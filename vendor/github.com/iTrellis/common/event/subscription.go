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

package event

import "fmt"

// Subscriber 消费者
type Subscriber interface {
	GetID() string
	Publish(values ...interface{})
	Stop()
}

// NewDefSubscriber 生成默认的消费者
func NewDefSubscriber(sub interface{}) (Subscriber, error) {
	var subscriber Subscriber
	switch s := sub.(type) {
	case func(...interface{}):
		subscriber = &defSubscriber{
			id: GenSubscriberID(),
			fn: s,
		}
	case Subscriber:
		subscriber = s
	default:
		return nil, fmt.Errorf("unkown subscriber type: %+v", s)
	}
	return subscriber, nil
}

// Subscriber is returned from the Subscribe function.
//
// This value and can be passed to Unsubscribe when the observer is no longer interested in receiving messages
type defSubscriber struct {
	id string
	fn func(values ...interface{})
}

// GetID return Subscriber's id
func (p *defSubscriber) GetID() string {
	return p.id
}

// Publish 发布信息
func (p *defSubscriber) Publish(values ...interface{}) {
	p.fn(values...)
}

// Stop do nothing
func (*defSubscriber) Stop() {}
