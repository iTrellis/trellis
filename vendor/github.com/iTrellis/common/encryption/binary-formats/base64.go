/*
Copyright © 2016 Henry Huang <hhh@rutcode.com>

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

package bsf

import (
	"encoding/base64"
	"fmt"
	"sync"
)

// default base64 encoders
const (
	EncodeStd    = "trellis::algo::encodeStd"
	EncodeRawStd = "trellis::algo::encodeRawStd"
	EncodeURL    = "trellis::algo::encodeURL"
	EncodeRawURL = "trellis::algo::encodeRawURL"
)

var defaultBase *base

type base struct {
	mapEncoders map[string]*base64.Encoding
	locker      sync.RWMutex
}

func (p *base) getEncoding(key string) (*base64.Encoding, bool) {
	p.locker.RLock()
	encoding, ok := p.mapEncoders[key]
	p.locker.RUnlock()
	return encoding, ok
}

func (p *base) setEncoding(key string, encoding *base64.Encoding) {
	p.locker.Lock()
	p.mapEncoders[key] = encoding
	p.locker.Unlock()
}

func init() {
	defaultBase = &base{
		mapEncoders: map[string]*base64.Encoding{
			EncodeStd:    base64.StdEncoding,
			EncodeRawStd: base64.RawStdEncoding,
			EncodeURL:    base64.URLEncoding,
			EncodeRawURL: base64.RawURLEncoding,
		},
	}
}

// NewEncoding get base64 encoding with input encoder
func NewEncoding(encoder string) *base64.Encoding {
	encoding, ok := defaultBase.getEncoding(encoder)
	if ok {
		return encoding
	}
	if len(encoder) != 64 {
		return nil
	}

	encoding = base64.NewEncoding(encoder)
	defaultBase.setEncoding(encoder, encoding)
	return encoding
}

// NewEncodingWithPadding get encoding with encoder and padding
func NewEncodingWithPadding(encoder string, padding rune) *base64.Encoding {
	key := fmt.Sprintf("%s::%d", encoder, padding)
	encoding, ok := defaultBase.getEncoding(key)
	if ok {
		return encoding
	}
	if len(encoder) != 64 {
		return nil
	}

	encoding = base64.NewEncoding(encoder).WithPadding(padding)
	defaultBase.setEncoding(key, encoding)
	return encoding
}

// Encode encode bytes with encoder
func Encode(encoder string, src []byte) string {
	encoding, ok := defaultBase.getEncoding(encoder)
	if ok {
		return encoding.EncodeToString(src)
	}
	return ""
}

// EncodeString encode string with encoder
func EncodeString(encoder string, src string) string {
	return Encode(encoder, []byte(src))
}

// Decode decode bytes with encoder
func Decode(encoder string, src []byte) ([]byte, error) {
	return DecodeString(encoder, string(src))
}

// DecodeString decode string with encoder
func DecodeString(encoder string, s string) ([]byte, error) {
	if encoding, ok := defaultBase.getEncoding(encoder); ok {
		bs, err := encoding.DecodeString(s)
		if err != nil {
			return nil, err
		}
		return bs, nil
	}
	return nil, nil
}
