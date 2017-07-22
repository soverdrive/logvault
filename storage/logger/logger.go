package logger

import (
	"fmt"
	"log"
	"os"
	"sync"

	pb "github.com/albert-widi/logvault/pb"
)

// FileLog of logger
type FileLog struct {
	logs map[string]*logger
	lock sync.Mutex
	dir  string
}

type logger struct {
	prefix   string
	filename string
	Log      *log.Logger
	file     *os.File
}

// NewFileLogger to create new logger object
func NewFileLogger(dir string) *FileLog {
	if dir != "" && dir[len(dir)-1:] != "/" {
		dir += "/"
	}
	fLog := &FileLog{
		logs: make(map[string]*logger),
		dir:  dir,
	}
	return fLog
}

func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func createDir(dir string) error {
	dirExists, err := dirExists(dir)
	if err != nil {
		return err
	}
	if !dirExists {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

func (flog *FileLog) createNewLogger(group, fn string) (*logger, error) {
	flog.lock.Lock()
	defer flog.lock.Unlock()

	group = fmt.Sprintf("%s%s", flog.dir, group)
	err := createDir(group)
	if err != nil {
		return nil, err
	}

	fileName := fn
	if group != "" {
		fileName = group + "/" + fileName
	}

	if l, ok := flog.logs[fileName]; ok {
		return l, nil
	}

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	lg := log.New(f, "", 0)
	l := &logger{
		prefix:   group,
		filename: fn,
		file:     f,
		Log:      lg,
	}
	flog.logs[fileName] = l
	return l, nil
}

// WriteLog for writing log
func (flog *FileLog) WriteLog(req *pb.IngestRequest) error {
	var err error
	l, ok := flog.logs[req.GetFilename()]
	if !ok {
		l, err = flog.createNewLogger(req.GetPrefix(), req.GetFilename())
		if err != nil {
			return err
		}
	}
	logContent := fmt.Sprintf("%s::%s", req.GetHostname(), req.GetLog())
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
