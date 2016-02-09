package util

func CreateOkResponse(data interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	m["ok"] = true
	if data != nil {
		m["data"] = data
	}
	return m
}

func CreateErrResponse(err error) map[string]interface{} {
	m := make(map[string]interface{})
	m["ok"] = false
	m["error"] = err.Error()
	return m
}
