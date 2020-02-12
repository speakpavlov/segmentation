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
	Segments []Segment `json:"segments"`
}

//Segment struct
type Segment struct {
	Expression string `json:"expression"`
	Value      string `json:"value"`
}

type SegmentationGetRequest struct {
	SegmentationId int         `json:"segmentation_id"`
	Data           interface{} `json:"data"`
}

func segmentationHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getHandler(w, r)
		case http.MethodPut:
			putHandler(w, r)
		}
	})
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	entry, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var segmentationInput SegmentationPutInput
	mErr := json.Unmarshal(entry, &segmentationInput)

	if mErr != nil {
		errorResponse(mErr, w, http.StatusBadRequest)
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

	segmentationId, sErr := db.PublishSegmentation(0, segments)

	if sErr != nil {
		errorResponse(sErr, w, http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"segmentation_id": segmentationId,
		"db":              db.SegmentationList,
	}

	jsonData, jErr := json.Marshal(data)

	if jErr != nil {
		errorResponse(jErr, w, http.StatusInternalServerError)
		return
	}

	successResponse(w, jsonData)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	entry, rErr := ioutil.ReadAll(r.Body)
	if rErr != nil {
		errorResponse(rErr, w, http.StatusBadRequest)
		return
	}

	var segmentationInput SegmentationGetRequest
	mErr := json.Unmarshal(entry, &segmentationInput)

	if mErr != nil {
		errorResponse(mErr, w, http.StatusBadRequest)
		return
	}

	segment, sErr := db.GetSegment(segmentationInput.SegmentationId, segmentationInput.Data)

	if sErr != nil {
		errorResponse(sErr, w, http.StatusBadRequest)
		return
	}

	if segment == nil {
		errorResponse(errors.New("Segment was not found."), w, http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"segment": segment,
	}

	jsonData, jErr := json.Marshal(data)

	if jErr != nil {
		errorResponse(jErr, w, http.StatusInternalServerError)
		return
	}

	successResponse(w, jsonData)
}

func successResponse(w http.ResponseWriter, jsonData []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func errorResponse(err error, w http.ResponseWriter, request int) {
	log.Print(err)

	data := map[string]interface{}{
		"error": err.Error(),
	}

	jsonData, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(request)
	w.Write(jsonData)
}
