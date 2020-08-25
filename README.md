# go-awslogs

[![Go Report Card](https://goreportcard.com/badge/github.com/dzhg/go-awslogs)](https://goreportcard.com/report/github.com/dzhg/go-awslogs)
![Github Workflow: Go](https://github.com/dzhg/go-awslogs/workflows/Go/badge.svg)

## Introduction

Inspired by https://github.com/jorgebastida/awslogs.

`go-awslogs` is a command line utility for easier access of AWS CloudWatch logs.

TL;DR:

```shell script
# Show all log groups
go-awslogs groups

# Show all log streams of a group
go-awslogs streams <group_name>

# Watch log stream
go-awslogs get <group_name> --watch

# Show help
go-awslogs help
```

## Usage

### Get Log Groups

```shell script
go-awslogs groups

# Filter the result with group name prefix
go-awslogs groups --prefix <prefix>
```

### Get Log Streams

```shell script
go-awslogs streams <group_name>
```

### Get Log Events

```shell script
go-awslogs get <group_name> [stream_name1 [, stream_name2 [,...]]]
```

Additional arguments:

`--watch`, `-w`: Watch the log streams (like `tail -f`)

`--start`, `-s`: Filter the log events by starting timestamp (Show log events after the timestamp)

`--end`, `-e`: Filter the log events by end timestamp (Show log events before the timestamp)

`--stream-name`, `-N`: Print the stream name as prefix of each line

### Time Parsing

For the timestamp arguments like `--start` and `--end`, `go-awslogs` accepts 3 forms
of time string:

**Relative Time**: `5min ago`, `1.5hours ago`, etc. (`ago` can be omitted)

**RFC3339**: `2020-07-23T15:23:35Z` or `2020-07-23T08:23:35-07:00`

**Human**: `01/02/2006 15:04:05 -07:00`

### AWS Profile

`go-awslogs` internally uses [`aws-sdk-go-v2`](https://github.com/aws/aws-sdk-go-v2). It loads
the default AWS profile without any setting.

A different AWS profile can be used by setting an environment variable:

```shell script
AWS_PROFILE=my_other_profile go-awslogs groups
```

Also, you can set environment variables directly for `AWS_SECRET_ACCESS_KEY`, `AWS_ACCESS_KEY` and `AWS_REGION` etc.
