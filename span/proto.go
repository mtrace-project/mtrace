package span

import (
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Span) ToProto() *SpanProto {
	if s == nil {
		return nil
	}

	protoSpan := &SpanProto{
		SpanId:        s.SpanId,
		ParentId:      s.ParentId,
		ServiceName:   s.ServiceName,
		OperationName: s.OperationName,
		SpanKind:      s.SpanKind,
		SpanStatus:    s.SpanStatus,
	}

	protoSpan.StartTime = timestamppb.New(s.StartTime)
	protoSpan.EndTime = timestamppb.New(s.EndTime)
	protoSpan.Duration = durationpb.New(s.Duration)

	if len(s.Attributes) > 0 {
		protoSpan.Attributes = make(map[string]*structpb.Value)
		for k, v := range s.Attributes {
			if val, err := structpb.NewValue(v); err == nil {
				protoSpan.Attributes[k] = val
			}
		}
	}

	return protoSpan
}
