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

	if mq.data[topic] != nil {
		if len(mq.data[topic]) == 0 {
			return nil
		}
		data = mq.data[topic][0].Body
		mq.Lock()
		mq.data[topic] = mq.data[topic][1:]
		mq.Unlock()
		return &data
	}

	return nil
}

func (mq *SimpleMessageQueueRepository) Clear(topic string) {
	if mq.data[topic] != nil {
		if len(mq.data[topic]) != 0 {
			mq.data[topic] = nil
		}
	}
}

func (mq *SimpleMessageQueueRepository) PeekAll(topic string) *[]byte {
	var res []byte

	if mq.data[topic] != nil {
		if len(mq.data[topic]) == 0 {
			return nil
		}
		mq.Lock()
		for _, data := range mq.data[topic] {
			res = append(res, data.Body...)
		}
		mq.Unlock()
		return &res
	}
	return nil
}

func (mq *SimpleMessageQueueRepository) Peek(topic string) *[]byte {
	var data []byte

	if mq.data[topic] != nil {
		if len(mq.data[topic]) == 0 {
			return nil
		}
		mq.Lock()
		data = mq.data[topic][0].Body
		mq.Unlock()
		return &data
	}

	return nil
}

func (mq *SimpleMessageQueueRepository) Push(topic string, obj []byte) {
	mq.Lock()
	mq.data[topic] = append(mq.data[topic], entity.Response{Body: obj})
	mq.Unlock()
}

func (mq *SimpleMessageQueueRepository) GetTopics() []string {
	var topics []string
	for key, _ := range mq.data {
		topics = append(topics, key)
	}
	return topics
}
