package entity

type Msg struct {
	MsgType     int64  `json:"msgType" bson:"msgType"`
	MsgSendType int64  `json:"msgSendType" bson:"msgSendType"`
	MsgBody     string `json:"msgBody" bson:"msgBody"`
	FromId      string `json:"fromId" bson:"fromId"`
	ToId        string `json:"toId" bson:"toId"`
}

/*

{"msgType":1,"msgSendType":1,"msgBody":"注册“,"fromId":"123"}

{"msgType":1,"msgSendType":1,"msgBody":"注册“,"fromId":"234"}

{"msgType":1,"msgSendType":1,"msgBody":"注册“,"fromId":"345"}

{"msgType":3,"msgSendType":1,"msgBody":"123->234“,"fromId":"123","toId":"234"}

{"msgType":4,"msgSendType":1,"msgBody":"123->all“,"fromId":123}

*/
