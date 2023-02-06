package config

type JWTConfig struct {
	Realm          string `mapstructure:"realm" json:"realm" yaml:"realm"`
	Issuer         string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`
	SecretKey      string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
	Timeout        int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	MaxRefreshTime int    `mapstructure:"max-refresh-time" json:"max-refresh-time" yaml:"max-refresh-time"`
}

type SecConfig struct {
	SignSecret      string `mapstructure:"sign-secret" json:"sign-secret" yaml:"sign-secret"`
	SignKeyName string `mapstructure:"sign-key-name" json:"sign-key-name" yaml:"sign-key-name"`
}