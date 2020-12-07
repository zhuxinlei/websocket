package test

import (
	"encoding/json"
	"testing"
)

func TestTopicPush(t *testing.T) {
	TopicPush("otc", nil, WsPushData{
		Type: "test",
		Data: json.RawMessage("1231"),
	})
}
