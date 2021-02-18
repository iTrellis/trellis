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
	"io"
	"os"
	"reflect"

	"github.com/iTrellis/common/event"
)

type chanWriter struct {
	logger   Logger
	stopChan chan bool
	logChan  chan *Event
	out      io.Writer

	subscriber event.Subscriber

	level     Level
	separator string
	buffer    int
}

// OptionChanWriter 操作配置函数
type OptionChanWriter func(*chanWriter)

// ChanWiterLevel 设置等级
func ChanWiterLevel(lvl Level) OptionChanWriter {
	return func(c *chanWriter) {
		c.level = lvl
	}
}

// ChanWiterBuffer 设置Chan的大小
func ChanWiterBuffer(buffer int) OptionChanWriter {
	return func(c *chanWriter) {
		c.buffer = buffer
	}
}

// ChanWiterSeparator 设置打印分隔符
func ChanWiterSeparator(separator string) OptionChanWriter {
	return func(c *chanWriter) {
		c.separator = separator
	}
}

// ChanWriter 标准窗体的输出对象
func ChanWriter(log Logger, opts ...OptionChanWriter) (Writer, error) {
	c := &chanWriter{
		logger:   log,
		out:      os.Stdout,
		stopChan: make(chan bool, 1),
	}
	c.init(opts...)

	c.looperLog()

	var err error
	c.subscriber, err = event.NewDefSubscriber(c.Publish)
	if err != nil {
		return nil, err
	}

	_, err = log.Subscriber(c.subscriber)
	if err != nil {
		c.stopChan <- true
		return nil, err
	}

	return c, nil
}

func (p *chanWriter) init(opts ...OptionChanWriter) {

	for _, o := range opts {
		o(p)
	}

	if p.buffer == 0 {
		p.logChan = make(chan *Event, defaultChanBuffer)
	} else {
		p.logChan = make(chan *Event, p.buffer)
	}

	if len(p.separator) == 0 {
		p.separator = "\t"
	}
}

func (p *chanWriter) Publish(evts ...interface{}) {
	for _, evt := range evts {
		switch eType := evt.(type) {
		case Event:
			p.logChan <- &eType
		case *Event:
			p.logChan <- eType
		default:
			panic(fmt.Errorf("unsupported event type: %s", reflect.TypeOf(evt).Name()))
		}
	}
}

func (p *chanWriter) looperLog() {
	go func() {
		for {
			select {
			case log := <-p.logChan:
				if log.Level < p.Level() {
					continue
				}
				data := []byte(generateLogs(log, p.separator))

				if color := LevelColors[log.Level]; len(color) != 0 {
					data = []byte(fmt.Sprintf(color, string(data)))
				}
				_, _ = p.Write(data)
			case <-p.stopChan:
				return
			}
		}
	}()
}

func (p chanWriter) Write(bs []byte) (int, error) {
	return p.out.Write(bs)
}

func (p *chanWriter) Level() Level {
	return p.level
}

func (p *chanWriter) GetID() string {
	return p.subscriber.GetID()
}

func (p *chanWriter) Stop() {
	if err := p.logger.RemoveSubscriber(p.subscriber.GetID()); err != nil {
		p.logger.Criticalf("failed remove Chan Writer: %s", err.Error())
	}
	p.stopChan <- true
}
