/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

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
