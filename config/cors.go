package config

type CorsConfig struct {
	AllowOrigins []string `mapstructure:"allow-origins" json:"allow-origins" yaml:"allow-origins"`
	AllowMethods []string `mapstructure:"allow-methods" json:"allow-methods" yaml:"allow-methods"`
	AllowHeaders []string `mapstructure:"allow-headers" json:"allow-headers" yaml:"allow-headers"`
}
