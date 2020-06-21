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
	return &PersistentStorage{dumpDirPath}
}

func (s PersistentStorage) Load() *SegmentationMap {
	seg := &SegmentationMap{}

	files, fErr := ioutil.ReadDir(s.dumpDirPath)
	if fErr != nil {
		//logger.Print("Dump was not loaded " + fErr.Error())

		os.Mkdir(s.dumpDirPath, 0777)
	}

	var segments []Segment

	for _, f := range files {
		tagId := f.Name()
		err := Load(s.dumpDirPath+tagId, &segments)
		if err != nil {
			log.Fatal(err)
		}

		sErr := seg.UpdateSegments(tagId, segments)
		if sErr != nil {
			log.Fatal(sErr)
		}

		seg.Segments[tagId] = segments

		//logger.Print("Tag: " + tagId + " was loaded")
	}

	return seg
}

func (s PersistentStorage) Save(segmentation *SegmentationMap) {

}
