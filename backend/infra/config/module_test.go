package config

import (
	"io/ioutil"
	"os"
	"testing"

	"code-shooting/infra/x/test"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"
)

type GoMarsConf struct {
	Apiversion string    `yaml:"apiversion,omitempty"`
	Strategy   StrategyA `yaml:"instrategy"`
}
type StrategyA struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
}

func TestConfigModule(t *testing.T) {
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

		var conf Config
		app := fx.New(
			NewModule(confFile.Name()),
			fx.Invoke(func(c Config) {
				conf = c
			}),
		)

		So(test.StartFxApp(app), ShouldBeNil)
		defer test.StopFxApp(app)

		Convey("When get the value with key of code-shooting", func() {
			gomarsConf := GoMarsConf{}
			err := conf.Get("code-shooting", &gomarsConf)

			Convey("Then invokeResult should equal expect ", func() {
				expect := GoMarsConf{
					Apiversion: "v1",
					Strategy: StrategyA{
						Name: "RoundRobin",
						Type: "Recreate",
					},
				}

				So(err, ShouldBeNil)
				So(gomarsConf, ShouldResemble, expect)
			})
		})
	})

	Convey("Given a empty config path ", t, func() {
		gomarsConf := GoMarsConf{}
		app := fx.New(
			NewModule(""),
			fx.Invoke(func(c Config) {
				c.Get("go.mars", &gomarsConf)
			}),
		)

		Convey("When start fx app", func() {
			err := test.StartFxApp(app)
			Convey("Then an err contains cannot unmarshal should return", func() {
				So(err.Error(), ShouldContainSubstring, "config path is empty")
			})
		})
	})
}
