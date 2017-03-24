package crawler

import (
	"github.com/cloudconsole/cloudspider/log"
	"github.com/cloudconsole/ultradns"
	"github.com/spf13/viper"
	"strings"
	"sync"

	"github.com/cloudconsole/cloudspider/storage"
)

// Extract only the required fields and discard the unwanted fields
func PruneUltraDnsFields(dnsRec ultradns.RRSet) storage.DNSDoc {
	doc := new(storage.DNSDoc)

	doc.CloudProvider = "UltraDNS"
	doc.ID = strings.TrimSuffix(dnsRec.RecName, ".")
	doc.Name = strings.TrimSuffix(dnsRec.RecName, ".")
	rt := dnsRec.RType
	doc.Type = rt[0 : len(rt)-4]

	if dnsRec.RRecords != nil {
		for _, record := range dnsRec.RRecords {
			doc.Records = append(doc.Records, strings.TrimSuffix(record, "."))
		}
	}

	return *doc
}

func CrawlUltraDns(cwg *sync.WaitGroup) {
	userName := viper.GetString("ultradns.username")
	password := viper.GetString("ultradns.password")
	var wwg sync.WaitGroup // Writer wait group
	totAdd := 0
	totUpd := 0

	conn := ultradns.NewSession()
	err := conn.Authenticate(userName, password)
	if err != nil {
		log.Error(map[string]interface{}{}, err.Error())
	}

	zones, err := conn.GetAllZones()
	if err != nil {
		log.Error(map[string]interface{}{}, err.Error())
	}

	for _, zone := range zones {
		zoneName := strings.TrimSuffix(zone.Property.ZName, ".")
		offset := 0
		limit := 100
		resCount := 0

		for {
			recs, resInfo, err := conn.GetRRsets(zoneName, offset, limit)
			if err != nil {
				log.Error(map[string]interface{}{}, err.Error())
			}
			resCount += resInfo.RetCount
			offset += limit

			var aDocs []interface{}
			var uDocs []interface{}

			for _, rec := range recs {
				rt := rec.RType
				if rt[0:len(rt)-4] == "CNAME" || rt[0:len(rt)-4] == "A" {
					dID := strings.TrimSuffix(rec.RecName, ".")
					if storage.DocExists(dID, "dns") {
						// update if dns record already in DB
						uDocs = append(uDocs, PruneUltraDnsFields(rec))
					} else {
						aDocs = append(aDocs, PruneUltraDnsFields(rec))
					}
				}
			}

			noADocs := len(aDocs)
			noUDocs := len(uDocs)
			totAdd += noADocs
			totUpd += noUDocs

			if noADocs != 0 {
				wwg.Add(noADocs)
				go storage.InsertMany(aDocs, "dns", &wwg)
			}

			if resCount == resInfo.Total {
				break
			}
		}
	}

	wwg.Wait() // wait for all the writer to finish

	log.Debug(map[string]interface{}{
		"servciename":  "ultradns",
		"totalAdded":   totAdd,
		"totalUpdated": totUpd,
	}, "Crawler finished")
	cwg.Done() // say UltraDNS crawler is done
}
