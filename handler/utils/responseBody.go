package handlerUtils

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

func ResponseBody(c echo.Context) (map[string]interface{}, error) {
    body := c.Request().Body
    data := make(map[string]interface{})

    err := json.NewDecoder(body).Decode(&data)

    if err != nil {
        return data, err
    }

    return data, nil
}
