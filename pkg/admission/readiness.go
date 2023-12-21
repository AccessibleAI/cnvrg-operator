package admission

import (
	"net/http"
)

type ReadinessHandler struct{}

func NewReadinessHandler() *ReadinessHandler {
	return &ReadinessHandler{}
}

func (h *ReadinessHandler) Handler(w http.ResponseWriter, r *http.Request) {

	resp := []byte("ready")
	endWithOk(resp, w)
}

func (h *ReadinessHandler) HandlerPath() string {
	return "/ready"
}
