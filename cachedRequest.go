package segmentation

import (
	"encoding/json"
	"io/ioutil"
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

	cachedRequest[url] = makeGetRequest(url)

	return cachedRequest[url]
}

func makeGetRequest(url string) interface{} {
	response, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic("Cannot make request: " + url + ", err:" + err.Error())
	}

	var resp HttpResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic("Cannot unmarshal: " + url + ", err:" + err.Error())
	}

	return resp.Result
}
