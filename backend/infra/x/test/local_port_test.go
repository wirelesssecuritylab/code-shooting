package test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAvailablePort(t *testing.T) {
	Convey("When get a available port", t, func() {
		port, err := GetAvailablePort()
		Convey("Then success", func() {
			So(port, ShouldNotEqual, 0)
			So(err, ShouldBeNil)
		})

	})
}
