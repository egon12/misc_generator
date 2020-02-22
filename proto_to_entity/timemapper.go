package proto_to_entity

const timemapper = `
package mapper

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func FromTimestamp(input *timestamp.Timestamp) *time.Time {
	r, _ := ptypes.Timestamp(input)
	return &r
}

func ToTimestamp(input *time.Time) *timestamp.Timestamp {
	r, _ := ptypes.TimestampProto(*input)
	return r
}

`
