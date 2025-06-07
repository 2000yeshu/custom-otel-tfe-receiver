package awscloudwatchmetricsreceiver

import (
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

func MetricsData(cfmetrics *TFECloudwatchMetrics, logger *zap.Logger) []pmetric.Metrics {
	acc := &metricDataAccumulator{}
	acc.getMetricsData(cfmetrics, logger)

	return acc.mds
}
