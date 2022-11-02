package config

func NewConfig() Config {
	return Config{
		VWAP: NewVWAP(),
		Feed: NewFeed(),
	}
}

type Config struct {
	VWAP
	Feed
}
