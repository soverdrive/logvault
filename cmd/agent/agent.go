package main

import (
	"context"
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

	pb "github.com/albert-widi/logvault/pb"
	"google.golang.org/grpc"
)

var (
	// flag list
	vaultFlag    = flag.String("vault", "", "vault host to locate vault service")
	hostnameFlag = flag.String("hostname", "", "host name for agent")
	groupFlag    = flag.String("group", "", "group name of agent")
	fileFlag     = flag.String("files", "", "file lists")
	dirFlag      = flag.String("dir", "", "directory of log file")
	delimiter    = ","
	// prefix for writer
	prefix     string
	hostname   string
	workerList []*worker
	client     pb.LogvaultClient
)

func getFileList(list, dir string) ([]string, error) {
	var files []string
	if list != "" {
		fileList := strings.Split(list, delimiter)
		for _, val := range fileList {
			files = append(files, val)
		}
	} else {
		if err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if !f.IsDir() {
				files = append(files, f.Name())
			}
			return err
		}); err != nil {
			return nil, err
		}
	}
	return files, nil
}

func main() {
	var err error
	flag.Parse()
	prefix = *groupFlag
	if prefix == "" {
		log.Fatal("Group is empty, please set the host group")
	}
	hostname = *hostnameFlag
	if hostname == "" {
		hostname, err = os.Hostname()
		if err != nil {
			log.Fatal("Cannot get hostname, hostname cannot be empty. Please set it through hostname flag. ", err.Error())
		}
	}

	list := *fileFlag
	dir := *dirFlag
	vault := *vaultFlag

	if dir == "" {
		log.Fatal("No directory detected")
	}
	// make sure dir have / char
	if dir[len(dir)-1:] != "/" {
		dir += "/"
	}

	if vault != "" {
		conn, err := grpc.Dial(vault, grpc.WithInsecure())
		if err != nil {
			log.Fatal("Cannot connect to gRPC host ", err.Error())
		}
		client = pb.NewLogvaultClient(conn)
	}

	files, err := getFileList(list, dir)
	if err != nil {
		log.Fatal("Cannot get file list ", err.Error())
	}
	workerList = make([]*worker, len(files))
	iterator := 0
	for key := range files {
		w := newWorker(dir, files[key])
		workerList[iterator] = w
		log.Printf("Tailing log for %s \n", w.fileName)
		go w.run(context.TODO())
		iterator++
	}

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		// exit all worker
		for _, w := range workerList {
			w.done()
		}
		log.Println("Agent exited")
	}
}

type worker struct {
	prefix    string
	directory string
	fileName  string
	cmd       *exec.Cmd
	pushChan  chan string
	doneChan  chan bool
}

const bufferSize = 100

func newWorker(directory, fileName string) *worker {
	completePath := directory + fileName
	w := &worker{
		directory: directory,
		fileName:  fileName,
		cmd:       exec.CommandContext(context.TODO(), "tail", "-f", completePath),
		pushChan:  make(chan string, bufferSize),
		doneChan:  make(chan bool),
	}
	w.cmd.Stdout = newWriter(w.pushChan)
	w.cmd.Stderr = newWriter(w.pushChan)
	return w
}

func (w *worker) run(ctx context.Context) {
	if err := w.cmd.Run(); err != nil {
		log.Fatal("Failed to run command ", err.Error())
		return
	}

	for {
		select {
		case content := <-w.pushChan:
			if client != nil {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
				defer cancel()
				_, err := client.PushLog(ctx, &pb.PushRequest{Log: content, Prefix: prefix, Hostname: hostname})
				if err != nil {
					log.Println("Failed to push to logee service ", err.Error())
				}
			}
		// exit when done
		case <-w.doneChan:
			return
		}
	}
}

func (w *worker) done() {
	w.doneChan <- true
}

// TODO: support filename eg: nginx.access.log
type writer struct {
	prefix   string
	pushChan chan string
}

func newWriter(pushChan chan string) *writer {
	w := new(writer)
	w.prefix = prefix
	w.pushChan = pushChan
	return w
}

// redirected writer
func (w *writer) Write(b []byte) (int, error) {
	fmt.Println(string(b))
	w.pushChan <- string(b)
	return len(b), nil
}
