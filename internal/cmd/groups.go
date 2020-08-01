package cmd

import (
	"fmt"
	"github.com/dzhg/go-awslogs/internal"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "get log groups",
	RunE:  groups,
}

func initGroupCmd() *cobra.Command {
	groupsCmd.Flags().StringP("prefix", "p", "", "log group prefix")
	return groupsCmd
}

func groups(cmd *cobra.Command, args []string) error {
	c, err := internal.NewAwsLogsClient()
	if err != nil {
		return errors.Wrap(err, "new AWS logs client")
	}

	prefix, err := cmd.Flags().GetString("prefix")
	if err != nil {
		return errors.Wrap(err, "prefix")
	}

	var prefixPtr *string
	if prefix != "" {
		prefixPtr = &prefix
	}

	limit := int64(50)

	res, err := c.DescribeLogGroups(&limit, prefixPtr, nil)

	result := make([]cloudwatchlogs.LogGroup, 0, 100)
	for {
		if err != nil {
			return errors.Wrap(err, "get groups")
		}

		result = append(result, res.LogGroups...)

		if res.NextToken != nil {
			res, err = c.DescribeLogGroups(&limit, nil, res.NextToken)
		} else {
			break
		}
	}

	for _, g := range result {
		fmt.Println(*g.LogGroupName)
	}

	return nil
}
