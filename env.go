package segmentation

import (
	"strconv"
	"strings"
)

type Env map[string]interface{}
type CurlFunc func(url string) interface{}

func NewEnv(in map[string]interface{}) Env {
	//init cache
	cache := &CachedRequest{}

	//add cached curl
	env := Env{
		"Curl": func(url string) interface{} {
			return cache.curl(url)
		},
	}

	//merge
	for key, value := range in {
		env[key] = value
	}

	return env
}

// Curl to external source
func (e Env) CurlWithUserId(url string, userId int) interface{} {
	replacedUrl := strings.Replace(url, "{userId}", strconv.Itoa(userId), 1)

	return e["Curl"].(CurlFunc)(replacedUrl)
}
