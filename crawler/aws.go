// Package crawler provides a spider which crawls all the cloud services
// and learn about the infrastructure configs
package crawler

import (
	"sync"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/route53"

	"github.com/cloudconsole/cloudspider/storage"
	"github.com/cloudconsole/cloudspider/log"
)


// Extract only the required fields and discard the unwanted fields
func PruneEC2Fields(inst *ec2.Instance) storage.MachineDoc {
	doc := new(storage.MachineDoc)

	doc.CloudProvider = "Amazon"
	doc.ID = *inst.InstanceId
	doc.State = *inst.State.Name
	doc.Virtualization = *inst.VirtualizationType
	doc.Architecture = *inst.Architecture
	doc.RootDevice = *inst.RootDeviceType
	doc.Type = *inst.InstanceType
	doc.DataCenter = *inst.Placement.AvailabilityZone
	doc.Tags = inst.Tags

	if inst.SecurityGroups != nil {
		for i := 0; i < len(inst.SecurityGroups); i++ {
			grpName := *inst.SecurityGroups[i].GroupName
			doc.SecurityGroup = append(doc.SecurityGroup, grpName)
		}
	}

	if inst.KeyName != nil {
		doc.SshKeyName = *inst.KeyName
	}

	if inst.IamInstanceProfile != nil {
		doc.IamProfile = *inst.IamInstanceProfile.Arn
	}

	if inst.LaunchTime != nil {
		doc.LaunchTime = *inst.LaunchTime
	}

	if inst.PrivateDnsName != nil {
		doc.PrivateDns = *inst.PrivateDnsName
	}

	if inst.PrivateIpAddress != nil {
		doc.PrivateIp = *inst.PrivateIpAddress
	}

	if inst.PublicDnsName != nil {
		doc.PublicDns = *inst.PublicDnsName
	}

	if inst.PublicIpAddress != nil {
		doc.PublicIp = *inst.PublicIpAddress
	}

	//b, err := json.Marshal(doc)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(b))

	return *doc
}

// Extract only the required fields and discard the unwanted fields
func PruneELBFields(elb *elb.LoadBalancerDescription) storage.LoadBalancerDoc {
	doc := new(storage.LoadBalancerDoc)

	doc.CloudProvider = "Amazon"
	doc.ID = *elb.LoadBalancerName
	doc.Name = *elb.LoadBalancerName
	doc.PublicDns = *elb.DNSName

	for _, dc := range elb.AvailabilityZones {
		doc.DataCenter = append(doc.DataCenter, *dc)
	}

	for _, instance := range elb.Instances {
		doc.Backends = append(doc.Backends, *instance.InstanceId)
	}

	return *doc
}

// Extract only the required fields and discard the unwanted fields
func PruneR53Fields(dnsRec *route53.ResourceRecordSet) storage.DNSDoc {
	doc := new(storage.DNSDoc)

	doc.CloudProvider = "Amazon"
	doc.ID = strings.TrimSuffix(*dnsRec.Name, ".")
	doc.Name = strings.TrimSuffix(*dnsRec.Name, ".")
	doc.Type = *dnsRec.Type

	if dnsRec.AliasTarget != nil {
		doc.Records = append(doc.Records, *dnsRec.AliasTarget.DNSName)
	}

	if dnsRec.ResourceRecords != nil {
		for _, record := range dnsRec.ResourceRecords {
			doc.Records = append(doc.Records, strings.TrimSuffix(*record.Value, "."))
		}
	}

	return *doc
}

// Checks if the document already exists in the list
func CheckDNSRecordExists(docs []storage.DNSDoc, newDoc storage.DNSDoc) bool {
	var exists = false

	for _, doc := range docs {
		if newDoc.ID == doc.ID {
			exists = true
			break
		}
	}

	return exists
}

// Merge the redundant/duplicate document records
func MergeDnsRecords(docs []storage.DNSDoc, dupDocs []storage.DNSDoc) []interface{} {
	var fdocs []interface{}

	for _, dup := range dupDocs {
		for _, doc := range docs {
			if dup.ID == doc.ID {
				mdoc := doc
				for _, record := range dup.Records {
					mdoc.Records = append(mdoc.Records, record)
				}
				fdocs = append(fdocs, mdoc)
				break
			} else {
				fdocs = append(fdocs, doc)
			}
		}
	}

	return fdocs
}

