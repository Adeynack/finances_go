package routes

import "github.com/labstack/echo/v4"

type HandlerGenerator[SpecificContext any] func(func(echo.Context, SpecificContext) error) echo.HandlerFunc

func WithPainter[SpecificContext any](
	group *echo.Group,
	generateHandler HandlerGenerator[SpecificContext],
	painterFunc func(r RequestContextPainter[SpecificContext]),
) {
	painter := RequestContextPainter[SpecificContext]{
		group:           group,
		generateHandler: generateHandler,
	}
	painterFunc(painter)
}

type RequestContextPainter[SpecificContext any] struct {
	group           *echo.Group
	generateHandler func(func(echo.Context, SpecificContext) error) echo.HandlerFunc
}

func (painter *RequestContextPainter[SpecificContext]) GET(
	path string,
	handler func(echo.Context, SpecificContext) error,
) {
	painter.group.GET(path, painter.generateHandler(handler))
}
