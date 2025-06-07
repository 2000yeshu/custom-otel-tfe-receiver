package awscloudwatchmetricsreceiver

const (
	ec2Prefix         = "vm1."
	rdsPrefix         = "db."
	elastiCachePrefix = "cache."
	ebsPrefix         = "storage."

	// TODO: Attribute not showing up for some reason
	attributeEc2InstanceId        = "aws.ec2.instance.id"
	attributeEBSVolumeId          = "aws.ebs.volume.id"
	attributeElastiCacheClusterId = "aws.elasticache.cluster.id"
	attributeDBInstanceIdentifier = "aws.rds.db_instance.identifier"
	attributeAvailibilityZone     = "aws.availability_zone"

	attributeCPUUtilized       = "cpu.utilized"
	attributeNetworkThroughput = "network.throughput"
	attributeWriteThroughput   = "write.throughput"
	attributeReadThroughput    = "read.throughput"

	unitNone = "None"
)

const (
	RDSInstanceCPUUtilization    = "CPUUtilization"
	RDSInstanceNetworkThroughput = "NetworkThroughput"
	RDSInstanceWriteThroughput   = "WriteThroughput"
	RDSInstanceReadThroughput    = "ReadThroughput"
)

var (
	RDSInstanceMetrics         = []string{RDSInstanceCPUUtilization, RDSInstanceNetworkThroughput, RDSInstanceWriteThroughput, RDSInstanceReadThroughput}
	ElastiCacheInstanceMetrics = []string{"CPUUtilized", "NetworkThroughput"}
)
