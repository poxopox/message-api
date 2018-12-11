package message

type Messages []*Message

type MessageStore interface {
	PutMessage(message *Message) error
	ReadMessages(user string) ( messages Messages, err error)
}

type InMemoryMessageStore struct {
	messages Messages
}
func (imds *InMemoryMessageStore) PutMessage(message *Message) error {
	imds.messages = append(imds.messages, message)
	return nil
}
func (imds *InMemoryMessageStore) ReadMessages(user string) ( messages Messages, err error ) {
	messages = imds.messages
	return
}

func NewInMemoryDataStore () *InMemoryMessageStore {
	return &InMemoryMessageStore{
		messages: make([]*Message, 0),
	}
}