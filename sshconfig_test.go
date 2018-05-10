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
    user ROOT
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
		list, err := List(f.Name(), nil)
		So(err, ShouldBeNil)
		So(len(list), ShouldEqual, 3)

		list, err = List(f.Name(), []string{"77"})
		So(err, ShouldBeNil)
		So(len(list), ShouldEqual, 2)

		list, err = List(f.Name(), []string{"root"})
		So(err, ShouldBeNil)
		So(len(list), ShouldEqual, 2)

		list, err = List(f.Name(), []string{"root"}, true)
		So(err, ShouldBeNil)
		So(len(list), ShouldEqual, 3)

		list, err = List(f.Name(), []string{"test", "77", "30"})
		So(err, ShouldBeNil)
		So(len(list), ShouldEqual, 1)
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
		So(Add(f.Name(), &HostConfig{Aliases: "test1", Connect: "2.2.2.2"}, ""), ShouldNotBeNil)

		add := &HostConfig{Aliases: "test4", Connect: "root@2.2.2.2"}
		So(Add(f.Name(), add, ""), ShouldBeNil)
		So(add, ShouldResemble, &HostConfig{Aliases: "test4", Connect: "root@2.2.2.2:22", Config: map[string]string{}})
		list, _ := List(f.Name(), []string{"test4"})
		So(list, ShouldResemble, []*HostConfig{add})
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
		list, _ := List(f.Name(), []string{"test1"})
		So(list, ShouldResemble, []*HostConfig{update1})

		So(update2, ShouldResemble, &HostConfig{Aliases: "test2", Connect: "wendell@192.168.99.20:77", Config: map[string]string{}})
		list, _ = List(f.Name(), []string{"test2"})
		So(list, ShouldResemble, []*HostConfig{update2})

		So(update3, ShouldResemble, &HostConfig{Aliases: "test4", Connect: "ROOT@192.168.99.30:77", Config: map[string]string{}})
		list, _ = List(f.Name(), []string{"test3"})
		So(list, ShouldBeEmpty)
		list, _ = List(f.Name(), []string{"test4"})
		So(list, ShouldResemble, []*HostConfig{update3})
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
		list, _ := List(f.Name(), []string{"test1"})
		So(list, ShouldBeEmpty)

		list, _ = List(f.Name(), []string{"test2"})
		So(list, ShouldBeEmpty)

		list, _ = List(f.Name(), nil)
		So(len(list), ShouldEqual, 1)
	})
}
