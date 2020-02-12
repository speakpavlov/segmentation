package segmentation

import (
	"errors"
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

func NewSegmentationDb(expression Expression) Db {
	return Db{Expression: expression}
}

func (db *Db) PublishSegmentation(id int, segments []Segment) bool {
	//edit or create
	edited := false

	//compile expression
	for i := range segments {
		program, err := db.Expression.compile(segments[i].Expression)

		if err != nil {
			return false
		}

		segments[i].byteCode = program
	}

	for i := range db.SegmentationList {
		if db.SegmentationList[i].Id == id {
			db.SegmentationList[i].Segments = segments
			edited = true
			break
		}
	}

	if edited != true {
		db.SegmentationList = append(db.SegmentationList, Segmentation{id, segments})
	}

	return true
}

func (db *Db) GetSegment(id int, data interface{}) (*Segment, error) {
	for _, segmentation := range db.SegmentationList {
		if segmentation.Id == id {
			for _, segment := range segmentation.Segments {
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

	return nil, errors.New("Segment was not found.")
}
