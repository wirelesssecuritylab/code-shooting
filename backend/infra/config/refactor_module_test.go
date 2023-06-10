package config

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"code-shooting/infra/x/test"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"

	"code-shooting/infra/config/internal/log"
	"code-shooting/infra/config/model"

	uConfig "go.uber.org/config"
)

type TestConf struct {
	Apiversion string       `yaml:"apiversion,omitempty"`
	Strategy   TestStrategy `yaml:"instrategy"`
}
type TestStrategy struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
}

func TestRefactorConfigModule(t *testing.T) {

	Convey("Given a config file ", t, func() {
		content := `code-shooting:
  apiversion: v1
  instrategy:
    name: RoundRobin
    type: Recreate`

		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()
		SetLogger(&log.DefaultLog{})
		var conf Config
		app := fx.New(
			NewModule(confFile.Name()),

			fx.Invoke(func(c Config) {
				conf = c
			}),
		)

		So(test.StartFxApp(app), ShouldBeNil)

		Convey("When get config by struct", func() {
			testConf := TestConf{}
			err := conf.Get("code-shooting", &testConf)

			Convey("Then invokeResult should equal expect ", func() {
				expect := TestConf{
					Apiversion: "v1",
					Strategy: TestStrategy{
						Name: "RoundRobin",
						Type: "Recreate",
					},
				}

				So(err, ShouldBeNil)
				So(testConf, ShouldResemble, expect)
			})
		})
		Convey("When get config by map\n", func() {
			testConf := make(map[string]interface{})
			err := conf.Get("code-shooting", &testConf)

			Convey("Then map should equal expect \n", func() {
				expect := map[string]interface{}{
					"apiversion": "v1",
					"instrategy": map[interface{}]interface{}{
						"name": "RoundRobin",
						"type": "Recreate",
					},
				}

				So(err, ShouldBeNil)
				So(testConf, ShouldResemble, expect)
			})
		})

		Convey("When get exist config Item value \n", func() {
			v := conf.GetValue("code-shooting.apiversion")
			Convey("Then got config value \n", func() {
				So(v, ShouldNotBeNil)
				So(v, ShouldEqual, "v1")
			})
		})
		Convey("When get not exist config Item value \n", func() {
			v := conf.GetValue("code-shooting.test")
			Convey("Then got config value is nil \n", func() {
				So(v, ShouldBeNil)
			})
		})
		Convey("When register event handler \n", func() {
			var event *model.Event
			RegisterEventHandler("code-shooting", func(es []*model.Event) {
				if len(es) > 0 {
					event = es[0]
				}
			})
			ProcessConfigEvent([]*model.Event{})
			confFile.WriteString(`test`)
			confFile.Sync()
			time.Sleep(2 * time.Second)
			Convey("Then got config value is nil \n", func() {
				So(event, ShouldNotBeNil)
				So(event.Value, ShouldEqual, "Recreatetest")
			})
		})
	})
}

func TestConfigModuleComplexYam(t *testing.T) {

	Convey("Given a config file \n", t, func() {
		content := `
      code-shooting:
        app:
          name: mario
          version: v1
          scene: pict
        serviceCenter:
          msb:
            namespace: director
            transport:
              host: 127.0.0.1:10081
              basepath: /api/microservices/v1
              schemes:
              - http
        rest-servers:
        - name: mario
          addr: 0.0.0.0:8083
          readtimeout: 180s
          writetimeout: 180s
          maxheaderbytes: 16384
          rootpath: /api/v1.0/mario
        test:
          a.b: 1`

		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		var conf Config
		app := fx.New(
			fx.Provide(func() Config {
				config, _ := NewConfig(confFile.Name())
				return config
			}),
			fx.Invoke(func(c Config) {
				conf = c
			}),
		)

		So(test.StartFxApp(app), ShouldBeNil)

		Convey("When get the value with key of code-shooting \n", func() {
			testConf := make(map[string]interface{})
			err := conf.Get("code-shooting", &testConf)
			So(err, ShouldBeNil)
			Convey("Then invokeResult should equal expect \n", func() {
				yaml, err := uConfig.NewYAML(uConfig.Source(strings.NewReader(string(content))))

				So(err, ShouldBeNil)
				expect := make(map[string]interface{})
				value := yaml.Get("code-shooting")
				So(value.HasValue(), ShouldBeTrue)
				err = value.Populate(expect)
				So(err, ShouldBeNil)
				So(testConf, ShouldResemble, expect)
			})
		})

		Convey("When get the value with nil param \n", func() {
			var testConf interface{}
			err := conf.Get("code-shooting", testConf)

			Convey("Then invokeResult error \n", func() {

				So(err.Error(), ShouldContainSubstring, "invalid object supplied")
			})
		})
		Convey("When get by key the key has . \n", func() {

			v := conf.GetValue("code-shooting.test.a..b")

			Convey("Then should success \n", func() {

				So(v, ShouldNotBeNil)

			})
		})
	})
}
