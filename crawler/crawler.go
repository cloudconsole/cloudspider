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
		go crawlAWS() // start the AWS Crawler
	}

	if viper.GetBool("cloudservices.ultradns") {
		wg.Add(1)
		go crawlUltraDns() // start the UltrDNS Crawler
	}

	if viper.GetBool("cloudservices.dnsmadeeasy") {
		wg.Add(1)
		go crawlDnsMadeEasy() // start the DNSMadeEasy Crawler
	}

	wg.Wait() // wait for all the crawler to finish
}

func crawlAWS() {
	var cwg sync.WaitGroup // Crawler wait group

	log.Info(map[string]interface{}{
		"provider": "AWS",
	}, "Crawler started")

	regions := viper.GetStringSlice("regions")

	// start the AWS EC2 crawler
	for _, region := range regions {
		cwg.Add(1)
		go CrawlAllInstances(region, &cwg)
	}

	// start the AWS ELB crawler
	for _, region := range regions {
		cwg.Add(1)
		go CrawlAllElbs(region, &cwg)
	}

	// start the AWS Route53 crawler
	cwg.Add(1)
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
	CrawlUltraDNS(&cwg)

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
	CrawlDNSMadeEasy(&cwg)

	cwg.Wait()
	log.Info(map[string]interface{}{
		"provider": "DnsMadeEasy",
	}, "Crawler finished")
	wg.Done() // say DNSMadeEasy Crawler is done
}
