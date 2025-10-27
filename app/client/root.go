package client

import (
	"api-gateway-module/common"
	"api-gateway-module/config"
	"api-gateway-module/kafka"
	"sync"

	"github.com/go-resty/resty/v2" // rest builder
)

const (
	_defaultBatchTime = 2
)

type HttpClient struct {
	client *resty.Client
	cfg    config.App

	producer kafka.Producer

	batchTime float64

	// TODO: 아래 데이터는 client에서 구성하면안됨. 추후 redis로 관리해야함.
	fetchLock    sync.Mutex
	mapper       []ApiRequestTopic
	fetchChannel chan ApiRequestTopic
}

func NewHttpClient(
	cfg config.App,
	producer map[string]kafka.Producer,
) *HttpClient {
	batchTime := cfg.Producer.BatchTime

	if batchTime == 0 {
		batchTime = _defaultBatchTime
	}

	if cfg.Http.BaseUrl == "" {
		panic("BaseUrl not existed")
	}

	httpClient := HttpClient{
		cfg:          cfg,
		producer:     producer[cfg.App.Name],
		batchTime:    batchTime,
		mapper:       make([]ApiRequestTopic, 0),
		fetchChannel: make(chan ApiRequestTopic),
	}

	httpClient.client = resty.New().
		SetJSONMarshaler(common.JsonHandler.Marshal).     // sonic Marshal
		SetJSONUnmarshaler(common.JsonHandler.Unmarshal). // sonic Unmarshal
		SetBaseURL(cfg.Http.BaseUrl)

	if len(cfg.Producer.URL) > 0 {
		go func() {
			httpClient.loop()
		}()
	}

	return &httpClient
}

func (h *HttpClient) GET(url string, router config.Router) (interface{}, error) {
	var err error
	var req *resty.Request
	var resp *resty.Response

	common.CB.Execute(func() ([]byte, error) {
		req = getRequest(h.client, router)
		resp, err = req.Get(url)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	// defer 키워드를 사용하여 함수가 종료되기 직전에 실행시킴
	defer h.handleRequestDefer(resp, req.Body)

	if err != nil {
		return nil, err
	}

	return string(resp.Body()), nil
}

func (h *HttpClient) POST(url string, router config.Router, requestBody interface{}) (interface{}, error) {
	var err error
	var req *resty.Request
	var resp *resty.Response

	common.CB.Execute(func() ([]byte, error) {
		req = getRequest(h.client, router)
		resp, err = req.Post(url)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	// defer 키워드를 사용하여 함수가 종료되기 직전에 실행시킴
	defer h.handleRequestDefer(resp, req.Body)

	if err != nil {
		return nil, err
	}

	return string(resp.Body()), nil
}

func (h *HttpClient) PUT(url string, router config.Router, requestBody interface{}) (interface{}, error) {
	var err error
	var req *resty.Request
	var resp *resty.Response

	common.CB.Execute(func() ([]byte, error) {
		req = getRequest(h.client, router)
		resp, err = req.Put(url)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	// defer 키워드를 사용하여 함수가 종료되기 직전에 실행시킴
	defer h.handleRequestDefer(resp, req.Body)

	if err != nil {
		return nil, err
	}

	// defer 키워드를 사용하여 함수가 종료되기 직전에 실행시킴
	defer h.handleRequestDefer(resp, req.Body)

	return string(resp.Body()), nil
}

func (h *HttpClient) DELETE(url string, router config.Router, requestBody interface{}) (interface{}, error) {
	var err error
	var req *resty.Request
	var resp *resty.Response

	common.CB.Execute(func() ([]byte, error) {
		req = getRequest(h.client, router)
		resp, err = req.Delete(url)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	// defer 키워드를 사용하여 함수가 종료되기 직전에 실행시킴
	defer h.handleRequestDefer(resp, req.Body)

	if err != nil {
		return nil, err
	}

	// defer 키워드를 사용하여 함수가 종료되기 직전에 실행시킴
	defer h.handleRequestDefer(resp, req.Body)

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
