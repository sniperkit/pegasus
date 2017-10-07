package netamqp

import (
	"github.com/streadway/amqp"
	"strconv"
	"time"
)

// publishing embeds the amqp.Publishing struct in order to configure it properly.
type publishing struct {
	amqp.Publishing
}

// configurePublishing configures the publishing method, assigning the headers value in publishing properties.
func (p *publishing) configurePublishing(headers map[string]string) {

	if headers["Content-Type"] != "" {
		p.ContentType = headers["Content-Type"]
	}

	if headers["MQ-Content-Encoding"] != "" {
		p.ContentEncoding = headers["MQ-Content-Encoding"]
	}

	if headers["MQ-Delivery-Mode"] != "" {
		v, _ := strconv.Atoi(headers["MQ-Delivery-Mode"])
		p.DeliveryMode = uint8(v)
	}

	if headers["MQ-Priority"] != "" {
		v, _ := strconv.Atoi(headers["MQ-Priority"])
		p.Priority = uint8(v)
	}

	if headers["MQ-Correlation-Id"] != "" {
		p.CorrelationId = headers["MQ-Correlation-Id"]
	}

	if headers["MQ-Reply-To"] != "" {
		p.ReplyTo = headers["MQ-Reply-To"]
	}

	if headers["MQ-Expiration"] != "" {
		p.Expiration = headers["MQ-Expiration"]
	}

	if headers["MQ-Message-Id"] != "" {
		p.MessageId = headers["MQ-Message-Id"]
	}

	if headers["MQ-Timestamp"] != "" {
		layout := "2006-01-02T15:04:05.000Z"
		t, _ := time.Parse(layout, headers["MQ-Timestamp"])
		p.Timestamp = t
	}

	if headers["MQ-Type"] != "" {
		p.Type = headers["MQ-Type"]
	}

	if headers["MQ-User-Id"] != "" {
		p.UserId = headers["MQ-User-Id"]
	}

	if headers["MQ-App-Id"] != "" {
		p.AppId = headers["MQ-App-Id"]
	}

}
