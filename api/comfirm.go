package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/labstack/echo/v4"
)

func (cli *Client) Confirm(c echo.Context) error {
	key := c.QueryParam("key")
	fmt.Println(key)

	v, isThere := cli.Setting.Objects[key]
	if !isThere {
		return c.JSON(http.StatusOK, map[string]any{
			"msg": "key defined",
		})
	}

	var (
		diff_x = v.Area.MaxX() - v.Area.MinX()
		diff_y = v.Area.MaxY() - v.Area.MinY()
		rx, ry int
	)

	for i := 0; i < 10; i++ {
		rx = rand.Intn(diff_x) - 1
		ry = rand.Intn(diff_y) - 1

		robotgo.Move(v.Area.MinX()+rx, v.Area.MinY()+ry)
		time.Sleep(500 * time.Millisecond)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"msg": "done",
	})
}
