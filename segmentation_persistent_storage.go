package segmentation

import (
	"io/ioutil"
	"log"
	"os"
)

type PersistentStorage struct {
	dumpDirPath string
}

func NewPersistentStorage(dumpDirPath string) *PersistentStorage {
	os.Mkdir(dumpDirPath, 0777)

	return &PersistentStorage{dumpDirPath}
}

func (s PersistentStorage) Load() map[string][]string {
	files, err := ioutil.ReadDir(s.dumpDirPath)
	if err != nil {
		log.Fatal(err)
	}

	segmentation := make(map[string][]string)

	for _, f := range files {
		var expressions []string
		tagId := f.Name()
		err := Load(s.dumpDirPath+"/"+tagId, &expressions)
		if err != nil {
			log.Fatal(err)
		}

		segmentation[tagId] = expressions
	}

	return segmentation
}

func (s PersistentStorage) SaveNewSegment(tagId string, segments []string) error {
	return Save(s.dumpDirPath+"/"+tagId, segments)
}
