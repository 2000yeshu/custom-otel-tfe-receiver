package awscloudwatchmetricsreceiver

import (
	"fmt"

	"go.opentelemetry.io/collector/pdata/pcommon"
)

func getEbsResource(tm EBSMetadata) pcommon.Resource {
	resource := pcommon.NewResource()

	// Set volume id
	resource.Attributes().PutStr(attributeEBSVolumeId, tm.VolumeID)

	return resource
}

func getRDSResource(tm RDSMetadata) pcommon.Resource {
	resource := pcommon.NewResource()

	// Set cluster id
	resource.Attributes().PutStr(attributeDBInstanceIdentifier, tm.DBInstanceIdentifier)
	resource.Attributes().PutStr(attributeAvailibilityZone, tm.AvailabilityZone)

	return resource
}

func getElastiCacheResource(tm ElastiCacheMetadata) pcommon.Resource {
	resource := pcommon.NewResource()

	// Set cluster id
	resource.Attributes().PutStr(attributeElastiCacheClusterId, tm.ClusterID)

	return resource
}

func getEc2Resource(tm EC2Metadata) pcommon.Resource {

	resource := pcommon.NewResource()
	fmt.Println("getEc2Resource called with EC2Metadata:", tm.InstanceID)
	// Set instance id

	resource.Attributes().PutStr(attributeEc2InstanceId, tm.InstanceID)

	// region, accountID, taskID := getResourceFromARN(tm.TaskARN)
	// resource.Attributes().PutStr(attributeECSCluster, getNameFromCluster(tm.Cluster))
	// resource.Attributes().PutStr(string(conventions.AWSECSTaskARNKey), tm.TaskARN)
	// resource.Attributes().PutStr(attributeECSTaskID, taskID)
	// resource.Attributes().PutStr(string(conventions.AWSECSTaskFamilyKey), tm.Family)

	// // Task revision: aws.ecs.task.version and aws.ecs.task.revision
	// resource.Attributes().PutStr(attributeECSTaskRevision, tm.Revision)
	// resource.Attributes().PutStr(string(conventions.AWSECSTaskRevisionKey), tm.Revision)

	// resource.Attributes().PutStr(attributeECSServiceName, tm.ServiceName)

	// resource.Attributes().PutStr(string(conventions.CloudAvailabilityZoneKey), tm.AvailabilityZone)
	// resource.Attributes().PutStr(attributeECSTaskPullStartedAt, tm.PullStartedAt)
	// resource.Attributes().PutStr(attributeECSTaskPullStoppedAt, tm.PullStoppedAt)
	// resource.Attributes().PutStr(attributeECSTaskKnownStatus, tm.KnownStatus)

	// // Task launchtype: aws.ecs.task.launch_type (raw string) and aws.ecs.launchtype (lowercase)
	// resource.Attributes().PutStr(attributeECSTaskLaunchType, tm.LaunchType)
	// switch lt := strings.ToLower(tm.LaunchType); lt {
	// case "ec2":
	// 	resource.Attributes().PutStr(string(conventions.AWSECSLaunchtypeKey), conventions.AWSECSLaunchtypeEC2.Value.AsString())
	// case "fargate":
	// 	resource.Attributes().PutStr(string(conventions.AWSECSLaunchtypeKey), conventions.AWSECSLaunchtypeFargate.Value.AsString())
	// }

	// resource.Attributes().PutStr(string(conventions.CloudRegionKey), region)
	// resource.Attributes().PutStr(string(conventions.CloudAccountIDKey), accountID)

	return resource
}
