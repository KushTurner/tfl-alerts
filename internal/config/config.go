package config

type TflConfig struct {
	AppId string
}

type AppConfig struct {
	TflConfig *TflConfig
}

func LoadConfig() (*AppConfig, error) {
	return &AppConfig{
		&TflConfig{AppId: ""},
	}, nil
}
