package fakes

// FakeQueueDispatcher provides a string mapping so we can assert on which Messages
// have been dispatched by calling clients.
type FakeQueueDispatcher struct {
	Messages      map[string][]interface{}
	DispatchCount int
}

// DispatchMessage implementation of dispatch message interface method
func (q *FakeQueueDispatcher) DispatchMessage(queue string, message interface{}) (err error) {
	q.DispatchCount++
	return
}
