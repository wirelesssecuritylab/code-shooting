package internal

import (
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
	"testing"
)

func TestConfig(t *testing.T) {
	Convey("Given a info level", t, func() {
		level := InfoLevel
		Convey("When check equal with info level", func() {
			equal := level.Equal(Level("info"))
			Convey("Then is equal", func() {
				So(equal, ShouldBeTrue)
			})
		})

		Convey("When check with invalid level", func() {
			equal := level.Equal(Level("xxx"))
			Convey("Then is not equal", func() {
				So(equal, ShouldBeFalse)
			})
		})

		Convey("When transform to zap level", func() {
			zapLevel, err := level.ToZapLevel()
			Convey("Then the level should be", func() {
				So(err, ShouldBeNil)
				So(zapLevel, ShouldEqual, zap.InfoLevel)
			})
		})

		Convey("When transform to string", func() {
			str := level.String()
			Convey("Then result is info", func() {
				So(str, ShouldEqual, "info")
			})
		})
	})

	Convey("Given a invalid level", t, func() {
		level := Level("invalid")

		Convey("When check equal with info level", func() {
			equal := level.Equal(InfoLevel)
			Convey("Then is not equal", func() {
				So(equal, ShouldBeFalse)
			})
		})

		Convey("When check with invalid level", func() {
			equal := level.Equal(Level("xxx"))
			Convey("Then is not equal", func() {
				So(equal, ShouldBeFalse)
			})
		})

		Convey("When transform to zap level", func() {
			_, err := level.ToZapLevel()
			Convey("Then should tips: incompatible with zap level", func() {
				So(err.Error(), ShouldContainSubstring, "incompatible with zap level")
			})
		})

		Convey("When transform to string", func() {
			str := level.String()
			Convey("Then result should be Level(invalid)", func() {
				So(str, ShouldEqual, "Level(invalid)")
			})
		})
	})

	Convey("Given a json encoder", t, func() {
		encoder := JsonEncoder
		Convey("When check equal json encoder", func() {
			equal := encoder.Equal(Encoder("JSON"))
			Convey("Then is equal", func() {
				So(equal, ShouldBeTrue)
			})
		})

		Convey("When check with invalid encoder", func() {
			equal := encoder.Equal(Encoder("xxx"))
			Convey("Then is not equal", func() {
				So(equal, ShouldBeFalse)
			})
		})

		Convey("When transform to string", func() {
			str := encoder.String()
			Convey("Then result is info", func() {
				So(str, ShouldEqual, "json")
			})
		})
	})

	Convey("Given a invalid encoder", t, func() {
		encoder := Encoder("invalid")
		Convey("When check equal json encoder", func() {
			equal := encoder.Equal(Encoder("JSON"))
			Convey("Then is not equal", func() {
				So(equal, ShouldBeFalse)
			})
		})

		Convey("When check with invalid encoder", func() {
			equal := encoder.Equal(Encoder("xxx"))
			Convey("Then is not equal", func() {
				So(equal, ShouldBeFalse)
			})
		})

		Convey("When transform to string", func() {
			str := encoder.String()
			Convey("Then result is info", func() {
				So(str, ShouldEqual, "Encoder(invalid)")
			})
		})
	})
}
