// Implements a Liveness probe with LivenessHandler
package probes

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type LivenessHandler struct {
}


func NewLivenessHandler() *LivenessHandler {
	return &LivenessHandler{}
}

// Liveness probe .
func (p *LivenessHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Debug("Liveness Probe was called .")
	rw.WriteHeader(http.StatusOK)
}
