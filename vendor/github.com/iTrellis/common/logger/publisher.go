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

package logger

import (
	"fmt"
	"time"

	"github.com/iTrellis/common/event"
)

// Publisher publish some informations
type Publisher interface {
	Log(Level, ...interface{})
	Logf(Level, string, ...interface{})

	WithPrefix(kvs ...interface{}) Publisher

	SetLevel(lvl Level)

	event.SubscriberGroup
}

// NewPublisher new a publisher
func NewPublisher() Publisher {
	return &publisher{
		eventGroup: event.NewSubscriberGroup(),
	}
}

type publisher struct {
	prefixes   []interface{}
	hasCaller  bool
	eventGroup event.SubscriberGroup
}

// SetLevel set looger's level
func (p *publisher) SetLevel(lvl Level) {
	p.Publish(lvl)
}

// Log 打印
func (p *publisher) Log(lvl Level, kvs ...interface{}) {
	p.publishLog(lvl, kvs...)
}

func (p *publisher) Logf(lvl Level, msg string, kvs ...interface{}) {
	p.publishLog(lvl, "msg", fmt.Sprintf(msg, kvs...))
}

// WithPrefix 加载前缀
func (p *publisher) WithPrefix(prefixes ...interface{}) Publisher {
	if len(prefixes) == 0 {
		return p
	}
	newPrefixes := append(prefixes, p.prefixes...)
	return &publisher{
		prefixes:   newPrefixes,
		hasCaller:  p.hasCaller || containsCaller(newPrefixes),
		eventGroup: p.eventGroup,
	}
}

func (p *publisher) publishLog(lvl Level, kvs ...interface{}) {
	if len(kvs) == 0 {
		return
	}
	fields := append(p.prefixes, kvs...)
	if p.hasCaller {
		if len(fields) != 0 {
			bindCallers(fields)
		}
	}
	p.Publish(&Event{
		Time:   time.Now(),
		Level:  lvl,
		Fields: fields,
	})
}

func (p *publisher) Subscriber(s interface{}) (event.Subscriber, error) {
	return p.eventGroup.Subscriber(s)
}

func (p *publisher) RemoveSubscriber(ids ...string) error {
	return p.eventGroup.RemoveSubscriber(ids...)
}

func (p *publisher) Publish(values ...interface{}) error {
	return p.eventGroup.Publish(values...)
}

func (p *publisher) ClearSubscribers() {
	p.eventGroup.ClearSubscribers()
}
