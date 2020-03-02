package segmentation

import (
	"testing"
)

func TestSegmentation(t *testing.T) {
	segmentation := NewSegmentation()

	err := segmentation.UpdateSegments("seg_1", []Segment{
		{Index: 1, Expression: "A == 1", Value: "1"},
		{Index: 2, Expression: "A == 2", Value: "3"},
		{Index: 3, Expression: "A == 3", Value: "3"},
	})

	if err != nil {
		t.Errorf("Publish was not unsuccessefully")
	}

	_, sErr := segmentation.GetSegments("seg_0", map[string]interface{}{"A": 1})

	if sErr != nil {
		t.Errorf("Id seg_0 has segments")
	}

	segment2, err2 := segmentation.GetSegments("seg_1", map[string]interface{}{"A": 1})

	if err2 != nil || segment2 == nil {
		t.Error("seg_1 has not segments")
	}

	if len(segment2) != 1 {
		t.Error("Segments are incorrect", len(segment2))
	}
}

func BenchmarkDbSegments(b *testing.B) {
	db := NewSegmentation()

	err := db.UpdateSegments("1", []Segment{
		{Index: 1, Expression: "A == 1", Value: "1"},
		{Index: 2, Expression: "A == 2", Value: "3"},
		{Index: 3, Expression: "A == 3", Value: "3"},
	})

	if err != nil {
		b.Errorf("Publish was not unsuccessefully")
	}

	for i := 0; i < b.N; i++ {
		segment2, err2 := db.GetSegments("1", map[string]interface{}{"A": 1})

		if err2 != nil || segment2 == nil {
			b.Error("Id 1 has not segments")
		}

		if segment2 != nil && segment2[0].Value != "1" {
			b.Error("Segment Id 1 is incorrect")
		}
	}
}

func BenchmarkDbSegmentsEmpty(b *testing.B) {
	db := NewSegmentation()

	err := db.UpdateSegments("1", []Segment{
		{Index: 1, Expression: "A == 1", Value: "1"},
		{Index: 2, Expression: "A == 2", Value: "3"},
		{Index: 3, Expression: "A == 3", Value: "3"},
	})

	if err != nil {
		b.Errorf("Publish was not unsuccessefully")
	}

	for i := 0; i < b.N; i++ {
		segments, _ := db.GetSegments("0", map[string]interface{}{"A": 1})

		if len(segments) > 0 {
			b.Errorf("Id 0 has segments")
		}
	}
}
