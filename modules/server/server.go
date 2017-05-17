package server

import (
	"fmt"
	"github.com/amlun/linda/linda"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func Start(linda *linda.Linda) {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(Recovery())

	var api = api{
		linda: linda,
	}

	r.GET("/api/ping", api.ping())
	r.GET("/api/tasks", api.tasks())
	r.GET("/api/queues", api.queuesStatus())
	r.GET("/api/job", api.getJob())
	r.POST("/api/task", api.pushTask())

	r.Run(":8081")
}

func Recovery() gin.HandlerFunc {
	return RecoveryWithJson()
}

func RecoveryWithJson() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			var result Result
			if err := recover(); err != nil {
				result.Code = http.StatusInternalServerError
				result.Msg = fmt.Sprintf("%s", err)
				c.JSON(500, result)
			}
		}()
		c.Next()
	}
}
