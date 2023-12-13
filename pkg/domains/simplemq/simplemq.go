package simplemq

import (
	"smq/pkg/entity"
	"sync"
)

type SimpleMessageQueueRepository struct {
	data map[string][]entity.Response
	sync.Mutex
}

func NewSimpleMessageQueueRepository() *SimpleMessageQueueRepository {
	m := make(map[string][]entity.Response)
	mq := SimpleMessageQueueRepository{data: m}
	return &mq
}

func (mq *SimpleMessageQueueRepository) Pop(topic string) *[]byte {
	var data []byte

	mq.Lock()
	if mq.data[topic] != nil {
		if len(mq.data[topic]) == 0 {
			return nil
		}
		data = mq.data[topic][0].Body
		mq.data[topic] = mq.data[topic][1:]
		if len(mq.data[topic]) == 0 {
			delete(mq.data, topic)
		}
		return &data
	}
	mq.Mutex.Unlock()
	return nil

}

func (mq *SimpleMessageQueueRepository) Clear(topic string) {
	mq.Lock()
	if mq.data[topic] != nil {
		if len(mq.data[topic]) != 0 {
			mq.data[topic] = nil
		}
	}
	mq.Unlock()
}

func (mq *SimpleMessageQueueRepository) PeekAll(topic string) *[]byte {
	var res []byte

	mq.Lock()
	if mq.data[topic] != nil {
		if len(mq.data[topic]) == 0 {
			return nil
		}
		for _, data := range mq.data[topic] {
			res = append(res, data.Body...)
		}
		return &res
	}
	mq.Unlock()
	return nil
}

func (mq *SimpleMessageQueueRepository) Peek(topic string) *[]byte {
	var data []byte
	mq.Lock()
	if mq.data[topic] != nil {
		if len(mq.data[topic]) == 0 {
			return nil
		}
		data = mq.data[topic][0].Body
		return &data
	}
	mq.Unlock()
	return nil
}

func (mq *SimpleMessageQueueRepository) Push(topic string, obj []byte) {
	mq.Lock()
	mq.data[topic] = append(mq.data[topic], entity.Response{Body: obj})
	mq.Unlock()
}

func (mq *SimpleMessageQueueRepository) GetTopics() []string {
	var topics []string
	for key := range mq.data {
		topics = append(topics, key)
	}
	return topics
}
