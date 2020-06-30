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

func (s PersistentStorage) Load() *SegmentationMap {
	seg := &SegmentationMap{}
	files, err := ioutil.ReadDir(s.dumpDirPath)
	if err != nil {
		log.Fatal(err)
	}

	var segments []string

	for _, f := range files {
		tagId := f.Name()
		err := Load(s.dumpDirPath+"/"+tagId, &segments)
		if err != nil {
			log.Fatal(err)
		}

		sErr := seg.UpdateSegments(tagId, segments)
		if sErr != nil {
			log.Fatal(sErr)
		}

		//logger.Print("Tag: " + tagId + " was loaded")
	}

	return seg
}

func (s PersistentStorage) SaveNewSegment(tagId string, segments []string) error {
	return Save(s.dumpDirPath+"/"+tagId, segments)
}
