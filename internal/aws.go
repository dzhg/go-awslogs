package internal

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"time"
)

// AWSError is a wrapper for any AWS related errors
type AWSError struct {
	context string
	cause   error
}

var _ error = (*AWSError)(nil)

func (e *AWSError) Error() string {
	return fmt.Sprintf("%s: %s", e.context, e.cause.Error())
}

func awsError(context string, err error) error {
	return &AWSError{context, err}
}

// AWSLogsClient wraps a cloudwatchlogs.Client
type AWSLogsClient struct {
	client *cloudwatchlogs.Client
}

// NewAwsLogsClient creates a new AWSLogsClient
func NewAwsLogsClient() (*AWSLogsClient, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, awsError("load default AWS config", err)
	}

	return &AWSLogsClient{cloudwatchlogs.New(cfg)}, nil
}

// DescribeLogGroups is a wrapper of AWS API describe-log-groups
func (c *AWSLogsClient) DescribeLogGroups(limit *int64, prefix, nextToken *string) (*cloudwatchlogs.DescribeLogGroupsResponse, error) {
	req := c.client.DescribeLogGroupsRequest(&cloudwatchlogs.DescribeLogGroupsInput{
		Limit:              limit,
		LogGroupNamePrefix: prefix,
		NextToken:          nextToken,
	})

	res, err := req.Send(context.TODO())

	if err != nil {
		return nil, awsError("describe log groups", err)
	}

	return res, nil
}

// LogStreamSorting represents the arguments for sorting of log streams
type LogStreamSorting struct {
	Descending *bool
	OrderBy    cloudwatchlogs.OrderBy
}

// DescribeLogStreams is a wrapper of AWS API describe-log-streams
func (c *AWSLogsClient) DescribeLogStreams(limit *int64, logGroupName *string, prefix *string, sorting *LogStreamSorting, nextToken *string) (*cloudwatchlogs.DescribeLogStreamsResponse, error) {
	req := c.client.DescribeLogStreamsRequest(&cloudwatchlogs.DescribeLogStreamsInput{
		Descending:          sorting.Descending,
		Limit:               limit,
		LogGroupName:        logGroupName,
		LogStreamNamePrefix: prefix,
		NextToken:           nextToken,
		OrderBy:             sorting.OrderBy,
	})

	res, err := req.Send(context.TODO())

	if err != nil {
		return nil, awsError("describe log streams", err)
	}

	return res, nil
}

// GetLogsParams represents the parameters of function StreamLogEvents
type GetLogsParams struct {
	Limit          *int64
	StartTime      *int64
	EndTime        *int64
	LogGroupName   *string
	LogStreamNames []string
}

// LogEvent represents an AWS log event
type LogEvent struct {
	Message    string
	EventID    string
	GroupName  string
	StreamName string
	Error      error
}

// StreamLogEvents streams the log events to the given channel
func (c *AWSLogsClient) StreamLogEvents(params *GetLogsParams, ch chan *LogEvent, watch *bool) {
	var originalInput *cloudwatchlogs.FilterLogEventsInput

	originalInput = &cloudwatchlogs.FilterLogEventsInput{
		EndTime:             params.EndTime,
		FilterPattern:       nil,
		Interleaved:         nil,
		Limit:               params.Limit,
		LogGroupName:        params.LogGroupName,
		LogStreamNamePrefix: nil,
		LogStreamNames:      params.LogStreamNames,
		NextToken:           nil,
		StartTime:           params.StartTime,
	}

	nextTokenInput := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName: params.LogGroupName,
		Limit:        params.Limit,
		NextToken:    nil,
	}

	lru, err := NewLRU(int(*params.Limit))
	if err != nil {
		ch <- &LogEvent{Error: err}
		return
	}

	input := originalInput

	for {
		req := c.client.FilterLogEventsRequest(input)

		response, err := req.Send(context.TODO())

		if err != nil {
			ch <- &LogEvent{Error: err}
			return
		}

		var cnt = 0
		for _, e := range response.FilterLogEventsOutput.Events {
			if _, ok := lru.Get(*e.EventId); ok {
				cnt++
				continue
			} else {
				lru.Add(*e.EventId, true)
				ch <- &LogEvent{
					Message:    *e.Message,
					EventID:    *e.EventId,
					GroupName:  *params.LogGroupName,
					StreamName: *e.LogStreamName,
					Error:      nil,
				}
			}
		}

		if response.NextToken == nil {
			if !*watch {
				return
			}
			// watch
			time.Sleep(time.Second)
			// fetch again from 5 minutes ago
			startTime := (time.Now().Add(-5 * time.Minute)).UnixNano() / int64(time.Millisecond)

			originalInput.StartTime = &startTime
			input = originalInput

			continue
		}

		nextTokenInput.NextToken = response.NextToken
		input = nextTokenInput
	}
}
