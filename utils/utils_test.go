package utils

import (
	"os/user"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFormatConnect(t *testing.T) {
	Convey("init", t, func() {
		So(FormatConnect("root", "1.1.1.1", "22"), ShouldEqual, "root@1.1.1.1:22")
	})
}

func TestParseConnct(t *testing.T) {
	Convey("init", t, func() {
		currentUser := ""
		tmp, err := user.Current()
		if err == nil {
			currentUser = tmp.Username
		}
		u, hostname, port := ParseConnct("root@1.1.1.1:77")
		So([]string{u, hostname, port}, ShouldResemble, []string{"root", "1.1.1.1", "77"})

		u, hostname, port = ParseConnct("1.1.1.1:77")
		So([]string{u, hostname, port}, ShouldResemble, []string{currentUser, "1.1.1.1", "77"})

		u, hostname, port = ParseConnct("1.1.1.1")
		So([]string{u, hostname, port}, ShouldResemble, []string{currentUser, "1.1.1.1", "22"})
	})
}

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
			So(Query(values, []string{"test", "2"}, false), ShouldBeTrue)
			So(Query(values, []string{"another", "2"}, false), ShouldBeTrue)

			So(Query(values, []string{"TEST", "1"}, false), ShouldBeFalse)
			So(Query(values, []string{"Another", "1"}, false), ShouldBeFalse)

			So(Query(values, []string{"TEST", "1"}, true), ShouldBeTrue)
			So(Query(values, []string{"Another", "1"}, true), ShouldBeTrue)
		})
	})
}

func TestCheckAlias(t *testing.T) {
	Convey("init", t, func() {
		aliasMap := map[string]*sshConfigHost{
			"test1": {},
			"test2": {},
		}
		Convey("check", func() {
			So(CheckAlias(aliasMap, true, "test1", "test2"), ShouldBeNil)
			So(CheckAlias(aliasMap, true, "test1", "test3"), ShouldNotBeNil)

			So(CheckAlias(aliasMap, false, "test1", "test2"), ShouldNotBeNil)
			So(CheckAlias(aliasMap, false, "test3", "test4"), ShouldBeNil)
		})
	})
}
