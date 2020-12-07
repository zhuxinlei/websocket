package ws

import (
	"context"
	"server-tokenhouse-ws/config"
	"server-tokenhouse-ws/models"
	"server-tokenhouse-ws/service/api"
	"strings"
	"time"
)

type Topic string

var ValidTopicTrie *models.Trie

func InitValidTopicTrie() {
	ValidTopicTrie = models.NewTrie("")
	for _, topic := range config.GetServer().Topic.Valid {
		if strings.Contains(topic, "$") {
			topics := replaceWildcard(topic)
			for _, topic := range topics {
				ValidTopicTrie.Add(topic)
			}
		} else {
			ValidTopicTrie.Add(topic)
		}
	}
}

func replaceWildcard(str string) []string {
	arr := strings.Split(str, models.TrieSeparator)
	var variable string
	for _, val := range arr {
		if strings.HasPrefix(val, "$") {
			variable = strings.TrimPrefix(val, "$")
			break
		}
	}
	values := config.GetServer().Topic.Dict[variable]
	results := make([]string, 0, len(values))
	for _, value := range values {
		newStr := strings.ReplaceAll(str, "$"+variable, value)
		if strings.Contains(newStr, "$") {
			results = append(results, replaceWildcard(newStr)...)
		} else {
			results = append(results, newStr)
		}
	}
	return results
}

// topic must format like 'prefix.suffix.extra'. e.g. 'eth.kline.1', 'eth.order'
func (t Topic) IsValid() bool {
	return ValidTopicTrie.Exist(string(t))
}
func (t Topic) IsDynamic() bool {
	id, ok := ValidTopicTrie.IsDynamic(string(t))
	if ok == false {
		return false
	}

	//验证house_id是否合法
	cfg := config.GetServer()
	host := cfg.TokenHouseHost + "/house_check?house_id=" + id.(string)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	res, err := api.DoRequest(ctx, "GET", host, nil, nil)
	if err != nil {
		return false
	}
	if res.Code != 200 {
		return false
	}
	return true
}
