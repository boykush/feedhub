package main

type config struct {
	HTTPPort             string `env:"BFF_HTTP_PORT" envDefault:"8080"`
	FeedServiceAddr      string `env:"FEED_SERVICE_ADDR" envDefault:"feed-service:50052"`
	CollectorServiceAddr string `env:"COLLECTOR_SERVICE_ADDR" envDefault:"collector-service:50053"`
}
