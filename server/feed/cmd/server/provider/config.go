package provider

// Config holds the feed service configuration.
type Config struct {
	Port        string `env:"FEED_SERVICE_PORT" envDefault:"50052"`
	DatabaseURL string `env:"DATABASE_URL,required"`
}
