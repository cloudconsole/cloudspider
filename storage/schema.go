package storage

import (
	"time"
)

// Describes an Tags
type Tags struct {
	Key string
	Value string
}

// Describes an machine
type NodeDoc struct {
	ID             string     `bson:"_id"`
	State          string     `bson:"state"`
	Virtualization string     `bson:"virtualization"`
	Architecture   string     `bson:"arch"`
	RootDevice     string     `bson:"root_device"`
	Type           string     `bson:"type"`
	DataCenter     string     `bson:"data_center"`
	SecurityGroup  []string   `bson:"security_group"`
	CloudProvider  string     `bson:"cloud_provider"`
	SshKeyName     string     `bson:"ssh_key_name"`
	Tags           []Tags      `bson:"tags"`
	IamProfile     string     `bson:"iam_profile,omitempty"`
	LaunchTime     time.Time  `bson:"launch_time"`
	PublicDns      string     `bson:"public_dns"`
	PrivateDns     string     `bson:"private_dns"`
	PublicIp       string     `bson:"public_ip"`
	PrivateIp      string     `bson:"private_ip"`
}

// Describes an load balancer
type LBDoc struct {
	ID            string    `bson:"_id"`
	Name          string    `bson:"name"`
	LaunchTime    time.Time `bson:"launch_time"`
	DataCenter    []string  `bson:"data_center"`
	CloudProvider string    `bson:"cloud_provider"`
	Backends      []string  `bson:"backends"`
	PublicDns     string    `bson:"public_dns"`
}

// Describes an DNS record
type DNSDoc struct {
	ID            string   `bson:"_id"`
	Name          string   `bson:"name"`
	Type          string   `bson:"type"`
	Records       []string `bson:"records"`
	CloudProvider string   `bson:"cloud_provider"`
}
