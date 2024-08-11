package utils

import echo "github.com/labstack/echo/v4"

func SetChunkedEnc(ectx echo.Context) {
	resp := ectx.Response()
	header := resp.Header()
	if header.Get("Transfer-Encoding") != "" {
		header.Set("Transfer-Encoding", "chunked")
	}
}
