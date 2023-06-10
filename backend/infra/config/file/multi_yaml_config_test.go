package file

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"syscall"
	"testing"
	"time"

	"code-shooting/infra/config/model"

	. "github.com/agiledragon/gomonkey"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

var yamlContentX = `
loadbalance:
    name: lb
    strategy:
        name: RoundRobin
        type: Recreate
    retryEnabled: true
    retryOnNext: 2
    scope:
    members:
        - systemA
        - systemB
`
var yamlContentY = `
loadbalance:
    strategy:
        name: RoundRobin
        type: Recreate
    retryEnabled: false
    retryOnNext: 3
    scope:
    members:
        - systemA
        - systemC
`
var yamlFileX = "xxxx.yaml"
var yamlFileY = "yyyy.yaml"

func NewConfig(filePath string) (model.ConfigSource, error) {
	source, err := NewFileSource(DefaultYamlSource, ParseYamlFile, filePath)
	if err != nil {
		return nil, err
	}
	var ch <-chan struct{}
	err = source.Start(ch)
	if err != nil {
		return nil, err
	}
	return source, nil
}

type FileStat struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	sysStat syscall.Stat_t
}

func (fs *FileStat) IsDir() bool        { return fs.mode&os.ModeDir != 0 }
func (fs *FileStat) Name() string       { return fs.name }
func (fs *FileStat) Size() int64        { return fs.size }
func (fs *FileStat) Mode() os.FileMode  { return fs.mode }
func (fs *FileStat) ModTime() time.Time { return fs.modTime }
func (fs *FileStat) Sys() interface{}   { return &fs.sysStat }

func TestMultiYamlConfig(t *testing.T) {

	file := &FileStat{}
	patches := ApplyFunc(os.Stat, func(_ string) (os.FileInfo, error) {
		return file, nil
	})
	defer patches.Reset()

	patches.ApplyMethod(reflect.TypeOf(file), "IsDir", func(_ *FileStat) bool {
		return true
	})

	patches.ApplyFunc(ioutil.ReadFile, func(file string) ([]byte, error) {
		var yamlContent string
		if file == yamlFileX {
			yamlContent = yamlContentX
		} else if file == yamlFileY || file == "yyyy.yml" {
			yamlContent = yamlContentY
		} else {
			return []byte{}, errors.New("config format not supported")
		}
		return []byte(yamlContent), nil
	})

	os.Setenv("RETRY_ENABLED", "false")
	defer os.Unsetenv("RETRY_ENABLED")

	Convey("Given a not exiseted directory ", t, func() {

		Convey("When new a configurator", func() {

			configurator, err := NewConfig("/xxxx/config/res")

			Convey("Then an error contained 'no such file or directory' should return and configurator should be nil", func() {
				So(err.Error(), ShouldContainSubstring, "no such file or directory")
				So(configurator, ShouldBeNil)
			})
		})
	})

	Convey("Given a config directory in which there are multiple yaml files and a configurator", t, func() {

		patches := ApplyFunc(getAllFilesPathBy, func(_ string, _ []string) ([]string, error) {
			confFileList := []string{yamlFileX, yamlFileY}
			return confFileList, nil
		})
		defer patches.Reset()
		err := os.Mkdir("/tmp/config", os.ModePerm)
		err = writeToFile("/tmp/config"+yamlFileX, yamlContentX)
		So(err, ShouldBeNil)
		err = writeToFile("/tmp/config"+yamlFileY, yamlContentY)
		So(err, ShouldBeNil)

		defer func() {
			os.RemoveAll("/tmp/config")
		}()

		configurator, err := NewConfig("/tmp/config")
		So(err, ShouldEqual, nil)

		Convey("When get the value with the key of loadbalance.name only exists in fileX, but not exists in fileY ", func() {

			name := configurator.Get("loadbalance.name")

			Convey("Then the value should be lb", func() {
				So(name, ShouldNotBeNil)
				So(name.Value, ShouldEqual, "lb")
			})
		})

		Convey("When get the value with the key of loadbalance.retryOnNext existed in mutiple files (value is number) ", func() {

			retryOnNext := configurator.Get("loadbalance.retryOnNext")

			Convey("Then the value should be 3", func() {
				So(retryOnNext, ShouldNotBeNil)
				So(retryOnNext.Value, ShouldEqual, 3)
			})
		})

		Convey("When get the value with the key of loadbalance.members existed in mutiple files (value is list)", func() {

			members := configurator.Get("loadbalance.members")

			Convey("Then the value should be 3", func() {
				So(members, ShouldNotBeNil)
				So(members.Value, ShouldContain, "systemA")
				So(members.Value, ShouldContain, "systemC")
				So(members.Value, ShouldNotContain, "systemB")
			})
		})

		Convey("When get the value with the key of loadbalance.retryEnabled exists in mutiple files (value is environment variable)", func() {

			retryEnabled := configurator.Get("loadbalance.retryEnabled")

			Convey("Then the value should be false", func() {
				So(retryEnabled, ShouldNotBeNil)
				So(retryEnabled.Value, ShouldEqual, false)
			})
		})
	})
}

