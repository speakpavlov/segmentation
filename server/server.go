package main

import (
	"flag"
	"github.com/speakpavlov/segmentation"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	apiVersion  = "v1"
	apiBasePath = "/api/" + apiVersion + "/"

	//path
	segmentationPath = apiBasePath + "segmentation"
	statusPath       = apiBasePath + "status"
)

var (
	port        int
	logfile     string
	dumpDirPath string

	persistentStorage *segmentation.PersistentStorage
	segmentationList  *segmentation.SegmentationMap
	logger            *log.Logger
)

func init() {
	flag.StringVar(&logfile, "logfile", "", "Location of the logfile")
	flag.StringVar(&dumpDirPath, "dumpDirPath", "./dump/", "Dump dir")
	flag.IntVar(&port, "port", 9090, "The port to listen on.")
}

func main() {
	flag.Parse()

	initLogger()

	persistentStorage = segmentation.NewPersistentStorage(dumpDirPath)
	segmentationList = persistentStorage.Load()

	http.Handle(segmentationPath, segmentationHandler(logger))
	http.Handle(statusPath, statusHandler(logger))
	logger.Printf("starting sever on :%d", port)

	strPort := ":" + strconv.Itoa(port)
	log.Fatal("ListenAndServe: ", http.ListenAndServe(strPort, nil))
}

func initLogger() {
	if logfile == "" {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		logger = log.New(f, "", log.LstdFlags)
	}
}
