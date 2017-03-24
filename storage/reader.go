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
	"github.com/cloudconsole/cloudspider/log"
	"gopkg.in/mgo.v2/bson"
)

// Checks if the document exists with the given ID
func DocExists(id string, collName string) bool {
	coll := DBConnection().C(collName)

	query := coll.FindId(id)
	count, err := query.Count()
	if err != nil {
		log.Fatal(map[string]interface{}{
			"DBError": err,
		}, "Database error")
	}

	if count >= 1 {
		return true
	} else {
		return false
	}
}

// Get the host document by id
func GetHosts(query string) (*[]NodeDoc, error) {
	coll := DBConnection().C("hosts")

	result := []NodeDoc{}
	coll.Find(bson.M{"$text": bson.M{"$search": query}}).All(&result)

	return &result, nil
}

// Get the host document by id
func GetHostById(id string) (*NodeDoc, error) {
	coll := DBConnection().C("hosts")

	result := NodeDoc{}
	coll.Find(bson.M{"_id": id}).One(&result)

	return &result, nil
}

// Get the host document by FQDN
func GetHostByFqdn(fqdn string) (*NodeDoc, error) {
	coll := DBConnection().C("hosts")

	result := NodeDoc{}
	coll.Find(bson.M{"public_dns": fqdn}).One(&result)

	return &result, nil
}

// Get the loadbalancer by host id
func GetLbsByHostId(id string) (*[]LBDoc, error) {
	coll := DBConnection().C("lbs")

	result := []LBDoc{}
	coll.Find(bson.M{"backends": id}).All(&result)

	return &result, nil
}

// Get the DNS by host fqdn
func GetDnsByFqdn(fqdn string) (*[]DNSDoc, error) {
	coll := DBConnection().C("dns")

	result := []DNSDoc{}
	coll.Find(bson.M{"records": fqdn}).All(&result)

	return &result, nil
}
