package awscloudwatchmetricsreceiver

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type StatsProvider struct {
	region string
}

func NewStatsProvider(region string) *StatsProvider {
	return &StatsProvider{region: region}
}

// Return an umbrella struct of all cloudwatch metrics
func (sp *StatsProvider) GetStats(ctx context.Context) (*EC2Stats, EC2Metadata, error) {
	var ec2Stats EC2Stats
	var metadata EC2Metadata
	metadata.InstanceID = "i-053db4dfa7c7463ad" // Replace with your instance ID

	out, err := sp.fetchCloudwatchMetrics(ctx)
	if err != nil {
		return nil, EC2Metadata{}, err
	}

	// fmt.Println("length of metric data results:", len(out.MetricDataResults), out.MetricDataResults)

	for _, metric := range out.MetricDataResults {
		ec2Stats.CPUUtilized = metric.Values
		ec2Stats.Timestamps = metric.Timestamps
	}
	return &ec2Stats, metadata, nil

}

// TODO: Prevent reinitialization of AWS SDK config
func (sp *StatsProvider) fetchCloudwatchMetrics(ctx context.Context) (*cloudwatch.GetMetricDataOutput, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(sp.region))
	if err != nil {
		return nil, err
	}

	// Create CloudWatch client
	cw := cloudwatch.NewFromConfig(cfg)

	// // Example: List metrics (you can customize this as needed)
	// input := &cloudwatch.ListMetricsInput{
	// 	Namespace:  aws.String("AWS/EC2"),        // Replace with your desired namespace
	// 	MetricName: aws.String("CPUUtilization"), // Replace with your desired metric name
	// 	Dimensions: []*cloudwatch.DimensionFilter{
	// 		{
	// 			Name:  aws.String("InstanceId"), // Replace with your desired dimension name
	// 			Value: aws.String("*"),          // Replace with your desired dimension value
	// 		},
	// 	},
	// }
	// result, err := cw.ListMetricsWithContext(ctx, input)
	// if err != nil {
	// 	receiver.logger.Error("failed to list CloudWatch metrics", zap.Error(err))
	// 	return err
	// }

	// for _, metric := range result.Metrics {
	// 	receiver.logger.Info("Found metric",
	// 		zap.String("namespace", aws.StringValue(metric.Namespace)),
	// 		zap.String("metricName", aws.StringValue(metric.MetricName)),
	// 	)
	// }

	// You can now use cw.GetMetricData or cw.GetMetricStatistics to fetch actual metric values as needed.
	// Implementation
	// For example, you can fetch metric data like this:
	inputData := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []cloudwatchtypes.MetricDataQuery{
			cloudwatchtypes.MetricDataQuery{
				Id: aws.String("m1"),
				MetricStat: &cloudwatchtypes.MetricStat{
					Metric: &cloudwatchtypes.Metric{
						Namespace:  aws.String("AWS/EC2"),
						MetricName: aws.String("CPUUtilization"),
						// Should be set in OTEL metric attributes
						Dimensions: []cloudwatchtypes.Dimension{
							{
								Name:  aws.String("InstanceId"),
								Value: aws.String("i-053db4dfa7c7463ad"), // Replace with your instance ID
							},
						},
					},
					Period: aws.Int32(60),         // 60 seconds
					Stat:   aws.String("Average"), // Replace with your desired statistic
				},
			},
		},
		StartTime: aws.Time(time.Now().Add(-time.Hour)), // 1 hour ago
		EndTime:   aws.Time(time.Now()),                 // Now
	}

	outputData, err := cw.GetMetricData(ctx, inputData)
	if err != nil {
		return nil, err
	}

	// print out metadata for debugging
	// fmt.Println("Metadata:", outputData.ResultMetadata)

	return outputData, nil
}
