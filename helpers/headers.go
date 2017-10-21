package helpers

const (
	// HTTPHeaderKey is the HTTP header prefix
	HTTPHeaderKey = "HP-"
	// GRPCHeaderKey is the GRPC header prefix
	GRPCHeaderKey = "GR-"
	// AMQPHeaderKey is the AMQP header prefix
	AMQPHeaderKey = "MQ-"
	// AMQPPramKey is the AMQP param prefix
	AMQPPramKey = "MP-"
)

// IsHTTPValidHeader returns true if the key is less than 3 characters or is not equal to GRPCHeaderKey & AMQPHeaderKey
func IsHTTPValidHeader(key string) bool {
	if len(key) > 3 && key[0:3] != GRPCHeaderKey && key[0:3] != AMQPHeaderKey && key[0:3] != AMQPPramKey {
		return true
	}
	return false
}

// IsGRPCValidHeader returns true if the key is less than 3 characters or is not equal to HTTPHeaderKey & AMQPHeaderKey
func IsGRPCValidHeader(key string) bool {
	if len(key) > 3 && key[0:3] != HTTPHeaderKey && key[0:3] != AMQPHeaderKey && key[0:3] != AMQPPramKey {
		return true
	}
	return false
}

// IsAMQPValidHeader returns true if the key is less than 3 characters or is not equal to GRPCHeaderKey & HTTPHeaderKey
func IsAMQPValidHeader(key string) bool {
	if len(key) > 3 && key[0:3] != GRPCHeaderKey && key[0:3] != HTTPHeaderKey && key[0:3] != AMQPPramKey {
		return true
	}
	return false
}

// AMQPParam returns true if the key is less than 3 characters or is not equal to GRPCHeaderKey & HTTPHeaderKey
func AMQPParam(key string) string {
	if len(key) > 3 && key[0:3] == AMQPPramKey {
		return key[3:]
	}
	return ""
}
