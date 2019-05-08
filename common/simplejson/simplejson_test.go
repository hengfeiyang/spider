package simplejson

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJson1(t *testing.T) {
	s := `{"code":1, "message":"ok", "items":[{"id":1,"name":"1xx"}, {"id":2, "name":"2xx"}]}`
	Convey("测试JSON解析", t, func() {
		j, err := New([]byte(s))
		So(err, ShouldBeNil)
		Convey("测试Get方法", func() {
			code := j.Get("code").MustInt()
			So(code, ShouldEqual, 1)
			message := j.Get("message").MustString()
			So(message, ShouldEqual, "ok")
			id := j.Get("items").GetIndex(0).Get("id").MustInt()
			So(id, ShouldEqual, 1)
		})
		Convey("测试Path方法", func() {
			p := j.Path("$.items[0].name").MustString()
			So(p, ShouldEqual, "1xx")
		})
	})
}
