package segmentation

import (
	"os"
	"testing"
)

func TestPersistentStorage(t *testing.T) {
	storage := NewPersistentStorage("testDump")
	defer os.RemoveAll("testDump")

	storage.SaveNewSegment("1",
		[]string{
			"A == 1",
			"A == 2",
		})

	segmentsMap := storage.Load()

	if segmentsMap["1"][0] != "A == 1" {
		t.Error("Data of first segment is incorrect")
	}
}
