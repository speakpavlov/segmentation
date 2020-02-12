package segmentation

import (
	"errors"
	"log"
	"strconv"
)

//Db struct
type Db struct {
	Expression       Expression
	SegmentationList []Segmentation
}

//Segmentation struct
type Segmentation struct {
	Id       int
	Segments []Segment
}

//Segment struct
type Segment struct {
	Index      int
	Expression string
	Value      string
	byteCode   interface{}
}

type Expression interface {
	compile(expression string) (interface{}, error)
	execute(program interface{}, env interface{}) (bool, error)
}

func NewSegmentationDb(expression Expression) *Db {
	return &Db{Expression: expression}
}

func (db *Db) PublishSegmentation(id int, segments []Segment) (int, error) {
	//edit or create
	edited := false

	//compile expression
	for i := range segments {
		if segments[i].Expression == "" {
			//don't need compile bytecode
			continue
		}

		program, err := db.Expression.compile(segments[i].Expression)

		if err != nil {
			log.Print(err)
			return 0, errors.New("Compilation segment #" + strconv.Itoa(i) + " was failed")
		}

		segments[i].byteCode = program
	}

	if id > 0 {
		for i := range db.SegmentationList {
			if db.SegmentationList[i].Id == id {
				db.SegmentationList[i].Segments = segments
				edited = true
				break
			}
		}
	}

	if edited != true {
		//next id, starts with 1
		id = len(db.SegmentationList) + 1

		db.SegmentationList = append(db.SegmentationList, Segmentation{id, segments})
	}

	return id, nil
}

func (db *Db) GetSegment(id int, data interface{}) (*Segment, error) {
	for _, segmentation := range db.SegmentationList {
		if segmentation.Id == id {
			for _, segment := range segmentation.Segments {
				//segment without rules
				if segment.byteCode == nil {
					return &segment, nil
				}

				result, err := db.Expression.execute(segment.byteCode, data)
				if err != nil {
					return nil, err
				}

				if result {
					return &segment, nil
				}
			}
		}
	}

	return nil, nil
}
