package segmentation

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type CachedRequest map[string]interface{}

type HttpResponse struct {
	Response interface{} `json:"response"`
}

func (cache CachedRequest) curl(url string) interface{} {
	if value, ok := cache[url]; ok {
		return value
	}

	//make request
	response, err := http.Get(url)
	if err != nil {
		//todo log
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	var resp HttpResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic(err)
	}

	//set to cache
	cache[url] = resp.Response

	return resp.Response
}
