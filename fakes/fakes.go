package fakes

// FakeQueueDispatcher provides a string mapping so we can assert on which Messages
// have been dispatched by calling clients.
type FakeQueueDispatcher struct {
	Messages      map[string][]interface{}
	DispatchCount int
}

// NewFakeQueueDispatcher creates a fake dispatcher
func NewFakeQueueDispatcher() (dispatcher *FakeQueueDispatcher) {
	dispatcher = &FakeQueueDispatcher{}
	dispatcher.Messages = make(map[string][]interface{})
	dispatcher.Messages["telemetry"] = make([]interface{}, 0)
	dispatcher.Messages["alert"] = make([]interface{}, 0)
	dispatcher.Messages["position"] = make([]interface{}, 0)
	return
}

// DispatchMessage implementation of dispatch message interface method
func (q *FakeQueueDispatcher) DispatchMessage(queue string, message interface{}) (err error) {
	q.DispatchCount++
	q.Messages[queue] = append(q.Messages[queue], message)
	return
}
