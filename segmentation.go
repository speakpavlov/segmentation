package segmentation

import (
	"errors"
	"log"
	"strconv"
)

//Segmentation struct
type Segmentation struct {
	Segments map[string][]Segment
}

//Segment struct
type Segment struct {
	Index      int
	Expression string
	Value      string
	ByteCode   interface{}
}

var expression MedvExpression

func NewSegmentation() *Segmentation {
	return &Segmentation{}
}

func (seg *Segmentation) UpdateSegments(tag string, segments []Segment) error {
	//compile expression
	for i := range segments {
		if segments[i].Expression == "" {
			//don't need compile byte code
			continue
		}

		program, err := expression.compile(segments[i].Expression)

		if err != nil {
			log.Print(err)
			return errors.New("Compilation segment #" + strconv.Itoa(i) + " was failed")
		}

		//update byte code
		segments[i].ByteCode = program
	}

	//init
	if seg.Segments == nil {
		seg.Segments = map[string][]Segment{}
	}

	seg.Segments[tag] = segments

	return nil
}

func (seg *Segmentation) GetSegments(tag string, data interface{}) ([]Segment, error) {
	var segments []Segment

	if segmentation, ok := seg.Segments[tag]; ok {
		for _, segment := range segmentation {
			//segment without rules
			if segment.ByteCode == nil {
				segments = append(segments, segment)
			}

			result, err := expression.execute(segment.ByteCode, data)
			if err != nil {
				return nil, err
			}

			if result {
				segments = append(segments, segment)
			}
		}
	}

	return segments, nil
}
