package httpapi

import "github.com/labstack/echo/v4"

func NewRouter(handler *Handler) *echo.Echo {
	e := echo.New()

	e.GET("/health", handler.Health)

	v1 := e.Group("/api/v1")
	v1.POST("/storyboards", handler.CreateStoryboard)
	v1.GET("/storyboards/:project_id", handler.GetStoryboard)
	v1.GET("/storyboards/:project_id/jobs/:job_id", handler.GetJobStatus)

	return e
}
