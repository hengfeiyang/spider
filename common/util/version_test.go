package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVersionDiff(t *testing.T) {
	Convey("测试版本对比", t, func() {
		So(VersionDiff("1.1.1", "1.2.1"), ShouldBeFalse)
	})
}
