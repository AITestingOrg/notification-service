package model

import "encoding/json"

type Message struct {
	RoutingKey string
	Body       map[string]*json.RawMessage
}
