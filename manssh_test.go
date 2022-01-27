package manssh

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	mainConfigContent = `
Include %s/config.d/*
Host home1
    hostname 192.168.1.11
Host main1
    hostname 192.168.1.10
Host main2
    hostname 192.168.1.20
    user wen
    port 77
Host main3
    hostname 192.168.1.30
    user ROOT
	port 77
`
	testConfigContent = `
Host *
	port 22022
Host test1
    hostname 192.168.2.10
    user root
    port 22
Host test2 main2
    hostname 192.168.2.20
    port 77
Host Test3
    hostname 192.168.2.30
    user ROOT
	port 77
`
	homeConfigContent = `
Host home1
    hostname 192.168.3.10
    user ROOT
    port 77
Host home2
    hostname 192.168.3.20
    user root
    port 77
Host home3
    hostname 192.168.3.30
    user ROOT
	port 77
`
)

var (
	configRootDir  = filepath.Join(os.TempDir(), "manssh")
	mainConfigPath = filepath.Join(configRootDir, "config")
	testConfigPath = filepath.Join(configRootDir, "config.d", "test")
	homeConfigPath = filepath.Join(configRootDir, "config.d", "home")
)

func initConfig() {
	_ = os.MkdirAll(configRootDir, os.ModePerm)
	_ = os.MkdirAll(filepath.Join(configRootDir, "config.d"), os.ModePerm)
	if err := ioutil.WriteFile(mainConfigPath, []byte(fmt.Sprintf(mainConfigContent, configRootDir)), 0644); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(testConfigPath, []byte(testConfigContent), 0644); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(homeConfigPath, []byte(homeConfigContent), 0644); err != nil {
		panic(err)
	}
}

func TestList(t *testing.T) {
	initConfig()
	defer os.Remove(configRootDir)

	hosts, err := List(mainConfigPath, ListOption{})
	require.Nil(t, err)
	require.Equal(t, 10, len(hosts))
	hostMap := map[string]*HostConfig{}
	for _, host := range hosts {
		hostMap[host.Alias] = host
	}

	main2 := hostMap["main2"]
	require.NotNil(t, main2)
	require.Equal(t, "wen@192.168.2.20:22022", main2.Connection)
	require.Equal(t, 1, len(main2.ImplicitConfig))
	require.Equal(t, "22022", main2.ImplicitConfig["port"])

	hosts, err = List(mainConfigPath, ListOption{
		Keywords: []string{"Test"},
	})
	require.Nil(t, err)
	require.Equal(t, 1, len(hosts))

	hosts, err = List(mainConfigPath, ListOption{
		Keywords:   []string{"Test"},
		IgnoreCase: true,
	})
	require.Nil(t, err)
	require.Equal(t, 3, len(hosts))
}

func TestAdd(t *testing.T) {
	initConfig()
	defer os.Remove(configRootDir)

	_, err := Add(mainConfigPath, &AddOption{
		Path:    testConfigPath,
		Alias:   "test1",
		Connect: "xxx@1.2.3.4:11",
		Config:  map[string]string{},
	})
	require.NotNil(t, err)

	host, err := Add(mainConfigPath, &AddOption{
		Path:    testConfigPath,
		Alias:   "test4",
		Connect: "xxx@1.2.3.4",
		Config:  map[string]string{},
	})
	require.Nil(t, err)
	require.Equal(t, "22022", host.ImplicitConfig["port"])
	require.Equal(t, "xxx@1.2.3.4:22022", host.Connection)

	hosts, err := List(mainConfigPath, ListOption{})
	require.Nil(t, err)
	require.Equal(t, 11, len(hosts))
}

func TestUpdate(t *testing.T) {
	initConfig()
	defer os.Remove(configRootDir)

	_, err := Update(mainConfigPath, &UpdateOption{
		Alias:   "test4",
		Connect: "xxx@1.2.3.4:11",
	})
	require.NotNil(t, err)

	host, err := Update(mainConfigPath, &UpdateOption{
		Alias:    "test1",
		NewAlias: "test4",
		Connect:  "xxx@1.2.3.4:11",
		Config: map[string]string{
			"IdentifyFile": "~/.ssh/test4",
		},
	})
	require.Nil(t, err)
	require.Equal(t, "xxx@1.2.3.4:22022", host.Connection)
	require.Equal(t, "~/.ssh/test4", host.OwnConfig["identifyfile"])
	require.Equal(t, "22022", host.ImplicitConfig["port"])

	hosts, err := List(mainConfigPath, ListOption{})
	require.Nil(t, err)
	require.Equal(t, 10, len(hosts))
	hostMap := map[string]*HostConfig{}
	for _, host := range hosts {
		hostMap[host.Alias] = host
	}
	require.Nil(t, hostMap["test1"])
	require.NotNil(t, hostMap["test4"])
}

func TestDelete(t *testing.T) {
	initConfig()
	defer os.Remove(configRootDir)

	_, err := Delete(mainConfigPath, "home1", "test1", "main4")
	require.NotNil(t, err)

	hosts, err := Delete(mainConfigPath, "home1", "test1", "*")
	require.Nil(t, err)
	require.Equal(t, 3, len(hosts))

	hosts, err = List(mainConfigPath, ListOption{})
	require.Nil(t, err)
	require.Equal(t, 8, len(hosts))
	hostMap := map[string]*HostConfig{}
	for _, host := range hosts {
		hostMap[host.Alias] = host
	}
	require.Nil(t, hostMap["home1"])
	require.Nil(t, hostMap["test1"])
	require.NotNil(t, hostMap["*"])
	require.Equal(t, "22", hostMap["main1"].ImplicitConfig["port"])
}

func TestGetFilePaths(t *testing.T) {
	initConfig()
	defer os.Remove(configRootDir)
	paths, err := GetFilePaths(mainConfigPath)
	require.Nil(t, err)
	require.Equal(t, 3, len(paths))
}
