package test

import (
	"context"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"
	"testing"
)

func TestFxApp(t *testing.T) {
	Convey("Given a normal app", t, func() {
		app := fx.New(fx.Invoke(func() {}))

		Convey("When start fx app", func() {
			err := StartFxApp(app)
			defer func() {
				StopFxApp(app)
			}()

			Convey("Then nothing will happen", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given a app with abnormal start callback ", t, func() {
		app := fx.New(fx.Invoke(func(lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return errors.New("abnormal start fx app")
				},
			})
		}))

		Convey("When start fx app", func() {
			err := StartFxApp(app)
			defer func() {
				StopFxApp(app)
			}()

			Convey("Then should tips: abnormal start fx app", func() {
				So(err.Error(), ShouldContainSubstring, "abnormal start fx app")
			})
		})
	})
}
