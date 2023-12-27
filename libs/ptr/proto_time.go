package ptr

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TimestamppbToTime(timestamp *timestamppb.Timestamp) *time.Time {
	if timestamp == nil {
		return nil
	}
	return Time(timestamp.AsTime())
}

func TimeToTimestamppb(time *time.Time) *timestamppb.Timestamp {
	if time == nil {
		return nil
	}
	return timestamppb.New(*time)
}
