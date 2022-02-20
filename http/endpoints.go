package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetVersion is the endpoint to get the current application version information
// @Summary get current version of the app
// @Description gets the timestamp, build version, release version, commit, branch, etc from current binary
// @Success 200 {object} Version
// @Router /version [get]
func (s *Server) GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, s.version)
}

func (s *Server) GetRefs(c *gin.Context) {
	c.JSON(http.StatusOK, s.archive.GetRefs())
}
