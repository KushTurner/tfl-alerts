package config

type TflConfig struct {
	Url string
}

type AppConfig struct {
	TflConfig *TflConfig
}

func LoadConfig() (*AppConfig, error) {
	return &AppConfig{
		&TflConfig{Url: "https://api.tfl.gov.uk"},
	}, nil
}