func TestExistedDirConfig(t *testing.T) {
	Convey("Given a config directory without files ", t, func() {

		configDir, _ := ioutil.TempDir("./", "res")
		defer func() {
			os.RemoveAll(configDir)
		}()

		Convey("When new a configurator", func() {

			configurator, err := NewConfig(configDir)

			Convey("Then an error contains 'no supported config file' should rerun and configurator should be nil", func() {
				So(err.Error(), ShouldContainSubstring, "no supported config file")
				So(configurator, ShouldBeNil)
			})
		})
	})

	Convey("Given a config directory with an empty folder ", t, func() {

		configDir, _ := ioutil.TempDir("./", "res")
		ioutil.TempDir(configDir, "yaml")
		defer func() {
			os.RemoveAll(configDir)
		}()

		Convey("When new a configurator ", func() {

			configurator, err := NewConfig(configDir)

			Convey("Then an error contains 'no supported config file' should return and configurator should be nil", func() {
				So(err.Error(), ShouldContainSubstring, "no supported config file")
				So(configurator, ShouldBeNil)
			})
		})
	})

	Convey("Given a config directory with a yaml file ", t, func() {

		configDir, _ := ioutil.TempDir("./", "res")
		var content = `loadbalance:
        retryOnNext: 2`
		appFile := path.Join(configDir, "app.yaml")
		ioutil.WriteFile(appFile, []byte(content), 0640)

		defer func() {
			os.RemoveAll(configDir)
		}()

		configurator, err := NewConfig(configDir)
		So(err, ShouldEqual, nil)

		Convey("When getting retryOnNext from the directory ", func() {

			retryOnNext := configurator.Get("loadbalance.retryOnNext")

			Convey("Then err should be nil and retryOnNext should be 2", func() {
				So(retryOnNext, ShouldNotBeNil)
				So(retryOnNext.Value, ShouldEqual, 2)
			})
		})
	})
	Convey("Given a config directory where has yaml and yml files both without content  ", t, func() {
		configDir, _ := ioutil.TempDir("./", "res")
		ioutil.TempFile(configDir, "app-*.yaml")
		ioutil.TempFile(configDir, "app-*.yml")

		defer func() {
			os.RemoveAll(configDir)
		}()

		configurator, err := NewConfig(configDir)
		So(err, ShouldEqual, nil)

		Convey("When getting retryOnNext from the directory ", func() {

			retryOnNext := configurator.Get("loadbalance.retryOnNext")

			Convey("Then a 'is not exist' error should return and retryOnNext should be 0", func() {
				So(retryOnNext, ShouldBeNil)
			})
		})
	})
}

func TestMultiFormatConfig(t *testing.T) {

	Convey("Given a config directory in which there are multiple files including json and yaml", t, func() {

		configDir, _ := ioutil.TempDir("./", "res")
		defer func() {
			os.RemoveAll(configDir)
		}()

		var bin_buf bytes.Buffer
		binary.Write(&bin_buf, binary.LittleEndian, 123456)
		ioutil.WriteFile(path.Join(configDir, "app.exe"), bin_buf.Bytes(), 0640)

		content := `apiversion: v1
instrategy:
    name: RoundRobin
    type: Recreate`
		ioutil.WriteFile(path.Join(configDir, "app.yml"), []byte(content), 0640)

		Convey("When get the value with the key of instrategy.name ", func() {

			configurator, err := NewConfig(configDir)
			So(err, ShouldEqual, nil)

			name := configurator.Get("instrategy.name")

			Convey("Then the value should be RoundRobin", func() {
				So(name, ShouldNotBeNil)
				So(name.Value, ShouldEqual, "RoundRobin")
			})
		})
	})
}

func TestOsStatError(t *testing.T) {
	file := &FileStat{}
	patches := ApplyFunc(os.Stat, func(_ string) (os.FileInfo, error) {
		return file, errors.Errorf("os.Stat error")
	})
	defer patches.Reset()

	Convey("Given a config directory has error ", t, func() {

		Convey("When new a configurator", func() {

			configurator, err := NewConfig("/xxxx/config/res")

			Convey("Then an error contains 'has error' and configurator should be nil", func() {
				So(err.Error(), ShouldContainSubstring, "no such file or directory")
				So(configurator, ShouldBeNil)
			})
		})
	})
}
