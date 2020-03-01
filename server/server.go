package main

import (
	"flag"
	"github.com/speakpavlov/segmentation"
	"io/ioutil"
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
)

var (
	port        int
	logfile     string
	dumpDirPath string

	seg    *segmentation.Segmentation
	logger *log.Logger
)

func init() {
	flag.StringVar(&logfile, "logfile", "", "Location of the logfile")
	flag.StringVar(&dumpDirPath, "dumpDirPath", "./dump/", "Dump dir")
	flag.IntVar(&port, "port", 9090, "The port to listen on.")
}

func main() {
	flag.Parse()

	initLogger()
	initializeSeg(logger)

	http.Handle(segmentationPath, segmentationHandler(logger))
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

func initializeSeg(logger *log.Logger) {
	seg = segmentation.NewSegmentation()

	files, fErr := ioutil.ReadDir(dumpDirPath)
	if fErr != nil {
		logger.Print("Dump was not loaded " + fErr.Error())

		os.Mkdir(dumpDirPath, 0777)
	}

	var segments []segmentation.Segment

	for _, f := range files {
		tagId := f.Name()
		err := Load(dumpDirPath+tagId, &segments)
		if err != nil {
			log.Fatal(err)
		}

		sErr := seg.UpdateSegments(tagId, segments)
		if sErr != nil {
			log.Fatal(sErr)
		}

		seg.Segments[tagId] = segments

		logger.Print("Tag: " + tagId + " was loaded")
	}
}
