package rabbitmq_transport

import (
	"bitbucket.org/code_horse/pegasus/network"
	"strconv"
)

type Properties struct {
	network.Properties
}

func NewProperties() *Properties {
	properties := &Properties{}
	properties.Fields = make(map[string]map[string]string)
	return properties
}

func NewSendProperties() *Properties {
	properties := &Properties{}
	properties.Fields = make(map[string]map[string]string)
	properties.SetQueueDurable(true).
		SetQueueAutoDelete(false).
		SetQueueExclusive(false).
		SetQueueNoWait(false).
		SetPublishMandatory(false).
		SetPublishImmediate(false).
		SetPublishExchange("")
	return properties
}

func NewListenProperties() *Properties {
	properties := &Properties{}
	properties.Fields = make(map[string]map[string]string)
	properties.
		SetQueueDurable(true).
		SetQueueAutoDelete(false).
		SetQueueExclusive(false).
		SetQueueNoWait(false).
		SetQueueBindNoWait(false).
		SetConsumeName("").
		SetConsumeAutoAct(false).
		SetConsumeExclusive(false).
		SetConsumeNoLocal(false).
		SetConsumeNoWait(false).
		SetQosPrefetchCount(1).
		SetQosPrefetchSize(0).
		SetQosGlobal(false)
	return properties
}


func BuildRabbitMQProperties(properties *network.Properties) *Properties {
	return &Properties{Properties: *properties}
}

func (p *Properties) GetProperties() *network.Properties {
	return &p.Properties
}

func (p *Properties) initMapper(key string) {
	if p.Fields[key] == nil {
		p.Fields[key] = make(map[string]string)
	}
}

func (p *Properties) SetPath(value string) *Properties {
	p.Path = value
	p.SetQueueName(value)
	return p
}

func (p *Properties) GetPath() string {
	return p.Path
}

func (p *Properties) SetKey(value string) *Properties {
	p.initMapper("KEY")
	p.Fields["KEY"]["VALUE"] = value
	return p
}

func (p *Properties) GetKey() string {
	return p.Fields["KEY"]["VALUE"]
}


func (p *Properties) SetPublishExchange(value string) *Properties {
	p.initMapper("PUBLISH-EXCHANGE")
	p.Fields["PUBLISH-EXCHANGE"]["VALUE"] = value
	return p
}

func (p *Properties) GetPublishExchange() string {
	return p.Fields["PUBLISH-EXCHANGE"]["VALUE"]
}

func (p *Properties) SetPublishMandatory(value bool) *Properties {
	p.initMapper("MANDATORY")
	p.Fields["MANDATORY"]["VALUE"] = strconv.FormatBool(value)
	return p
}

func (p *Properties) GetPublishMandatory() bool {
	if p.Fields["MANDATORY"]["VALUE"] == "true" {
		return true
	}
	return false
}

func (p *Properties) SetPublishImmediate(value bool) *Properties {
	p.initMapper("IMMEDIATE")
	p.Fields["IMMEDIATE"]["VALUE"] = strconv.FormatBool(value)
	return p
}

func (p *Properties) GetPublishImmediate() bool {
	if p.Fields["IMMEDIATE"]["VALUE"] == "true" {
		return true
	}
	return false
}

func (p *Properties) SetDeclarationName(name string) *Properties {
	p.initMapper("DECLARATION-DECLARE")
	p.Fields["DECLARATION-DECLARE"]["NAME"] = name
	return p
}

func (p *Properties) GetDeclarationName() string {
	return p.Fields["DECLARATION-DECLARE"]["NAME"]
}

func (p *Properties) SetDeclarationTypeToTopic() *Properties {
	p.initMapper("DECLARATION-DECLARE")
	p.Fields["DECLARATION-DECLARE"]["TYPE"] = "topic"
	return p
}

func (p *Properties) SetDeclarationTypeToDirect() *Properties {
	p.initMapper("DECLARATION-DECLARE")
	p.Fields["DECLARATION-DECLARE"]["TYPE"] = "direct"
	return p
}

func (p *Properties) SetDeclarationTypeToHeaders() *Properties {
	p.initMapper("DECLARATION-DECLARE")
	p.Fields["DECLARATION-DECLARE"]["TYPE"] = "headers"
	return p
}

