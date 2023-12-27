package router

const (
	// AWSLambdaSyncInvocationPayloadMax exceed response and response Lambda invocation payload quota
	// https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-limits.html
	AWSLambdaSyncInvocationPayloadMax = 6291456 // 6MB
)

const (
	AWSLambdaExceedInvocationPayloadCode = "AWS_LAMBDA_EXCEEDED_PAYLOAD_SIZE"
	GRPCResourceExhausted                = "GRPC_RESOURCE_EXHAUSTED"
)

// ErrorStatus presents the error status model
type ErrorStatus struct {
	// Code indicate the unique error identifier code
	Code string `json:"code"`
	// Message presents the error message description
	Message string `json:"message"`
}

// NewInternalErr returns the internal server error status
func NewInternalErr() ErrorStatus {
	return ErrorStatus{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Internal Server Error",
	}
}

// NewExceededInvocationPayloadErr returns the Lambda invocation payload exceed payload error status.
func NewExceededInvocationPayloadErr() ErrorStatus {
	return ErrorStatus{
		Code:    AWSLambdaExceedInvocationPayloadCode,
		Message: "body size is too long",
	}
}

// NewgRPCResourceExhaustedErr returns the error resource exhausted status.
func NewgRPCResourceExhaustedErr() ErrorStatus {
	return ErrorStatus{
		Code:    GRPCResourceExhausted,
		Message: "received resource larger than max size",
	}
}
