package netamqp

import (
	"github.com/cpapidas/pegasus/peg"
	"github.com/streadway/amqp"
	"strconv"
	"time"
)

// publishing embeds the amqp.Publishing struct in order to configure it properly.
type publishing struct {
	amqp.Publishing
	optionHeaders map[string]string
}

// NewPublishing create and returns a new object
func buildPublishing(options *peg.Options, body []byte) publishing {
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
		if peg.IsAMQPValidHeader(k) {
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
	p.setContentType()
	p.setContent()
	p.setPriority()
	p.setCorrelation()
	p.setDelivery()
	p.setReply()
	p.setExpiration()
	p.setMessage()
	p.setTimestamp()
	p.setType()
	p.setUser()
	p.setApp()
}

// setContentType set header
func (p *publishing) setContentType() {
	if p.optionHeaders["Content-Type"] != "" {
		p.ContentType = p.optionHeaders["Content-Type"]
	}
}

// setContent set header
func (p *publishing) setContent() {
	if p.optionHeaders["MQ-Content-Encoding"] != "" {
		p.ContentEncoding = p.optionHeaders["MQ-Content-Encoding"]
	}
}

// setDelivery set header
func (p *publishing) setDelivery() {
	if p.optionHeaders["MQ-Delivery-Mode"] != "" {
		v, _ := strconv.Atoi(p.optionHeaders["MQ-Delivery-Mode"])
		p.DeliveryMode = uint8(v)
	}
}

// setPriority set header
func (p *publishing) setPriority() {
	if p.optionHeaders["MQ-Priority"] != "" {
		v, _ := strconv.Atoi(p.optionHeaders["MQ-Priority"])
		p.Priority = uint8(v)
	}
}

// setCorrelation set header
func (p *publishing) setCorrelation() {
	if p.optionHeaders["MQ-Correlation-Id"] != "" {
		p.CorrelationId = p.optionHeaders["MQ-Correlation-Id"]
	}
}

// setReply set header
func (p *publishing) setReply() {
	if p.optionHeaders["MQ-Reply-To"] != "" {
		p.ReplyTo = p.optionHeaders["MQ-Reply-To"]
	}
}

// setExpiration set header
func (p *publishing) setExpiration() {
	if p.optionHeaders["MQ-Expiration"] != "" {
		p.Expiration = p.optionHeaders["MQ-Expiration"]
	}
}

// setMessage set header
func (p *publishing) setMessage() {
	if p.optionHeaders["MQ-Message-Id"] != "" {
		p.MessageId = p.optionHeaders["MQ-Message-Id"]
	}
}

// setTimestamp set header
func (p *publishing) setTimestamp() {
	if p.optionHeaders["MQ-Timestamp"] != "" {
		layout := "2006-01-02T15:04:05.000Z"
		t, _ := time.Parse(layout, p.optionHeaders["MQ-Timestamp"])
		p.Timestamp = t
	}
}

// setType set header
func (p *publishing) setType() {
	if p.optionHeaders["MQ-Type"] != "" {
		p.Type = p.optionHeaders["MQ-Type"]
	}
}

// setUser set header
func (p *publishing) setUser() {
	if p.optionHeaders["MQ-User-Id"] != "" {
		p.UserId = p.optionHeaders["MQ-User-Id"]
	}
}

func (p *publishing) setApp() {
	if p.optionHeaders["MQ-App-Id"] != "" {
		p.AppId = p.optionHeaders["MQ-App-Id"]
	}
}
