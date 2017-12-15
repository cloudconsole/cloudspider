// handlers.common_test.go
//
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
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

// Test that a GET request to the home page returns the home page with
// the HTTP code 200 for an unauthenticated user
func TestShowIndexPageUnauthenticated(t *testing.T) {
    r := getRouter(true)

    r.GET("/", showIndexPage)

    // Create a request to send to the above route
    req, _ := http.NewRequest("GET", "/", nil)

    testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
        // Test that the http status code is 200
        statusOK := w.Code == http.StatusOK

        // Test that the page title is "Home Page"
        // You can carry out a lot more detailed tests using libraries that can
        // parse and process HTML pages
        p, err := ioutil.ReadAll(w.Body)
        pageOK := err == nil && strings.Index(string(p), "<title>Cloud Spider</title>") > 0

        return statusOK && pageOK
    })
}
