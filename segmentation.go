package segmentation

import (
	"errors"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
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

		program, err := expr.Compile(segments[i].Expression)

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

func (seg *Segmentation) GetSegments(tag string, data map[string]interface{}) ([]Segment, error) {
	var segments []Segment
	env := NewEnv(data)

	if segmentation, ok := seg.Segments[tag]; ok {
		for _, segment := range segmentation {
			//segment without rules
			if segment.ByteCode == nil {
				segments = append(segments, segment)
			}

			result, err := expr.Run(segment.ByteCode.(*vm.Program), env)
			if err != nil {
				return nil, err
			}

			if result.(bool) {
				segments = append(segments, segment)
			}
		}
	}

	return segments, nil
}
