package utils

import (
	"fmt"

	"github.com/labstack/echo"
)

type RequestError struct {
	Code int
	Data string
}

func (err *RequestError) Error() string {
	return fmt.Sprintf("%d: %s", err.Code, err.Data)
}

func (err *RequestError) SendResponse(c echo.Context) error {
	return c.String(err.Code, err.Data)
}

type Sender interface {
	SendResponse(c echo.Context) error
}

func HandleError(c echo.Context, err error) error {
	sender, ok := err.(Sender)
	fmt.Println(err)
	if ok {
		return sender.SendResponse(c)
	}
	return c.String(500, "Internal Server Error")
}
