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

// NewCodec Takes in a connection/buffer and returns a new Codec
type NewCodec func() Codec

// Codec is a simple encoding interface used for the broker/transport
// where headers are not supported by the underlying implementation.
type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(bs []byte, v interface{}) error
	String() string
}
