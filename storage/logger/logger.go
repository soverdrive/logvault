package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// FileLog of logger
type FileLog struct {
	logs map[string]*logger
	lock sync.Mutex
}

type logger struct {
	prefix   string
	filename string
	Log      *log.Logger
	file     *os.File
}

// NewFileLogger to create new logger object
func NewFileLogger() *FileLog {
	fLog := &FileLog{
		logs: make(map[string]*logger),
	}
	return fLog
}

func (flog *FileLog) createNewLogger(prefix string) (*logger, error) {
	flog.lock.Lock()
	defer flog.lock.Unlock()
	if l, ok := flog.logs[prefix]; ok {
		return l, nil
	}
	fileName := prefix + ".log"
	f, err := os.OpenFile(prefix+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	lg := log.New(f, "", 0)
	l := &logger{
		prefix:   prefix,
		filename: fileName,
		file:     f,
		Log:      lg,
	}
	flog.logs[prefix] = l
	return l, nil
}

// WriteLog for writing log
func (flog *FileLog) WriteLog(prefix, hostname, content string) error {
	var err error
	l, ok := flog.logs[prefix]
	if !ok {
		l, err = flog.createNewLogger(prefix)
		if err != nil {
			return err
		}
	}
	logContent := fmt.Sprintf("%s::%s", hostname, content)
	fmt.Printf("Write to FileLog: %s", logContent)
	l.Log.Print(logContent)
	return nil
}

// Close all files
func (flog *FileLog) Close() error {
	for _, val := range flog.logs {
		val.file.Close()
	}
	return nil
}
