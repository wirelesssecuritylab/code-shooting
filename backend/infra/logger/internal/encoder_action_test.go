package internal

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEncoderActionRegister(t *testing.T) {
	Convey("Given a encode actions map", t, func() {
		encodes := map[string]func() string{
			"TC": func() string {
				return "trace_id_1"
			},
			"SP": func() string {
				return "span_id_1"
			},
		}
		Convey("When register encoder actions", func() {
			err := RegisterEncoderActions(encodes)

			Convey("Then register succ and resitered encodes can be queryed", func() {
				So(err, ShouldBeNil)

				_, ok := encodeActions["TC"]
				So(ok, ShouldBeTrue)

				_, ok2 := encodeActions["SP"]
				So(ok2, ShouldBeTrue)
			})
		})
	})
}

func TestEncoderActionRegisterWhichEncoderKeyExist(t *testing.T) {
	Convey("Given a encode actions map with encode key which already exist in encodeActions", t, func() {
		encodes := map[string]func() string{
			"T": func() string {
				return "fake-timestamp"
			},
		}
		Convey("When register encoder actions", func() {
			err := RegisterEncoderActions(encodes)

			Convey("Then register failed,err should tips: key: T already exist.", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "key: T already exist.")
			})
		})
	})
}

func TestEncoderActionRegisterWhichSomeEncoderRegisterFailed(t *testing.T) {
	Convey("Given a encode actions map which one key already exist in encodeActions and the others are normal", t, func() {
		encodes := map[string]func() string{
			"T": func() string {
				return "fake-timestamp"
			},
			"TCID": func() string {
				return "trace_id_1"
			},
			"SPID": func() string {
				return "span_id_1"
			},
		}
		Convey("When register encoder actions", func() {
			err := RegisterEncoderActions(encodes)

			Convey("Then register failed and given encodes should not be found in encodeActions", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "key: T already exist.")

				_, ok := encodeActions["TCID"]
				So(ok, ShouldBeFalse)

				_, ok2 := encodeActions["SPID"]
				So(ok2, ShouldBeFalse)
			})
		})
	})
}

func TestEncoderActionRegisterWhichEncoderKeyIsEmpty(t *testing.T) {
	Convey("Given a encode actions map with empty encode key", t, func() {
		encodes := map[string]func() string{
			"": func() string {
				return "error"
			},
		}
		Convey("When register encoder actions", func() {
			err := RegisterEncoderActions(encodes)

			Convey("Then register failed,err should tips: encodeActions has empty key.", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "encodeActions has empty key.")
			})
		})
	})
}
