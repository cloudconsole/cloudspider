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

// routes.go

package ui

func initializeRoutes() {

	// Handle the index route
	router.GET("/", showIndexPage)

	// Handle the settings route
	router.GET("/settings", showIndexPage)

	// Handle GET requests at /host/search_query
	router.GET("/host/:s_query", describeHost)

	// Handle GET requests at /search?search_query=search_string
	router.GET("/search", showSearchResults)

	//// Group user related routes together
	//userRoutes := router.Group("/u")
	//{
	//	// Handle the GET requests at /u/register
	//	// Show the registration page
	//	// Ensure that the user is not logged in by using the middleware
	//	userRoutes.GET("/register", showRegistrationPage)
	//
	//	// Handle POST requests at /u/register
	//	// Ensure that the user is not logged in by using the middleware
	//	userRoutes.POST("/register", register)
	//}

}
