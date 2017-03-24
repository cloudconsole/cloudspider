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

// handlers.article.go

package ui

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cloudconsole/cloudspider/storage"
	"strings"
)

func describeHost(c *gin.Context) {
	sQuery := c.Param("s_query")

	// Check if host exists
	if host, err := getHost(sQuery); err == nil {
		// fetch all the lbs & dns mapped with the host
		lbs, err := storage.GetLbsByHostId(host.ID)
		dns, err := storage.GetDnsByFqdn(host.PublicDns)

		// fetch the dns mapping for the host dns
		// this hack is required to fetch multiple cname mapping
		for _, dn := range *dns {
			mdns, _ := storage.GetDnsByFqdn(dn.Name)
			for _, val := range *mdns {
				*dns = append(*dns, val)
			}
		}

		// fetch all the dns name mapping with the lbs
		for _, lb := range *lbs {
			lbdns, _ := storage.GetDnsByFqdn(lb.PublicDns)
			for _, val := range *lbdns {
				*dns = append(*dns, val)
				// further fetch the dns mapping for the lbs dns
				// this hack is required to fetch multiple cname mapping
				mdns, _ := storage.GetDnsByFqdn(val.Name)
				for _, val := range *mdns {
					*dns = append(*dns, val)
				}
			}
		}

		if err == nil {
			// Call the render function with the payload and the name of the
			// template
			render(c, gin.H{
				"s_query": sQuery,
				"host": host,
				"lbs": lbs,
				"dns": dns,
			}, "describe-instance.html")
		} else {
			// If the elb/dns is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
		}
	} else {
		// If the host is not found, abort with an error
		c.AbortWithError(http.StatusNotFound, err)
	}
}

// Fetch host information
func getHost(query string) (*storage.NodeDoc, error) {
	switch  {
	case strings.HasPrefix(query, "i-"):
		return storage.GetHostById(query)
	case strings.HasPrefix(query, "ec2-"):
		return storage.GetHostByFqdn(query)
	}

	return nil, nil
}
