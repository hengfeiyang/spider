package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStripTags(t *testing.T) {
	Convey("测试HTML标签过滤", t, func() {
		s := []byte("啊<div class=\"a\">是<span>打</span><!--注释-->发<img src=\"http://asf.jpg\">是<br /></div>")
		So(string(KeepTags(s)), ShouldEqual, "啊是打发是")
		So(string(KeepTags(s, "br", "div")), ShouldEqual, "啊<div class=\"a\">是打发是<br /></div>")
		So(string(StripTags(s)), ShouldEqual, string(s))
		So(string(StripTags(s, "br", "div")), ShouldEqual, `啊是<span>打</span><!--注释-->发<img src="http://asf.jpg">是`)
	})
}
