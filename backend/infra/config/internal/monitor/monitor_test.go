package monitor

import (
	"testing"

	"code-shooting/infra/config/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMonitor(t *testing.T) {

	Convey("Given a config monitor ", t, func() {
		monitor := NewConfigstoreMonitor()
		receiveEvent := false
		monitor.RegisterEventHandler("a.b", func([]*model.Event) {
			receiveEvent = true
		})
		Convey("When process event key containt a.b \n", func() {
			monitor.ProcessConfigEvent([]*model.Event{
				{
					ConfigItem: model.ConfigItem{Key: "a.b.c"},
				},
			})
			Convey("Then process event . \n", func() {
				So(receiveEvent, ShouldBeTrue)
			})
		})

		Convey("When process event key not containt a.b \n", func() {
			monitor.ProcessConfigEvent([]*model.Event{
				{
					ConfigItem: model.ConfigItem{Key: "a.d.c"},
				},
			})
			Convey("Then not process event . \n", func() {
				So(receiveEvent, ShouldBeFalse)
			})
		})
	})
}
