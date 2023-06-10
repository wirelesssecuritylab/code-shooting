package alarmproc

import (
	"fm-mgt/constant"
	"fm-mgt/model"
	"fm-mgt/utils"
	"sqm/library/opssdk/log/plog"

	"encoding/json"

	"net/http"
	"time"
)

func (s *EventController) ReceiveAlarm() {
	receivedAlarm := model.SouthAlarm{}
	err := json.Unmarshal(s.Ctx.Input.RequestBody, &receivedAlarm)
	if err != nil {
		plog.Error(constant.Module, "", constant.PROCESS_NAME, "unmarshal request body err: ",
			utils.BytesToString(s.Ctx.Input.RequestBody))
		s.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	msgSendAndResponse := model.MsgSendAndResponse{
		SouthAlarm: &receivedAlarm,
		Response:   make(chan error, 1),
	}
	err := TaskQueue.TaskQueueInstance().PushBack(msgSendAndResponse)
	if err != nil {
		plog.Error(constant.Module, "", constant.PROCESS_NAME, "receivedAlarm alarm queue push error",
			"alarm is ", constant.GetStructMsgInfo(receivedAlarm))
		s.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	timerTimeOut := time.NewTimer(time.Duration(constant.SendInfoTimeOutTime) * time.Minute)
	select {
	case <-msgSendAndResponse.Response:
		s.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
	case <-timerTimeOut.C:
		plog.Error(constant.Module, "", constant.PROCESS_NAME, "receivedAlarm handle the alarm time out",
			"alarm is ", constant.GetStructMsgInfo(receivedAlarm))
		s.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
	}
	close(msgSendAndResponse.Response)
	return
}
