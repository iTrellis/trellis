package configure

import (
	"github.com/iTrellis/common/logger"
	"github.com/iTrellis/config"
)

type Logger struct {
	Type    string         `json:"type" yaml:"type"` // || std | file | logrus
	Level   logger.Level   `json:"level" yaml:"level"`
	Options config.Options `json:"options" yaml:"options"`
}
