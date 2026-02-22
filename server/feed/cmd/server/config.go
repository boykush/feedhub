package main

type config struct {
	Port string `env:"FEED_SERVICE_PORT" envDefault:"50052"`
}
