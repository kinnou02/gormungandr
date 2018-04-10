package main

import (
	"net/http"

	"github.com/CanalTP/gormungandr"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NoRouteHandler(kraken *gormungandr.Kraken) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		filter, err := gormungandr.ParsePath(c.Param("filter"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		if filter.Api == "route_schedules" {
			request := NewRouteScheduleRequest()
			if err := c.ShouldBindQuery(&request); err != nil {
				logrus.Debugf("%+v\n", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
			request.ForbiddenUris = append(request.ForbiddenUris, c.QueryArray("forbidden_uris[]")...)
			request.Filters = append(request.Filters, filter.Filters...)
			RouteSchedule(c, kraken, &request)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "api not found"})

		}
	}
	return gin.HandlerFunc(fn)
}
