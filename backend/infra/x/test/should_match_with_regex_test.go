package test

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestShouldRegexMatchWith(t *testing.T) {
	Convey("Given a regex expression: ^abc[\\d]+$", t, func() {
		regex := "^abc[\\d]+$"

		Convey("When matched it with abc1234", func() {
			res := ShouldMatchWithRegex("abc1234", regex)

			Convey("Then result is match", func() {
				So(res, ShouldBeEmpty)
			})
		})

		Convey("When matched it with xxx", func() {
			res := ShouldMatchWithRegex("xxx", regex)

			Convey("Then should tips: but it didn't", func() {
				So(res, ShouldContainSubstring, "but it didn't")
			})
		})

		Convey("When matched it with integer: 10", func() {
			res := ShouldMatchWithRegex(10, regex)

			Convey("Then should tips: Actual argument to this assertion must be string", func() {
				So(res, ShouldContainSubstring, "Actual argument to this assertion must be string")
			})
		})
	})

	Convey("Given a string abc1234", t, func() {
		str := "abc1234"

		Convey("When match it with regex: a, abc, ^abc", func() {
			res := ShouldMatchWithRegex(str, "a", "abc", "^abc")

			Convey("Then result is match", func() {
				So(res, ShouldBeEmpty)
			})
		})

		Convey("When match it with regex: a, 1", func() {
			res := ShouldMatchWithRegex(str, "a", 1)

			Convey("Then should tips: Expected argument to this assertion must be string", func() {
				So(res, ShouldContainSubstring, "Expected argument to this assertion must be string")
			})
		})

		Convey("When match it with nothing", func() {
			res := ShouldMatchWithRegex(str)

			Convey("Then should tips: Expected value to this assertion must be specified", func() {
				So(res, ShouldContainSubstring, "Expected value to this assertion must be specified")
			})
		})
	})
}
