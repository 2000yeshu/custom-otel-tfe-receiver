package awscloudwatchmetricsreceiver

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"

	"go.uber.org/zap"
)

type metricDataAccumulator struct {
	mds []pmetric.Metrics
}

func (acc *metricDataAccumulator) getMetricsData(ec2Stats *EC2Stats, metadata EC2Metadata, _ *zap.Logger) {
	cloudwatchMetric := CloudWatchMetrics{}
	ec2Resource := getEc2Resource(metadata)

	lengthOfDataPoints := len(ec2Stats.CPUUtilized)

	for idx := range lengthOfDataPoints {
		timestamp := pcommon.NewTimestampFromTime(ec2Stats.Timestamps[idx])
		cloudwatchMetric.CPUUtilized = ec2Stats.CPUUtilized[idx]
		acc.accumulate(convertToOTLPMetrics(ec2Prefix, cloudwatchMetric, ec2Resource, timestamp))
	}

}

func (acc *metricDataAccumulator) accumulate(md pmetric.Metrics) {
	if acc.mds == nil {
		acc.mds = make([]pmetric.Metrics, 0)
	}
	acc.mds = append(acc.mds, md)
}
