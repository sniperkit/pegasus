package dummies

import (
	"bitbucket.org/code_horse/pegasus/transport"
	"bitbucket.org/code_horse/pegasus/transport/http_transport"
	"net/url"
)

func SendLogs(content string) {

	httpTransporter := transport.NewHttpTransporter(nil)

	params := url.Values{}
	params.Add("content", content)

	_, err := httpTransporter.Send(
		http_transport.NewProperties().
			SetPath("http://http_log_service:8800/log?content="+params.Encode()).
			SetGetMethod().
			GetProperties(),
		http_transport.NewOptions().
			GetOptions(),
		nil,
	)

	if err != nil {
		panic(err)
	}

}
