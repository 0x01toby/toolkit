package transport

import (
	"encoding/json"
	"fmt"
	"github.com/taorzhang/toolkit/client/jsonrpc/codec"
	"github.com/valyala/fasthttp"
)

type Http struct {
	addr    string
	headers map[string]string
	client  *fasthttp.Client
}

func NewHttp(addr string) *Http {
	return &Http{addr: addr, client: &fasthttp.Client{}, headers: make(map[string]string)}
}

// SetHeaders 设置header头
func (h *Http) SetHeaders(headers map[string]string) {
	if len(headers) > 0 {
		h.headers = headers
	}
}

// SetMaxConnPerHost 设置最大连接
func (h *Http) SetMaxConnPerHost(count int) {
	h.client.MaxConnsPerHost = count
}

// Close 关闭连接
func (h *Http) Close() error {
	return nil
}

// Call standard eth method call
func (h *Http) Call(method string, out interface{}, params ...interface{}) error {
	request := codec.Request{
		ID:      1,
		JsonRPC: "2.0",
		Method:  method,
	}
	if len(params) > 0 {
		data, err := json.Marshal(params)
		if err != nil {
			return fmt.Errorf("marshal request params:%v", err)
		}
		request.Params = data
	}
	raw, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal request:%v", err)
	}
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	req.SetRequestURI(h.addr)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	for k, v := range h.headers {
		req.Header.Add(k, v)
	}
	req.SetBody(raw)
	if err = h.client.Do(req, resp); err != nil {
		return fmt.Errorf("clinet do:%v", err)
	}
	var response codec.Response
	if err = json.Unmarshal(resp.Body(), &response); err != nil {
		return fmt.Errorf("json unmarshal response.body:%v", err)
	}
	if response.Error != nil {
		return fmt.Errorf("response error:%v", response.Error)
	}

	if err = json.Unmarshal(response.Result, out); err != nil {
		return fmt.Errorf("json unmarshal response.result:%v", err)
	}
	return nil
}

// EthCall 单独一个eth call
func (h *Http) EthCall() error {
	return nil
}

// BatchEthCall 批量执行 eth call
func (h *Http) BatchEthCall() error {
	return nil
}
