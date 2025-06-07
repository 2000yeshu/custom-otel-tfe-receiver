package awscloudwatchmetricsreceiver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/rds"
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
	ctx, m.cancel = context.WithCancel(ctx)

	go func() {
		ticker := time.NewTicker(m.config.PollInterval)
		for {
			select {
			case <-ticker.C:
				if err := m.collectDataFromCloudWatch(ctx); err != nil {
					fmt.Println(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (m *awsCloudWatchMetricsReceiver) collectDataFromCloudWatch(ctx context.Context) error {
	// fetch stats
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		return err
	}

	// Create CloudWatch client
	cw := cloudwatch.NewFromConfig(cfg)
	rds := rds.NewFromConfig(cfg)
	m.provider = NewStatsProvider("us-west-2", m.config, cw, rds, &TFECloudwatchMetrics{})
	if err := m.provider.GetStats(ctx); err != nil {
		return err
	}

	// convert to otel metrics
	mds := MetricsData(m.provider.cloudwatchMetrics, m.logger)

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
