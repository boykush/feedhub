package main

type config struct {
	Port        string `env:"COLLECTOR_SERVICE_PORT" envDefault:"50053"`
	DatabaseURL string `env:"DATABASE_URL,required"`
}
