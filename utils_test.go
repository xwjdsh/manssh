package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestArgumentsCheck(t *testing.T) {
	Convey("init", t, func() {
		So(argumentsCheck(2, 3, 4), ShouldNotBeNil)
		So(argumentsCheck(2, 1, 1), ShouldNotBeNil)
		So(argumentsCheck(2, 1, 4), ShouldBeNil)
		So(argumentsCheck(2, 2, 2), ShouldBeNil)
		So(argumentsCheck(2, -1, -1), ShouldBeNil)
	})
}

func TestContains(t *testing.T) {
	Convey("init", t, func() {
		values := []string{"test1", "test2", "another1", "another2"}
		Convey("check", func() {
			So(contains(values, "es"), ShouldBeTrue)
			So(contains(values, "ther"), ShouldBeTrue)
			So(contains(values, "test2"), ShouldBeTrue)
			So(contains(values, "another1"), ShouldBeTrue)

			So(contains(values, "tt"), ShouldBeFalse)
			So(contains(values, "thr"), ShouldBeFalse)
			So(contains(values, "test8"), ShouldBeFalse)
			So(contains(values, "test3"), ShouldBeFalse)
		})
	})
}

func TestQuery(t *testing.T) {
	Convey("init", t, func() {
		values := []string{"test1", "test2", "another1", "another2"}
		Convey("check", func() {
			So(query(values, []string{"test", "2"}), ShouldBeTrue)
			So(query(values, []string{"another", "1"}), ShouldBeTrue)

			So(query(values, []string{"test", "3"}), ShouldBeFalse)
			So(query(values, []string{"another", "3"}), ShouldBeFalse)
		})
	})
}