// Crawl all the machines/instances/hosts/servers running on a AWS cloud
func CrawlAllInstances(region string, cwg *sync.WaitGroup) {
	conn := ec2.New(session.New(), &aws.Config{Region: &region})
	totAdd := 0
	totUpd := 0
	var wwg sync.WaitGroup  // Writer wait group

	resp, err := conn.DescribeInstances(nil)
	if err != nil {
		log.Fatal(map[string]interface{}{}, err.Error())
	}

	for _, instances := range resp.Reservations {
		var aDocs []interface{}
		var uDocs []interface{}
		for _, instance := range instances.Instances {
			if instance.InstanceLifecycle == nil {
				if *instance.State.Name == "running" {
					if storage.DocExists(*instance.InstanceId, "machines") {
						uDocs = append(uDocs, PruneEC2Fields(instance))
					} else {
						aDocs = append(aDocs, PruneEC2Fields(instance))
					}
				}
			}
		}

		noADocs := len(aDocs)
		noUDocs := len(uDocs)
		totAdd += noADocs
		totUpd += noUDocs
		// totalInstance += numUpdateDocs
		if noADocs != 0 {
			wwg.Add(noADocs)
			go storage.InsertMany(aDocs, "machines", &wwg)
		}
	}

	wwg.Wait()  // wait for all the writer to finish

	log.Debug(map[string]interface{}{
		"servciename": "aws_ec2",
		"totalAdded": totAdd,
		"totalUpdated": totUpd,
		"region": region,
	}, "Crawled finished")

	// Ensure machines collection index has been created
	eerr := storage.EnsureMachinesIndex()
	if eerr != nil {
		log.Error(map[string]interface{}{}, eerr.Error())
	}

	cwg.Done()  // say AWS EC2 crawler is done
}

//Crawl all the loadbalncers in a AWS cloud
func CrawlAllElbs(region string, cwg *sync.WaitGroup) {
	conn := elb.New(session.New(), aws.NewConfig().WithRegion(region))
	totAdd := 0
	totUpd := 0
	var wwg sync.WaitGroup  // Writer wait group

	// get all the ELBS
	resp, err := conn.DescribeLoadBalancers(nil)
	if err != nil {
		log.Fatal(map[string]interface{}{}, err.Error())
	}

	var aDocs []interface{}
	var uDocs []interface{}
	for _, elb := range resp.LoadBalancerDescriptions {
		if storage.DocExists(*elb.LoadBalancerName, "loadbalancers") {
			// update if elbs already in DB
			uDocs = append(uDocs, PruneELBFields(elb))
		} else {
			// Insert new elbs into DB
			aDocs = append(aDocs, PruneELBFields(elb))
		}
	}

	noADocs := len(aDocs)
	noUDocs := len(uDocs)
	totAdd += noADocs
	totUpd += noUDocs

	if noADocs != 0 {
		wwg.Add(noADocs)
		go storage.InsertMany(aDocs, "loadbalancers", &wwg)
	}

	wwg.Wait()  // wait for all the writer to finish

	log.Debug(map[string]interface{}{
		"servciename": "aws_elb",
		"totalAdded": totAdd,
		"totalUpdated": totUpd,
		"region": region,
	}, "Crawler finished")
	cwg.Done()  // say AWS ELB crawler is done
}

//Crawl all the DNS records in a AWS cloud
func CrawlAllRoute53(cwg *sync.WaitGroup) {
	conn := route53.New(session.New())
	var zones []string
	totAdd := 0
	totUpd := 0
	var wwg sync.WaitGroup  // Writer wait group

	resp, err := conn.ListHostedZones(nil)
	if err != nil {
		log.Fatal(map[string]interface{}{}, err.Error())
	}

	// get all the zones
	for _, zone := range resp.HostedZones {
		zones = append(zones, *zone.Id)
	}

	for _, zone := range zones {
		params := &route53.ListResourceRecordSetsInput{
			HostedZoneId: aws.String(zone),
		}
		// get all the records for a zone
		resp, err := conn.ListResourceRecordSets(params)
		if err != nil {
			log.Fatal(map[string]interface{}{}, err.Error())
		}

		var aDocs []interface{}
		var uDocs []interface{}
		var docsToAdd []storage.DNSDoc
		var docsToUpd []storage.DNSDoc
		var dupDocs []storage.DNSDoc
		for _, rec := range resp.ResourceRecordSets {
			if *rec.Type == "CNAME" || *rec.Type == "A" {
				nDoc := PruneR53Fields(rec)
				dID := strings.TrimSuffix(*rec.Name, ".")
				if storage.DocExists(dID, "dns") {
					// update if dns record already in DB
					docsToUpd = append(docsToUpd, PruneR53Fields(rec))
				} else {
					if len(docsToAdd) == 0 {
						// Insert new dns record into DB
						docsToAdd = append(docsToAdd, nDoc)
					} else {
						// Check for duplicate record entries on the zone
						if CheckDNSRecordExists(docsToAdd, nDoc) {
							// Insert new dns record into DB
							dupDocs = append(dupDocs, nDoc)
						} else {
							// Insert new dns record into DB
							docsToAdd = append(docsToAdd, nDoc)
						}
					}

				}
			}
		}

		// Check if there are duplicate docs and merge them
		if len(dupDocs) != 0 {
			aDocs = MergeDnsRecords(docsToAdd, dupDocs)
		}

		noADocs := len(aDocs)
		noUDocs := len(uDocs)
		totAdd += noADocs
		totUpd += noUDocs

		if noADocs != 0 {
			wwg.Add(noADocs)
			go storage.InsertMany(aDocs, "dns", &wwg)
		}
	}

	wwg.Wait()  // wait for all the writer to finish

	log.Debug(map[string]interface{}{
		"servciename": "aws_route53",
		"totalAdded": totAdd,
		"totalUpdated": totUpd,
	}, "Crawled finished")
	cwg.Done()  // say AWS Route53 crawler is done
}
