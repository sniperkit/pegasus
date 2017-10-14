package netamqp

import (
	"github.com/streadway/amqp"
	"strconv"
	"time"
	"github.com/cpapidas/pegasus/network"
	"github.com/cpapidas/pegasus/helpers"
)

// publishing embeds the amqp.Publishing struct in order to configure it properly.
type publishing struct {
	amqp.Publishing
	optionHeaders map[string]string
}

// NewPublishing create and returns a new object
func buildPublishing(options *network.Options, body []byte) publishing {
	pub := publishing{}
	pub.DeliveryMode = amqp.Persistent
	pub.ContentType = "text/plain"
	pub.optionHeaders = options.GetHeaders()
	pub.configurePublishing()
	pub.Headers = make(map[string]interface{})
	pub.setHeaders()
	pub.Body = body
	pub.setParams(options.GetParams())
	return pub
}

// setHeaders sets the publishing Headers
func (p *publishing) setHeaders() {
	for k, v := range p.optionHeaders {
		if helpers.IsAMQPValidHeader(k) {
			p.Headers[k] = v
		}
	}
}

// setParams sets the publishing Params using the headers flag. Params are traveling through headers using the prefix of
// MP-<whatever>. Those headers are temporary and removed in Listen method.
func (p *publishing) setParams(params map[string]string) {
	for k, v := range params {
		p.Headers["MP-"+k] = v
	}
}

// configurePublishing configures the publishing method, assigning the headers value in publishing properties.
func (p *publishing) configurePublishing() {

	if p.optionHeaders["Content-Type"] != "" {
		p.ContentType = p.optionHeaders["Content-Type"]
	}

	if p.optionHeaders["MQ-Content-Encoding"] != "" {
		p.ContentEncoding = p.optionHeaders["MQ-Content-Encoding"]
	}

	if p.optionHeaders["MQ-Delivery-Mode"] != "" {
		v, _ := strconv.Atoi(p.optionHeaders["MQ-Delivery-Mode"])
		p.DeliveryMode = uint8(v)
	}

	if p.optionHeaders["MQ-Priority"] != "" {
		v, _ := strconv.Atoi(p.optionHeaders["MQ-Priority"])
		p.Priority = uint8(v)
	}

	if p.optionHeaders["MQ-Correlation-Id"] != "" {
		p.CorrelationId = p.optionHeaders["MQ-Correlation-Id"]
	}

	if p.optionHeaders["MQ-Reply-To"] != "" {
		p.ReplyTo = p.optionHeaders["MQ-Reply-To"]
	}

	if p.optionHeaders["MQ-Expiration"] != "" {
		p.Expiration = p.optionHeaders["MQ-Expiration"]
	}

	if p.optionHeaders["MQ-Message-Id"] != "" {
		p.MessageId = p.optionHeaders["MQ-Message-Id"]
	}

	if p.optionHeaders["MQ-Timestamp"] != "" {
		layout := "2006-01-02T15:04:05.000Z"
		t, _ := time.Parse(layout, p.optionHeaders["MQ-Timestamp"])
		p.Timestamp = t
	}

	if p.optionHeaders["MQ-Type"] != "" {
		p.Type = p.optionHeaders["MQ-Type"]
	}

	if p.optionHeaders["MQ-User-Id"] != "" {
		p.UserId = p.optionHeaders["MQ-User-Id"]
	}

	if p.optionHeaders["MQ-App-Id"] != "" {
		p.AppId = p.optionHeaders["MQ-App-Id"]
	}

}
