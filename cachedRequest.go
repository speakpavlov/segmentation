package segmentation

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type CachedRequest map[string]interface{}

type HttpResponse struct {
	Result interface{} `json:"result"`
}

func (cachedRequest CachedRequest) makeCachedRequest(url string) interface{} {
	//if exist in cachedRequest
	if value, ok := cachedRequest[url]; ok {
		return value
	}

	cachedRequest[url] = makeRequest(url)

	return cachedRequest[url]
}

func makeRequest(url string) interface{} {
	response, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Print("Cannot load, url: " + url + ", err: " + err.Error())

		return false
	}

	var resp HttpResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Print("Cannot unmarshal, url: " + url + ", err: " + err.Error())

		return false
	}

	return resp.Result
}
