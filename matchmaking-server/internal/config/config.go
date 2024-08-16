package config

type Config struct {
	Redis Redis `yaml:"redis"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}
