package awscloudwatchmetricsreceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"

	"go.opentelemetry.io/collector/receiver"
)

var (
	typeStr = component.MustNewType("tfeawscloudwatchmetricsreceiver")
)

func NewFactory() receiver.Factory {
	return receiver.NewFactory(typeStr, createDefaultConfig, receiver.WithMetrics(createMetricsReceiver, component.StabilityLevelUnmaintained))
}

func createMetricsReceiver(_ context.Context, params receiver.Settings, baseCfg component.Config, consumer consumer.Metrics) (receiver.Metrics, error) {
	cfg := baseCfg.(*Config)
	rcvr := newMetricReceiver(cfg, params.Logger, consumer)
	return rcvr, nil
}

func createDefaultConfig() component.Config {
	return &Config{
		PollInterval: defaultPollInterval,
		Metrics:      &MetricsConfig{},
	}
}
