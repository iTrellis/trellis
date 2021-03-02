package configure

import (
	"github.com/iTrellis/common/logger"
	"github.com/iTrellis/config"
)

var (
	DefaultLevel = logger.InfoLevel

	DefaultStackSkip = 6
)

type Logger struct {
	Type      string         `json:"type" yaml:"type"` // || std | file | logrus
	Level     *logger.Level  `json:"level" yaml:"level"`
	StackSkip *int           `json:"stack_skip" yaml:"stack_skip"`
	Options   config.Options `json:"options" yaml:"options"`
}
