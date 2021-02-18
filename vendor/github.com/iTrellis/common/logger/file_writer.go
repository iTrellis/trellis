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
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/iTrellis/common/event"
	"github.com/iTrellis/common/files"
)

type fileWriter struct {
	logger Logger

	opts fileWriterOptions

	// filePath  string

	stopChan chan bool
	logChan  chan *Event

	subscriber event.Subscriber
}

type fileWriterOptions struct {
	level Level

	goNum int

	separator  string
	fileName   string
	maxLength  int64
	chanBuffer int

	moveFileType MoveFileType
	// 最大保留日志个数，如果为0则全部保留
	maxBackupFile int
}

type routingFileWriter struct {
	locker        sync.Mutex
	opts          fileWriterOptions
	fileName      string
	writeFileTime time.Time
	lastMoveFlag  int
	ticker        *time.Ticker
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

// OptionFileWriter 操作配置函数
type OptionFileWriter func(*fileWriterOptions)

// FileWiterLevel 设置等级
func FileWiterLevel(lvl Level) OptionFileWriter {
	return func(f *fileWriterOptions) {
		f.level = lvl
	}
}

// FileWiterRoutings 设置Gouting数量
func FileWiterRoutings(num int) OptionFileWriter {
	return func(f *fileWriterOptions) {
		if num < 1 {
			num = 1
		}
		f.goNum = num
	}
}

// FileWiterBuffer 设置Chan的大小
func FileWiterBuffer(buffer int) OptionFileWriter {
	return func(f *fileWriterOptions) {
		f.chanBuffer = buffer
	}
}

// FileWiterSeparator 设置打印分隔符
func FileWiterSeparator(separator string) OptionFileWriter {
	return func(f *fileWriterOptions) {
		f.separator = separator
	}
}

// FileWiterFileName 设置文件名
func FileWiterFileName(name string) OptionFileWriter {
	return func(f *fileWriterOptions) {
		f.fileName = name
	}
}

// FileWiterMaxLength 设置最大文件大小
func FileWiterMaxLength(length int64) OptionFileWriter {
	return func(f *fileWriterOptions) {
		f.maxLength = length
	}
}

// FileWiterMaxBackupFile 文件最大数量
func FileWiterMaxBackupFile(num int) OptionFileWriter {
	return func(f *fileWriterOptions) {
		f.maxBackupFile = num
	}
}

// FileWiterMoveFileType 设置移动文件的类型
func FileWiterMoveFileType(typ MoveFileType) OptionFileWriter {
	return func(f *fileWriterOptions) {
		f.moveFileType = typ
	}
}

// FileWriter 标准窗体的输出对象
func FileWriter(log Logger, opts ...OptionFileWriter) (Writer, error) {
	fw := &fileWriter{
		logger:   log,
		stopChan: make(chan bool, 1),
	}

	err := fw.init(opts...)
	if err != nil {
		return nil, err
	}

	fw.subscriber, err = event.NewDefSubscriber(fw.Publish)
	if err != nil {
		fw.Stop()
		return nil, err
	}

	_, err = log.Subscriber(fw.subscriber)
	if err != nil {
		fw.stopChan <- true
		return nil, err
	}
	return fw, err
}

var fileExecutor = files.New()

func (p *fileWriter) init(opts ...OptionFileWriter) error {

	for _, o := range opts {
		o(&p.opts)
	}

	if len(p.opts.fileName) == 0 {
		return errors.New("file name not exist")
	}

	if p.opts.chanBuffer == 0 {
		p.logChan = make(chan *Event, defaultChanBuffer)
	} else {
		p.logChan = make(chan *Event, p.opts.chanBuffer)
	}

	if len(p.opts.separator) == 0 {
		p.opts.separator = "\t"
	}

	for i := 0; i < p.opts.goNum; i++ {
		rfw := routingFileWriter{
			opts:          p.opts,
			writeFileTime: time.Now(),
			ticker:        time.NewTicker(time.Second * 30),
		}

		if p.opts.goNum == 1 {
			rfw.fileName = p.opts.fileName
		} else {
			rfw.fileName = fmt.Sprintf("%s.%d", p.opts.fileName, i)
		}

		fi, err := fileExecutor.FileInfo(rfw.fileName)
		if err == nil {
			// 说明文件存在
			rfw.writeFileTime = fi.ModTime()
		} else {
			// 没有文件创建文件
			_, err = fileExecutor.WriteAppend(rfw.fileName, "")
			if err != nil {
				return err
			}
		}

		switch p.opts.moveFileType {
		case MoveFileTypePerMinite:
			rfw.lastMoveFlag = rfw.writeFileTime.Minute()
		case MoveFileTypeHourly:
			rfw.lastMoveFlag = rfw.writeFileTime.Hour()
		case MoveFileTypeDaily:
			rfw.lastMoveFlag = rfw.writeFileTime.Day()
		}

		rfw.looperLog(p)
	}

	return nil
}

func (p *fileWriter) Publish(evts ...interface{}) {
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

func (p *fileWriter) Write(bs []byte) (int, error) {
	return fileExecutor.WriteAppendBytes(p.opts.fileName, bs)
}

func (p *routingFileWriter) looperLog(fw *fileWriter) {
	go func() {
		for {
			select {
			case log := <-fw.logChan:
				if log.Level >= p.opts.level {
					_, _ = p.innerLog(log)
				}
			case t := <-p.ticker.C:
				flag := 0
				switch p.opts.moveFileType {
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
				p.locker.Lock()
				_ = p.judgeMoveFile()
				p.locker.Unlock()
			case <-fw.stopChan:
				return
			}
		}
	}()
}

func (p *routingFileWriter) judgeMoveFile() error {

	timeNow, flag := time.Now(), 0
	switch p.opts.moveFileType {
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

func (p *routingFileWriter) moveFile() error {
	var timeStr string
	switch p.opts.moveFileType {
	case MoveFileTypePerMinite:
		timeStr = time.Now().Format("200601021504-05.999999999")
	case MoveFileTypeHourly:
		timeStr = time.Now().Format("2006010215-0405.999999999")
	case MoveFileTypeDaily:
		timeStr = time.Now().Format("20060102-150405.999999999")
	}

	err := fileExecutor.Rename(p.fileName, fmt.Sprintf("%s_%s", p.fileName, timeStr))
	if err != nil {
		return err
	}

	if err = p.removeOldFiles(); err != nil {
		return err
	}

	_, err = fileExecutor.Write(p.fileName, "")

	return err
}

func (p *routingFileWriter) removeOldFiles() error {
	if 0 == p.opts.maxBackupFile {
		return nil
	}

	path := filepath.Dir(p.fileName)

	// 获取日志文件列表
	dirLis, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	// 根据文件名过滤日志文件
	fileSort := FileSort{}
	fileNameSplit := strings.Split(p.fileName, "/")
	filePrefix := fmt.Sprintf("%s_", fileNameSplit[len(fileNameSplit)-1])
	for _, f := range dirLis {
		if strings.Contains(f.Name(), filePrefix) {
			fileSort = append(fileSort, f)
		}
	}

	if len(fileSort) <= int(p.opts.maxBackupFile) {
		return nil
	}

	// 根据文件修改日期排序，保留最近的N个文件
	sort.Sort(fileSort)
	for _, f := range fileSort[p.opts.maxBackupFile:] {
		err := os.Remove(path + "/" + f.Name())
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *routingFileWriter) innerLog(evt *Event) (n int, err error) {
	p.locker.Lock()
	defer p.locker.Unlock()

	if err = p.judgeMoveFile(); err != nil {
		return
	}

	n, err = p.Write([]byte(generateLogs(evt, p.opts.separator)))

	if p.opts.maxLength == 0 {
		return
	}

	fi, e := fileExecutor.FileInfo(p.fileName)
	if e != nil {
		return 0, e
	}

	if p.opts.maxLength > fi.Size() {
		return
	}

	err = p.moveFile()

	return
}

func (p *routingFileWriter) Write(bs []byte) (int, error) {
	return fileExecutor.WriteAppendBytes(p.fileName, bs)
}

func (p *fileWriter) Level() Level {
	return p.opts.level
}

func (p *fileWriter) GetID() string {
	return p.subscriber.GetID()
}

func (p *fileWriter) Stop() {
	if err := p.logger.RemoveSubscriber(p.subscriber.GetID()); err != nil {
		p.logger.Criticalf("failed remove Chan Writer: %s", err.Error())
	}
	p.stopChan <- true
}

// FileSort 文件排序
type FileSort []os.FileInfo

func (fs FileSort) Len() int {
	return len(fs)
}

func (fs FileSort) Less(i, j int) bool {
	return fs[i].Name() > fs[j].Name()
}

func (fs FileSort) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}
