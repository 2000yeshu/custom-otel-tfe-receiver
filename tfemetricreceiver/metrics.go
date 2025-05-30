package awscloudwatchmetricsreceiver

import (
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

func MetricsData(ec2Stats *EC2Stats, metadata EC2Metadata, logger *zap.Logger) []pmetric.Metrics {
	acc := &metricDataAccumulator{}
	acc.getMetricsData(ec2Stats, metadata, logger)

	return acc.mds
}
