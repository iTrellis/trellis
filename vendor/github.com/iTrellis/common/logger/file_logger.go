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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iTrellis/common/files"
)

const (
	// 默认的通道大小
	defaultChanBuffer int = 100000
)

type fileLogger struct {
	id string

	options FileOptions

	logChan  chan *Event
	stopChan chan bool

	writeFileTime time.Time
	lastMoveFlag  int

	ticker *time.Ticker

	prefixes []interface{}
}

// FileOptions file options
type FileOptions struct {
	level Level

	separator  string
	fileName   string
	maxLength  int64
	chanBuffer int

	moveFileType MoveFileType
	// 最大保留日志个数，如果为0则全部保留
	maxBackupFile int

	publisher Publisher
}

// MoveFileType move file type
type MoveFileType int

// MoveFileTypes
const (
	MoveFileTypeNothing   MoveFileType = iota // 不移动
	MoveFileTypePerMinite                     // 按分钟移动
	MoveFileTypeHourly                        // 按小时移动
	MoveFileTypeDaily                         // 按天移动
)

// NewFileLogger 标准窗体的输出对象
func NewFileLogger(opts ...FileOption) (Logger, error) {
	fw := &fileLogger{
		id:       uuid.NewString(),
		ticker:   time.NewTicker(time.Second * 30),
		stopChan: make(chan bool, 1),
	}

	err := fw.init(opts...)
	if err != nil {
		return nil, err
	}

	if fw.options.publisher != nil {
		_, err := fw.options.publisher.Subscriber(fw)
		if err != nil {
			return nil, err
		}
	}

	return fw, err
}

var fileExecutor = files.New()

func (p *fileLogger) init(opts ...FileOption) error {

	for _, o := range opts {
		o(&p.options)
	}

	if len(p.options.fileName) == 0 {
		return errors.New("file name not exist")
	}

	if p.options.chanBuffer == 0 {
		p.logChan = make(chan *Event, defaultChanBuffer)
	} else {
		p.logChan = make(chan *Event, p.options.chanBuffer)
	}

	if len(p.options.separator) == 0 {
		p.options.separator = "\t"
	}

	fi, err := fileExecutor.FileInfo(p.options.fileName)
	if err == nil {
		// 说明文件存在
		p.writeFileTime = fi.ModTime()
	} else {
		// 没有文件创建文件
		_, err = fileExecutor.WriteAppend(p.options.fileName, "")
		if err != nil {
			return err
		}
	}

	switch p.options.moveFileType {
	case MoveFileTypePerMinite:
		p.lastMoveFlag = p.writeFileTime.Minute()
	case MoveFileTypeHourly:
		p.lastMoveFlag = p.writeFileTime.Hour()
	case MoveFileTypeDaily:
		p.lastMoveFlag = p.writeFileTime.Day()
	}

	go p.looperLog()

	return nil
}

func (p *fileLogger) SetLevel(lvl Level) {
	p.options.level = lvl
}

func (p *fileLogger) write(bs []byte) (int, error) {
	return fileExecutor.WriteAppendBytes(p.options.fileName, bs)
}

func (p *fileLogger) Log(kvs ...interface{}) error {
	var ss []string
	for _, v := range kvs {
		ss = append(ss, toString(v))
	}
	_, err := p.write([]byte(strings.Join(ss, p.options.separator)))
	return err
}

func (p *fileLogger) looperLog() {
	for {
		select {
		case log := <-p.logChan:
			if log.Level >= p.options.level {
				_, _ = p.innerLog(log)
			}
		case t := <-p.ticker.C:
			flag := 0
			switch p.options.moveFileType {
			case MoveFileTypePerMinite:
				flag = t.Minute()
			case MoveFileTypeHourly:
				flag = t.Hour()
			case MoveFileTypeDaily:
				flag = t.Day()
			}
			if p.lastMoveFlag == flag {
				continue
			}
			_ = p.judgeMoveFile()
		case <-p.stopChan:
			p.ticker.Stop()
			return
		}
	}
}

func (p *fileLogger) judgeMoveFile() error {

	timeNow, flag := time.Now(), 0
	switch p.options.moveFileType {
	case MoveFileTypePerMinite:
		flag = timeNow.Minute()
	case MoveFileTypeHourly:
		flag = timeNow.Hour()
	case MoveFileTypeDaily:
		flag = time.Now().Day()
	default:
		return nil
	}

	if flag == p.lastMoveFlag {
		return nil
	}
	p.lastMoveFlag = flag
	p.writeFileTime = time.Now()
	return p.moveFile()
}

func (p *fileLogger) moveFile() error {
	var timeStr string
	switch p.options.moveFileType {
	case MoveFileTypePerMinite:
		timeStr = time.Now().Format("200601021504-05.999999999")
	case MoveFileTypeHourly:
		timeStr = time.Now().Format("2006010215-0405.999999999")
	case MoveFileTypeDaily:
		timeStr = time.Now().Format("20060102-150405.999999999")
	}

	err := fileExecutor.Rename(p.options.fileName, fmt.Sprintf("%s_%s", p.options.fileName, timeStr))
	if err != nil {
		return err
	}

	if err = p.removeOldFiles(); err != nil {
		return err
	}

	_, err = fileExecutor.Write(p.options.fileName, "")

	return err
}

