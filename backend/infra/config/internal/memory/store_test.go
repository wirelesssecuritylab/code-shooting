package memory

import (
	"testing"
	"time"

	"code-shooting/infra/config/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMemoryStore(t *testing.T) {

	Convey("Given a memory store. \n", t, func() {
		store := NewMemoryStore()

		Convey("When Create configItem. \n", func() {
			ver, err := store.Create("key", 1)
			So(err, ShouldBeNil)
			So(ver, ShouldBeLessThan, time.Now().String())
			Convey("Then get the configItem success. \n", func() {
				item := store.Get("key")
				So(item, ShouldNotBeNil)
				So(item, ShouldResemble, &model.ConfigItem{Value: 1, Key: "key", Version: ver})
			})
		})
		Convey("When Create configItem key is exist. \n", func() {
			store.Create("key", 1)
			ver, err := store.Create("key", 2)
			Convey("Then create error. \n", func() {
				So(err.Error(), ShouldContainSubstring, "already exists")
				So(ver, ShouldEqual, "")
			})
		})

		Convey("When get not exist configItem. \n", func() {
			item := store.Get("key2")
			Convey("Then should return nil. \n", func() {
				So(item, ShouldBeNil)
			})
		})

		Convey("When update a exist configItem. \n", func() {
			store.Create("key", 1)
			ver, err := store.Update("key", 2)
			Convey("Then should success. \n", func() {
				So(err, ShouldBeNil)
				So(ver, ShouldBeLessThan, time.Now().String())
			})
		})

		Convey("When delete a  exist configItem. \n", func() {
			store.Create("key", 1)
			err := store.Delete("key")
			Convey("Then should success. \n", func() {
				So(err, ShouldBeNil)
			})
		})
		Convey("When update a not exist configItem. \n", func() {
			ver, err := store.Update("key", 2)
			Convey("Then should return err. \n", func() {
				So(err.Error(), ShouldContainSubstring, "not found")
				So(ver, ShouldEqual, "")
			})
		})

		Convey("When get all configItems. \n", func() {
			store.Create("key1", 1)
			store.Create("key2", 1)
			configs := store.GetAll()
			Convey("Then should success. \n", func() {

				So(len(configs), ShouldEqual, 2)
			})
		})
	})
}
