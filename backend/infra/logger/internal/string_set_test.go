package internal

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestStringSet(t *testing.T) {
	Convey("Given a empty string set", t, func() {
		a := NewStringSet()
		Convey("When check equal with a empty string set", func() {
			e := a.Equal(NewStringSet())
			Convey("Then the result is equal", func() {
				So(e, ShouldBeTrue)
			})
		})

		Convey("When check with a non empty string set(a)", func() {
			n := NewStringSet()
			n.Add("a")
			e := a.Equal(n)
			Convey("Then should not equal", func() {
				So(e, ShouldBeFalse)
			})
		})

		Convey("When add element a to it", func() {
			a.Add("a")
			Convey("Then a should in it", func() {
				So(a.Contains("a"), ShouldBeTrue)
			})
		})
	})

	Convey("Given a string set set(a, b)", t, func() {
		a := NewStringSet()
		a.Add("a")
		a.Add("b")
		Convey("When delete a from set", func() {
			a.Delete("a")
			Convey("Then set just has b element", func() {
				So(a.Size(), ShouldEqual, 1)
				So(a.Contains("b"), ShouldBeTrue)
				So(a.Contains("a"), ShouldBeFalse)
			})
		})
	})

	Convey("Given two set a(a, b, c), b(b, d, e)", t, func() {
		a := NewStringSet()
		a.Add("a")
		a.Add("b")
		a.Add("c")

		b := NewStringSet()
		b.Add("b")
		b.Add("d")
		b.Add("e")

		Convey("When union the two set", func() {
			u := a.Union(b)
			Convey("Then the result is set(a, b, c, d, e)", func() {
				expect := NewStringSet()
				expect.Add("a", "b", "c", "d", "e")
				So(u.Equal(expect), ShouldBeTrue)
			})
		})

		Convey("When intersection the two set", func() {
			u := a.Intersection(b)
			Convey("Then the result is set(b)", func() {
				expect := NewStringSet()
				expect.Add("b")
				So(u.Equal(expect), ShouldBeTrue)
			})
		})

		Convey("When difference the two set", func() {
			u := a.Difference(b)
			Convey("Then the result is set(a, c)", func() {
				expect := NewStringSet()
				expect.Add("a", "c")
				So(u.Equal(expect), ShouldBeTrue)
			})
		})

		Convey("When check equal with each other", func() {
			e := b.Equal(a)
			Convey("Then is not equal", func() {
				So(e, ShouldBeFalse)
			})
		})
	})
}
