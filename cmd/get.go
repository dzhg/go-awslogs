package cmd

import (
	"fmt"
	"github.com/dzhg/go-awslogs/internal"
	"github.com/dzhg/go-awslogs/internal/tsp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <group> [stream1 [stream2] ...]",
	Short: "get log events from AWS CloudWatch log streams",
	RunE:  getLogs,
	Args:  cobra.MinimumNArgs(1),
}

func initGetCmd() *cobra.Command {
	getCmd.Flags().BoolP("watch", "w", false, "watch the log stream")
	getCmd.Flags().StringP("start", "s", "1m", "start timestamp")
	getCmd.Flags().StringP("end", "e", "", "end timestamp")
	getCmd.Flags().BoolP("stream-name", "N", false, "print stream name")

	return getCmd
}

func getLogs(cmd *cobra.Command, args []string) error {
	c, err := internal.NewAwsLogsClient()
	if err != nil {
		return errors.Wrap(err, "new AWS log client")
	}

	watch, err := cmd.Flags().GetBool("watch")

	if err != nil {
		return errors.Wrap(err, "watch")
	}

	start, err := cmd.Flags().GetString("start")

	if err != nil {
		return errors.Wrap(err, "start")
	}

	startTime, err := tsp.ParseString(start)

	if err != nil {
		return errors.Wrap(err, "start")
	}

	end, err := cmd.Flags().GetString("end")
	if err != nil {
		return errors.Wrap(err, "end")
	}

	var endTimePtr *int64
	if end != "" {
		endTime, err := tsp.ParseString(end)
		if err != nil {
			return errors.Wrap(err, "end")
		}
		endTimePtr = &endTime
	}

	limit := int64(1000)
	logGroupName := args[0]
	var logStreamNames []string
	if len(args) == 1 {
		logStreamNames = nil
	} else {
		logStreamNames = args[1:]
	}

	params := internal.GetLogsParams{
		Limit:          &limit,
		StartTime:      &startTime,
		EndTime:        endTimePtr,
		LogGroupName:   &logGroupName,
		LogStreamNames: logStreamNames,
	}

	events := make(chan *internal.LogEvent)

	go c.StreamLogEvents(&params, events, &watch)

	for e := range events {
		if e.Error != nil {
			return e.Error
		}
		printEvent(e, cmd)
	}

	return nil
}

func printEvent(e *internal.LogEvent, cmd *cobra.Command) {
	showStream, _ := cmd.Flags().GetBool("stream-name")
	if showStream {
		fmt.Printf("%s ", e.StreamName)
	}
	fmt.Println(e.Message)
}
