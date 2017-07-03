package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	pb "github.com/albert-widi/logvault/cmd/logee/pb"
	"google.golang.org/grpc"
)

var (
	// flag list
	hostFlag  = flag.String("host", "", "host to locate logee service")
	groupFlag = flag.String("group", "", "group name of agent")
	fileFlag  = flag.String("files", "", "file lists")
	dirFlag   = flag.String("dir", "", "directory of log file")
	delimiter = ","
	// prefix for writer
	prefix     string
	workerList []*worker
	client     pb.LogeeClient
)

func getFileList(list, dir string) ([]string, error) {
	if dir == "" {
		return nil, errors.New("No directory detected")
	}
	// make sure dir have / char
	if dir[len(dir)-1:] != "/" {
		dir += "/"
	}

	var files []string
	if list != "" {
		fileList := strings.Split(list, delimiter)
		for _, val := range fileList {
			files = append(files, dir+val)
		}
	} else {
		if err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if !f.IsDir() {
				files = append(files, path)
			}
			return err
		}); err != nil {
			return nil, err
		}
	}
	return files, nil
}

func main() {
	flag.Parse()
	prefix = *groupFlag
	if prefix == "" {
		log.Fatal("Group is empty, please set the host group")
	}
	list := *fileFlag
	dir := *dirFlag
	host := *hostFlag

	if host != "" {
		conn, err := grpc.Dial(host, grpc.WithInsecure())
		if err != nil {
			log.Fatal("Cannot connect to gRPC host ", err.Error())
		}
		client = pb.NewLogeeClient(conn)
	}

	files, err := getFileList(list, dir)
	if err != nil {
		log.Fatal("Cannot get file list ", err.Error())
	}
	workerList = make([]*worker, len(files))
	iterator := 0
	for key := range files {
		w := newWorker(files[key])
		workerList[iterator] = w
		log.Printf("Tailing log for %s \n", w.fileName)
		go w.run(context.TODO())
		iterator++
	}

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Println("Agent exited")
	}
}

type worker struct {
	fileName string
	cmd      *exec.Cmd
}

func newWorker(fileName string) *worker {
	w := &worker{
		fileName: fileName,
		cmd:      exec.CommandContext(context.TODO(), "tail", "-f", fileName),
	}
	w.cmd.Stdout = newWriter(prefix)
	w.cmd.Stderr = newWriter(prefix)
	return w
}

func (w *worker) run(ctx context.Context) {
	if err := w.cmd.Run(); err != nil {
		log.Fatal("Failed to run command ", err.Error())
		return
	}
}

type writer struct {
	prefix string
}

func newWriter(prefix string) *writer {
	w := new(writer)
	w.prefix = prefix
	return w
}

// redirected writer
func (w *writer) Write(b []byte) (int, error) {
	fmt.Println(string(b))
	// TODO: using worker to pipe into buffered channel and queue, smaller point of failure
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		_, err := client.PushLog(ctx, &pb.PushRequest{Log: string(b), Prefix: w.prefix})
		if err != nil {
			log.Println("Failed to push to logee service ", err.Error())
		}
	}
	return len(b), nil
}
