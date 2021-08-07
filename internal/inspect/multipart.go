package inspect

import (
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	MAX_UPLOAD_SIZE = 5 * 1024 * 1024
)

func MultipartForm(c *gin.Context, resp gin.H) {
	// handle multi-part form
	err := c.Request.ParseMultipartForm(MAX_UPLOAD_SIZE)
	if err != nil {
		log.Trace(err.Error())
	}

	if c.Request.MultipartForm != nil {
		resp["multipart"] = c.Request.MultipartForm.Value
		log.Debug("multipart detected")
		if len(c.Request.MultipartForm.Value) > 0 {
			for k, values := range c.Request.MultipartForm.Value {
				for _, v := range values {
					log.WithField(k, v).Debug("multipart form value")
				}
			}

		}

		if len(c.Request.MultipartForm.File) > 0 {
			for k, values := range c.Request.MultipartForm.File {
				for i, v := range values {
					log.WithFields(log.Fields{
						"key":     k,
						"i":       i,
						"summary": multipartFileSummary(v),
					}).Debug("multipart file detected")
					resp[fmt.Sprintf("multipart-form-file-%s-%d", k, i)] = multipartFileSummary(v)
				}
			}

		}
	}

	return
}

func multipartFileSummary(fh *multipart.FileHeader) string {
	return fmt.Sprintf("%s %s %d", fh.Filename, fh.Header, fh.Size)
}
