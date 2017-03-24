// Copyright © 2016 Ashok Raja <ashokraja.r@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// main.go

package ui

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/ekyoung/gin-nice-recovery"
)

var router *gin.Engine

func Run() {

	// Set the router as the default one provided by Gin
	router = gin.New()

	router.Use(gin.Logger()) // Install the default logger, not required

	// Install nice.Recovery, passing the handler to call after recovery
	router.Use(nice.Recovery(recoveryHandler))

	// Serving static files
	router.Static("/static", "./ui/static")

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob("ui/templates/*")

	// Initialize the routes
	initializeRoutes()

	// Start serving the application
	router.Run()

}

func recoveryHandler(c *gin.Context, err interface{}) {
	c.HTML(http.StatusInternalServerError, "panic.html", gin.H{
		"err": err,
	})
}
