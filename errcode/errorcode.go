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

package errcode

import "github.com/go-trellis/common/errors"

const namespace = "trellis"

// 错误集合
var (
	ErrBadRequest      = errors.TN(namespace, 10, "bad request: {{.err}}")
	ErrAPINotFound     = errors.TN(namespace, 11, "api not found")
	ErrGetService      = errors.TN(namespace, 12, "{{.err}}")
	ErrGetServiceTopic = errors.TN(namespace, 13, "failed get service's topic")
	ErrCallService     = errors.TN(namespace, 14, "failed call service: {{.err}}")
)
