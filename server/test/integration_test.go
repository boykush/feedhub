package test

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Disable Ryuk to avoid Docker Desktop connectivity issues on macOS.
	// Cleanup is handled by t.Cleanup instead.
	t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	ctx := context.Background()

	// Get project root (server/test -> project root is ../../)
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("failed to get project root: %v", err)
	}

	// 1. Start PostgreSQL container
	pgContainer, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("collector_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	t.Cleanup(func() { pgContainer.Terminate(ctx) })

	collectorDBURL, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	// Create feed_test database for feed service
	feedDBURL := createDatabase(t, collectorDBURL, "feed_test")

	// 2. Run Atlas migrations
	collectorMigrationsDir := filepath.Join(projectRoot, "server/collector/migrations")
	atlasCmd := exec.Command("mise", "exec", "--", "atlas", "migrate", "apply",
		"--dir", fmt.Sprintf("file://%s", collectorMigrationsDir),
		"--url", collectorDBURL,
	)
	atlasCmd.Dir = projectRoot
	atlasCmd.Stdout = os.Stdout
	atlasCmd.Stderr = os.Stderr
	if err := atlasCmd.Run(); err != nil {
		t.Fatalf("failed to run collector atlas migrations: %v", err)
	}

	feedMigrationsDir := filepath.Join(projectRoot, "server/feed/migrations")
	feedAtlasCmd := exec.Command("mise", "exec", "--", "atlas", "migrate", "apply",
		"--dir", fmt.Sprintf("file://%s", feedMigrationsDir),
		"--url", feedDBURL,
	)
	feedAtlasCmd.Dir = projectRoot
	feedAtlasCmd.Stdout = os.Stdout
	feedAtlasCmd.Stderr = os.Stderr
	if err := feedAtlasCmd.Run(); err != nil {
		t.Fatalf("failed to run feed atlas migrations: %v", err)
	}

	// 3. Start services using free ports to avoid conflicts
	feedPort := freePort(t)
	collectorPort := freePort(t)
	bffPort := freePort(t)

	// Start feed service
	feedCmd := exec.Command(filepath.Join(projectRoot, "server/feed/bin/server"))
	feedCmd.Env = append(os.Environ(),
		fmt.Sprintf("FEED_SERVICE_PORT=%s", feedPort),
		fmt.Sprintf("DATABASE_URL=%s", feedDBURL),
	)
	feedCmd.Stdout = os.Stdout
	feedCmd.Stderr = os.Stderr
	if err := feedCmd.Start(); err != nil {
		t.Fatalf("failed to start feed service: %v", err)
	}
	t.Cleanup(func() { feedCmd.Process.Kill(); feedCmd.Wait() })

	// Start collector service
	collectorCmd := exec.Command(filepath.Join(projectRoot, "server/collector/bin/server"))
	collectorCmd.Env = append(os.Environ(),
		fmt.Sprintf("COLLECTOR_SERVICE_PORT=%s", collectorPort),
		fmt.Sprintf("DATABASE_URL=%s", collectorDBURL),
	)
	collectorCmd.Stdout = os.Stdout
	collectorCmd.Stderr = os.Stderr
	if err := collectorCmd.Start(); err != nil {
		t.Fatalf("failed to start collector service: %v", err)
	}
	t.Cleanup(func() { collectorCmd.Process.Kill(); collectorCmd.Wait() })

	// Start BFF service
	bffCmd := exec.Command(filepath.Join(projectRoot, "server/bff/bin/server"))
	bffCmd.Env = append(os.Environ(),
		fmt.Sprintf("BFF_HTTP_PORT=%s", bffPort),
		fmt.Sprintf("FEED_SERVICE_ADDR=localhost:%s", feedPort),
		fmt.Sprintf("COLLECTOR_SERVICE_ADDR=localhost:%s", collectorPort),
	)
	bffCmd.Stdout = os.Stdout
	bffCmd.Stderr = os.Stderr
	if err := bffCmd.Start(); err != nil {
		t.Fatalf("failed to start bff service: %v", err)
	}
	t.Cleanup(func() { bffCmd.Process.Kill(); bffCmd.Wait() })

	// 4. Wait for BFF health
	healthURL := fmt.Sprintf("http://localhost:%s/api/v1/feeds/health", bffPort)
	waitForHealth(t, healthURL, 30*time.Second)

	// 5. Run hurl tests
	hurlCmd := exec.Command("mise", "exec", "--", "hurl",
		"--variable", fmt.Sprintf("bff_port=%s", bffPort),
		"--test",
		filepath.Join(projectRoot, "server/test/collector_health.hurl"),
		filepath.Join(projectRoot, "server/test/collector_operations.hurl"),
		filepath.Join(projectRoot, "server/test/feed_health.hurl"),
		filepath.Join(projectRoot, "server/test/feed_operations.hurl"),
		filepath.Join(projectRoot, "server/test/feed_list.hurl"),
	)
	hurlCmd.Dir = projectRoot
	hurlCmd.Stdout = os.Stdout
	hurlCmd.Stderr = os.Stderr
	if err := hurlCmd.Run(); err != nil {
		t.Fatalf("hurl tests failed: %v", err)
	}
}

func createDatabase(t *testing.T, baseURL, dbName string) string {
	t.Helper()

	db, err := sql.Open("postgres", baseURL)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		t.Fatalf("failed to create database %s: %v", dbName, err)
	}

	return replaceDBName(baseURL, dbName)
}

func replaceDBName(connStr, newDB string) string {
	// Connection string format: postgres://user:pass@host:port/dbname?params
	qIdx := strings.Index(connStr, "?")
	query := ""
	base := connStr
	if qIdx >= 0 {
		query = connStr[qIdx:]
		base = connStr[:qIdx]
	}
	lastSlash := strings.LastIndex(base, "/")
	if lastSlash < 0 {
		return connStr
	}
	return base[:lastSlash+1] + newDB + query
}

func freePort(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return fmt.Sprintf("%d", port)
}

func waitForHealth(t *testing.T, url string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("health check at %s did not pass within %v", url, timeout)
}