func (p *Properties) SetDeclarationTypeToFanout() *Properties {
	p.initMapper("DECLARATION-DECLARE")
	p.Fields["DECLARATION-DECLARE"]["TYPE"] = "fanout"
	return p
}

func (p *Properties) GetDeclarationType() string {
	return p.Fields["DECLARATION-DECLARE"]["TYPE"]
}

func (p *Properties) SetDeclarationDurable(durable bool) *Properties {
	p.initMapper("DECLARATION-DECLARE")
	p.Fields["DECLARATION-DECLARE"]["DURABLE"] = strconv.FormatBool(durable)
	return p
}

func (p *Properties) GetDeclarationDurable() bool {
	if p.Fields["DECLARATION-DECLARE"]["DURABLE"] == "false" {
		return false
	}
	return true
}

func (p *Properties) SetDeclarationAutoDelete(durable bool) *Properties {
	p.initMapper("DECLARATION-DECLARE")
	p.Fields["DECLARATION-DECLARE"]["AUTO-DELETE"] = strconv.FormatBool(durable)
	return p
}

func (p *Properties) GetDeclarationAutoDelete() bool {
	if p.Fields["DECLARATION-DECLARE"]["AUTO-DELETE"] == "true" {
		return true
	}
	return false
}

func (p *Properties) SetDeclarationExclusive(durable bool) *Properties {
	p.initMapper("DECLARATION-DECLARE")
	p.Fields["DECLARATION-DECLARE"]["EXCLUSIVE"] = strconv.FormatBool(durable)
	return p
}

func (p *Properties) GetDeclarationExclusive() bool {
	if p.Fields["DECLARATION-DECLARE"]["EXCLUSIVE"] == "true" {
		return true
	}
	return false
}

func (p *Properties) GetDeclarationNoWait() bool {
	if p.Fields["DECLARATION-DECLARE"]["NO-WAIT"] == "true" {
		return true
	}
	return false
}

func (p *Properties) SetDeclarationNoWait(durable bool) *Properties {
	p.initMapper("DECLARATION-DECLARE")
	p.Fields["DECLARATION-DECLARE"]["NO-WAIT"] = strconv.FormatBool(durable)
	return p
}

// ----------------------------- QUEUE ----------------------------- //

func (p *Properties) SetQueueName(value string) *Properties {
	p.initMapper("DECLARATION-QUEUE-DECLARE")
	p.Fields["DECLARATION-QUEUE-DECLARE"]["NAME"] = value
	return p
}

func (p *Properties) GetQueueName() string {
	return p.Fields["DECLARATION-QUEUE-DECLARE"]["NAME"]
}

func (p *Properties) SetQueueDurable(durable bool) *Properties {
	p.initMapper("DECLARATION-QUEUE-DECLARE")
	p.Fields["DECLARATION-QUEUE-DECLARE"]["DURABLE"] = strconv.FormatBool(durable)
	return p
}

func (p *Properties) GetQueueDurable() bool {
	if p.Fields["DECLARATION-QUEUE-DECLARE"]["DURABLE"] == "false" {
		return false
	}
	return true
}

func (p *Properties) SetQueueAutoDelete(durable bool) *Properties {
	p.initMapper("DECLARATION-QUEUE-DECLARE")
	p.Fields["DECLARATION-QUEUE-DECLARE"]["AUTO-DELETE"] = strconv.FormatBool(durable)
	return p
}

func (p *Properties) GetQueueAutoDelete() bool {
	if p.Fields["DECLARATION-QUEUE-DECLARE"]["AUTO-DELETE"] == "true" {
		return true
	}
	return false
}

func (p *Properties) SetQueueExclusive(durable bool) *Properties {
	p.initMapper("DECLARATION-QUEUE-DECLARE")
	p.Fields["DECLARATION-QUEUE-DECLARE"]["EXCLUSIVE"] = strconv.FormatBool(durable)
	return p
}

func (p *Properties) GetQueueExclusive() bool {
	if p.Fields["DECLARATION-QUEUE-DECLARE"]["EXCLUSIVE"] == "true" {
		return true
	}
	return false
}

