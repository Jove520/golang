package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

func TrimSpaceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		trimQuery(c)
		trimJSONBody(c)

		c.Next()
	}
}

func trimQuery(c *gin.Context) {
	params := c.Request.URL.Query()

	for key, values := range params {
		for i, value := range values {
			params[key][i] = strings.TrimSpace(value)
		}
	}

	c.Request.URL.RawQuery = params.Encode()
}

func trimJSONBody(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil || len(bodyBytes) == 0 {
		return
	}

	var body any
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		return
	}

	trimJSONStrings(body)

	newBody, err := json.Marshal(body)
	if err != nil {
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(newBody))
}

func trimJSONStrings(value any) {
	switch v := value.(type) {
	case map[string]any:
		for key, item := range v {
			if text, ok := item.(string); ok {
				v[key] = strings.TrimSpace(text)
				continue
			}

			trimJSONStrings(item)
		}
	case []any:
		for i, item := range v {
			if text, ok := item.(string); ok {
				v[i] = strings.TrimSpace(text)
				continue
			}

			trimJSONStrings(item)
		}
	}
}
