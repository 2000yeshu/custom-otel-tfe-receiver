package awscloudwatchmetricsreceiver

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"

	"go.uber.org/zap"
)

type awsCloudWatchMetricsReceiver struct {
	config        *Config
	nextStartTime time.Time
	logger        *zap.Logger
	cancel        context.CancelFunc
	nextConsumer  consumer.Metrics
	wg            *sync.WaitGroup
	provider      *StatsProvider
	doneChan      chan bool
}

func newMetricReceiver(cfg *Config, logger *zap.Logger, consumer consumer.Metrics) *awsCloudWatchMetricsReceiver {
	return &awsCloudWatchMetricsReceiver{
		config:        cfg,
		nextStartTime: time.Now().Add(-cfg.PollInterval),
		logger:        logger,
		wg:            &sync.WaitGroup{},
		nextConsumer:  consumer,
		doneChan:      make(chan bool),
	}
}

func (m *awsCloudWatchMetricsReceiver) Start(ctx context.Context, _ component.Host) error {
	m.logger.Debug("starting to poll for CloudWatch metrics")
	ctx, m.cancel = context.WithCancel(ctx)

	go func() {
		ticker := time.NewTicker(m.config.PollInterval)
		for {
			select {
			case <-ticker.C:
				_ = m.collectDataFromCloudWatch(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (m *awsCloudWatchMetricsReceiver) collectDataFromCloudWatch(ctx context.Context) error {
	// fetch stats
	m.provider = NewStatsProvider("us-west-2")
	stats, metadata, err := m.provider.GetStats(ctx)
	if err != nil {
		return err
	}

	// convert to otel metrics
	mds := MetricsData(stats, metadata, m.logger)

	for _, md := range mds {
		err := m.nextConsumer.ConsumeMetrics(ctx, md)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *awsCloudWatchMetricsReceiver) Shutdown(_ context.Context) error {
	if m.cancel != nil {
		m.cancel()
	}
	return nil
}
