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

// Option used by NewMessage
type Option func(*Options)

// Options parameters
type Options struct {
	Service *service.Service
	Topic   string
	Payload *BasePayload

	ContentType string
}

func ContentType(ct string) Option {
	return func(o *Options) {
		o.ContentType = ct
	}
}

func Topic(topic string) Option {
	return func(o *Options) {
		o.Topic = topic
	}
}

func MessagePayload(payload *BasePayload) Option {
	return func(o *Options) {
		o.Payload = payload
	}
}

func Service(s *service.Service) Option {
	return func(o *Options) {
		o.Service = s
	}
}
