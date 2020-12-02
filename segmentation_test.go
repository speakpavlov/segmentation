package segmentation

import (
	"testing"
)

func TestSegmentationExist(t *testing.T) {
	segmentation := initBaseSegmentation(t)

	_, sErr := segmentation.GetSegments("seg_0", map[string]interface{}{"A": 1})

	if sErr != nil {
		t.Errorf("Id seg_0 has segments")
	}
}

func TestSegmentationIndex(t *testing.T) {
	segmentation := initBaseSegmentation(t)

	segment, err := segmentation.GetSegments("seg_1", map[string]interface{}{"A": 2})

	if err != nil {
		t.Error("seg_1 has not segments")
	}

	if segment == nil || segment[0] != 1 {
		t.Error("seg_1 with index 1 has incorrect value")
	}

	if len(segment) != 2 {
		t.Error("byteCodes are incorrect", len(segment))
	}
}
func TestSegmentationEmpty(t *testing.T) {
	segmentation := initBaseSegmentation(t)

	segment, err := segmentation.GetSegments("seg_1", map[string]interface{}{"A": 5})

	if err != nil {
		t.Error("seg_1 has not segments")
	}

	if segment == nil || segment[0] != 2 {
		t.Error("seg_1 with index 2 has incorrect value")
	}

	if len(segment) != 1 {
		t.Error("byteCodes are incorrect", len(segment))
	}
}

func TestNotCompiledExpression(t *testing.T) {
	segmentation := NewSegmentationMap()

	err := segmentation.UpdateSegments("seg_1", []string{
		"a = 1",
	})

	if err == nil {
		t.Errorf("Incorrect expression should case error")
	}
}

func initBaseSegmentation(t *testing.T) *SegmentationMap {
	segmentation := NewSegmentationMap()

	err := segmentation.UpdateSegments("seg_1", []string{
		"A == 1",
		"A == 2",
		"",
	})

	if err != nil {
		t.Errorf("Publish was not unsuccessefully")
	}

	return segmentation
}
