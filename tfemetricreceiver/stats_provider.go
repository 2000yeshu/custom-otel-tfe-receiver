package awscloudwatchmetricsreceiver

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	aurorards "github.com/aws/aws-sdk-go-v2/service/rds"
	aurorardstypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type StatsProvider struct {
	region            string
	cfg               *Config
	cw                *cloudwatch.Client
	rds               *aurorards.Client
	cloudwatchMetrics *TFECloudwatchMetrics
}

func NewStatsProvider(region string, config *Config, cw *cloudwatch.Client, rds *aurorards.Client, cwm *TFECloudwatchMetrics) *StatsProvider {
	return &StatsProvider{region: region, cfg: config, cw: cw, rds: rds, cloudwatchMetrics: cwm}
}

// Takes the arn and type of resource
// Also takes the metrics requested by user
// Uses the repective single dimension for the resource type and
// Fetches all the possible values for that dimension using the arn GET query
// returns the MetricDataQuery for eachmetric+eachdimension for that resource arn
// Also validates that the user has not provided an invalid metric name
func (sp *StatsProvider) getMetricsFromARN(ctx context.Context, cloudwatchMetrics *TFECloudwatchMetrics, nameConfig *NamedConfig) error {
	switch nameConfig.Type {
	case rdsResource:
		// ARN format: arn:aws:rds:region:account-id:cluster:db-cluster-identifier
		parts := strings.Split(nameConfig.ARN, ":")
		if len(parts) != 7 {
			return fmt.Errorf("invalid ARN: %s", nameConfig.ARN)
		}

		clusterId := parts[6]
		fmt.Println("Fetching metrics for RDS cluster:", clusterId)
		// var page *aurorards.DescribeDBInstancesOutput
		paginator := aurorards.NewDescribeDBInstancesPaginator(sp.rds, &aurorards.DescribeDBInstancesInput{
			Filters: []aurorardstypes.Filter{
				{
					Name:   aws.String("db-cluster-id"),
					Values: []string{clusterId},
				},
			},
		})
		if paginator == nil {
			return fmt.Errorf("failed to create paginator for RDS instances")
		}

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return err
			}

			for _, instance := range page.DBInstances {
				if instance.DBClusterIdentifier != nil && instance.DBInstanceIdentifier != nil && *instance.DBClusterIdentifier == clusterId {
					sp.cloudwatchMetrics.RDSStats[*instance.DBInstanceIdentifier] = RDSStats{
						MetricsData: make(map[string][]Float64DataPoint, 0),
						// Insert metadata, metrics will be populated later
						RDSMetadata: RDSMetadata{
							DBInstanceIdentifier: *instance.DBInstanceIdentifier,
							Region:               sp.region,
							AvailabilityZone:     *instance.AvailabilityZone,
						},
					}
					// fmt.Println("Found RDS instance:", *instance.DBInstanceIdentifier, "in cluster:", clusterId)

				}
			}

		}

		for instanceId, _ := range sp.cloudwatchMetrics.RDSStats {
			for _, metric := range nameConfig.MetricNames {
				validMetric := false
				for _, validMetricName := range RDSInstanceMetrics {
					if metric == validMetricName {
						validMetric = true
						break
					}
				}
				if !validMetric {
					return fmt.Errorf("invalid metric name: %s for RDS instance: %s", metric, instanceId)
				}

				// TODO: Validate metric name against a predefined list of valid RDS metrics
				// fmt.Println("Fetching metric data for:", metric, "for RDS instance:", *instance.DBInstanceIdentifier)
				thisMetricQueryForThisDbInstance := []cloudwatchtypes.MetricDataQuery{{
					Id: aws.String(rdsResource),
					MetricStat: &cloudwatchtypes.MetricStat{
						Metric: &cloudwatchtypes.Metric{
							Namespace:  aws.String(rdsNamespace),
							MetricName: aws.String(metric),
							Dimensions: []cloudwatchtypes.Dimension{
								{
									Name:  aws.String("DBInstanceIdentifier"),
									Value: &instanceId,
								},
							},
						},
						Period: aws.Int32(int32(defaultPeriod.Seconds())),
						Stat:   aws.String("Average"), // Example aggregation, can be parameterized
					},
				}}
				outputData, err := sp.cw.GetMetricData(ctx, &cloudwatch.GetMetricDataInput{
					MetricDataQueries: thisMetricQueryForThisDbInstance,
					StartTime:         aws.Time(time.Now().Add(-time.Hour)), // 1 hour ago
					EndTime:           aws.Time(time.Now()),                 // Now
				})
				if err != nil {
					return err
				}

				// fpdp := make([]Float64DataPoint, 0)
				// cloudwatchMetrics.RDSStats[*instance.DBInstanceIdentifier].MetricsData = make(map[string]fpdp)
				// fmt.Println("Is map nil", cloudwatchMetrics.RDSStats[*instance.DBInstanceIdentifier].MetricsData == nil)
				cloudwatchMetrics.RDSStats[instanceId].MetricsData[metric] = make([]Float64DataPoint, 0)
				for idx, value := range outputData.MetricDataResults[0].Values {
					cloudwatchMetrics.RDSStats[instanceId].MetricsData[metric] = append(
						cloudwatchMetrics.RDSStats[instanceId].MetricsData[metric],
						Float64DataPoint{
							Value:     value,
							Timestamp: outputData.MetricDataResults[0].Timestamps[idx],
						},
					)
				}

			}
			fmt.Println("Fetched metric data for RDS instance:", instanceId, "with length of CPUUtilization values:", len(cloudwatchMetrics.RDSStats[instanceId].MetricsData[RDSInstanceCPUUtilization]), "with length of NetworkThroughput values:",
				len(sp.cloudwatchMetrics.RDSStats[instanceId].MetricsData[RDSInstanceNetworkThroughput]),
				"with length of WriteThroughput values:", len(sp.cloudwatchMetrics.RDSStats[instanceId].MetricsData[RDSInstanceWriteThroughput]),
				"with length of ReadThroughput values:", len(sp.cloudwatchMetrics.RDSStats[instanceId].MetricsData[RDSInstanceReadThroughput]))
		}
		return nil
	case elastiCacheResource:
	case ebsResource:

	}
	return fmt.Errorf("unsupported resource type: %s", nameConfig.Type)
}

// TODO: Prevent reinitialization of AWS SDK config
// TODO: Check AWS api limitations and handle them
func (sp *StatsProvider) fetchCloudwatchMetrics(ctx context.Context) error {
	sp.cloudwatchMetrics.RDSStats = make(map[string]RDSStats)

	for _, names := range sp.cfg.Metrics.Names {
		if names == nil {
			continue
		}

		if names.Period.Seconds() == 0 {
			names.Period = defaultPeriod
		}

		err := sp.getMetricsFromARN(ctx, sp.cloudwatchMetrics, names)
		if err != nil {
			return err
		}
	}

	return nil
}

// Return an umbrella struct of all cloudwatch metrics
func (sp *StatsProvider) GetStats(ctx context.Context) error {
	// var ec2Stats []EC2Stats
	if err := sp.fetchCloudwatchMetrics(ctx); err != nil {
		return err
	}

	return nil
}
