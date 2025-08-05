package config

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	JWT struct {
		Secret string
	}
}

func Load() *Config {
	cfg := &Config{}

	cfg.DB.Host = "localhost"
	cfg.DB.Port = "5432"
	cfg.DB.User = "gorent"
	cfg.DB.Password = "gorentpass"
	cfg.DB.Name = "gorent"
	cfg.DB.SSLMode = "disable"

	cfg.JWT.Secret = "jwt-key"

	return cfg
}
