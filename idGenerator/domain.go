package idgenerator

const (
	TRACE_ID_LENGTH = 32
	SPAN_ID_LENGTH  = 16
)

type IdGenerator interface {
	Generate(length int) (string, error)
}
