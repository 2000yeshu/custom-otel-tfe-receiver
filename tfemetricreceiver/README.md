To add a new metric for a specific type
1. Add the types to constant.go
2. Add the block to accumulator

## Requirements
Our requirement for creating a config file is that:
1. We can't specify the specific instance ids (for e.g. for Amazon RDS, ElastiCache, etc.)
2. We want to produce raw metrics from this receiver. And not aggregated(or cluster) level metrics. 

## Possible Strategies
With these requirements, we had two ways to implement this. 

### Method 1 
A dumb implementation wherein a user specifies the NS, metrics and the regex for dimension values that they have to monitor and we will just fetch those metrics from cloudwatch + fetch all the dimensions for those metrics matching the regex and pass to the next consumer.
First we will list all the dimensions related to a metric. 
E.g. 
```
Metric: CPUUtilization idx: 0 Dimension: DatabaseClass Value: db.r6i.xlarge
Metric: CPUUtilization idx: 1 Dimension: DBInstanceIdentifier Value: sandbox-tfe-rds-cluster-instance-2
Metric: CPUUtilization idx: 2 Dimension: DBInstanceIdentifier Value: sandbox-tfe-rds-cluster-instance-1
Metric: CPUUtilization idx: 4 Dimension: DBInstanceIdentifier Value: sandbox-tfe-rds-cluster-instance-0
Metric: CPUUtilization idx: 5 Dimension: EngineName Value: aurora-postgresql
Metric: CPUUtilization idx: 6 Dimension: Role Value: READER
Metric: CPUUtilization idx: 6 Dimension: DBClusterIdentifier Value: sandbox-tfe-rds-cluster-us-west-2
Metric: CPUUtilization idx: 7 Dimension: Role Value: WRITER
Metric: CPUUtilization idx: 7 Dimension: DBClusterIdentifier Value: sandbox-tfe-rds-cluster-us-west-2
Metric: CPUUtilization idx: 8 Dimension: DBClusterIdentifier Value: sandbox-tfe-rds-cluster-us-west-2
```
Then, filter the dimensions which matches with the regex user has given. There could be multiple dimension key that matches them.
E.g. DBInstanceIdentifier Value: sandbox-tfe-rds-cluster-instance-*
```
Metric: CPUUtilization idx: 1 Dimension: DBInstanceIdentifier Value: sandbox-tfe-rds-cluster-instance-2
Metric: CPUUtilization idx: 2 Dimension: DBInstanceIdentifier Value: sandbox-tfe-rds-cluster-instance-1
Metric: CPUUtilization idx: 4 Dimension: DBInstanceIdentifier Value: sandbox-tfe-rds-cluster-instance-0
```
Then, append this metric for each dimensions value to the GetMetricsDataRequest. and fetch the values and timestamps. Each will be a different resource taggesd respectively.

If user has specified faulty regex and we 

Now it has complicacies:
1. An api call for each metric for listing dimensions for each metric.
2. Regex is not deterministic such that outputs can be changed if more resources matching the regex is added to the AWS account.
3. To even match the dimension values to a regex we have to select a dimension key. For e.g. if we want to get the CPUUtilization metrics of all the instances in a RDS cluster. Which dimension do we provide the 
```
- namespace: AWS/RDS
  metric_name: ["CPUUtilization", "WriteThroughput"]
  dimensions:
   - name: DBInstanceidentifier
     value: "sandbox-tfe-rds-cluster-instance-*"
```
### Method 2