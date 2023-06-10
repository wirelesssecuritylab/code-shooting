package logger

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"code-shooting/infra/config"
	"code-shooting/infra/x/test"

	"github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/fx"
)

func TestGlobalLogger(t *testing.T) {
	Convey("Given a redirect target file for stdout", t, func() {
		f, _ := ioutil.TempFile(".", "default-*.log")
		defer func() {
			f.Close()
			os.Remove(f.Name())
		}()

		Convey("When log messages with info level", func() {
			test.DoActionsInFileRedirectedContext(os.Stdout, f, func() {
				Named("mod_name").With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
					Debug("debug level: a", "b", 1, 2, 3)
				Named("mod_name").With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
					Debugf("debug level f: %s %s %d %d %d", "a", "b", 1, 2, 3)

				With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
					Info("info level: a", "b", 1, 2, 3)
				Named("mod_name").With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
					Info("info level: a", "b", 1, 2, 3)
				With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
					Infof("info level f: %s %s %d %d %d", "a", "b", 1, 2, 3)
				Named("mod_name").With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
					Infof("info level f: %s %s %d %d %d", "a", "b", 1, 2, 3)

				Sync()
			})

			Convey("Then the log contains only the info log", func() {
				content, _ := ioutil.ReadFile(f.Name())

				hostName, err := os.Hostname()
				if err != nil {
					hostName = "unknown"
				}
				timeFormat := "\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{3}(([+-]\\d{2}:\\d{2})|Z)"
				expect := timeFormat + "\tINFO\t" + hostName + "\t\t" +
					"info level: ab1 2 3\t{\"a\": \"b\", \"c\": 1, \"d\": \"1h0m0s\"}\t\\[out_of_the_box_log_test\\.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tINFO\t" + hostName + "\tmod_name\t" +
					"info level: ab1 2 3\t{\"a\": \"b\", \"c\": 1, \"d\": \"1h0m0s\"}\t\\[out_of_the_box_log_test\\.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tINFO\t" + hostName + "\t\t" +
					"info level f: a b 1 2 3\t{\"a\": \"b\", \"c\": 1, \"d\": \"1h0m0s\"}\t\\[out_of_the_box_log_test\\.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tINFO\t" + hostName + "\tmod_name\t" +
					"info level f: a b 1 2 3\t{\"a\": \"b\", \"c\": 1, \"d\": \"1h0m0s\"}\t\\[out_of_the_box_log_test\\.go\\]\\[\\d+\\]\n$"

				So(string(content), test.ShouldMatchWithRegex, expect)
			})
		})

		Convey("When query log level", func() {
			level := GetLevel()

			Convey("Then the level is info", func() {
				expect := "info"
				So(level, ShouldEqual, expect)
			})
		})
	})

	Convey("Given a new log level of debug", t, func() {
		oldLevel := GetLevel()
		defer func() {
			SetLevel(oldLevel)
		}()
		newLevel := "debug"
		SetLevel(newLevel)

		Convey("When query log level", func() {
			level := GetLevel()

			Convey("Then the level is debug", func() {
				expect := "debug"
				So(level, ShouldEqual, expect)
			})
		})
	})

	Convey("Given a logger of file sink", t, func() {
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
      fileMode: "0644"`

		conf, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			conf.Close()
			os.Remove(conf.Name())
		}()

		conf.WriteString(content)
		conf.Sync()

		removeFilesMatchWith("/tmp", "mars-json")
		removeFilesMatchWith("/tmp", "mars-another-json")
		defer func() {
			removeFilesMatchWith("/tmp", "mars-json")
			removeFilesMatchWith("/tmp", "mars-another-json")
		}()

		removeFilesMatchWith("/tmp", "mars-json")
		removeFilesMatchWith("/tmp", "mars-another-json")
		defer func() {
			removeFilesMatchWith("/tmp", "mars-json")
			removeFilesMatchWith("/tmp", "mars-another-json")
		}()

		logger, err := NewLogger(conf.Name())
		So(err, ShouldBeNil)

		old := GetLogger()
		defer func() {
			SetLogger(old)
		}()

		SetLogger(logger)

		Convey("When log messages with dynamic adjust level", func() {
			Named("mod_name").With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				Debug("debug level: a", "b", 1, 2, 3)
			Named("mod_name").With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				Debugf("debug level f: %s %s %d %d %d", "a", "b", 1, 2, 3)
			Named("mod_name").With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				Info("info level: a", "b", 1, 2, 3)
			Named("mod_name").AddCallerSkip(1).With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				Infof("info level f: %s %s %d %d %d", "a", "b", 1, 2, 3)

			Sync()

			Convey("Then the log content should be expect(discard lower level log)", func() {
				timeFormat := "\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{3}(([+-]\\d{4})|Z)"
				expect := "{\"L\":\"INFO\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/out_of_the_box_log_test\\.go:\\d+\"," +
					"\"M\":\"info level: ab1 2 3\",\"a\":\"b\",\"c\":1,\"d\":\"1h0m0s\"}\n" +

					"{\"L\":\"INFO\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"convey/discovery\\.go:\\d+\"," +
					"\"M\":\"info level f: a b 1 2 3\",\"a\":\"b\",\"c\":1,\"d\":\"1h0m0s\"}\n$"

				for _, filePath := range []string{"/tmp/mars-json.log", "/tmp/mars-another-json.log"} {
					content, err := ioutil.ReadFile(filePath)
					So(err, ShouldBeNil)
					So(string(content), test.ShouldMatchWithRegex, expect)
				}
			})
		})
	})

	Convey("Given a logger direct by config path(main scene of fx application)", t, func() {
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
      fileMode: "0644"`

		conf, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			conf.Close()
			os.Remove(conf.Name())
		}()

		conf.WriteString(content)
		conf.Sync()

		removeFilesMatchWith("/tmp", "mars-json")
		removeFilesMatchWith("/tmp", "mars-another-json")
		defer func() {
			removeFilesMatchWith("/tmp", "mars-json")
			removeFilesMatchWith("/tmp", "mars-another-json")
		}()

		removeFilesMatchWith("/tmp", "mars-json")
		removeFilesMatchWith("/tmp", "mars-another-json")
		defer func() {
			removeFilesMatchWith("/tmp", "mars-json")
			removeFilesMatchWith("/tmp", "mars-another-json")
		}()

		patches := gomonkey.ApplyFunc(os.Exit, func(code int) {
			panic(code)
		})
		defer patches.Reset()

		logger, err := NewLogger(conf.Name())
		So(err, ShouldBeNil)
		logger = logger.Named("app")

		old := GetLogger()
		defer func() {
			SetLogger(old)
		}()

		SetLogger(logger)

		Convey("When use fx start application", func() {
			app := fx.New(
				fx.Logger(GetLogger().CreateStdLogger()),
				fx.Invoke(func() {
					logger.Info("Hello")
				}))
			test.StartFxApp(app)
			defer test.StopFxApp(app)

			Convey("And log app messages", func() {
				Named("mod_name").With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
					Debug("debug level: a", "b", 1, 2, 3)
				Named("mod_name").With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
					Debugf("debug level f: %s %s %d %d %d", "a", "b", 1, 2, 3)
				Named("mod_name").With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
					Info("info level: a", "b", 1, 2, 3)
				Named("mod_name").AddCallerSkip(1).With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
					Infof("info level f: %s %s %d %d %d", "a", "b", 1, 2, 3)

				Debug("debug level: ", 1)
				Debugf("debug level f: %d", 1)
				Info("info level: ", 1)
				Infof("info level f: %d", 1)
				Warn("warn level: ", 1)
				Warnf("warn level f: %d", 1)
				Error("error level: ", 1)
				Errorf("error level f: %d", 1)

				So(func() {
					Panic("panic level: ", 1)
				}, ShouldPanic)
				So(func() {
					Panicf("panic level f: %d", 1)
				}, ShouldPanic)
				So(func() {
					Fatal("fatal level: ", 1)
				}, ShouldPanic)
				So(func() {
					Fatalf("fatal level f: %d", 1)
				}, ShouldPanic)

				Sync()

				Convey("Then the log content should be expect(discard lower level log)", func() {
					timeFormat := "\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{3}(([+-]\\d{4})|Z)"
					expect := "{\"L\":\"INFO\",\"T\":\"" + timeFormat + "\",\"N\":\"app\"," +
						"\"C\":\"fx[^\"]*\\.go:\\d+\"," +
						"\"M\":\"\\[Fx\\]\\s+PROVIDE[^\"]+\"}\n" +

						"({\"L\":\"INFO\",\"T\":\"" + timeFormat + "\",\"N\":\"app\"," +
						"\"C\":\"[^\"]*\\.go:\\d+\"," +
						"\"M\":\"[^\"]+\"}\n)*" +

						"{\"L\":\"INFO\",\"T\":\"" + timeFormat + "\",\"N\":\"app\"," +
						"\"C\":\"logger/out_of_the_box_log_test\\.go:\\d+\"," +
						"\"M\":\"Hello\"}\n" +

						"({\"L\":\"INFO\",\"T\":\"" + timeFormat + "\",\"N\":\"app\"," +
						"\"C\":\"[^\"]*\\.go:\\d+\"," +
						"\"M\":\"[^\"]+\"}\n)*" +

						"{\"L\":\"INFO\",\"T\":\"" + timeFormat + "\",\"N\":\"app\"," +
						"\"C\":\"fx[^\"]*\\.go:\\d+\"," +
						"\"M\":\"\\[Fx\\]\\s+RUNNING\"}\n" +

						"{\"L\":\"INFO\",\"T\":\"" + timeFormat + "\",\"N\":\"app\\.mod_name\"," +
						"\"C\":\"logger/out_of_the_box_log_test\\.go:\\d+\"," +
						"\"M\":\"info level: ab1 2 3\",\"a\":\"b\",\"c\":1,\"d\":\"1h0m0s\"}\n" +

						"{\"L\":\"INFO\",\"T\":\"" + timeFormat + "\",\"N\":\"app\\.mod_name\"," +
						"\"C\":\"convey/discovery\\.go:\\d+\"," +
						"\"M\":\"info level f: a b 1 2 3\",\"a\":\"b\",\"c\":1,\"d\":\"1h0m0s\"}\n" +

						"{\"L\":\"INFO\",\"T\":\"" + timeFormat + "\",\"N\":\"app\"," +
						"\"C\":\"logger/out_of_the_box_log_test\\.go:\\d+\"," +
						"\"M\":\"info level: 1\"}\n" +

						"{\"L\":\"INFO\",\"T\":\"" + timeFormat + "\",\"N\":\"app\"," +
						"\"C\":\"logger/out_of_the_box_log_test\\.go:\\d+\"," +
						"\"M\":\"info level f: 1\"}\n" +

						"{\"L\":\"WARN\",\"T\":\"" + timeFormat + "\",\"N\":\"app\"," +
						"\"C\":\"logger/out_of_the_box_log_test\\.go:\\d+\"," +
						"\"M\":\"warn level: 1\"}\n" +

						"{\"L\":\"WARN\",\"T\":\"" + timeFormat + "\",\"N\":\"app\"," +
						"\"C\":\"logger/out_of_the_box_log_test\\.go:\\d+\"," +
						"\"M\":\"warn level f: 1\"}\n" +

						"{\"L\":\"ERROR\",\"T\":\"" + timeFormat + "\",\"N\":\"app\"," +
						"\"C\":\"logger/out_of_the_box_log_test\\.go:\\d+\"," +
						"\"M\":\"error level: 1\"}\n" +

						"{\"L\":\"ERROR\",\"T\":\"" + timeFormat + "\",\"N\":\"app\"," +
						"\"C\":\"logger/out_of_the_box_log_test\\.go:\\d+\"," +
						"\"M\":\"error level f: 1\"}\n" +

						"{\"L\":\"PANIC\",\"T\":\"" + timeFormat + "\",\"N\":\"app\"," +
						"\"C\":\"logger/out_of_the_box_log_test\\.go:\\d+\"," +
						"\"M\":\"panic level: 1\"}\n" +

						"{\"L\":\"PANIC\",\"T\":\"" + timeFormat + "\",\"N\":\"app\"," +
						"\"C\":\"logger/out_of_the_box_log_test\\.go:\\d+\"," +
						"\"M\":\"panic level f: 1\"}\n" +

						"{\"L\":\"FATAL\",\"T\":\"" + timeFormat + "\",\"N\":\"app\"," +
						"\"C\":\"logger/out_of_the_box_log_test\\.go:\\d+\"," +
						"\"M\":\"fatal level: 1\"}\n" +

						"{\"L\":\"FATAL\",\"T\":\"" + timeFormat + "\",\"N\":\"app\"," +
						"\"C\":\"logger/out_of_the_box_log_test\\.go:\\d+\"," +
						"\"M\":\"fatal level f: 1\"}\n$"

					for _, filePath := range []string{"/tmp/mars-json.log", "/tmp/mars-another-json.log"} {
						content, err := ioutil.ReadFile(filePath)
						So(err, ShouldBeNil)
						So(string(content), test.ShouldMatchWithRegex, expect)
					}
				})
			})
		})
	})
}

func TestLogInFx(t *testing.T) {
	Convey("Given a logger and config \n", t, func() {
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
      fileMode: "0644"`

		conf, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			conf.Close()
			os.Remove(conf.Name())
		}()

		conf.WriteString(content)
		conf.Sync()

		removeFilesMatchWith("/tmp", "mars-json")

		defer func() {
			removeFilesMatchWith("/tmp", "mars-json")
		}()

		logger, err := NewLogger(conf.Name())
		So(err, ShouldBeNil)
		SetLogger(logger)
		config.SetLogger(logger)
		app := fx.New(
			config.NewModule(conf.Name()),
			fx.Invoke(func(c config.Config) {

			}),
		)
		So(test.StartFxApp(app), ShouldBeNil)
		defer test.StartFxApp(app)

		Convey("When log messages with dynamic adjust level \n", func() {
			SetLevel("Debug")
			defer SetLevel("Info")
			Debug("test debug log")
			Convey("Then the log content should be expect(discard lower level log) \n", func() {
				content, err := ioutil.ReadFile("/tmp/mars-json.log")
				So(err, ShouldBeNil)
				So(string(content), ShouldContainSubstring, "test debug log")
				So(string(content), ShouldContainSubstring, "DEBUG")
			})
		})
	})
}
