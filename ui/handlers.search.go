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

package ui

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/cloudconsole/cloudspider/storage"
)

func showSearchResults(c *gin.Context) {
	q := c.Request.URL.Query()
	sQuery := q["search_query"][0]

	// Check if search docs exists
	if hosts, err := storage.GetHosts(sQuery); err == nil {
		// Call the render function with the payload and the name of the
		// template
		render(c, gin.H{
			"s_query": sQuery,
			"hosts": hosts,
		}, "search-results.html")
	} else {
		// If the host is not found, abort with an error
		c.AbortWithError(http.StatusNotFound, err)
	}
}
