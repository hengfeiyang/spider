package util

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCompare1(t *testing.T) {
	Convey("测试相等比较", t, func() {
		Convey("字符相等测试", func() {
			So(Eq("123", "123"), ShouldBeTrue)
			So(Eq("123", "1223"), ShouldBeFalse)
		})
		Convey("数字相等测试", func() {
			So(Eq(123, 123), ShouldBeTrue)
			So(Eq(123, 456), ShouldBeFalse)
		})
		Convey("错误类型比较", func() {
			So(Eq("123", 123), ShouldBeFalse)
		})
		Convey("interface比较", func() {
			var v1 interface{} = "123"
			var v2 = "123"
			So(Eq(v1, v2), ShouldBeTrue)
		})
	})

	Convey("测试不相等比较", t, func() {
		Convey("字符不相等测试", func() {
			So(Neq("123", "123"), ShouldBeFalse)
			So(Neq("123", "1223"), ShouldBeTrue)
		})
		Convey("数字不相等测试", func() {
			So(Neq(123, 123), ShouldBeFalse)
			So(Neq(123, 456), ShouldBeTrue)
		})
		Convey("错误类型比较", func() {
			So(Neq("123", 123), ShouldBeTrue)
		})
		Convey("interface比较", func() {
			var v1 interface{} = "123"
			var v2 = "123"
			So(Neq(v1, v2), ShouldBeFalse)
		})
	})

	Convey("测试大于小于比较", t, func() {
		So(Gt(123, 456), ShouldBeFalse)
		So(Gt(123, 121), ShouldBeTrue)
		So(Lt(123, 456), ShouldBeTrue)
	})

	Convey("测试空值", t, func() {
		So(Empty(""), ShouldBeTrue)
		So(Empty(0), ShouldBeTrue)
		var v int64 = 0
		So(Empty(v), ShouldBeTrue)
		So(Empty(float32(0.0)), ShouldBeTrue)
		So(Empty(float64(0.0)), ShouldBeTrue)
		So(Empty(nil), ShouldBeTrue)
		So(Empty(false), ShouldBeTrue)
		So(Empty("1"), ShouldBeFalse)
		So(Empty(1), ShouldBeFalse)
		So(Empty(1.1), ShouldBeFalse)
		So(Empty(errors.New("x")), ShouldBeFalse)
	})
}
