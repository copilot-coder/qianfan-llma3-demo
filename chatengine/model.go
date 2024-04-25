package chatengine

import "errors"

var (
	ErrNoResponse = errors.New("no response")
	ErrRateLimit  = errors.New("rate limit")
)

type Config struct {
	AccessKey string
	SecretKey string
	Model     string
}

type ChatReq struct {
	Messages        []Message `json:"messages"`
	MaxOutputTokens int       `json:"max_output_tokens"`
	Stream          bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	ErrCode int    `json:"error_code"`
	ErrMsg  string `json:"error_msg"`
	Id      string `json:"id"`
	IsEnd   bool   `json:"is_end"`
	Result  string `json:"result"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Err error `json:"-"`
}
