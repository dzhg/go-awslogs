package internal

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"time"
)

type AWSError struct {
	context string
	cause   error
}

func (e *AWSError) Error() string {
	return fmt.Sprintf("%s: %s", e.context, e.cause.Error())
}

func awsError(context string, err error) error {
	return &AWSError{context, err}
}

type AWSLogsClient struct {
	client *cloudwatchlogs.Client
}

func NewAwsLogsClient() (*AWSLogsClient, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, awsError("load default AWS config", err)
	}

	cfg.Region = "us-west-2"

	return &AWSLogsClient{cloudwatchlogs.New(cfg)}, nil
}

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

type LogStreamSorting struct {
	Descending *bool
	OrderBy    cloudwatchlogs.OrderBy
}

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

type GetLogsParams struct {
	Limit          *int64
	StartTime      *int64
	EndTime        *int64
	LogGroupName   *string
	LogStreamNames []string
}

type LogEvent struct {
	Message    string
	EventId    string
	GroupName  string
	StreamName string
	Error      error
}

func (c *AWSLogsClient) StreamLogEvents(params *GetLogsParams, ch chan *LogEvent, watch *bool) {
	var input *cloudwatchlogs.FilterLogEventsInput

	input = &cloudwatchlogs.FilterLogEventsInput{
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

	lru, err := NewLRU(int(*params.Limit))
	if err != nil {
		ch <- &LogEvent{Error: err}
		return
	}

	for {
		req := c.client.FilterLogEventsRequest(input)

		response, err := req.Send(context.TODO())

		if err != nil {
			ch <- &LogEvent{Error: err}
			return
		}

		for _, e := range response.FilterLogEventsOutput.Events {
			if _, ok := lru.Get(*e.EventId); ok {
				continue
			} else {
				lru.Add(*e.EventId, *e.Message)
				ch <- &LogEvent{
					Message:    *e.Message,
					EventId:    *e.EventId,
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
			continue
		}

		input = &cloudwatchlogs.FilterLogEventsInput{
			LogGroupName: params.LogGroupName,
			NextToken:    response.NextToken,
		}
	}
}
