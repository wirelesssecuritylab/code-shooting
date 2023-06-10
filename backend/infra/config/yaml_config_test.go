package config

import (
	"errors"

	"io/ioutil"
	"os"

	"reflect"
	"testing"

	. "github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
)

type Template struct {
	Metadata string `yaml:"apiversion,omitempty"`
}
type Strategy struct {
	Name string  `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
	Template
}

type Member struct {
	Name string  `yaml:"name,omitempty"`
	Type *string `yaml:"type,omitempty"`
}

type MemberName struct {
	Name string `yaml:"name,omitempty"`
}

type InStrategy struct {
	Strategy `json:",inline"`
}

type app struct {
	Apiversion string `json:"apiversion" yaml:"apiversion"`
	Name       string `yaml:"name" json:"name"`
	Type       string
}

type AppJsonInline struct {
	Apiversion string `json:"apiversion,omitempty"`
	Strategy   `json:",inline"`
}

type AppJsonInline2 struct {
	Apiversion string `json:"apiversion,omitempty"`
	InStrategy `json:"instrategy"`
}

type AppYamlInline struct {
	Apiversion string `yaml:"apiversion,omitempty"`
	Member     `yaml:",inline"`
}

type AppYamlInline2 struct {
	Apiversion string     `yaml:"apiversion,omitempty"`
	Member     MemberName `yaml:",inline"`
}

var appContent = `
apiversion: v1
name: RoundRobin
type: Recreate
`

var app2Content = `
apiversion: v1
instrategy:
    name: RoundRobin
    type: Recreate
`

var yamlContent = `
apiversion: v1
loadbalance:
    name: lb
    strategy:
        name: RoundRobin
        type: Recreate
    retryEnabled: false
    retryOnNext: 2
    scope:
    members:
        - systemA
        - systemB
`

type JsonTagCaseSensitive struct {
	AA  string `json:"aa"`
	BB  string `json:"BB,omitempty"`
	CC  string
	DD  string `json:"d_d"`
	D_D string
}

var caseSensitiveJsonContent = `
aa: this is aa
bb: this is bb
cc: this is cc
d_d: this is d_d
qq: this is qq
`

type YamlTagCaseSensitive struct {
	AA string `yaml:"aa"`
	BB string `yaml:"BB"`
	CC string `yaml:"cc"`
	DD string `yaml:"CC"`
	EE string
	FF string
}

var caseSensitiveYamlContent = `
aa: this is aa
BB: this is BB
cc: this is cc
CC: this is CC
ee: this is ee
qq: this is qq
`

func TestBasicTypeInSingleYamlConfig(t *testing.T) {
	file := &FileStat{}
	patches := ApplyFunc(os.Stat, func(_ string) (os.FileInfo, error) {
		return file, nil
	})
	defer patches.Reset()

	patches.ApplyMethod(reflect.TypeOf(file), "IsDir", func(_ *FileStat) bool {
		return false
	})

	patches.ApplyFunc(ioutil.ReadFile, func(_ string) ([]byte, error) {
		return []byte(yamlContent), nil
	})

	Convey("Given a existed yaml config file and a configurator", t, func() {

		yamlFile, _ := ioutil.TempFile(".", "*.yaml")
		defer func() {
			yamlFile.Close()
			os.Remove(yamlFile.Name())
		}()

		configurator, err := NewConfig(yamlFile.Name())
		So(err, ShouldEqual, nil)

		Convey("When get the value with the key of loadbalance.retryOnNext(number) ", func() {

			retryOnNext := 0
			err := configurator.Get("loadbalance.retryOnNext", &retryOnNext)

			Convey("Then the value should be 2", func() {
				So(err, ShouldEqual, nil)
				So(retryOnNext, ShouldEqual, 2)
			})
		})

		Convey("When get the value with the key of loadbalance.name(string) ", func() {

			name := ""
			err := configurator.Get("loadbalance.name", &name)

			Convey("Then the value should be lb", func() {
				So(err, ShouldEqual, nil)
				So(name, ShouldEqual, "lb")
			})
		})

		Convey("When get the value with the key of loadbalance.retryEnabled(bool) ", func() {

			retryEnabled := true
			err := configurator.Get("loadbalance.retryEnabled", &retryEnabled)

			Convey("Then the value should be false", func() {
				So(err, ShouldEqual, nil)
				So(retryEnabled, ShouldEqual, false)
			})
		})

		Convey("When get the value with the key of loadbalance.scope(null) ", func() {

			var scope *string
			err := configurator.Get("loadbalance.scope", &scope)

			Convey("Then the value should be null", func() {
				So(err, ShouldEqual, nil)
				So(scope, ShouldEqual, nil)
			})
		})

		Convey("When get the value with the key of loadbalance.strategy(map), and value-type is struct", func() {

			strategy := Strategy{}
			err := configurator.Get("loadbalance.strategy", &strategy)

			Convey("Then the value should be RoundRobin", func() {
				So(err, ShouldEqual, nil)
				So(strategy.Name, ShouldEqual, "RoundRobin")
			})
		})

		Convey("When get the value with the key of loadbalance.strategy(map), and value-type is struct which has a pointer member", func() {

			strategy := Strategy{}
			err := configurator.Get("loadbalance.strategy", &strategy)

			Convey("Then value should be Recreate", func() {
				So(err, ShouldEqual, nil)
				So(*strategy.Type, ShouldEqual, "Recreate")
			})
		})

		Convey("When get the value with the key of loadbalance.strategy(map), and value-type is struct which has a member not in yml", func() {

			strategy := Strategy{}
			err := configurator.Get("loadbalance.strategy", &strategy)

			Convey("Then the value should be expect", func() {
				So(err, ShouldEqual, nil)
				sType := "Recreate"
				So(strategy, ShouldResemble, Strategy{Name: "RoundRobin", Type: &sType, Template: Template{Metadata: ""}})
			})
		})

		Convey("When get the value with the key of loadbalance.strategy(map), and value-type is map", func() {

			var strategy map[string]string
			err := configurator.Get("loadbalance.strategy", &strategy)

			Convey("Then the value with the key of \"name\" should be RoundRobin", func() {
				So(err, ShouldEqual, nil)
				So(strategy["name"], ShouldEqual, "RoundRobin")
			})
		})

		Convey("When get the value with the key of loadbalance.members(list)", func() {

			var members [2]string
			err := configurator.Get("loadbalance.members", &members)

			Convey("Then the value should containe systemA and systemB", func() {
				So(err, ShouldEqual, nil)
				So(members, ShouldContain, "systemA")
				So(members, ShouldContain, "systemB")
			})
		})

		Convey("When get the value with the key of loadbalance.strategy.id(not existed)", func() {

			id := ""
			err := configurator.Get("loadbalance.strategy.id", &id)
			Convey("Then an 'is not exist' error should return and the value should be empty", func() {
				So(err, ShouldNotBeNil)
				So(IsNotExist(err), ShouldBeTrue)
				So(id, ShouldBeEmpty)
			})
		})

		Convey("When input the value-type with the key of loadbalance.strategy is inconsistent with the value-type in yml", func() {

			strategy := ""
			err := configurator.Get("loadbalance.strategy", &strategy)

			Convey("Then an error should return, and the value should be empty", func() {
				So(err, ShouldBeError)
				So(strategy, ShouldBeEmpty)
			})
		})
	})

}

func TestEmbededTypeInSingleYamlConfig(t *testing.T) {

	file := &FileStat{}
	patches := ApplyFunc(os.Stat, func(_ string) (os.FileInfo, error) {
		return file, nil
	})
	defer patches.Reset()

	patches.ApplyMethod(reflect.TypeOf(file), "IsDir", func(_ *FileStat) bool {
		return false
	})

	Convey("Given a existed yaml config file and a configurator", t, func() {

		yamlFile, _ := ioutil.TempFile(".", "*.yaml")
		defer func() {
			yamlFile.Close()
			os.Remove(yamlFile.Name())
		}()

		patches := ApplyFunc(ioutil.ReadFile, func(_ string) ([]byte, error) {
			return []byte(appContent), nil
		})
		defer patches.Reset()

		configurator, err := NewConfig(yamlFile.Name())
		So(err, ShouldEqual, nil)

		Convey("When input the configureObj with a embedded json struct ", func() {

			appInline := AppJsonInline{}
			err := configurator.Get("", &appInline)

			Convey("Then the configureObj should equal expect", func() {
				So(err, ShouldBeNil)
				typeValue := "Recreate"
				So(appInline, ShouldResemble, AppJsonInline{Apiversion: "v1", Strategy: Strategy{Name: "RoundRobin", Type: &typeValue, Template: Template{Metadata: ""}}})
			})
		})

		Convey("When input the configureObj with a embedded yaml struct ", func() {

			appInline := AppYamlInline{}
			err := configurator.Get("", &appInline)

			Convey("Then the configureObj should equal expect", func() {
				So(err, ShouldBeNil)
				typeValue := "Recreate"
				So(appInline, ShouldResemble, AppYamlInline{Apiversion: "v1", Member: Member{Name: "RoundRobin", Type: &typeValue}})
			})
		})

		Convey("When input the configureObj with a yaml-inline-tag struct ", func() {

			appInline := AppYamlInline2{}
			err := configurator.Get("", &appInline)

			Convey("Then the configureObj should equal expect", func() {
				So(err, ShouldBeNil)
				So(appInline, ShouldResemble, AppYamlInline2{Apiversion: "v1", Member: MemberName{Name: "RoundRobin"}})
			})
		})

		Convey("When input the configureObj has both the json and yaml tags, and no tag ", func() {

			appConf := app{}
			err := configurator.Get("", &appConf)

			Convey("Then the configureObj should be all configuration", func() {
				So(err, ShouldEqual, nil)
				So(appConf, ShouldResemble, app{Apiversion: "v1", Name: "RoundRobin", Type: "Recreate"})
			})
		})
	})

	Convey("Given a existed yaml config file and a configurator ", t, func() {

		yamlFile, _ := ioutil.TempFile(".", "*.yaml")
		defer func() {
			yamlFile.Close()
			os.Remove(yamlFile.Name())
		}()

		patches := ApplyFunc(ioutil.ReadFile, func(_ string) ([]byte, error) {
			return []byte(app2Content), nil
		})
		defer patches.Reset()

		configurator, err := NewConfig(yamlFile.Name())
		So(err, ShouldEqual, nil)

		Convey("When input the configureObj with multi-layer embedded json struct ", func() {

			appInline := AppJsonInline2{}
			err := configurator.Get("", &appInline)

			Convey("Then the configureObj should equal expect", func() {
				So(err, ShouldBeNil)
				typeValue := "Recreate"
				So(appInline, ShouldResemble, AppJsonInline2{Apiversion: "v1", InStrategy: InStrategy{Strategy: Strategy{Name: "RoundRobin", Type: &typeValue, Template: Template{Metadata: ""}}}})
			})
		})
	})

	Convey("Given a existed yaml config file and a configurator ", t, func() {

		yamlFile, _ := ioutil.TempFile(".", "*.yaml")
		defer func() {
			yamlFile.Close()
			os.Remove(yamlFile.Name())
		}()

		patches := ApplyFunc(ioutil.ReadFile, func(_ string) ([]byte, error) {
			return []byte(caseSensitiveJsonContent), nil
		})
		defer patches.Reset()

		configurator, err := NewConfig(yamlFile.Name())
		So(err, ShouldEqual, nil)

		Convey("When input the configureObj with json case sensitive tag ", func() {

			jsonConf := JsonTagCaseSensitive{}
			jsonErr := configurator.Get("", &jsonConf)

			Convey("Then the configureObjshould equal expect", func() {
				// json 不允许配置内容多于结构定义，反之可以
				// json 不按tag解析，按小写变量名解析
				So(jsonErr, ShouldEqual, nil)
				So(jsonConf, ShouldResemble, JsonTagCaseSensitive{AA: "this is aa", BB: "this is bb", CC: "this is cc", DD: "", D_D: "this is d_d"})
			})
		})
	})

	Convey("Given a existed yaml config file and a configurator ", t, func() {

		yamlFile, _ := ioutil.TempFile(".", "*.yaml")
		defer func() {
			yamlFile.Close()
			os.Remove(yamlFile.Name())
		}()

		patches := ApplyFunc(ioutil.ReadFile, func(_ string) ([]byte, error) {
			return []byte(caseSensitiveYamlContent), nil
		})
		defer patches.Reset()

		configurator, err := NewConfig(yamlFile.Name())
		So(err, ShouldEqual, nil)

		Convey("When input the configureObj with yaml case sensitive tag ", func() {
			patches := ApplyFunc(ioutil.ReadFile, func(_ string) ([]byte, error) {
				return []byte(caseSensitiveYamlContent), nil
			})
			defer patches.Reset()

			yamlConf := YamlTagCaseSensitive{}
			yamlErr := configurator.Get("", &yamlConf)

			Convey("Then the configureObj should equal expect", func() {
				// yaml 不允许配置内容多于结构定义，反之可以
				// yaml 按tag解析，没有tag 按小写变量名解析，大小写敏感
				So(yamlErr, ShouldEqual, nil)
				So(yamlConf, ShouldResemble, YamlTagCaseSensitive{AA: "this is aa", BB: "this is BB", CC: "this is cc", DD: "this is CC", EE: "this is ee", FF: ""})
			})
		})
	})
}

func TestAbnormalYamlConfig(t *testing.T) {
	file := &FileStat{}
	patches := ApplyFunc(os.Stat, func(_ string) (os.FileInfo, error) {
		return file, nil
	})
	defer patches.Reset()

	patches.ApplyFunc(ioutil.ReadFile, func(_ string) ([]byte, error) {
		return []byte{}, errors.New("The file is not existed")
	})

	Convey("Given an empty config file path ", t, func() {
		Convey("When new a configurator ", func() {

			configurator, err := NewConfig("")

			Convey("Then an error should return and contain 'config path is empty' error info", func() {
				So(err.Error(), ShouldContainSubstring, "config path is empty")
				So(configurator, ShouldBeNil)
			})
		})
	})

	Convey("Given a config file with the suffix of conf ", t, func() {

		confFile := "/tmp/config"
		os.Mkdir(confFile, os.ModePerm)
		defer os.RemoveAll(confFile)

		Convey("When new a configurator ", func() {

			configurator, err := NewConfig(confFile)

			Convey("Then an error should return and contain 'no supported config file' error info", func() {
				So(err.Error(), ShouldContainSubstring, "no supported config file")
				So(configurator, ShouldBeNil)
			})
		})
	})

	Convey("Given a not existed yaml config file ", t, func() {

		notExistedFile := "notExisted.yml"

		Convey("When new a configurator ", func() {

			configurator, err := NewConfig(notExistedFile)

			Convey("Then an error should return and contain 'no such file or directory' error info", func() {
				So(err.Error(), ShouldContainSubstring, "no such file or directory")
				So(configurator, ShouldBeNil)
			})
		})
	})

}
