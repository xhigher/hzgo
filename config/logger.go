package config

type LoggerConfig struct {
	Filename string `yaml:"filename"`
	Level         string `yaml:"level"`
	Format        string `yaml:"format"`
	Prefix        string `yaml:"prefix"`
	Director      string `yaml:"director"`
	ShowLine      bool   `yaml:"showLine"`
	EncodeLevel   string `yaml:"encode-level"`
	StacktraceKey string `yaml:"stacktrace-key"`
	LogInConsole  bool   `yaml:"logger-in-console"`
}
