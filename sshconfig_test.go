package manssh

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const configContent = `
Host test1
    hostname 192.168.99.10
    user root
    port 22

Host test2
    hostname 192.168.99.20
    user root
    port 77

Host test3
    hostname 192.168.99.30
    user root
		port 77
`

func TestList(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return
	}
	defer os.Remove(f.Name())
	if _, err = f.WriteString(configContent); err != nil {
		return
	}
	Convey("init", t, func() {
		So(len(List(f.Name())), ShouldEqual, 3)
		So(len(List(f.Name(), "77")), ShouldEqual, 2)
		So(len(List(f.Name(), "test", "77", "30")), ShouldEqual, 1)
	})
}

func TestAdd(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return
	}
	defer os.Remove(f.Name())
	if _, err = f.WriteString(configContent); err != nil {
		return
	}
	Convey("init", t, func() {
		So(Add(f.Name(), &HostConfig{Aliases: "test1", Connect: "2.2.2.2"}), ShouldNotBeNil)

		add := &HostConfig{Aliases: "test4", Connect: "root@2.2.2.2"}
		So(Add(f.Name(), add), ShouldBeNil)
		So(add, ShouldResemble, &HostConfig{Aliases: "test4", Connect: "root@2.2.2.2:22", Config: map[string]string{}})
		So(List(f.Name(), "test4"), ShouldResemble, []*HostConfig{add})
	})
}

func TestUpdate(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return
	}
	defer os.Remove(f.Name())
	if _, err = f.WriteString(configContent); err != nil {
		return
	}
	Convey("init", t, func() {
		So(Update(f.Name(), &HostConfig{Aliases: "test4", Connect: "2.2.2.2"}, ""), ShouldNotBeNil)

		update1 := &HostConfig{Aliases: "test1", Connect: "root@2.2.2.2:77"}
		update2 := &HostConfig{Aliases: "test2", Config: map[string]string{"user": "wendell", "port": "77"}}
		update3 := &HostConfig{Aliases: "test3"}
		So(Update(f.Name(), update1, ""), ShouldBeNil)
		So(Update(f.Name(), update2, ""), ShouldBeNil)
		So(Update(f.Name(), update3, "test4"), ShouldBeNil)

		So(update1, ShouldResemble, &HostConfig{Aliases: "test1", Connect: "root@2.2.2.2:77", Config: map[string]string{}})
		So(List(f.Name(), "test1"), ShouldResemble, []*HostConfig{update1})

		So(update2, ShouldResemble, &HostConfig{Aliases: "test2", Connect: "wendell@192.168.99.20:77", Config: map[string]string{}})
		So(List(f.Name(), "test2"), ShouldResemble, []*HostConfig{update2})

		So(update3, ShouldResemble, &HostConfig{Aliases: "test4", Connect: "root@192.168.99.30:77", Config: map[string]string{}})
		So(List(f.Name(), "test3"), ShouldBeEmpty)
		So(List(f.Name(), "test4"), ShouldResemble, []*HostConfig{update3})
	})
}

func TestDelete(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return
	}
	defer os.Remove(f.Name())
	if _, err = f.WriteString(configContent); err != nil {
		return
	}
	Convey("init", t, func() {
		So(Delete(f.Name(), "test4"), ShouldNotBeNil)

		So(Delete(f.Name(), "test1", "test2"), ShouldBeNil)
		So(List(f.Name(), "test1"), ShouldBeEmpty)
		So(List(f.Name(), "test2"), ShouldBeEmpty)
		So(len(List(f.Name())), ShouldEqual, 1)
	})
}
