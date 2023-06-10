package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConvertConfigKey(t *testing.T) {

	Convey("Given a config key has a..b \n", t, func() {

		Convey("When convert key \n", func() {
			covertKey := ConvertOutKeyToInner("a....b")
			Convey("Then k covert result k a,,b \n", func() {
				So(covertKey, ShouldEqual, "a..b")

			})
		})
	})

	Convey("Given a config key has a..b#c \n", t, func() {

		Convey("When convert key \n", func() {
			covertKey := ConvertInnerKeyToOut("a..b#c")
			Convey("Then k covert result k a,,b \n", func() {
				So(covertKey, ShouldEqual, "a..b.c")

			})
		})
	})
}
