package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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

func main() {
	err := createDir("/var/log/producer")
	if err != nil {
		log.Fatal("Cannot create producer log directory")
	}

	f1, err := os.OpenFile("/var/log/producer/log1.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Cannot create f1 ", err.Error())
	}
	defer f1.Close()
	f2, err := os.OpenFile("/var/log/producer/log2.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Cannot create f2 ", err.Error())
	}
	defer f2.Close()
	log1 := log.New(f1, "F1:", log.Ldate|log.Ltime|log.Lshortfile)
	log2 := log.New(f2, "F2:", log.Ldate|log.Ltime|log.Lshortfile)

	ticker := time.NewTicker(time.Second * 1)
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case t := <-ticker.C:
			log.Println("Logging something happened ", t.String())
			log1.Println("Something happened in F1 ", t.String())
			log2.Println("Something happened in F2 ", t.String())
		case <-term:
			log.Println("Producer exited")
			return
		}
	}
}
