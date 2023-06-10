package repository

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"code-shooting/domain/entity/spec"
)

func TestSpecToQuery(t *testing.T) {
	Convey("test institute equal", t, func() {
		query, args := specToQuery(spec.Institute.Equal("无线研究院"))
		So(query, ShouldEqual, "institute = ?")
		So(args, ShouldResemble, []interface{}{"无线研究院"})
	})
	Convey("test import time until", t, func() {
		n := time.Now()
		query, args := specToQuery(spec.ImportTime.Until(n))
		So(query, ShouldEqual, "import_time <= ?")
		So(args, ShouldResemble, []interface{}{n})
	})
	Convey("test and", t, func() {
		n := time.Now()
		andSpec := spec.NewAndSpec(spec.Center.Equal("虚拟化中心"), spec.ImportTime.Since(n))
		query, args := specToQuery(andSpec)
		So(query, ShouldEqual, "(center = ?) AND (import_time >= ?)")
		So(args, ShouldResemble, []interface{}{"虚拟化中心", n})
	})
}
