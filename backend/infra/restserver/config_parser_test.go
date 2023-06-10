package restserver

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	marsconfig "code-shooting/infra/config"
)

func TestConfigParser(t *testing.T) {
	Convey("Given no code-shooting config", t, func() {

		var goMarsConf marsconfig.Config

		Convey("When parse restserver config", func() {
			_, err := GetConfigParser().Parse(goMarsConf)
			Convey("Then err should not be nil and contains 'config does not exist'", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "config does not exist")
			})
		})
	})
	Convey("Given a no-rest-server config file", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: json
    outputPaths: 
    - stdout
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

		goMarsConf, _ := marsconfig.NewConfig(conf.Name())

		Convey("When parse restserver config", func() {
			_, err := GetConfigParser().Parse(goMarsConf)
			Convey("Then err should not be nil and contains 'the key does not exist'", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "the key does not exist")
			})
		})
	})

	Convey("Given a two-rest-server with the same name config file ", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: json
    outputPaths: 
    - stdout
  rest-servers:
  - name: myserver
    addr: 127.0.0.1:9081
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384
    rootpath: /mypath
    middlewares:
    - name: stats
  - name: myserver
    addr: 127.0.0.1:9082
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384
    rootpath: /urpath
    middlewares:
    - name: stats`

		conf, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			conf.Close()
			os.Remove(conf.Name())
		}()

		conf.WriteString(content)
		conf.Sync()

		goMarsConf, _ := marsconfig.NewConfig(conf.Name())

		Convey("When parse restserver config", func() {
			_, err := GetConfigParser().Parse(goMarsConf)
			Convey("Then err should not be nil and contains 'name myserver is duplicated'", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "name myserver is duplicated")
			})
		})
	})

	Convey("Given a restserver config file with error addr fomart", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: json
    outputPaths: 
    - stdout
  rest-servers:
  - name: myserver
    addr: 127.0.0.1,9081
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384
    rootpath: /mypath
    middlewares:
    - name: stats
`
		conf, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			conf.Close()
			os.Remove(conf.Name())
		}()
		conf.WriteString(content)
		conf.Sync()
		goMarsConf, _ := marsconfig.NewConfig(conf.Name())
		Convey("When parse restserver config", func() {
			_, err := GetConfigParser().Parse(goMarsConf)
			Convey("Then should return addr format error '", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "split host port from address")
			})
		})
	})

	Convey("Given a restserver config file with error ip fomart", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: json
    outputPaths: 
    - stdout
  rest-servers:
  - name: myserver
    addr: 1270.0.1:9081
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384
    rootpath: /mypath
    middlewares:
    - name: stats
`
		conf, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			conf.Close()
			os.Remove(conf.Name())
		}()
		conf.WriteString(content)
		conf.Sync()
		goMarsConf, _ := marsconfig.NewConfig(conf.Name())
		Convey("When parse restserver config", func() {
			_, err := GetConfigParser().Parse(goMarsConf)
			Convey("Then should return ip format error '", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "not a valid textual representation of an IP address")
			})
		})
	})
}
