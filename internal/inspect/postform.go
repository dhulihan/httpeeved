package inspect

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func PostForm(c *gin.Context, resp gin.H) {
	// form data
	err := c.Request.ParseForm()
	if err != nil {
		log.Trace(err.Error())
	}

	if len(c.Request.PostForm) > 0 {
		for k, values := range c.Request.PostForm {
			for _, v := range values {
				log.WithField(k, v).Debug("form value")
			}
		}

		resp["form"] = c.Request.PostForm
	}
}
