package http

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(handler *Handler) *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	}))

	e.POST("/team/add", handler.CreateTeam)
	e.GET("/team/get", handler.GetTeam)

	e.POST("/users/setIsActive", handler.SetUserIsActive)
	e.GET("/users/getReview", handler.GetReviewerPRs)

	e.POST("/pullRequest/create", handler.CreatePR)
	e.POST("/pullRequest/merge", handler.MergePR)
	e.POST("/pullRequest/reassign", handler.ReassignReviewer)

	e.GET("/stats/users", handler.GetUserStats)
	e.GET("/stats/prs", handler.GetPRStats)
	e.GET("/stats/workload", handler.GetReviewerWorkload)

	return e
}
