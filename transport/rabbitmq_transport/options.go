package rabbitmq_transport

import (
	"bitbucket.org/code_horse/pegasus/network"
	"github.com/streadway/amqp"
	"strconv"
	"time"
)

type Options struct {
	network.Options
}

func NewOptions() *Options {
	options := &Options{}
	options.Fields = make(map[string]map[string]string)
	return options
}

func NewSendOptions() *Options {
	options := &Options{}
	options.Fields = make(map[string]map[string]string)
	options.SetPublishContentType("text/plain").
		SetPublishDeliveryModePersistent().
		SetPriority(1)
	return options
}

func BuildOptions(networkOptions *network.Options) *Options {
	return &Options{Options: *networkOptions}
}

func BuildOptionsByBytes(data []byte) *Options {
	options := &Options{}
	return BuildOptions(options.Unmarshal(data))
}

func (o *Options) GetOptions() *network.Options {
	return &o.Options
}

func (o *Options) SetPublishContentType(value string) *Options {
	o.initMapper("CONTENT-TYPE")
	o.Fields["CONTENT-TYPE"]["VALUE"] = value
	return o
}

func (o *Options) GetPublishContentType() string {
	return o.Fields["CONTENT-TYPE"]["VALUE"]
}

func (o *Options) SetPublishContentEncoding(value string) *Options {
	o.initMapper("CONTENT-TYPE")
	o.Fields["CONTENT-TYPE"]["VALUE"] = value
	return o
}

func (o *Options) GetPublishContentEncoding() string {
	return o.Fields["CONTENT-TYPE"]["VALUE"]
}

func (o *Options) SetPublishDeliveryMode(value uint8) *Options {
	o.initMapper("DELIVERY-MODE")
	o.Fields["DELIVERY-MODE"]["VALUE"] = strconv.Itoa(int(value))
	return o
}

func (o *Options) SetPublishDeliveryModePersistent() *Options {
	o.initMapper("DELIVERY-MODE")
	o.Fields["DELIVERY-MODE"]["VALUE"] = strconv.Itoa(int(amqp.Persistent))
	return o
}

func (o *Options) SetPublishDeliveryModeTransient() *Options {
	o.initMapper("DELIVERY-MODE")
	o.Fields["DELIVERY-MODE"]["VALUE"] = strconv.Itoa(int(amqp.Transient))
	return o
}

func (o *Options) GetPublishDeliveryMode() uint8 {
	value, err := strconv.Atoi(o.Fields["DELIVERY-MODE"]["VALUE"])
	if err != nil {
		panic(err)
	}
	return uint8(value)
}

func (o *Options) SetPriority(value uint8) *Options {
	o.initMapper("PRIORITY")
	o.Fields["PRIORITY"]["VALUE"] = strconv.Itoa(int(value))
	return o
}

func (o *Options) GetPriority() uint8 {
	value, err := strconv.Atoi(o.Fields["PRIORITY"]["VALUE"])
	if err != nil {
		panic(err)
	}
	return uint8(value)
}

func (o *Options) SetPublishCorrelationId(value string) *Options {
	o.initMapper("CORRELATION")
	o.Fields["CORRELATION"]["VALUE"] = value
	return o
}

func (o *Options) GetPublishCorrelationId() string {
	return o.Fields["CORRELATION"]["VALUE"]
}

func (o *Options) SetPublishReplyTo(value string) *Options {
	o.initMapper("REPLY-TO")
	o.Fields["REPLY-TO"]["VALUE"] = value
	return o
}

func (o *Options) GetPublishReplyTo() string {
	return o.Fields["REPLY-TO"]["VALUE"]
}

func (o *Options) SetPublishExpiration(value string) *Options {
	o.initMapper("EXPIRATION")
	o.Fields["EXPIRATION"]["VALUE"] = value
	return o
}

func (o *Options) GetPublishExpiration() string {
	return o.Fields["EXPIRATION"]["VALUE"]
}

func (o *Options) SetPublishMessageId(value string) *Options {
	o.initMapper("MESSAGE-ID")
	o.Fields["MESSAGE-ID"]["VALUE"] = value
	return o
}

func (o *Options) GetPublishMessageId() string {
	return o.Fields["MESSAGE-ID"]["VALUE"]
}

func (o *Options) SetPublishTimestamp(value time.Time) *Options {
	o.initMapper("TIMESTAMP")
	o.Fields["TIMESTAMP"]["VALUE"] = value.String()
	return o
}

func (o *Options) GetPublishTimestamp() time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, o.Fields["TIMESTAMP"]["VALUE"])
	if err != nil {
		return time.Now()
	}
	return t
}

func (o *Options) SetPublishType(value string) *Options {
	o.initMapper("TYPE")
	o.Fields["TYPE"]["VALUE"] = value
	return o
}

func (o *Options) GetPublishType() string {
	return o.Fields["TYPE"]["VALUE"]
}

func (o *Options) SetPublishUserId(value string) *Options {
	o.initMapper("USER-ID")
	o.Fields["USER-ID"]["VALUE"] = value
	return o
}

func (o *Options) GetPublishUserId() string {
	return o.Fields["USER-ID"]["VALUE"]
}

func (o *Options) SetPublishAppId(value string) *Options {
	o.initMapper("APP-ID")
	o.Fields["APP-ID"]["VALUE"] = value
	return o
}

func (o *Options) GetPublishAppId() string {
	return o.Fields["APP-ID"]["VALUE"]
}

func (o *Options) initMapper(key string) {
	if o.Fields[key] == nil {
		o.Fields[key] = make(map[string]string)
	}
}
