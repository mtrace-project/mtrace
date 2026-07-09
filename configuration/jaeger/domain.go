package jaeger

import "github.com/spf13/viper"

type JaegerConfig struct {
	BaseURL string `mapstructure:"base_url"` // "http://localhost:16686"
}

func NewJaegerConfig(v *viper.Viper) *JaegerConfig {
	return &JaegerConfig{
		BaseURL: v.GetString("jaeger.base_url"),
	}
}
