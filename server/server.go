package main

import (
	"flag"
	"fmt"
	"github.com/speakpavlov/segmentation"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	apiVersion  = "v1"
	apiBasePath = "/api/" + apiVersion + "/"

	version = "1.0.0"

	//path
	segmentationPath = apiBasePath + "segmentation"
)

var (
	port     int
	logfile  string
	dumpfile string
	ver      bool

	db *segmentation.Db
)

func init() {
	flag.StringVar(&logfile, "logfile", "", "Location of the logfile")
	flag.StringVar(&dumpfile, "dumpfile", "dump.json", "Dump file")
	flag.BoolVar(&ver, "version", false, "Print server version.")
	flag.IntVar(&port, "port", 9090, "The port to listen on.")
}

func main() {
	flag.Parse()

	if ver {
		fmt.Printf("SegmentationPutInput v%s", version)
		os.Exit(0)
	}

	var logger *log.Logger

	if logfile == "" {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		logger = log.New(f, "", log.LstdFlags)
	}

	//var err error
	db = segmentation.NewSegmentationDb(segmentation.AntonMedvExpression{})

	//err := Load(dumpfile, &db.SegmentationList)
	//
	//if err == nil {
	//	logger.Print("Dump was loaded")
	//}

	http.Handle(segmentationPath, serviceLoader(segmentationHandler(logger), requestMetrics(logger)))

	logger.Printf("starting sever on :%d", port)

	strPort := ":" + strconv.Itoa(port)
	log.Fatal("ListenAndServe: ", http.ListenAndServe(strPort, nil))
}
