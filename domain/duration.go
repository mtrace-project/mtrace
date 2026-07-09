package domain

import (
	"time"

	"go.yaml.in/yaml/v3"
)

type Duration time.Duration

func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}

	parsed, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*d = Duration(parsed)
	return nil
}

func (d *Duration) ToTimeDuration() *time.Duration {
	if d == nil {
		return nil
	}
	dur := time.Duration(*d)
	return &dur
}

func FromTimeDuration(d time.Duration) Duration {
	return Duration(d)
}
