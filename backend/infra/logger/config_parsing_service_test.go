package logger

import (
	"io/ioutil"
	"os"
	"testing"

	"code-shooting/infra/config"
	"code-shooting/infra/logger/internal"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParsingLoggerConfig(t *testing.T) {
	Convey("Given a json encoder config", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: json
    outputPaths:
    - /tmp/mars-json.log
    - /tmp/mars-another-json.log
    rotateConfig:
      maxSize: 1 # MB
      maxBackups: 2 #
      maxAge: 7 # days
      compress: false
      fileMode: "0600"`

		conf, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			conf.Close()
			os.Remove(conf.Name())
		}()

		conf.WriteString(content)
		conf.Sync()

		goMarsConf, _ := config.NewConfig(conf.Name())

		Convey("When parse logger config", func() {
			parsedConfig, err := getConfigParsingService().Parse(goMarsConf)
			Convey("Then parsed config should be expect", func() {
				So(err, ShouldBeNil)
				So(parsedConfig.Level, ShouldEqual, "info")
				So(parsedConfig.Encoder, ShouldEqual, internal.JsonEncoder)
				So(parsedConfig.OutputPaths, ShouldResemble, []string{"/tmp/mars-json.log", "/tmp/mars-another-json.log"})
				So(parsedConfig.RotateConfig, ShouldResemble, internal.RotateConfig{MaxSize: 1, MaxBackups: 2, MaxAge: 7, Compress: false, FileMode: "0600"})
			})
		})
	})

	Convey("Given a plain encoder config", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths: 
    - /tmp/mars-plain.log
    - /tmp/mars-another-plain.log
    rotateConfig:
      maxSize: 1 # MB
      maxBackups: 2 # 
      maxAge: 7 # days
      compress: false`

		conf, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			conf.Close()
			os.Remove(conf.Name())
		}()

		conf.WriteString(content)
		conf.Sync()

		goMarsConf, _ := config.NewConfig(conf.Name())

		Convey("When parse logger config", func() {
			parsedConfig, err := getConfigParsingService().Parse(goMarsConf)
			Convey("Then parsed config should be expect", func() {
				So(err, ShouldBeNil)
				So(parsedConfig.Level, ShouldEqual, "info")
				So(parsedConfig.Encoder, ShouldEqual, internal.PlainEncoder)
				So(parsedConfig.OutputPaths, ShouldResemble, []string{"/tmp/mars-plain.log", "/tmp/mars-another-plain.log"})
				So(parsedConfig.RotateConfig, ShouldResemble, internal.RotateConfig{MaxSize: 1, MaxBackups: 2, MaxAge: 7, Compress: false, FileMode: "0640"})
			})
		})
	})

	Convey("Given a invalid encoder config with unsupported encoder", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: xxx
    outputPaths: 
    - /tmp/mars-plain.log
    - /tmp/mars-another-plain.log
    rotateConfig:
      maxSize: 1 # MB
      maxBackups: 2 # 
      maxAge: 7 # days
      compress: false`

		conf, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			conf.Close()
			os.Remove(conf.Name())
		}()

		conf.WriteString(content)
		conf.Sync()

		goMarsConf, _ := config.NewConfig(conf.Name())

		Convey("When parse logger config", func() {
			_, err := getConfigParsingService().Parse(goMarsConf)
			Convey("Then parsed config should failed and show tips: unsupported encoder: xxx", func() {
				So(err.Error(), ShouldContainSubstring, "unsupported encoder: xxx")
			})
		})
	})

	Convey("Given a invalid encoder config with unexpect type of maxsize", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths: 
    - /tmp/mars-plain.log
    - /tmp/mars-another-plain.log
    rotateConfig:
      maxSize: 1MB # MB
      maxBackups: 2 # 
      maxAge: 7 # days
      compress: false`

		conf, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			conf.Close()
			os.Remove(conf.Name())
		}()

		conf.WriteString(content)
		conf.Sync()

		goMarsConf, _ := config.NewConfig(conf.Name())

		Convey("When parse logger config", func() {
			_, err := getConfigParsingService().Parse(goMarsConf)
			Convey("Then parsed config should failed and show tips: cannot unmarshal", func() {
				So(err.Error(), ShouldContainSubstring, "cannot unmarshal")
			})
		})
	})

	Convey("Given a invalid encoder config with empty output paths", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    rotateConfig:
      maxSize: 1 # MB
      maxBackups: 2 # 
      maxAge: 7 # days
      compress: false`

		conf, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			conf.Close()
			os.Remove(conf.Name())
		}()

		conf.WriteString(content)
		conf.Sync()

		goMarsConf, _ := config.NewConfig(conf.Name())

		Convey("When parse logger config", func() {
			_, err := getConfigParsingService().Parse(goMarsConf)
			Convey("Then parsed config should failed and show tips: output paths is empty", func() {
				So(err.Error(), ShouldContainSubstring, "output paths is empty")
			})
		})
	})

}
