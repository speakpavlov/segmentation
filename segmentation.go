package segmentation

import (
	"errors"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"log"
	"strconv"
)

//SegmentationMap struct
type SegmentationMap struct {
	byteCodes map[string][]*vm.Program
}

//Segment struct
type Segment struct {
	Expression string
	ByteCode   interface{}
}

func NewSegmentationMap() *SegmentationMap {
	return &SegmentationMap{}
}

func (seg *SegmentationMap) UpdateSegments(tag string, expressions []string) error {
	var byteCodes = make([]*vm.Program, len(expressions))

	for i := range expressions {
		if expressions[i] == "" {
			//empty is valid segment
			expressions[i] = "true"
		}

		//compile expression
		program, err := expr.Compile(expressions[i])

		if err != nil {
			log.Print(err)

			return errors.New("Compilation segment #" + strconv.Itoa(i) + " was failed")
		}

		//update byte code
		byteCodes[i] = program
	}

	if seg.byteCodes == nil {
		seg.byteCodes = map[string][]*vm.Program{}
	}

	seg.byteCodes[tag] = byteCodes

	return nil
}

func (seg *SegmentationMap) GetSegments(tag string, data map[string]interface{}) ([]int, error) {
	var segmentIndexes []int
	env := NewEnv(data)

	if segmentation, ok := seg.byteCodes[tag]; ok {
		for index, segment := range segmentation {
			//segment without rules

			result, err := expr.Run(segment, env)
			if err != nil {
				return nil, err
			}

			if result.(bool) {
				segmentIndexes = append(segmentIndexes, index)
			}
		}
	}

	return segmentIndexes, nil
}
