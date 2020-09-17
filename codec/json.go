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

import (
	"encoding/json"
)

type jsonCodec struct{}

func (*jsonCodec) Unmarshal(bytes []byte, obj interface{}) error {
	return json.Unmarshal(bytes, obj)
}

func (*jsonCodec) Marshal(body interface{}) ([]byte, error) {
	return json.Marshal(body)
}

func (*jsonCodec) String() string {
	return JSON
}

func newJSONCodec() (Codec, error) {
	return (*jsonCodec)(nil), nil
}
