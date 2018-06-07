package schedules

import (
	"github.com/CanalTP/gonavitia"
	"github.com/CanalTP/gormungandr"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

type Publisher interface {
	PublishRouteSchedule(request RouteScheduleRequest, response gonavitia.RouteScheduleResponse, c gin.Context) error
}

type AuthOption func(*gin.RouterGroup)

func Auth(authMiddleware gin.HandlerFunc) AuthOption {
	return func(group *gin.RouterGroup) {
		group.Use(authMiddleware)
	}
}
func SkipAuth() AuthOption {
	return func(group *gin.RouterGroup) {}
}

func SetupApi(router *gin.Engine, kraken *gormungandr.Kraken, statPublisher Publisher, auth AuthOption) {
	cov := router.Group("/v1/coverage/:coverage")
	auth(cov)
	cov.GET("/*filter", NoRouteHandler(kraken, statPublisher))

	router.Use(location.New(location.Config{
		Scheme: "http",
		Host:   "navitia.io",
	}))
}
