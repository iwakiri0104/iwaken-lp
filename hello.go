// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hello

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"net/http"
	"os"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func MethodOverride(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method == "POST" {
			method := c.Request().PostFormValue("_method")
			if method == "PUT" || method == "PATCH" || method == "DELETE" {
				c.Request().Method = method
			}
		}
		return next(c)
	}
}

func main() {

	list, err := template.New("t").ParseGlob("views/template/*.html")

	t := &Template{
		templates: template.Must(list, err),
	}

	e := echo.New()

	//メソッド対応処理
	e.Pre(MethodOverride)
	e.Renderer = t

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	//自動マイグレーション
	//models.Init()

	//db := models.DatabaseConnection()
	//routerf
	e.GET("/", Index)
	e.GET("/liff", Liff)
	e.GET("/detail", Detail)

	e.Static("/assets/css", "./views/assets/css")
	e.Static("/assets/js", "./views/assets/js")
	e.Static("/assets/scss", "./views/assets/scss")
	e.Static("/assets/img", "./views/assets/img")
	e.Static("/assets/vendor", "./views/assets/vendor")

	// Determine port for HTTP service...
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	// start server
	e.Logger.Fatal(e.Start(port))
}

// index
func Index(c echo.Context) error {

	return c.Render(http.StatusOK, "index", nil)
}

// index
func Detail(c echo.Context) error {
	return c.Render(http.StatusOK, "detail", nil)
}

// index
func Liff(c echo.Context) error {
	return c.Render(http.StatusOK, "liff", nil)
}
