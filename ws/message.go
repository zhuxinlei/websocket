package ws

const (
	SubscribeStatusOk     = "ok"
	SubscribeStatusFailed = "failed"
)

// 单人消息
type SingleMessage struct {
	Target  ClientKey
	Content []byte
}

// 主题消息
type TopicMessage struct {
	Topic   Topic
	Targets []ClientKey
	Content []byte
}

// 订阅请求
type SubscribeRequest struct {
	ID    string `json:"id"`
	Sub   Topic  `json:"sub"`
	Unsub Topic  `json:"unsub"`
}

// 订阅响应
type SubscribeResponse struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Subbed   Topic  `json:"subbed,omitempty"`
	Unsubbed Topic  `json:"unsubbed,omitempty"`
	Ts       int64  `json:"ts"`                  // millisecond timestamp
	ErrorMsg string `json:"error_msg,omitempty"` // 错误信息
}
