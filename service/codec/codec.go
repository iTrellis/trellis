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

package codec

// Codec
type Codec interface {
	Reader
	Writer
	Close() error
	String() string
}

// Reader
type Reader interface {
	ReadHeader(*Payload) error
	ReadBody(interface{}) error
}

// Writer
type Writer interface {
	Write(*Payload, interface{}) error
}

// Marshaler is a simple encoding interface used for the broker/transport
// where headers are not supported by the underlying implementation.
type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	String() string
}

// Payload payload between services
type Payload struct {
	// The values read from the socket
	Header map[string]string
	Body   []byte
}
