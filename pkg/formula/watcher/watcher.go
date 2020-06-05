package watcher

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
)

type WatchManager struct {
	watcher *fsnotify.Watcher
}

func New() *WatchManager {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	return &WatchManager{watcher: watcher}
}

func (w *WatchManager) Watch(dir string) {
	defer w.watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)

				if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Rename == fsnotify.Rename {
					if Exists(event.Name) && isDir(event.Name) {
						w.watch(event.Name)
					}
				}

				// TODO: run Makefile or .bat
			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	w.start(dir)
	fmt.Printf("Watching dir %s \n", dir)
	<-done
}

func (w *WatchManager) watch(dirName string) {
	if err := w.watcher.Add(dirName); err != nil {
		log.Println(err)
		return
	}
	return
}

func (w *WatchManager) start(name string) {
	if !isDir(name) {
		return
	}

	w.watch(name)

	open, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	names, err := open.Readdirnames(-1)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, n := range names { // Watch all directories in current dir
		dir := fmt.Sprintf("%s/%s", name, n)
		if !isDir(dir) {
			continue
		}

		w.watch(dir)
		w.start(dir)
	}

	return
}

func isDir(name string) bool {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return false
	}

	info, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return false
	}

	if !info.IsDir() {
		return false
	}
	return true
}

func Exists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}

	return true
}