func (p *Properties) GetQueueNoWait() bool {
	if p.Fields["DECLARATION-QUEUE-DECLARE"]["NO-WAIT"] == "true" {
		return true
	}
	return false
}

func (p *Properties) SetQueueNoWait(durable bool) *Properties {
	p.initMapper("DECLARATION-QUEUE-DECLARE")
	p.Fields["DECLARATION-QUEUE-DECLARE"]["NO-WAIT"] = strconv.FormatBool(durable)
	return p
}

func (p *Properties) GetQueueBindNoWait() bool {
	if p.Fields["DECLARATION-QUEUE-DECLARE"]["NO-WAIT-BIND"] == "true" {
		return true
	}
	return false
}

func (p *Properties) SetQueueBindNoWait(durable bool) *Properties {
	p.initMapper("DECLARATION-QUEUE-DECLARE")
	p.Fields["DECLARATION-QUEUE-DECLARE"]["NO-WAIT-BIND"] = strconv.FormatBool(durable)
	return p
}

// ----------------------------- CONSUMER ----------------------------- //

func (p *Properties) SetConsumeName(value string) *Properties {
	p.initMapper("CONSUME")
	p.Fields["CONSUME"]["NAME"] = value
	return p
}

func (p *Properties) GetConsumeName() string {
	return p.Fields["CONSUME"]["NAME"]
}

func (p *Properties) SetConsumeAutoAct(durable bool) *Properties {
	p.initMapper("CONSUME")
	p.Fields["CONSUME"]["AUTO-ACT"] = strconv.FormatBool(durable)
	return p
}

func (p *Properties) GetConsumeAutoAct() bool {
	if p.Fields["CONSUME"]["AUTO-ACT"] == "true" {
		return true
	}
	return false
}

func (p *Properties) SetConsumeExclusive(value bool) *Properties {
	p.initMapper("CONSUME")
	p.Fields["CONSUME"]["EXCLUSIVE"] = strconv.FormatBool(value)
	return p
}

func (p *Properties) GetConsumeExclusive() bool {
	if p.Fields["CONSUME"]["EXCLUSIVE"] == "true" {
		return true
	}
	return false
}

func (p *Properties) SetConsumeNoLocal(value bool) *Properties {
	p.initMapper("CONSUME")
	p.Fields["CONSUME"]["NO-LOCAL"] = strconv.FormatBool(value)
	return p
}

func (p *Properties) GetConsumeNoLocal() bool {
	if p.Fields["CONSUME"]["NO-LOCAL"] == "true" {
		return true
	}
	return false
}

func (p *Properties) SetConsumeNoWait(value bool) *Properties {
	p.initMapper("CONSUME")
	p.Fields["CONSUME"]["NO-WAIT"] = strconv.FormatBool(value)
	return p
}

func (p *Properties) GetConsumeNoWait() bool {
	if p.Fields["CONSUME"]["NO-WAIT"] == "true" {
		return true
	}
	return false
}

// ----------------------------- QOS ----------------------------- //

func (p *Properties) SetQosPrefetchCount(value int) *Properties {
	p.initMapper("QOS")
	p.Fields["QOS"]["PREFETCH-COUNT"] = strconv.Itoa(value)
	return p
}

func (p *Properties) GetQosPrefetchCount() int {
	value, err := strconv.Atoi(p.Fields["QOS"]["PREFETCH-COUNT"])
	if err != nil {
		return 1
	}
	return value
}

func (p *Properties) SetQosPrefetchSize(value int) *Properties {
	p.initMapper("QOS")
	p.Fields["QOS"]["PREFETCH-SIZE"] = strconv.Itoa(value)
	return p
}

func (p *Properties) GetQosPrefetchSize() int {
	value, err := strconv.Atoi(p.Fields["QOS"]["PREFETCH-SIZE"])
	if err != nil {
		return 0
	}
	return value
}

func (p *Properties) SetQosGlobal(value bool) *Properties {
	p.initMapper("QOS")
	p.Fields["QOS"]["GLOBAL"] = strconv.FormatBool(value)
	return p
}

func (p *Properties) GetQosGlobal() bool {
	if p.Fields["QOS"]["GLOBAL"] == "true" {
		return true
	}
	return false
}