package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"syscall"
	"testing"
	"time"

	. "github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
)

type FileStat struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	sys     syscall.Stat_t
}

func (fs *FileStat) IsDir() bool        { return fs.mode&os.ModeDir != 0 }
func (fs *FileStat) Name() string       { return fs.name }
func (fs *FileStat) Size() int64        { return fs.size }
func (fs *FileStat) Mode() os.FileMode  { return fs.mode }
func (fs *FileStat) ModTime() time.Time { return fs.modTime }
func (fs *FileStat) Sys() interface{}   { return &fs.sys }

func TestEnvConfig(t *testing.T) {

	file := &FileStat{}
	patches := ApplyFunc(os.Stat, func(_ string) (os.FileInfo, error) {
		return file, nil
	})
	defer patches.Reset()

	patches.ApplyMethod(reflect.TypeOf(file), "IsDir", func(_ *FileStat) bool {
		return false
	})

	Convey("Given a yaml config file including environment variables in ${VAR} format ", t, func() {

		content := `
service:
  kafka:
    address: ${kafkaIp}
`
		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		Convey("When the environment variables is configured and get the value with the key of service.kafka", func() {

			os.Setenv("kafkaIp", "127.0.0.1")
			defer os.Unsetenv("kafkaIp")

			configurator, err := NewConfig(confFile.Name())
			So(err, ShouldBeNil)

			var kafka = struct {
				Address string
			}{}
			err = configurator.Get("service.kafka", &kafka)

			Convey("Then the value should be 127.0.0.1", func() {
				So(err, ShouldEqual, nil)
				So(kafka.Address, ShouldEqual, "127.0.0.1")
			})
		})

		Convey("When the environment variables is not configured ", func() {

			configurator, err := NewConfig(confFile.Name())

			Convey("Then an error contain 'default is empty' should return and the configurator should be nil ", func() {
				So(err.Error(), ShouldContainSubstring, "default is empty")
				So(configurator, ShouldBeNil)
			})
		})
	})

	Convey("Given a yaml config file including environment variables in $VAR format", t, func() {

		content := `
service:
  kafka:
    address: $kafkaIp
`
		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		Convey("When the environment variables is configured and get the value with the key of service.kafka", func() {

			os.Setenv("kafkaIp", "127.0.0.1")
			defer os.Unsetenv("kafkaIp")

			configurator, err := NewConfig(confFile.Name())
			So(err, ShouldBeNil)

			var kafka = struct {
				Address string
			}{}
			err = configurator.Get("service.kafka", &kafka)

			Convey("Then the value should be 127.0.0.1", func() {
				So(err, ShouldEqual, nil)
				So(kafka.Address, ShouldEqual, "127.0.0.1")
			})
		})
	})

	Convey("Given a yaml config file including environment variables in ${VAR:default} format", t, func() {

		content := `
service:
  kafka:
    address: ${kafkaIp:10.10.10.10}
`
		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		Convey("When the environment variables is configured and get the value with the key of service.kafka", func() {

			os.Setenv("kafkaIp", "127.0.0.1")
			defer os.Unsetenv("kafkaIp")

			configurator, err := NewConfig(confFile.Name())
			So(err, ShouldBeNil)

			kafka := struct {
				Address string
			}{}
			err = configurator.Get("service.kafka", &kafka)

			Convey("Then the value should be 127.0.0.1 ", func() {
				So(err, ShouldEqual, nil)
				So(kafka.Address, ShouldEqual, "127.0.0.1")
			})
		})

		Convey("When the environment variables is not configured and get the value with the key of service.kafka", func() {

			configurator, err := NewConfig(confFile.Name())
			So(err, ShouldBeNil)

			kafka := struct {
				Address string
			}{}
			err = configurator.Get("service.kafka", &kafka)

			Convey("Then the value should be default value:10.10.10.10", func() {
				So(err, ShouldEqual, nil)
				So(kafka.Address, ShouldEqual, "10.10.10.10")
			})
		})
	})

	Convey("Given a yaml config file including multiple environment variables in ${VAR1}:${VAR2} format", t, func() {

		content := `
service:
  kafka:
    address: ${kafkaIp}:${kafkaPort}
`
		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		Convey("When the environment variables is configured and get the value with the key of service.kafka", func() {

			os.Setenv("kafkaIp", "127.0.0.1")
			defer os.Unsetenv("kafkaIp")
			os.Setenv("kafkaPort", "1111")
			defer os.Unsetenv("kafkaPort")

			configurator, err := NewConfig(confFile.Name())
			So(err, ShouldBeNil)

			kafka := struct {
				Address string
			}{}
			err = configurator.Get("service.kafka", &kafka)

			Convey("Then the value should be 127.0.0.1:1111 ", func() {
				So(err, ShouldEqual, nil)
				So(kafka.Address, ShouldEqual, "127.0.0.1:1111")
			})
		})

		Convey("When not all the environment variables is configured ", func() {

			os.Setenv("kafkaIp", "127.0.0.1")
			defer os.Unsetenv("kafkaIp")

			configurator, err := NewConfig(confFile.Name())

			Convey("Then an error contained 'default is empty' should rerurn and the configurator should be nil", func() {
				So(err.Error(), ShouldContainSubstring, "default is empty")
				So(configurator, ShouldBeNil)
			})
		})
	})

	Convey("Given a yaml config file including multiple environment variables in $VAR1:$VAR2 format", t, func() {

		content := `
service:
  kafka:
    address: $kafkaIp:$kafkaPort
`
		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		Convey("When the environment variables is configured ", func() {

			os.Setenv("kafkaIp", "127.0.0.1")
			defer os.Unsetenv("kafkaIp")
			os.Setenv("kafkaPort", "1111")
			defer os.Unsetenv("kafkaPort")

			configurator, err := NewConfig(confFile.Name())
			So(err, ShouldBeNil)

			kafka := struct {
				Address string
			}{}
			err = configurator.Get("service.kafka", &kafka)

			Convey("Then the value should be 127.0.0.1:1111 ", func() {
				So(err, ShouldEqual, nil)
				So(kafka.Address, ShouldEqual, "127.0.0.1:1111")
			})
		})
	})

	Convey("Given a yaml config file including multiple environment variables in $(VAR1):${VAR2} format", t, func() {

		content := `
service:
  kafka:
    address: $(kafkaIp):${kafkaPort}
`
		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		Convey("When the environment variables is configured ", func() {

			os.Setenv("kafkaIp", "127.0.0.1")
			defer os.Unsetenv("kafkaIp")
			os.Setenv("kafkaPort", "1111")
			defer os.Unsetenv("kafkaPort")

			configurator, err := NewConfig(confFile.Name())
			So(err, ShouldBeNil)

			kafka := struct {
				Address string
			}{}
			err = configurator.Get("service.kafka", &kafka)

			Convey("Then the value should be $(kafkaIp):1111 ", func() {
				So(err, ShouldEqual, nil)
				So(kafka.Address, ShouldEqual, "$(kafkaIp):1111")
			})
		})
	})

	Convey("Given a yaml config file including multiple environment variables in ${VAR:default1}:${VAR2:default2} format", t, func() {

		content := `
service:
  kafka:
    address: ${kafkaIp:10.10.10.10}:${kafkaPort:1010}
`
		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		Convey("When the environment variables is configured ", func() {

			os.Setenv("kafkaIp", "127.0.0.1")
			defer os.Unsetenv("kafkaIp")
			os.Setenv("kafkaPort", "1111")
			defer os.Unsetenv("kafkaPort")

			configurator, err := NewConfig(confFile.Name())
			So(err, ShouldBeNil)

			kafka := struct {
				Address string
			}{}
			err = configurator.Get("service.kafka", &kafka)

			Convey("Then the value should be 127.0.0.1:1111 ", func() {
				So(err, ShouldEqual, nil)
				So(kafka.Address, ShouldEqual, "127.0.0.1:1111")
			})
		})

		Convey("When not all the environment variables is configured ", func() {

			os.Setenv("kafkaIp", "127.0.0.1")
			defer os.Unsetenv("kafkaIp")

			configurator, err := NewConfig(confFile.Name())
			So(err, ShouldBeNil)

			kafka := struct {
				Address string
			}{}
			err = configurator.Get("service.kafka", &kafka)

			Convey("Then the value should be 127.0.0.1:1010", func() {
				So(err, ShouldEqual, nil)
				So(kafka.Address, ShouldEqual, "127.0.0.1:1010")
			})
		})

		Convey("When none of the environment variables is configured ", func() {

			configurator, err := NewConfig(confFile.Name())
			So(err, ShouldBeNil)

			kafka := struct {
				Address string
			}{}
			err = configurator.Get("service.kafka", &kafka)

			Convey("Then the value should be 10.10.10.10:1010", func() {
				So(err, ShouldEqual, nil)
				So(kafka.Address, ShouldEqual, "10.10.10.10:1010")
			})
		})
	})

	Convey("Given a yaml config file including multiple environment variables in $VAR:default1:$VAR2:default2 format", t, func() {

		content := `
service:
  kafka:
    address: $kafkaIp:10.10.10.10:$kafkaPort:1010
`
		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		Convey("When the environment variables is configured ", func() {

			os.Setenv("kafkaIp", "127.0.0.1")
			defer os.Unsetenv("kafkaIp")
			os.Setenv("kafkaPort", "1111")
			defer os.Unsetenv("kafkaPort")

			configurator, err := NewConfig(confFile.Name())
			So(err, ShouldBeNil)

			kafka := struct {
				Address string
			}{}
			err = configurator.Get("service.kafka", &kafka)

			Convey("Then the value should be 127.0.0.1:10.10.10.10:1111:1010 ", func() {
				So(err, ShouldEqual, nil)
				So(kafka.Address, ShouldEqual, "127.0.0.1:10.10.10.10:1111:1010")
			})
		})

		Convey("When not all the environment variables is configured ", func() {

			os.Setenv("kafkaIp", "127.0.0.1")
			defer os.Unsetenv("kafkaIp")

			configurator, err := NewConfig(confFile.Name())

			Convey("Then err should contain 'default is empty' and value should be empty", func() {
				So(err.Error(), ShouldContainSubstring, "default is empty")
				So(configurator, ShouldBeNil)
			})
		})
	})
}
