package model

type Message struct {
	RoutingKey 	string
	UserID 		string 	`json:"userId"`
	Data 		[]byte 	`json:"data"`
}
