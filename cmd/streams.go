package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/dzhg/go-awslogs/internal"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var streamsCmd = &cobra.Command{
	Use:   "streams <group>",
	Short: "get log streams",
	Long:  `get log streams of a given group, sorted by last event time (desc)`,
	Args:  cobra.ExactArgs(1),
	RunE:  streams,
}

func initStreamCmd() *cobra.Command {
	return streamsCmd
}

func streams(cmd *cobra.Command, args []string) error {
	c, err := internal.NewAwsLogsClient()
	if err != nil {
		return errors.Wrap(err, "new AWS logs client")
	}

	descending := true
	limit := int64(50)
	logGroupName := args[0]
	var orderBy cloudwatchlogs.OrderBy = "LastEventTime"
	var sorting = internal.LogStreamSorting{
		Descending: &descending,
		OrderBy: orderBy,
	}

	res, err := c.DescribeLogStreams(&limit, &logGroupName, nil, &sorting, nil)

	result := make([]cloudwatchlogs.LogStream, 0, 100)
	for {
		if err != nil {
			return errors.Wrap(err, "get streams")
		}

		result = append(result, res.LogStreams...)

		if res.NextToken != nil {
			res, err = c.DescribeLogStreams(&limit, &logGroupName, nil, &sorting, res.NextToken)
		} else {
			break
		}
	}

	for _, s := range result {
		fmt.Println(*s.LogStreamName)
	}

	return nil
}
