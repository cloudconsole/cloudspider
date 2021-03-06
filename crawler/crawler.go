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

package crawler

import (
	"sync"

	"github.com/cloudconsole/cloudspider/log"
	"github.com/spf13/viper"
)

var wg sync.WaitGroup // Runner wait group

func Run() {
	if viper.GetBool("cloudservices.aws") {
		wg.Add(1)
		go crawlAws() // start the AWS Crawler
	}

	if viper.GetBool("cloudservices.ultradns") {
		wg.Add(1)
		go crawlUltraDns() // start the UltrDNS Crawler
	}

	if viper.GetBool("cloudservices.dnsmadeeasy") {
		wg.Add(1)
		go crawlDnsMadeEasy() // start the DNSMadeEasy Crawler
	}

	//if viper.GetBool("cloudservices.akamai") {
	//	wg.Add(1)
	//	go crawlAkamaiDns() // start the Akamai Crawler
	//}

	wg.Wait() // wait for all the crawler to finish
}

func crawlAws() {
	var cwg sync.WaitGroup // Crawler wait group

	log.Info(map[string]interface{}{
		"provider": "AWS",
	}, "Crawler started")

	regions := viper.GetStringSlice("regions")

	// start the AWS EC2 crawler
	for _, region := range regions {
		cwg.Add(1)
		log.Info(map[string]interface{}{
			"servciename": "aws_ec2",
			"region":      region,
		}, "Crawler started")
		go CrawlAllInstances(region, &cwg)
	}

	// start the AWS ELB crawler
	for _, region := range regions {
		cwg.Add(1)
		log.Info(map[string]interface{}{
			"servciename": "aws_elb",
			"region":      region,
		}, "Crawler started")
		go CrawlAllElbs(region, &cwg)
	}

	// start the AWS Route53 crawler
	cwg.Add(1)
	log.Info(map[string]interface{}{
		"servciename": "aws_route53",
		"region":      "global",
	}, "Crawler started")
	go CrawlAllRoute53(&cwg)

	cwg.Wait()
	log.Info(map[string]interface{}{
		"provider": "AWS",
	}, "Crawler finised")
	wg.Done() // say all the AWS Crawler is done
}

func crawlUltraDns() {
	var cwg sync.WaitGroup // Crawler wait group

	log.Info(map[string]interface{}{
		"provider": "UltraDns",
	}, "Crawler started")

	// start the UltraDNS crawler
	cwg.Add(1)
	CrawlUltraDns(&cwg)

	cwg.Wait()
	log.Info(map[string]interface{}{
		"provider": "UltraDns",
	}, "Crawler finished")
	wg.Done() // say UltraDns Crawler is done
}

func crawlDnsMadeEasy() {
	var cwg sync.WaitGroup // Crawler wait group

	log.Info(map[string]interface{}{
		"provider": "DnsMadeEasy",
	}, "Crawler started")

	// start the UltraDNS crawler
	cwg.Add(1)
	CrawlDnsMadeEasy(&cwg)

	cwg.Wait()
	log.Info(map[string]interface{}{
		"provider": "DnsMadeEasy",
	}, "Crawler finished")
	wg.Done() // say DNSMadeEasy Crawler is done
}

//func crawlAkamaiDns() {
//	var cwg sync.WaitGroup // Crawler wait group
//
//	log.Info(map[string]interface{}{
//		"provider": "AkamaiDNS",
//	}, "Crawler started")
//
//	// start the UltraDNS crawler
//	cwg.Add(1)
//	CrawlAkamaiDns(&cwg)
//
//	cwg.Wait()
//	log.Info(map[string]interface{}{
//		"provider": "AkamaiDNS",
//	}, "Crawler finished")
//	wg.Done() // say AkamaiDNS Crawler is done
//}
