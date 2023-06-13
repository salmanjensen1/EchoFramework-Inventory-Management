package Response

import "github.com/labstack/echo/v4"

type SystemResponse struct {
	Status  int       `json:"Status"`
	Message string    `json:"message"`
	Data    *echo.Map `json:"data"`
}
