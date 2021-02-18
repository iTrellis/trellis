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
	"time"

	"github.com/iTrellis/config"
	"github.com/iTrellis/trellis/service"
)

type Registry struct {
	Name string               `json:"name" yaml:"name"`
	Type service.RegisterType `json:"type" yaml:"type"`

	Address []string      `json:"address" yaml:"address"`
	Secure  bool          `json:"secure" yaml:"secure"`
	Timeout time.Duration `json:"timeout" yaml:"timeout"`

	Watchers []Watcher `json:"watchers" yaml:"watchers"`
}

type Watcher struct {
	service.Service `json:",inline" yaml:",inline"`

	Options config.Options `json:"options" yaml:"options"`
}
