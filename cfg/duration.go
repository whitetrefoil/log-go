package cfg

import "time"

type Duration struct {
	time.Duration
}

func (d *Duration) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
