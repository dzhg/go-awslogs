package tsp

import (
	"log"
	"reflect"
	"testing"
	"time"
)

func safeParse(layout, s string) time.Time {
	t, err := time.Parse(layout, s)
	if err != nil {
		log.Println(err.Error())
	}
	return t
}

func rfc3339(s string) time.Time {
	return safeParse(time.RFC3339, s)
}

func TestParseRelative(t *testing.T) {
	type args struct {
		s string
		t time.Time
	}
	newArgs := func(s string, t string) args {
		return args{s, rfc3339(t)}
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			"1m",
			newArgs("1m", "2020-07-21T10:49:10-07:00"),
			rfc3339("2020-07-21T10:48:10-07:00"),
			false,
		},
		{
			"45 seconds",
			newArgs("45 seconds", "2020-07-21T10:49:50-07:00"),
			rfc3339("2020-07-21T10:49:05-07:00"),
			false,
		},
		{
			"1.5 hour",
			newArgs("1.5 hour", "2020-07-21T10:49:50-07:00"),
			rfc3339("2020-07-21T09:19:50-07:00"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRelative(tt.args.s, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRelative() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRelative() got = %v, want %v", got, tt.want)
			}
		})
	}
}