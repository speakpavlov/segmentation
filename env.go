package segmentation

type Env map[string]interface{}

func NewEnv(env map[string]interface{}) Env {
	//init fresh cachedRequest
	env["cachedRequestObj"] = &CachedRequest{}

	return env
}

// Request func
func (e Env) GetRequest(url string) interface{} {
	return e["cachedRequestObj"].(*CachedRequest).makeCachedRequest(url)
}
