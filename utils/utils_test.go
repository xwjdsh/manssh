package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestArgumentsCheck(t *testing.T) {
	Convey("init", t, func() {
		So(ArgumentsCheck(2, 3, 4), ShouldNotBeNil)
		So(ArgumentsCheck(2, 1, 1), ShouldNotBeNil)
		So(ArgumentsCheck(2, 1, 4), ShouldBeNil)
		So(ArgumentsCheck(2, 2, 2), ShouldBeNil)
		So(ArgumentsCheck(2, -1, -1), ShouldBeNil)
	})
}

func TestQuery(t *testing.T) {
	Convey("init", t, func() {
		values := []string{"test1", "test2", "another1", "another2"}
		Convey("check", func() {
			So(Query(values, []string{"test", "2"}), ShouldBeTrue)
			So(Query(values, []string{"another", "1"}), ShouldBeTrue)

			So(Query(values, []string{"test", "3"}), ShouldBeFalse)
			So(Query(values, []string{"another", "3"}), ShouldBeFalse)
		})
	})
}
