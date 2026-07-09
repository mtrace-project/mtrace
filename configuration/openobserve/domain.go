package openobserve

import "github.com/spf13/viper"

type OpenObserveConfig struct {
	BaseURL    string `mapstructure:"base_url"`    // "http://localhost:5080"
	OrgName    string `mapstructure:"org_name"`    // "default"
	StreamName string `mapstructure:"stream_name"` // "default"
	Username   string `mapstructure:"username"`    // "admin@example.com"
	Password   string `mapstructure:"password"`    // "admin"
}

func NewOpenObserveConfig(v *viper.Viper) *OpenObserveConfig {
	return &OpenObserveConfig{
		BaseURL:    v.GetString("openobserve.base_url"),
		OrgName:    v.GetString("openobserve.org_name"),
		StreamName: v.GetString("openobserve.stream_name"),
		Username:   v.GetString("openobserve.username"),
		Password:   v.GetString("openobserve.password"),
	}
}
