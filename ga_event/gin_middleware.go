package ga_event

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GALogger(trackingId string, domainName string) gin.HandlerFunc {
	gaContext := NewGAContext(trackingId, domainName)
	return func(c *gin.Context) {
		cid := ""

		for _, cookie := range c.Request.Cookies() {
			if cookie.Name == "GAT" {
				cid = cookie.Value
				break
			}
		}

		if cid == "" {
			generateUUID(&cid)
			cookie := &http.Cookie{"GAT", cid, "/", "", time.Now().Add(18000 * time.Hour), "", 18000 * 60 * 60, false, false, "", nil}
			http.SetCookie(c.Writer, cookie)
		}

		c.Next()
		err := gaContext.NewPageView(c.Request.UserAgent(), cid, c.ClientIP(), c.Request.URL.Path, "", nil).Log()

		if err != nil {
			log.Printf("Unable to log GA PageView: %v\n", err)
		}
	}
}
