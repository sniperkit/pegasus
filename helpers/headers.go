package helpers

const (
	// HTTPHeaderKey is the header HTTP header prefix
	HTTPHeaderKey = "HP-"
	// GRPCHeaderKey is the header GRPC header prefix
	GRPCHeaderKey = "GR-"
	// AMQPHeaderKey is the header AMQP header prefix
	AMQPHeaderKey = "MQ-"
)

// IsHTTPValidHeader returns true if the key is less than 3 characters or is not equal to GRPCHeaderKey & AMQPHeaderKey
func IsHTTPValidHeader(key string) bool {
	if len(key) < 3 {
		return true
	} else if key[0:3] != GRPCHeaderKey && key[0:3] != AMQPHeaderKey {
		return true
	}
	return false
}

// IsGRPCValidHeader returns true if the key is less than 3 characters or is not equal to HTTPHeaderKey & AMQPHeaderKey
func IsGRPCValidHeader(key string) bool {
	if len(key) < 3 {
		return true
	} else if key[0:3] != HTTPHeaderKey && key[0:3] != AMQPHeaderKey {
		return true
	}
	return false
}

// IsAMQPValidHeader returns true if the key is less than 3 characters or is not equal to GRPCHeaderKey & HTTPHeaderKey
func IsAMQPValidHeader(key string) bool {
	if len(key) < 3 {
		return true
	} else if key[0:3] != GRPCHeaderKey && key[0:3] != HTTPHeaderKey {
		return true
	}
	return false
}
