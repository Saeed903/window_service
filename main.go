package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kardianos/service"
	"github.com/saeed903/windows_service/config"
	"github.com/saeed903/windows_service/pkg/utils"
)

type myService struct{}

func (s *myService) Start(svc service.Service) error {
    go s.run()
    return nil
}

func (s *myService) Stop(svc service.Service) error {
    return nil
}

func (s *myService) run() {
    var path string

    cfg, err := config.InitConfig()
    if err != nil {
        log.Fatalln("unable read config file", err)
    }

    flag.StringVar(&path, "folder path", cfg.FolderWatchPath, "folder path to watching...")

    // if path == "" {
    //     path = "C:/data"
    // }

    if !filepath.IsAbs(path) {
        if err := os.Mkdir(path, os.ModeAppend); err != nil {
            fmt.Printf("Unable Create Folder in path: %s, err is %s \n", path, err)
        }
    }

    fmt.Printf("System start watch folder: %s and remove any file in this folder...", path)

	tempFiles := []string{}
	files := []string{}
	oldFiles := []string{}
	//oldTime := time.Duration(1)
    //var readFiles []string
	for {
		readFiles := []string{}
		files = []string{}
		if len(tempFiles) == 0 {
			readFiles, _ = utils.ReadFiles(path)
			// if err != nil {
			// 	fmt.Printf("can not read from path: %s and err: %v",path, err)
			// }

			if len(readFiles) == 0 {
				oldFiles = []string{}
			} else {
				if len(oldFiles) > 0 {
					exists := false
					for j, readFile := range readFiles {
						for i, oldFile := range oldFiles {
							if readFile[i] == oldFile[j] {
								exists = true
								break
							}
						}
						if !exists {
							files = append(files, readFile)
						}
						exists = false
					}
				} else {
					files = append(files, readFiles...)
				}

				tempFiles = append(tempFiles, files...)
			}

		}
		if len(tempFiles) == 0 {
			time.Sleep(time.Second)
			continue
		}
		//l := len(tempFiles)

		oldFiles = append(oldFiles, tempFiles...)
		//oldTime = time.Duration(3*l)

		for _, file := range tempFiles {
			//defer wg.Done()
			err := os.Remove(file)
            if err != nil {
                log.Fatalf("Remove %s err: %s", filepath.Base(file), err)
                return
            }
            fmt.Printf("file: %s remove from folder: %s\n", filepath.Base(file), filepath.Dir(file))
		}

		//time.Sleep(time.Duration(l/2) * time.Second)

		tempFiles = []string{}

		fmt.Printf("Count files read directory: %d\n", len(files))
		//fmt.Printf("Time Sleep for read directory: %d", l+1)

		time.Sleep(time.Second)

	}
}

func main() {
    svcConfig := &service.Config{
        Name:        "RemoveFileService",
        DisplayName: "My Service",
        Description: "This is a sample Go service.",
    }

    prg := &myService{}
    svc, err := service.New(prg, svcConfig)
    if err != nil {
        log.Fatal(err)
    }

    // Start the service.
    if err := svc.Run(); err != nil {
        log.Fatal(err)
    }
}
