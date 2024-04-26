package chatengine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// API接口文档: https://cloud.baidu.com/doc/WENXINWORKSHOP/s/ilv62om62
type Engine struct {
	cfg        Config
	httpClient *http.Client
}

func NewEngine(cfg Config) *Engine {
	return &Engine{
		cfg:        cfg,
		httpClient: &http.Client{},
	}
}

// 流式输出
func (e *Engine) StreamRequest(ctx context.Context, req ChatReq) (chan *ChatResponse, error) {
	req.Stream = true
	httpResp, err := e.prepare(ctx, &req)
	if err != nil {
		return nil, err
	}

	code := httpResp.StatusCode
	if code != http.StatusOK {
		defer httpResp.Body.Close()
		bytes, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, err
		}
		log.Println("error: http status code:", code, ", body: "+string(bytes))
		if code == http.StatusTooManyRequests {
			return nil, ErrRateLimit
		}
		return nil, fmt.Errorf("invalid http status %v", code)
	}

	// httpResp.Body will be closed by Stream
	ch := make(chan *ChatResponse)
	stream := NewStream(ch, httpResp.Body)
	go stream.Recv()
	return ch, nil
}

func (e *Engine) ChatRequest(ctx context.Context, req ChatReq) (*ChatResponse, error) {
	req.Stream = false
	httpResp, err := e.prepare(ctx, &req)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	code := httpResp.StatusCode
	if code != http.StatusOK {
		if code == http.StatusTooManyRequests {
			return nil, ErrRateLimit
		}
		return nil, fmt.Errorf("invalid http status %v", code)
	}

	bytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp ChatResponse
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	if resp.ErrCode != 0 {
		if resp.ErrCode == 18 {
			return nil, ErrRateLimit
		} else {
			return nil, fmt.Errorf(resp.ErrMsg)
		}
	}

	return &resp, nil
}

func (e *Engine) prepare(ctx context.Context, req *ChatReq) (*http.Response, error) {
	jsonValue, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	urlStr := "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/" + strings.ToLower(req.Model)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, urlStr, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(urlStr)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Host", u.Hostname())
	request.Header.Set("x-bce-date", formatISO8601Date(nowUTCSeconds()))

	headers := make(map[string]string)
	for k, v := range request.Header {
		headers[k] = v[0]
	}
	headersToSign := make(map[string]struct{})
	for k := range headers {
		headersToSign[strings.ToLower(k)] = struct{}{}
	}
	signStr := bceSign(e.cfg.AccessKey, e.cfg.SecretKey, request.URL.RequestURI(), request.Method, nil, headers, headersToSign)
	request.Header.Set(AUTHORIZATION, signStr)

	return e.httpClient.Do(request)
}
