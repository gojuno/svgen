package tests

type (
	// Supported types
	StringType string
	IntType    int

	// Non-supported type
	FloatType float64

	// Types without values
	String2Type string
)

const (
	StringValue1 StringType = "1"
	StringValue2 StringType = "2"

	IntValue1 IntType = 1
	IntValue2 IntType = 2

	FloatValue1 FloatType = 1.0
	FloatValue2 FloatType = 2.0
)
