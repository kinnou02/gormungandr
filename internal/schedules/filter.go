package schedules

import (
	"net/http"

	"github.com/CanalTP/gormungandr"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func NoRouteHandler(kraken *gormungandr.Kraken, publisher Publisher) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		request_id := uuid.NewV4()
		logger := logrus.WithField("request_id", request_id)
		filter, err := gormungandr.ParsePath(c.Param("filter"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if filter.API == "route_schedules" {
			request := NewRouteScheduleRequest()
			if err := c.ShouldBindQuery(&request); err != nil {
				logger.Debugf("%+v\n", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
			request.ID = request_id
			request.ForbiddenUris = append(request.ForbiddenUris, c.QueryArray("forbidden_uris[]")...)
			request.Filters = append(request.Filters, filter.Filters...)
			if user, ok := gormungandr.GetUser(c); ok {
				request.User = user
			}
			if len(request.Filters) < 1 {
				c.JSON(http.StatusNotFound, gin.H{"message": "at least one filter is required"})
				return
			}

			request.Coverage = c.Param("coverage")
			RouteSchedule(c, kraken, &request, publisher, logger)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "API not found"})
			return
		}
	}
	return gin.HandlerFunc(fn)
}
