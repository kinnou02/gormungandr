package schedules

import (
	"net/http"

	"github.com/CanalTP/gonavitia"
	"github.com/CanalTP/gormungandr"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type Publisher interface {
	PublishRouteSchedule(request RouteScheduleRequest, response gonavitia.RouteScheduleResponse, c echo.Context) error
}

func NoRouteHandler(kraken *gormungandr.Kraken, publisher Publisher) echo.HandlerFunc {
	fn := func(c echo.Context) error {
		request_id := uuid.NewV4()
		logger := logrus.WithField("request_id", request_id)
		filter, err := gormungandr.ParsePath(c.Param("*"))
		if err != nil {
			//c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if filter.API == "route_schedules" {
			request := NewRouteScheduleRequest()
			if err := c.Bind(&request); err != nil {
				logger.Debugf("%+v\n", err)
				//c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return err
			}
			query := c.Request().URL.Query() // Parse only once
			request.ID = request_id
			request.ForbiddenUris = append(request.ForbiddenUris, query["forbidden_uris[]"]...)
			request.Filters = append(request.Filters, filter.Filters...)
			/*
				if user, ok := gormungandr.GetUser(c); ok {
					request.User = user
				}
			*/
			if len(request.Filters) < 1 {
				return echo.NewHTTPError(http.StatusNotFound, "at least one filter is required")
			}

			request.Coverage = c.Param("coverage")
			return RouteSchedule(c, kraken, &request, publisher, logger)
		} else {
			//c.JSON(http.StatusNotFound, gin.H{"error": "API not found"})
			return echo.NewHTTPError(http.StatusNotFound, "API not Found")
		}
	}
	return echo.HandlerFunc(fn)
}
