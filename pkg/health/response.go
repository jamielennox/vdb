package health

import "net/http"

type CheckResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (r *CheckResult) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

type Response struct {
	Success bool                   `json:"success"`
	Checks  map[string]CheckResult `json:"checks"`
}

func (r *Response) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}