func (p *fileLogger) removeOldFiles() error {
	if 0 == p.options.maxBackupFile {
		return nil
	}

	path := filepath.Dir(p.options.fileName)

	// 获取日志文件列表
	dirLis, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	// 根据文件名过滤日志文件
	fileSort := FileSort{}
	fileNameSplit := strings.Split(p.options.fileName, "/")
	filePrefix := fmt.Sprintf("%s_", fileNameSplit[len(fileNameSplit)-1])
	for _, f := range dirLis {
		if strings.Contains(f.Name(), filePrefix) {
			fileSort = append(fileSort, f)
		}
	}

	if len(fileSort) <= int(p.options.maxBackupFile) {
		return nil
	}

	// 根据文件修改日期排序，保留最近的N个文件
	sort.Sort(fileSort)
	for _, f := range fileSort[p.options.maxBackupFile:] {
		err := os.Remove(path + "/" + f.Name())
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *fileLogger) genLogs(evt *Event) string {
	var logs []string

	kvs := genLogs(evt)
	for i := 0; i < len(kvs); i += 2 {
		logs = append(logs, fmt.Sprintf("%s=%s", kvs[i], kvs[i+1]))
	}

	gEnd := "\n"
	switch runtime.GOOS {
	case "windows":
		gEnd += "\r\n"
	}

	return fmt.Sprintf("%s%s", strings.Join(logs, p.options.separator), gEnd)
}
func (p *fileLogger) innerLog(evt *Event) (n int, err error) {

	if err = p.judgeMoveFile(); err != nil {
		return
	}

	n, err = p.write([]byte(p.genLogs(evt)))

	if p.options.maxLength == 0 {
		return
	}

	fi, e := fileExecutor.FileInfo(p.options.fileName)
	if e != nil {
		return 0, e
	}

	if p.options.maxLength > fi.Size() {
		return
	}

	err = p.moveFile()

	return
}

func (p *fileLogger) Level() Level {
	return p.options.level
}

func (p *fileLogger) GetID() string {
	return p.id
}

func (p *fileLogger) Stop() {
	p.stopChan <- true

	p.options.publisher = nil

	close(p.logChan)
}

func (p *fileLogger) Publish(evts ...interface{}) error {
	for _, evt := range evts {
		switch t := evt.(type) {
		case Event:
			p.logChan <- &t
		case *Event:
			evt := *t
			p.logChan <- &evt
		case Level:
			p.options.level = t
		default:
			return fmt.Errorf("unsupported event type: %+v", reflect.TypeOf(evt))
		}
	}
	return nil
}

func (p *fileLogger) pubLog(level Level, kvs ...interface{}) {
	p.Publish(&Event{
		Time:   time.Now(),
		Level:  level,
		Fields: kvs,
	})
}

// Debug 调试
func (p *fileLogger) Debug(kvs ...interface{}) {
	p.pubLog(DebugLevel, kvs...)
}

// Debugf 调试
func (p *fileLogger) Debugf(msg string, kvs ...interface{}) {
	p.Debug("msg", fmt.Sprintf(msg, kvs...))
}

// Info 信息
func (p *fileLogger) Info(kvs ...interface{}) {
	p.pubLog(InfoLevel, kvs...)
}

// Infof 信息
func (p *fileLogger) Infof(msg string, kvs ...interface{}) {
	p.Info("msg", fmt.Sprintf(msg, kvs...))
}

// Warn 警告
func (p *fileLogger) Warn(kvs ...interface{}) {
	p.pubLog(WarnLevel, kvs...)
}

// Warnf 警告
func (p *fileLogger) Warnf(msg string, kvs ...interface{}) {
	p.Warn("msg", fmt.Sprintf(msg, kvs...))
}

// Error 错误
func (p *fileLogger) Error(kvs ...interface{}) {
	p.pubLog(ErrorLevel, kvs...)
}

// Errorf 错误
func (p *fileLogger) Errorf(msg string, kvs ...interface{}) {
	p.Error("msg", fmt.Sprintf(msg, kvs...))
}

// Critical 严重的
func (p *fileLogger) Critical(kvs ...interface{}) {
	p.pubLog(CriticalLevel, kvs...)
}

// Criticalf 严重的
func (p *fileLogger) Criticalf(msg string, kvs ...interface{}) {
	p.Critical("msg", fmt.Sprintf(msg, kvs...))
}

// Panic panic
func (p *fileLogger) Panic(kvs ...interface{}) {
	p.pubLog(PanicLevel, kvs...)
}

// Panicf panic
func (p *fileLogger) Panicf(msg string, kvs ...interface{}) {
	p.Panic("msg", fmt.Sprintf(msg, kvs...))
}
