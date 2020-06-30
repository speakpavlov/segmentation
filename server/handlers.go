package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type SegmentationPutInput struct {
	TagId       string   `json:"tag_id"`
	Expressions []string `json:"expressions"`
}

type SegmentationGetRequest struct {
	TagId string                 `json:"tag_id"`
	Data  map[string]interface{} `json:"data"`
}

//status handler
func statusHandler(l *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeSuccess(w, l, "OK")
	})
}

func segmentationHandler(l *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			loadSegmentationHandler(w, r, l)
		case http.MethodPut:
			importSegmentationHandler(w, r, l)
		}
	})
}

func importSegmentationHandler(w http.ResponseWriter, r *http.Request, l *log.Logger) {
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

	sErr := segmentationList.UpdateSegments(segmentationInput.TagId, segmentationInput.Expressions)
	if sErr != nil {
		writeError(sErr, w, l, http.StatusBadRequest)
		return
	}

	dErr := persistentStorage.SaveNewSegment(segmentationInput.TagId, segmentationInput.Expressions)
	if dErr != nil {
		writeError(dErr, w, l, http.StatusInternalServerError)
		return
	}

	writeSuccess(w, l, "OK")
}

func loadSegmentationHandler(w http.ResponseWriter, r *http.Request, l *log.Logger) {
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

	segments, sErr := segmentationList.GetSegments(segmentationInput.TagId, segmentationInput.Data)
	if sErr != nil {
		writeError(sErr, w, l, http.StatusBadRequest)
		return
	}

	if len(segments) == 0 {
		writeError(errors.New("Segmentation was not found."), w, l, http.StatusNotFound)
		return
	}

	//var result []map[string]interface{}
	//for _, seg := range segments {
	//	result = append(result, map[string]interface{}{
	//		"index": seg,
	//	})
	//}

	writeSuccess(w, l, segments)
}

/////

func writeSuccess(w http.ResponseWriter, l *log.Logger, result interface{}) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"response": result,
	})
}

func writeError(err error, w http.ResponseWriter, l *log.Logger, request int) {
	if request != http.StatusNotFound {
		l.Print(err)
	}

	writeJSON(w, request, map[string]interface{}{
		"response": false,
		"error":    err.Error(),
	})
}

func writeJSON(w http.ResponseWriter, code int, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(data)
}
