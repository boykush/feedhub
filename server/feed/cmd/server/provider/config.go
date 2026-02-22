package provider

// Config holds the feed service configuration.
type Config struct {
	Port string `env:"FEED_SERVICE_PORT" envDefault:"50052"`
}
