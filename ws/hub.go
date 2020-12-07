package ws

import (
	log "github.com/sirupsen/logrus"
)

var (
	hubManager *hub
)

func init() {
	hubManager = newhub()
}

// 订阅信息
type subscription struct {
	client *Client
	topic  Topic
	done   chan struct{}
}

type hub struct {
	// Registered clients.
	clients map[ClientKey]*Client
	// Register requests from the clients.
	register chan *Client
	// Unregister requests from clients.
	unregister chan *Client

	// sub
	sub chan subscription
	// unsub
	unsub chan subscription
	// topics
	topics map[Topic]map[ClientKey]struct{}

	// SinglePush is a channel for pushing info to a specify client
	SinglePush chan *SingleMessage
	// BroadcastPush
	BroadcastPush chan []byte
	// TopicPush
	TopicPush chan *TopicMessage
}

func GetHubManager() *hub {
	return hubManager
}

func newhub() *hub {
	return &hub{
		clients:       make(map[ClientKey]*Client),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		sub:           make(chan subscription),
		unsub:         make(chan subscription),
		topics:        make(map[Topic]map[ClientKey]struct{}),
		SinglePush:    make(chan *SingleMessage, 1024),
		BroadcastPush: make(chan []byte, 1024),
		TopicPush:     make(chan *TopicMessage, 1024),
	}
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.key] = client
			log.Debugf("[hub] Register client: %+v", client)
		case client := <-h.unregister:
			//断开websocket连接或者取消订阅的时候查看主题下是否还有订阅用户，如果没有则删除内存中的该主题相关数据

			if _, ok := h.clients[client.key]; ok {
				delete(h.clients, client.key)
				close(client.send)
			}
			//判断所有主题下是否还有订阅用户，将没有订阅用户的主题删除
			for k, v := range h.topics {
				for clientKey, _ := range v {
					if clientKey == client.key {
						delete(h.topics[k], clientKey)
					}
				}
				len := len(v)

				if len == 0 {
					v = nil
				}
			}
			log.Debugf("[hub] Unregister client: %+v", client)
		case sub := <-h.sub:
			clients, ok := h.topics[sub.topic]
			if !ok {
				clients = make(map[ClientKey]struct{})
				h.topics[sub.topic] = clients
			}
			clients[sub.client.key] = struct{}{}
			close(sub.done)

			log.Debugf("[hub] sub topic: %s, client: %+v", sub.topic, sub.client.key)

		case unsub := <-h.unsub:
			//断开websocket连接或者取消订阅的时候查看主题下是否还有订阅用户，如果没有则删除内存中的该主题相关数据
			clients, ok := h.topics[unsub.topic]
			if ok {
				delete(clients, unsub.client.key)
			}
			close(unsub.done)
			//取消时判断如果该topic下没有关注的用户，则删除该topic
			len := len(h.topics[unsub.topic])

			if len == 0 {
				clients = nil
			}
			log.Debugf("[hub] unsub topic: %s, client: %+v", unsub.topic, unsub.client.key)

		case singleMessage := <-h.SinglePush:
			client, ok := h.clients[singleMessage.Target]
			if ok {
				client.send <- singleMessage.Content
				log.Debugf("[hub] SinglePush client: %+v, msg: %s", client, string(singleMessage.Content))
			}
		case message := <-h.BroadcastPush:
			for _, client := range h.clients {
				client.send <- message
			}
			log.Debugf("[hub] BroadcastPush msg: %s", string(message))
		case topicMessage := <-h.TopicPush:
			subscribers, ok := h.topics[topicMessage.Topic]
			if ok {
				if len(topicMessage.Targets) > 0 { // push specific targets
					for _, target := range topicMessage.Targets {
						_, ok := subscribers[target]
						if ok {
							client, ok := h.clients[target]
							if ok {
								client.send <- topicMessage.Content
							} else { // disconnected without unsub
								delete(subscribers, target)
							}
						}
					}
				} else { // push all
					for key := range subscribers {
						client, ok := h.clients[key]
						if ok {
							client.send <- topicMessage.Content
						} else { // disconnected without unsub
							delete(subscribers, key)
						}
					}
				}
			}
			log.Debugf("[hub] TopicPush topic: %s msg: %s", topicMessage.Topic, string(topicMessage.Content))
		}
	}
}
