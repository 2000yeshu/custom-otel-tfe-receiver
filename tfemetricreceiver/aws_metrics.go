package awscloudwatchmetricsreceiver

import (
	"time"

	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type EBSStats struct {
	VolumeWriteBytes float64   `json:"volume_write_bytes"`
	Timestamps       time.Time `json:"timestamps"`
}

type EBSMetadata struct {
	VolumeID string `json:"volume_id"`
}

type ElastiCacheStats struct {
	CPUUtilized float64   `json:"cpu_utilized"`
	Timestamps  time.Time `json:"timestamps"`
}

type ElastiCacheMetadata struct {
	ClusterID string `json:"cluster_id"`
}

// Multiple type of metrics for a single RDS instance
// and also we have some metadata
type RDSStats struct {
	// CPUUtilized       []Float64DataPoint `json:"cpu_utilized"`
	// NetworkThroughput []Float64DataPoint `json:"network_throughput"`

	// Key name is the metric name
	//  Allowed CPUUtilized, NetworkThroughput
	MetricsData map[string][]Float64DataPoint `json:"disk_io"`
	RDSMetadata RDSMetadata                   `json:"rds_metadata"`
}

type Float64DataPoint struct {
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type RDSMetadata struct {
	DBInstanceIdentifier string `json:"db_instance_identifier"`
	Region               string `json:"region"`
	AvailabilityZone     string `json:"availability_zone"`
}

type EC2Stats struct {
	CPUUtilized float64   `json:"cpu_utilized"`
	Timestamps  time.Time `json:"timestamps"`
}

type EC2Metadata struct {
	InstanceID string `json:"instance_id"`
}

const (
	ec2Resource         = "aws_ec2_instance"
	rdsResource         = "aws_rds_instance"
	elastiCacheResource = "aws_elasticache_cluster"
	ebsResource         = "aws_ebs_volume"

	ec2Namespace         = "AWS/EC2"
	elastiCacheNamespace = "AWS/ElastiCache"
	rdsNamespace         = "AWS/RDS"
	ebsNamespace         = "AWS/EBS"
)

type DimensionBasedMetric struct {
	Prefix string `json:"prefix"` // e.g., "vm1", "db", "cache", "storage"
	Key    string `json:"key"`    // Dimension key, e.g., "InstanceId"
	Value  string `json:"value"`  // Dimension value, e.g., "i-1234567890abcdef0"
	// key is metric name, value is MetricDataResult
	MetricValues map[string]cloudwatchtypes.MetricDataResult
}

type TFECloudwatchMetrics struct {
	EC2Stats            []EC2Stats          `json:"ec2_stats"`
	EC2Metadata         EC2Metadata         `json:"ec2_metadata"`
	ElastiCacheStats    []ElastiCacheStats  `json:"elasticache_stats"`
	ElastiCacheMetadata ElastiCacheMetadata `json:"elasticache_metadata"`

	// this we can have multiple instances
	// for a single RDS instance, we can have multiple metrics
	// map with instance identifier as key
	RDSStats map[string]RDSStats `json:"rds_stats"`

	EBSStats    []EBSStats  `json:"ebs_stats"`
	EBSMetadata EBSMetadata `json:"ebs_metadata"`

	// MetricName, Dimension Key, Dimension Value,
	DimensionBasedMetric []*DimensionBasedMetric `json:"dimension_based_metric"`
}
