package configuration

import (
	"github.com/labstack/echo"
	servestatic "github.com/ncpa0/htmx-framework/serve-static"
)

type Configuration struct {
	// When enabled, the `.html` extension will be stripped from the URL pathnames.
	StripExtension bool
	// When enabled, the server will print debug information to the console.
	DebugMode bool
	// The entrypoint file containing the JSX pages used to generate the views html files.
	//
	// Defaults to `index.tsx`.
	Entrypoint string
	// The directory to which output the generated html files, and from which those will be hosted.
	//
	// Defaults to `views`.
	ViewsDir string
	// The directory to which output the static files, and from which those will be hosted.
	//
	// Defaults to `static`.
	StaticDir string
	// The URL path from under which the static files will be hosted.
	//
	// Defaults to `/static`.
	StaticURL        string
	BeforeStaticSend func(resp *servestatic.StaticResponse, c echo.Context) error
}

var Current *Configuration = &Configuration{
	StripExtension: false,
	DebugMode:      false,
	Entrypoint:     "index.tsx",
	ViewsDir:       "views",
	StaticDir:      "static",
	StaticURL:      "/static",
}
