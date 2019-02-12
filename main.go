package main

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os"
)

func main() {
	certDir := os.Getenv("CERT_DIR")
	watchDir := os.Getenv("WATCH_DIR")
	certFileName := os.Getenv("CERT_FILE_NAME")

	if len(certDir) < 1 {
		log.Println("Missing cert dir")
		log.Println()
		flag.Usage()
		os.Exit(1)
	}

	if len(watchDir) < 1 {
		log.Println("Missing watch dir")
		log.Println()
		flag.Usage()
		os.Exit(1)
	}

	if len(certFileName) < 1 {
		log.Println("Missing file name")
		log.Println()
		flag.Usage()
		os.Exit(1)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case _ = <-watcher.Events:
				funcName(watchDir, certDir, certFileName)
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	funcName(watchDir, certDir, certFileName)

	log.Printf("Watching directory: %q", watchDir)
	err = watcher.Add(watchDir)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}

func funcName(watchDir string, certDir string, certFileName string) {
	log.Println("secret updated")
	crtFile := fmt.Sprintf("%s/tls.crt", watchDir)
	keyFile := fmt.Sprintf("%s/tls.key", watchDir)
	if Exists(crtFile) && Exists(keyFile) {

		crtIn, err := os.Open(crtFile)
		if err != nil {
			log.Fatalln("failed to open second file for reading:", err)
		}
		defer crtIn.Close()
		keyIn, err := os.Open(keyFile)
		if err != nil {
			log.Fatalln("failed to open second file for reading:", err)
		}
		defer keyIn.Close()

		out, err := os.OpenFile(fmt.Sprintf("%s/%s", certDir, certFileName), os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalln("failed to open outpout file:", err)
		}
		defer out.Close()

		_, err = io.Copy(out, crtIn)
		if err != nil {
			log.Fatalln("failed to append cert file to output:", err)
		}
		_, err = io.Copy(out, keyIn)
		if err != nil {
			log.Fatalln("failed to append key file to output:", err)
		}

	} else {
		log.Printf("missing one of tls.crt or tls.key skip")
	}
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
