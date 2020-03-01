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
	TagId string      `json:"tag_id"`
	Data  interface{} `json:"data"`
}

func segmentationHandler(l *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			getHandler(w, r, l)
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
		errorResponse(mErr, w, l, http.StatusBadRequest)
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
		errorResponse(dErr, w, l, http.StatusInternalServerError)
		return
	}

	sErr := seg.UpdateSegments(segmentationInput.TagId, segments)

	if sErr != nil {
		errorResponse(sErr, w, l, http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"response": "OK",
	}

	jsonData, jErr := json.Marshal(data)

	if jErr != nil {
		errorResponse(jErr, w, l, http.StatusInternalServerError)
		return
	}

	successResponse(w, jsonData)
}

func getHandler(w http.ResponseWriter, r *http.Request, l *log.Logger) {
	entry, rErr := ioutil.ReadAll(r.Body)
	if rErr != nil {
		errorResponse(rErr, w, l, http.StatusBadRequest)
		return
	}

	var segmentationInput SegmentationGetRequest

	mErr := json.Unmarshal(entry, &segmentationInput)
	if mErr != nil {
		errorResponse(mErr, w, l, http.StatusBadRequest)
		return
	}

	segments, sErr := seg.GetSegments(segmentationInput.TagId, segmentationInput.Data)
	if sErr != nil {
		errorResponse(sErr, w, l, http.StatusBadRequest)
		return
	}

	if len(segments) == 0 {
		errorResponse(errors.New("Segments was not found."), w, l, http.StatusNotFound)
		return
	}

	var result []map[string]interface{}
	for _, seg := range segments {
		result = append(result, map[string]interface{}{
			"index": seg.Index,
			"value": seg.Value,
		})
	}

	data := map[string]interface{}{
		"response": result,
	}

	jsonData, jErr := json.Marshal(data)
	if jErr != nil {
		errorResponse(jErr, w, l, http.StatusInternalServerError)
		return
	}

	successResponse(w, jsonData)
}

func successResponse(w http.ResponseWriter, jsonData []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func errorResponse(err error, w http.ResponseWriter, l *log.Logger, request int) {
	if request != http.StatusNotFound {
		l.Print(err)
	}

	data := map[string]interface{}{
		"response": false,
		"error":    err.Error(),
	}

	jsonData, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(request)
	w.Write(jsonData)
}
