package redis

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/poxopox/text-buffer/message"
)

type MemoryStore struct {
	RedisClient *redis.Client
}

func NewMessageStore(redisClient *redis.Client) *MemoryStore {
	return &MemoryStore{RedisClient: redisClient}
}

func (mem *MemoryStore) PutMessage(message *message.Message) error {
	messageJson, err :=  message.ToJson()
	if err != nil {
		return err
	}
	mem.RedisClient.LPush("messages", messageJson)
	return nil
}
func (mem *MemoryStore) ReadMessages (user string) ( messages message.Messages, err error) {
	messages = make(message.Messages, 0)
	results := mem.RedisClient.LRange("messages", 0, -1)
	for _, resJson := range results.Val() {
		marshaledRes := &message.Message{}
		err := json.Unmarshal([]byte(resJson), marshaledRes)
		if err != nil {
			return nil, err
		}
		messages = append(messages, marshaledRes)
	}
	return messages, nil
}