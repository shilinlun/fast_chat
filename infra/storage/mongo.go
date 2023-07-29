package storage

import (
	"context"
	"errors"
	"fast_chat/core/entity"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	session *mongo.Database
}

func ProvideMongo() *Mongo {
	// 设置客户端选项
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// 连接 MongoDB
	c, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	session := c.Database("fast_chat")
	return &Mongo{
		session: session,
	}
}

func (m *Mongo) Insert(ctx context.Context, msg *entity.Msg) error {
	conn := m.session.Collection("msg")
	if conn == nil {
		return errors.New("get mongo session failed")
	}
	resp, err := conn.InsertOne(ctx, msg)
	if err != nil {
		return err
	}
	if resp.InsertedID == nil {
		return errors.New("insert msg err")
	}
	return nil
}
