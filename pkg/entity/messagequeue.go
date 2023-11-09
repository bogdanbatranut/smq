package entity

import "sync"

type MessageQueue struct {
	data map[string][]Response
	sync.Mutex
}

func NewMessageQueue() *MessageQueue {
	m := make(map[string][]Response)
	mq := MessageQueue{data: m}
	return &mq
}

func (mq *MessageQueue) Pop(topic string) *[]byte {
	var data []byte

	if mq.data[topic] != nil {
		if len(mq.data[topic]) == 0 {
			return nil
		}
		data = mq.data[topic][0].body
		mq.Lock()
		mq.data[topic] = mq.data[topic][1:]
		mq.Unlock()
		return &data
	}

	return nil
}

func (mq *MessageQueue) Clear(topic string) {
	if mq.data[topic] != nil {
		if len(mq.data[topic]) != 0 {
			mq.data[topic] = nil
		}
	}
}

func (mq *MessageQueue) PeekAll(topic string) *[]byte {
	var res []byte

	if mq.data[topic] != nil {
		if len(mq.data[topic]) == 0 {
			return nil
		}
		mq.Lock()
		for _, data := range mq.data[topic] {
			res = append(res, data.body...)
		}
		mq.Unlock()
		return &res
	}
	return nil
}

func (mq *MessageQueue) Peek(topic string) *[]byte {
	var data []byte

	if mq.data[topic] != nil {
		if len(mq.data[topic]) == 0 {
			return nil
		}
		mq.Lock()
		data = mq.data[topic][0].body
		mq.Unlock()
		return &data
	}

	return nil
}

func (mq *MessageQueue) Push(topic string, obj []byte) {
	mq.Lock()
	mq.data[topic] = append(mq.data[topic], Response{body: obj})
	mq.Unlock()
}
