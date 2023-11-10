package messagequeue

type IMessageQueueRepository interface {
	Pop(string) *[]byte
	Clear(string)
	PeekAll(string) *[]byte
	Peek(string) *[]byte
	Push(string, []byte)
}
