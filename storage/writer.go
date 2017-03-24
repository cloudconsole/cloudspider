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

package storage

import (
	"sync"
	"strings"
	"gopkg.in/mgo.v2"
	"github.com/cloudconsole/cloudspider/log"
)

// Write many docs in to mongodb
func InsertMany(Docs []interface{}, collName string, wg *sync.WaitGroup) {
	coll := DBConnection().C(collName)

	for _, doc := range Docs {
		if err := coll.Insert(doc); err != nil {
			if !strings.Contains(err.Error(), "E11000") {
				log.Fatal(map[string]interface{}{
					"DBError": err,
				}, "Database error")
			}
		}
		wg.Done() // say write is done
	}
}

// Ensure index exits
func EnsureHostIndex() error {
	coll := DBConnection().C("hosts")

	return coll.EnsureIndex(
		mgo.Index{
			Key: []string{
				"$text:ssh_key_name",
				"$text:type",
				"$text:public_dns",
				"$text:private_dns",
				"$text:tags",
				"$text:security_group",
			},
			Name:       "host_index",
			Background: true,
		})
}
