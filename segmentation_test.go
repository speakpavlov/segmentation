package segmentation

import (
	"testing"
)

func TestSegmentation(t *testing.T) {
	db := NewSegmentationDb(&AntonMedvExpression{})

	_, err := db.PublishSegmentation(1, []Segment{
		{Index: 1, Expression: "A == 1", Value: "1"},
		{Index: 2, Expression: "A == 2", Value: "3"},
		{Index: 3, Expression: "A == 3", Value: "3"},
	})

	if err != nil {
		t.Errorf("Publish was not unsuccessefully")
	}

	_, sErr := db.GetSegment(0, map[string]interface{}{"A": 1})

	if sErr != nil {
		t.Errorf("Id 0 has segments")
	}

	segment2, err2 := db.GetSegment(1, map[string]interface{}{"A": 1})

	if err2 != nil || segment2 == nil {
		t.Error("Id 1 has not segments")
	}

	if segment2 != nil && segment2.Value != "1" {
		t.Error("Segment Id 1 is incorrect")
	}
}

func BenchmarkDbSegments(b *testing.B) {
	db := NewSegmentationDb(&AntonMedvExpression{})

	_, err := db.PublishSegmentation(1, []Segment{
		{Index: 1, Expression: "A == 1", Value: "1"},
		{Index: 2, Expression: "A == 2", Value: "3"},
		{Index: 3, Expression: "A == 3", Value: "3"},
	})

	if err != nil {
		b.Errorf("Publish was not unsuccessefully")
	}

	for i := 0; i < b.N; i++ {
		segment2, err2 := db.GetSegment(1, map[string]interface{}{"A": 1})

		if err2 != nil || segment2 == nil {
			b.Error("Id 1 has not segments")
		}

		if segment2 != nil && segment2.Value != "1" {
			b.Error("Segment Id 1 is incorrect")
		}
	}
}

func BenchmarkDbSegmentsEmpty(b *testing.B) {
	db := NewSegmentationDb(&AntonMedvExpression{})

	_, err := db.PublishSegmentation(1, []Segment{
		{Index: 1, Expression: "A == 1", Value: "1"},
		{Index: 2, Expression: "A == 2", Value: "3"},
		{Index: 3, Expression: "A == 3", Value: "3"},
	})

	if err != nil {
		b.Errorf("Publish was not unsuccessefully")
	}

	for i := 0; i < b.N; i++ {
		_, err := db.GetSegment(0, map[string]interface{}{"A": 1})

		if err == nil {
			b.Errorf("Id 0 has segments")
		}
	}
}
