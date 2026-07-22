package config

type assemblerConfig struct {
	MinBuildTimeSec int `yaml:"min_build_time_sec" env:"ASSEMBLER_MIN_BUILD_TIME_SEC" env-default:"5"`
	MaxBuildTimeSec int `yaml:"max_build_time_sec" env:"ASSEMBLER_MAX_BUILD_TIME_SEC" env-default:"15"`
}

func (c *assemblerConfig) MinSec() int {
	if c.MinBuildTimeSec <= 0 {
		return 5
	}
	return c.MinBuildTimeSec
}

func (c *assemblerConfig) MaxSec() int {
	if c.MaxBuildTimeSec < c.MinSec() {
		return c.MinSec()
	}
	return c.MaxBuildTimeSec
}
