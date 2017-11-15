package manssh

import (
	"testing"

	"github.com/kevinburke/ssh_config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFormatConnect(t *testing.T) {
	Convey("init", t, func() {
		So(FormatConnect("root", "1.1.1.1", "22"), ShouldEqual, "root@1.1.1.1:22")
	})
}

func TestParseConnct(t *testing.T) {
	Convey("init", t, func() {
		user, hostname, port := ParseConnct("root@1.1.1.1:77")
		So([]string{user, hostname, port}, ShouldResemble, []string{"root", "1.1.1.1", "77"})

		user, hostname, port = ParseConnct("1.1.1.1:77")
		So([]string{user, hostname, port}, ShouldResemble, []string{GetHomeDir(), "1.1.1.1", "77"})

		user, hostname, port = ParseConnct("1.1.1.1")
		So([]string{user, hostname, port}, ShouldResemble, []string{GetHomeDir(), "1.1.1.1", "22"})
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
			So(Query(values, []string{"test", "2"}), ShouldBeTrue)
			So(Query(values, []string{"another", "1"}), ShouldBeTrue)

			So(Query(values, []string{"test", "3"}), ShouldBeFalse)
			So(Query(values, []string{"another", "3"}), ShouldBeFalse)
		})
	})
}

func TestCheckAlias(t *testing.T) {
	Convey("init", t, func() {
		aliasMap := map[string]*ssh_config.Host{
			"test1": &ssh_config.Host{},
			"test2": &ssh_config.Host{},
		}
		Convey("check", func() {
			So(CheckAlias(aliasMap, true, "test1", "test2"), ShouldBeNil)
			So(CheckAlias(aliasMap, true, "test1", "test3"), ShouldNotBeNil)

			So(CheckAlias(aliasMap, false, "test1", "test2"), ShouldNotBeNil)
			So(CheckAlias(aliasMap, false, "test3", "test4"), ShouldBeNil)
		})
	})
}
