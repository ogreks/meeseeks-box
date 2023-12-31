package middleware

import (
	"github.com/gin-gonic/gin"
)

func SkippedPathPrefixes(c *gin.Context, prefixes ...string) bool {
	if len(prefixes) == 0 {
		return false
	}

	path := c.Request.URL.Path
	pathLen := len(path)
	for _, p := range prefixes {
		if pl := len(p); pl <= pathLen && path[:pl] == p {
			return true
		}
	}

	return false
}
