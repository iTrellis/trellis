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

package message

import "github.com/iTrellis/trellis/service"

// Message is the interface for publishing asynchronously
type Message interface {
	Service() *service.Service
	Topic() string
	Payload() *BasePayload
	ContentType() string
}

// Payload payload between services
type Payload struct {
	ID       string
	Target   string
	Method   string
	Endpoint string
	Error    string

	BasePayload
}

// BasePayload payload between services
type BasePayload struct {
	// The values read from the socket
	Header map[string]string
	Body   []byte
}
