package awscloudwatchmetricsreceiver

import (
	"fmt"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"

	"go.uber.org/zap"
)

type metricDataAccumulator struct {
	mds []pmetric.Metrics
}

func (acc *metricDataAccumulator) getMetricsDataUsingDimensions(cfmetrics *TFECloudwatchMetrics, _ *zap.Logger) {

	for _, dimBasedMetric := range cfmetrics.DimensionBasedMetric {
		fmt.Println("Processing dimension-based metric:", dimBasedMetric.Key, "with value:", dimBasedMetric.Value)
		resource := getResourceFromDimensions(dimBasedMetric.Key, dimBasedMetric.Value)
		for metricKey, metricData := range dimBasedMetric.MetricValues {

			for idx, metricValue := range metricData.Values {
				timestamp := pcommon.NewTimestampFromTime(metricData.Timestamps[idx])

				acc.accumulate(convertToOTLPMetrics(dimBasedMetric.Prefix+metricKey, metricValue, resource, timestamp))
				fmt.Println("Accumulated metric:", dimBasedMetric.Prefix+metricKey, "with value:", metricValue, "at timestamp:", metricData.Timestamps[idx])
			}

		}

	}

}

func (acc *metricDataAccumulator) getMetricsData(cfmetrics *TFECloudwatchMetrics, _ *zap.Logger) {
	// cloudwatchMetric := CloudWatchMetrics{}

	// ec2Resource := getEc2Resource(cfmetrics.EC2Metadata)

	for _, rdsinstance := range cfmetrics.RDSStats {
		// Tag resource with this particular RDS instance metadata
		rdsResource := getRDSResource(rdsinstance.RDSMetadata)
		fmt.Println("length of cpuutil data points:", len(rdsinstance.MetricsData[RDSInstanceCPUUtilization]))
		for _, metricData := range rdsinstance.MetricsData[RDSInstanceCPUUtilization] {
			timestamp := pcommon.NewTimestampFromTime(metricData.Timestamp)
			acc.accumulate(convertToOTLPMetrics(rdsPrefix+attributeCPUUtilized, metricData.Value, rdsResource, timestamp))
		}

		fmt.Println("length of network throughput data points:", len(rdsinstance.MetricsData[RDSInstanceNetworkThroughput]))
		for _, metricData := range rdsinstance.MetricsData[RDSInstanceNetworkThroughput] {
			timestamp := pcommon.NewTimestampFromTime(metricData.Timestamp)
			acc.accumulate(convertToOTLPMetrics(rdsPrefix+attributeNetworkThroughput, metricData.Value, rdsResource, timestamp))
		}

		fmt.Println("length of write throughput data points:", len(rdsinstance.MetricsData[RDSInstanceWriteThroughput]))
		for _, metricData := range rdsinstance.MetricsData[RDSInstanceWriteThroughput] {
			timestamp := pcommon.NewTimestampFromTime(metricData.Timestamp)
			acc.accumulate(convertToOTLPMetrics(rdsPrefix+attributeWriteThroughput, metricData.Value, rdsResource, timestamp))
		}

		fmt.Println("length of read throughput data points:", len(rdsinstance.MetricsData[RDSInstanceReadThroughput]))
		for _, metricData := range rdsinstance.MetricsData[RDSInstanceReadThroughput] {
			timestamp := pcommon.NewTimestampFromTime(metricData.Timestamp)
			acc.accumulate(convertToOTLPMetrics(rdsPrefix+attributeReadThroughput, metricData.Value, rdsResource, timestamp))
		}

	}
	// elastiCacheResource := getElastiCacheResource(cfmetrics.ElastiCacheMetadata)
	// ebsResource := getEbsResource(cfmetrics.EBSMetadata)

	// lengthOfEC2DataPoints := len(cfmetrics.EC2Stats)
	// for idx := range lengthOfEC2DataPoints {
	// 	timestamp := pcommon.NewTimestampFromTime(cfmetrics.EC2Stats[idx].Timestamps)
	// 	cloudwatchMetric.GaugeMetricValue = cfmetrics.EC2Stats[idx].CPUUtilized
	// 	acc.accumulate(convertToOTLPMetrics(ec2Prefix, cloudwatchMetric, ec2Resource, timestamp))
	// }

	// lengthOfElastiCacheDataPoints := len(cfmetrics.ElastiCacheStats)
	// for idx := range lengthOfElastiCacheDataPoints {
	// 	timestamp := pcommon.NewTimestampFromTime(cfmetrics.ElastiCacheStats[idx].Timestamps)
	// 	cloudwatchMetric.GaugeMetricValue = cfmetrics.ElastiCacheStats[idx].CPUUtilized
	// 	acc.accumulate(convertToOTLPMetrics(elastiCachePrefix, cloudwatchMetric, elastiCacheResource, timestamp))
	// }

	// lengthOfEBSDataPoints := len(cfmetrics.EBSStats)
	// for idx := range lengthOfEBSDataPoints {
	// 	timestamp := pcommon.NewTimestampFromTime(cfmetrics.EBSStats[idx].Timestamps)
	// 	cloudwatchMetric.GaugeMetricValue = cfmetrics.EBSStats[idx].VolumeWriteBytes
	// 	acc.accumulate(convertToOTLPMetrics(ebsPrefix, cloudwatchMetric, ebsResource, timestamp))
	// }
}

func (acc *metricDataAccumulator) accumulate(md pmetric.Metrics) {
	if acc.mds == nil {
		acc.mds = make([]pmetric.Metrics, 0)
	}
	acc.mds = append(acc.mds, md)
}
