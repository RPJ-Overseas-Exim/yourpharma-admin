package handlerUtils

import (
	"encoding/json"
	"io"

	"github.com/labstack/echo/v4"
)

func ResponseBody(c echo.Context) (map[string]interface{}, error) {
    body := c.Request().Body

    return GetJson(body)
}

func GetJson(jsonData io.ReadCloser) (map[string]interface{}, error) {
    data := make(map[string]interface{})
    err := json.NewDecoder(jsonData).Decode(&data)

    if err != nil {
        return data, err
    }

    return data, nil
}
