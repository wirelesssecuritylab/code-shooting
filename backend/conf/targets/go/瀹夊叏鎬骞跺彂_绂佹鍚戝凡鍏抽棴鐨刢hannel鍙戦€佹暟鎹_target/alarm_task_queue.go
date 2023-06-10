package taskqueue

func (s *TaskQueue) Delete(queueIndex int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	index := 0
	for index < queueIndex {
		s.queue[index].Response <- nil
		index = index + 1
	}
	s.queue = s.queue[queueIndex:]
}

func (s *TaskQueue) PushBack(msgSendAndResponse model.MsgSendAndResponse) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.queue = append(s.queue, msgSendAndResponse)
	return nil
}