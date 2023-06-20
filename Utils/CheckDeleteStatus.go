package Utils

import "github.com/labstack/echo/v4"

func CheckDeleteStatus(c echo.Context, deleteStatus bool) error {
	if !deleteStatus {
		return nil
	}
	return c.String(500, "User is deleted")
}
