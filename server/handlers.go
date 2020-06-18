package main

import (
	"encoding/json"
	"errors"
	"github.com/speakpavlov/segmentation"
	"io/ioutil"
	"log"
	"net/http"
)

type SegmentationPutInput struct {
	TagId    string    `json:"tag_id"`
	Segments []Segment `json:"segments"`
}

//Segment struct
type Segment struct {
	Expression string `json:"expression"`
	Value      string `json:"value"`
}

type SegmentationGetRequest struct {
	TagId string                 `json:"tag_id"`
	Data  map[string]interface{} `json:"data"`
}

//status handler
func statusHandler(l *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Print(r.RequestURI)

		writeSuccess(w, l, "OK")
	})
}

func segmentationHandler(l *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postHandler(w, r, l)
		case http.MethodPut:
			putHandler(w, r, l)
		}
	})
}

func putHandler(w http.ResponseWriter, r *http.Request, l *log.Logger) {
	entry, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l.Print(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var segmentationInput SegmentationPutInput
	mErr := json.Unmarshal(entry, &segmentationInput)

	if mErr != nil {
		writeError(mErr, w, l, http.StatusBadRequest)
		return
	}

	var segments []segmentation.Segment

	index := 0
	for _, v := range segmentationInput.Segments {
		segments = append(segments, segmentation.Segment{
			Index:      index,
			Expression: v.Expression,
			Value:      v.Value,
		})
		index++
	}

	dErr := Save(dumpDirPath+segmentationInput.TagId, segments)
	if dErr != nil {
		writeError(dErr, w, l, http.StatusInternalServerError)
		return
	}

	sErr := seg.UpdateSegments(segmentationInput.TagId, segments)
	if sErr != nil {
		writeError(sErr, w, l, http.StatusBadRequest)
		return
	}

	writeSuccess(w, l, "OK")
}

func postHandler(w http.ResponseWriter, r *http.Request, l *log.Logger) {
	entry, rErr := ioutil.ReadAll(r.Body)
	if rErr != nil {
		writeError(rErr, w, l, http.StatusBadRequest)
		return
	}

	var segmentationInput SegmentationGetRequest

	mErr := json.Unmarshal(entry, &segmentationInput)
	if mErr != nil {
		writeError(mErr, w, l, http.StatusBadRequest)
		return
	}

	segments, sErr := seg.GetSegments(segmentationInput.TagId, segmentationInput.Data)
	if sErr != nil {
		writeError(sErr, w, l, http.StatusBadRequest)
		return
	}

	if len(segments) == 0 {
		writeError(errors.New("Segments was not found."), w, l, http.StatusNotFound)
		return
	}

	var result []map[string]interface{}
	for _, seg := range segments {
		result = append(result, map[string]interface{}{
			"index": seg.Index,
			"value": seg.Value,
		})
	}

	writeSuccess(w, l, result)
}

/////

func writeSuccess(w http.ResponseWriter, l *log.Logger, result interface{}) {
	writeJSON(w, map[string]interface{}{
		"response": result,
	})
}

func writeError(err error, w http.ResponseWriter, l *log.Logger, request int) {
	if request != http.StatusNotFound {
		l.Print(err)
	}

	w.WriteHeader(request)

	writeJSON(w, map[string]interface{}{
		"response": false,
		"error":    err.Error(),
	})
}

func writeJSON(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
