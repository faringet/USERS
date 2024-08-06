package config

type Config struct {
	Postgre  Postgre `mapstructure:"POSTGRE"`
	Logger   Logger  `mapstructure:"LOGGER"`
	LocalURL string  `mapstructure:"THIS_APP_URL"`
}

type Logger struct {
	Production  string `mapstructure:"PRODUCTION"`
	Development string `mapstructure:"DEVELOPMENT"`
}

type Postgre struct {
	DbURL string `mapstructure:"DB_URL"`
}
