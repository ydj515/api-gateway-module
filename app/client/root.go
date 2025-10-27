package client

import (
	"api-gateway-module/common"
	"api-gateway-module/config"
	"api-gateway-module/kafka"
	"fmt"

	"github.com/go-resty/resty/v2" // rest builder
)

const (
	_defaultBatchTime = 2
)

type HttpClient struct {
	client   *resty.Client
	cfg      config.App
	producer kafka.Producer
}

func NewHttpClient(
	cfg config.App,
	producer map[string]kafka.Producer,
) HttpClient {
	batchTime := cfg.Producer.BatchTime

	if batchTime == 0 {
		batchTime = _defaultBatchTime
	}

	if cfg.Http.BaseUrl == "" {
		panic("BaseUrl not existed")
	}

	client := resty.New().
		SetJSONMarshaler(common.JsonHandler.Marshal).     // sonic Marshal
		SetJSONUnmarshaler(common.JsonHandler.Unmarshal). // sonic Unmarshal
		SetBaseURL(cfg.Http.BaseUrl)
	return HttpClient{
		cfg:      cfg,
		client:   client,
		producer: producer[cfg.App.Name],
	}
}

func (h HttpClient) GET(url string, router config.Router) (interface{}, error) {
	var err error
	var req *resty.Request
	var resp *resty.Response

	req = getRequest(h.client, router)
	resp, err = req.Get(url)

	fmt.Println(resp)

	if err != nil {
		return nil, err
	}
	return string(resp.Body()), nil
}

func (h HttpClient) POST(url string, router config.Router, requestBody interface{}) (interface{}, error) {
	var err error
	var req *resty.Request
	var resp *resty.Response

	req = getRequest(h.client, router).SetBody(requestBody)
	resp, err = req.Post(url)

	fmt.Println(resp)

	if err != nil {
		return nil, err
	}
	return string(resp.Body()), nil
}

func (h HttpClient) PUT(url string, router config.Router, requestBody interface{}) (interface{}, error) {
	var err error
	var req *resty.Request
	var resp *resty.Response

	// defer // 함수가 종료할 때 실행

	req = getRequest(h.client, router).SetBody(requestBody)
	resp, err = req.Put(url)

	fmt.Println(resp)

	if err != nil {
		return nil, err
	}
	return string(resp.Body()), nil
}

func (h HttpClient) DELETE(url string, router config.Router, requestBody interface{}) (interface{}, error) {
	var err error
	var req *resty.Request
	var resp *resty.Response

	// defer // 함수가 종료할 때 실행

	req = getRequest(h.client, router).SetBody(requestBody)
	resp, err = req.Delete(url)

	fmt.Println(resp)

	if err != nil {
		return nil, err
	}
	return string(resp.Body()), nil
}

func getRequest(client *resty.Client, router config.Router) *resty.Request {
	// h.client.R().SetAuthScheme().SetAuthToken().SetHeaders()
	req := client.R().EnableTrace()

	if router.Auth != nil {
		if len(router.Auth.Key) != 0 {
			req.SetAuthScheme(router.Auth.Key)
		}
		req.SetAuthScheme(router.Auth.Token)
	}
	if router.Header != nil {
		req.SetHeaders(router.Header)
	}
	return req
}
