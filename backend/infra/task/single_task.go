package task

import (
	"sync"

	"code-shooting/infra/logger"
)

type SingleTask struct {
	runningTaskTag sync.Map
	waitingTaskTag sync.Map
	finishTask     chan string
	stopTask       chan bool
	taskFunc       func(string, func())
}

func NewSingleTask(f func(string, func())) (*SingleTask, error) {
	t := &SingleTask{
		finishTask: make(chan string),
		stopTask:   make(chan bool),
		taskFunc:   f,
	}
	go t.loopTask()
	return t, nil
}

func (t *SingleTask) SubmitTask(taskName string) {
	if _, ok := t.waitingTaskTag.Load(taskName); ok {
		return
	}

	if _, ok := t.runningTaskTag.Load(taskName); ok {
		t.waitingTaskTag.Store(taskName, "")
		return
	}

	t.runningTaskTag.Store(taskName, "")
	t.runningTask(taskName)
	return
}

func (t *SingleTask) ShowTask() {
	t.runningTaskTag.Range(func(key, value interface{}) bool {
		logger.Debugf("running task : %s", key)
		return true
	})
	t.waitingTaskTag.Range(func(key, value interface{}) bool {
		logger.Debugf("waiting task : %s", key)
		return true
	})
}

func (t *SingleTask) finishTaskNotice(param string) func() {
	return func() {
		t.finishTask <- param
		logger.Debugf("start finish notice task:%s", param)
	}
}

func (t *SingleTask) ReleaseTask() {
	t.stopTask <- true
}

func (t *SingleTask) runningTask(param string) {
	go func() {
		t.taskFunc(param, t.finishTaskNotice(param))
	}()
}

func (t *SingleTask) loopTask() {
	for {
		select {
		case taskName := <-t.finishTask:
			logger.Debugf("deal finish notice task:%s", taskName)
			t.runningTaskTag.Delete(taskName)
			if _, ok := t.waitingTaskTag.LoadAndDelete(taskName); ok {
				t.runningTask(taskName)
			}
		case <-t.stopTask:
			break
		}
	}
}
