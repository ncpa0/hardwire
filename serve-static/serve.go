package servestatic

import (
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/ncpa0/htmx-framework/utils"
)

type StaticFile struct {
	Path              string
	RelPath           string
	Content           []byte
	ContentType       string
	LastModifiedAt    *time.Time
	LastModifiedAtRFC string
	Etag              string
}

var staticFiles []StaticFile = []StaticFile{}

func (f *StaticFile) Revalidate() error {
	// check if the file has changed since last time
	// and reload it if it has
	file, err := os.Open(f.Path)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		return err
	}

	modTime := info.ModTime()
	if modTime.Equal(*f.LastModifiedAt) {
		return nil
	}

	buff := make([]byte, info.Size())
	_, err = file.Read(buff)

	if err != nil {
		return err
	}

	f.Content = buff
	f.Etag = utils.HashBytes(buff)
	f.LastModifiedAt = &modTime
	f.LastModifiedAtRFC = modTime.Format(http.TimeFormat)
	f.ContentType = http.DetectContentType(buff)

	return nil
}

func getStaticFile(filepath string) ([]byte, string, *time.Time, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, "", nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, "", nil, err
	}

	buff := make([]byte, info.Size())
	_, err = file.Read(buff)

	if err != nil {
		return nil, "", nil, err
	}

	modTime := info.ModTime()
	return buff, http.DetectContentType(buff), &modTime, err
}

type Configuration struct {
	BeforeSend func(c echo.Context) error
}

func Serve(server *echo.Echo, baseUrl string, root string, conf *Configuration) {
	if root[len(root)-1] != '/' {
		root += "/"
	}

	utils.Walk(root, func(root string, dirs []string, files []string) error {
		for _, file := range files {
			filepath := path.Join(root, file)
			relativePath := filepath[len(root):]
			content, ctype, modTime, err := getStaticFile(filepath)

			if err == nil {
				staticFiles = append(staticFiles, StaticFile{
					Path:              filepath,
					RelPath:           relativePath,
					Content:           content,
					ContentType:       ctype,
					Etag:              utils.HashBytes(content),
					LastModifiedAt:    modTime,
					LastModifiedAtRFC: modTime.Format(http.TimeFormat),
				})
			}
		}
		return nil
	})

	server.GET(baseUrl+"/*", func(c echo.Context) error {
		routePath := c.Param("*")
		for _, file := range staticFiles {
			if file.RelPath == routePath {
				return sendFile(file, c, conf)
			}
		}

		// check if files exists in fs, and if it does load it into memory
		// and serve it
		filepath := path.Join(root, routePath)
		content, ctype, modTime, err := getStaticFile(filepath)

		if err == nil {
			staticFiles = append(staticFiles, StaticFile{
				Path:              filepath,
				RelPath:           routePath,
				Content:           content,
				ContentType:       ctype,
				Etag:              utils.HashBytes(content),
				LastModifiedAt:    modTime,
				LastModifiedAtRFC: modTime.Format(http.TimeFormat),
			})

			return sendFile(staticFiles[len(staticFiles)-1], c, conf)
		}

		return c.String(404, "Not found")
	})
}

func sendFile(file StaticFile, c echo.Context, conf *Configuration) error {
	file.Revalidate()

	if c.Request().Header.Get("If-None-Match") == file.Etag {
		return c.NoContent(304)
	}

	h := c.Response().Header()
	h.Set("Content-Type", file.ContentType)
	h.Set("Last-Modified", file.LastModifiedAtRFC)
	h.Set("Date", time.Now().Format(http.TimeFormat))
	h.Set("Cache-Control", "public, max-age=2592000")
	h.Set("Accept-Ranges", "bytes")
	h.Set("ETag", strconv.FormatInt(file.LastModifiedAt.Unix(), 10))

	if conf.BeforeSend != nil {
		err := conf.BeforeSend(c)
		if err != nil {
			return err
		}
	}

	requestedRange := utils.ParseRangeHeader(&h)
	if requestedRange != nil {
		contentLength := strconv.FormatInt(requestedRange.End-requestedRange.Start+1, 10)
		contentRange := ("bytes " +
			strconv.FormatInt(requestedRange.Start, 10) +
			"-" + strconv.FormatInt(requestedRange.End, 10) +
			"/" + strconv.FormatInt(int64(len(file.Content)), 10))
		h.Set("Content-Length", contentLength)
		h.Set("Content-Range", contentRange)

		if conf.BeforeSend != nil {
			err := conf.BeforeSend(c)
			if err != nil {
				return err
			}
		}

		return c.Blob(200, file.ContentType, file.Content[requestedRange.Start:requestedRange.End+1])
	} else {
		if conf.BeforeSend != nil {
			err := conf.BeforeSend(c)
			if err != nil {
				return err
			}
		}

		return c.Blob(200, file.ContentType, file.Content)
	}
}
