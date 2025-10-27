package client

import (
	"api-gateway-module/common"
	"time"

	"github.com/go-resty/resty/v2"
)

// 중복 발행을 방지하기 위해 lock
func (h *HttpClient) fetchToKafka() {
	h.fetchLock.Lock()
	defer h.fetchLock.Unlock()

	if len(h.mapper) > 0 {
		ent := h.mapper
		h.mapper = make([]ApiRequestTopic, 0)

		v, err := common.JsonHandler.Marshal(ent)

		if err == nil {
			h.producer.SendEvent(v)
		}

	}
}

func (h *HttpClient) loop() {
	ticker := time.NewTicker(time.Duration(h.batchTime) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.fetchToKafka()
		}
	}

}

func (h *HttpClient) handleRequestDefer(resp *resty.Response, request interface{}) {
	if len(h.cfg.Producer.URL) > 0 {
		h.mapper = append(h.mapper, NewApiRequestTopic(resp, request))
	}
}
