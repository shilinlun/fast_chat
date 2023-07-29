package entity

type Msg struct {
	MsgType     int64  `json:"msgType" bson:"msgType"`
	MsgSendType int64  `json:"msgSendType" bson:"msgSendType"`
	MsgBody     string `json:"msgBody" bson:"msgBody"`
	FromId      int64  `json:"fromId" bson:"fromId"`
	ToId        int64  `json:"toId" bson:"toId"`
}
