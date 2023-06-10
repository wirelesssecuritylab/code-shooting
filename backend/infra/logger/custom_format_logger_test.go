package logger

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"code-shooting/infra/x/test"

	"github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomFormatLogger(t *testing.T) {
	Convey("Given a custom format(plog style) logger(with plain encoder) config ", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    format: "$${T}\t$${L}\t$${H}\t$${N}\tTransactionID=$${TransactionID:null}\tInstanceID=$${InstanceID:null}\t[ObjectID=$${ObjectID:null},ObjectType=$${ObjectType:null}]\t$${M}\t[$${C_F}]:[$${C_L}]\n"
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

		logger, err := NewLogger(conf.Name())
		So(err, ShouldBeNil)

		patches := gomonkey.ApplyFunc(os.Exit, func(code int) {
			panic(code)
		})
		defer patches.Reset()

		Convey("When log messages with dynamic adjust level", func() {
			logger = logger.Named("mod_name")
			logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				Debug("debug level: a", "b", 1, 2, 3)
			logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				Debugf("debug level f: %s %s %d %d %d", "a", "b", 1, 2, 3)
			logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				With(StringField("TransactionID", "tid_1"), IntField("InstanceID", 1)).
				With(StringField("ObjectID", "obj_1"), StringField("ObjectType", "obj_t1")).
				Info("info level: a", "b", 1, 2, 3)
			logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				With(StringField("TransactionID", "tid_1"), IntField("InstanceID", 1)).
				With(StringField("ObjectID", "obj_1"), StringField("ObjectType", "obj_t1")).
				Infof("info level f: %s %s %d %d %d", "a", "b", 1, 2, 3)
			logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				Warn("warn level: a", "b", 1, 2, 3)
			logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				Warnf("warn level f: %s %s %d %d %d", "a", "b", 1, 2, 3)
			logger.With(StringField("s", "a")).With(IntField("b", 1)).
				With(DurationField("TransactionID", 1*time.Hour+1*time.Minute+1*time.Second), IntField("InstanceID", 1)).
				With(StringField("ObjectID", "obj_1"), StringField("ObjectType", "obj_t1")).
				Error("error level: a", 1, 2)
			logger.With(StringField("s", "a")).With(DurationField("d", time.Duration(1*time.Minute))).
				Errorf("error level f: %s %d %d", "a", 1, 2)
			So(func() {
				logger.With(StringField("a", "b")).
					Panic("panic level: a", "b", 1, 2)
			}, ShouldPanic)
			So(func() {
				logger.With(StringField("a", "b")).
					Panicf("panic level f: %s %s %d %d", "a", "b", 1, 2)
			}, ShouldPanic)
			So(func() {
				logger.With(StringField("f", "b")).
					Fatal("fatal level: a", "b", 1, 2)
			}, ShouldPanic)
			So(func() {
				logger.With(StringField("f", "b")).
					Fatalf("fatal level f: %s %s %d %d", "a", "b", 1, 2)
			}, ShouldPanic)

			logger.SetLevel("Warn")
			logger.With(StringField("l", "i")).
				Info("info level 1: a", "b", 1, 2)
			logger.With(StringField("l", "i")).
				Infof("info level f 1: %s %s %d %d", "a", "b", 1, 2)
			logger.With(StringField("l", "w")).
				Warn("warn level 1: a", "b", 1, 2)
			logger.With(StringField("l", "w")).
				Warnf("warn level f 1: %s %s %d %d", "a", "b", 1, 2)

			logger.SetLevel("Error")
			logger.With(StringField("l", "w")).
				Warn("warn level 2: a", "b", 1, 2)
			logger.With(StringField("l", "w")).
				Warnf("warn level f 1: %s %s %d %d", "a", "b", 1, 2)
			logger.With(StringField("l", "e")).
				Error("error level 2: a", "b", 1, 2)
			logger.With(StringField("l", "e")).
				Errorf("error level f 2: %s %s %d %d", "a", "b", 1, 2)

			logger.SetLevel("Panic")
			logger.With(StringField("l", "e")).
				Error("error level 3: a", "b", 1, 2)
			logger.With(StringField("l", "e")).
				Errorf("error level f 3: %s %s %d %d", "a", "b", 1, 2)
			So(func() {
				logger.With(StringField("l", "p")).
					Panic("panic level 3: a", "b", 1, 2)
			}, ShouldPanic)
			So(func() {
				logger.With(StringField("l", "p")).
					Panicf("panic level f 3: %s %s %d %d", "a", "b", 1, 2)
			}, ShouldPanic)

			logger.SetLevel("Fatal")
			logger.With(StringField("l", "p")).
				Panic("panic level 4: a", "b", 1, 2)
			logger.With(StringField("l", "p")).
				Panicf("panic level f 4: %s %s %d %d", "a", "b", 1, 2)
			So(func() {
				logger.With(StringField("l", "f")).
					Fatal("fatal level 4: a", "b", 1, 2)
			}, ShouldPanic)
			So(func() {
				logger.With(StringField("l", "f")).
					Fatalf("fatal level f 4: %s %s %d %d", "a", "b", 1, 2)
			}, ShouldPanic)

			logger.SetLevel("Debug")
			logger.With(StringField("l", "d")).
				Debug("debug level 5: a", "b", 1, 2)
			logger.With(StringField("l", "d")).
				Debugf("debug level f 5: %s %s %d %d", "a", "b", 1, 2)

			logger.Sync()

			Convey("Then the log content should be expect(discard lower level log)", func() {
				hostName, err := os.Hostname()
				if err != nil {
					hostName = "unknown"
				}
				timeFormat := "\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{3}(([+-]\\d{2}:\\d{2})|Z)"

				expect := timeFormat + "\tINFO\t" + hostName + "\tmod_name\t" +
					"TransactionID=tid_1\tInstanceID=1\t\\[ObjectID=obj_1,ObjectType=obj_t1\\]\t" +
					"info level: ab1 2 3\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tINFO\t" + hostName + "\tmod_name\t" +
					"TransactionID=tid_1\tInstanceID=1\t\\[ObjectID=obj_1,ObjectType=obj_t1\\]\t" +
					"info level f: a b 1 2 3\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tWARN\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"warn level: ab1 2 3\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tWARN\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"warn level f: a b 1 2 3\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tERROR\t" + hostName + "\tmod_name\t" +
					"TransactionID=1h1m1s\tInstanceID=1\t\\[ObjectID=obj_1,ObjectType=obj_t1\\]\t" +
					"error level: a1 2\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tERROR\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"error level f: a 1 2\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tPANIC\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"panic level: ab1 2\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tPANIC\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"panic level f: a b 1 2\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tFATAL\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"fatal level: ab1 2\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tFATAL\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"fatal level f: a b 1 2\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tWARN\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"warn level 1: ab1 2\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tWARN\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"warn level f 1: a b 1 2\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tERROR\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"error level 2: ab1 2\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tERROR\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"error level f 2: a b 1 2\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tPANIC\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"panic level 3: ab1 2\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tPANIC\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"panic level f 3: a b 1 2\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tFATAL\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"fatal level 4: ab1 2\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tFATAL\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"fatal level f 4: a b 1 2\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tDEBUG\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"debug level 5: ab1 2\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tDEBUG\t" + hostName + "\tmod_name\t" +
					"TransactionID=null\tInstanceID=null\t\\[ObjectID=null,ObjectType=null\\]\t" +
					"debug level f 5: a b 1 2\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n$"

				for _, filePath := range []string{"/tmp/mars-json.log", "/tmp/mars-another-json.log"} {
					info, statErr := os.Stat(filePath)
					So(statErr, ShouldBeNil)
					So(info.Mode(), ShouldEqual, os.FileMode(0644))
					content, readErr := ioutil.ReadFile(filePath)
					So(readErr, ShouldBeNil)
					So(string(content), test.ShouldMatchWithRegex, expect)
				}
			})
		})
	})

	Convey("Given a custom format(like plog style but has extends fields) logger(with plain encoder) config ", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain 
    format: "$${T}\t$${L}\t$${H}\t$${N}\t$${M} $${E}\t[$${C_F}]:[$${C_L}]\n"
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

		logger, err := NewLogger(conf.Name())
		So(err, ShouldBeNil)

		patches := gomonkey.ApplyFunc(os.Exit, func(code int) {
			panic(code)
		})
		defer patches.Reset()

		Convey("When log messages", func() {
			logger = logger.Named("mod_name")
			logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				With(StringField("TransactionID", "tid_1"), IntField("InstanceID", 1)).
				With(StringField("ObjectID", "obj_1"), StringField("ObjectType", "obj_t1")).
				Info("info level: a", "b", 1, 2, 3)
			logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				With(StringField("TransactionID", "tid_1"), IntField("InstanceID", 1)).
				With(StringField("ObjectID", "obj_1"), StringField("ObjectType", "obj_t1")).
				Infof("info level f: %s %s %d %d %d", "a", "b", 1, 2, 3)

			logger.Sync()

			Convey("Then the log content should be expect(discard lower level log)", func() {
				hostName, err := os.Hostname()
				if err != nil {
					hostName = "unknown"
				}
				timeFormat := "\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{3}(([+-]\\d{2}:\\d{2})|Z)"

				expect := timeFormat + "\tINFO\t" + hostName + "\tmod_name\t" +
					"info level: ab1 2 3 {\"a\": \"b\", \"c\": 1, \"d\": \"1h0m0s\", " +
					"\"TransactionID\": \"tid_1\", \"InstanceID\": 1, \"ObjectID\": \"obj_1\", \"ObjectType\": \"obj_t1\"}" +
					"\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tINFO\t" + hostName + "\tmod_name\t" +
					"info level f: a b 1 2 3 {\"a\": \"b\", \"c\": 1, \"d\": \"1h0m0s\", " +
					"\"TransactionID\": \"tid_1\", \"InstanceID\": 1, \"ObjectID\": \"obj_1\", \"ObjectType\": \"obj_t1\"}" +
					"\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n$"

				for _, filePath := range []string{"/tmp/mars-json.log", "/tmp/mars-another-json.log"} {
					content, err := ioutil.ReadFile(filePath)
					So(err, ShouldBeNil)
					So(string(content), test.ShouldMatchWithRegex, expect)
				}
			})
		})
	})

	Convey("Given a logger config(with rotate config) ", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain 
    outputPaths: 
    - /tmp/plain-rotate.log
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

		removeFilesMatchWith("/tmp", "plain-rotate")
		defer func() {
			removeFilesMatchWith("/tmp", "plain-rotate")
		}()

		logger, err := NewLogger(conf.Name())
		So(err, ShouldBeNil)

		patches := gomonkey.ApplyFunc(os.Exit, func(code int) {
			panic(code)
		})
		defer patches.Reset()

		Convey("When log message larger than 2MB", func() {
			msg := ""
			for i := 0; i <= 1024; i++ {
				msg = msg + strconv.Itoa(i%10)
			}

			logger = logger.Named("mod_name")
			for i := 0; i <= 2*1024; i++ {
				logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
					With(StringField("TransactionID", "tid_1"), IntField("InstanceID", 1)).
					With(StringField("ObjectID", "obj_1"), StringField("ObjectType", "obj_t1")).
					Info(msg)
			}
			logger.Sync()

			Convey("Then the log file number is 3, and file name format is plain-rotate[-2020-11-13T02-23-03.396].log", func() {
				files, err := listFilesMatchWith("/tmp", "plain-rotate")
				So(err, ShouldBeNil)
				So(files, ShouldHaveLength, 3)
				for _, file := range files {
					So(file, test.ShouldMatchWithRegex, "^/tmp/plain-rotate(-\\d{4}-\\d{2}-\\d{2}T\\d{2}-\\d{2}-\\d{2}\\.\\d{3})?\\.log$")
				}
			})
		})
	})

	Convey("Given a logger config(without rotate config) ", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain 
    outputPaths: 
    - /tmp/plain-rotate.log`

		conf, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			conf.Close()
			os.Remove(conf.Name())
		}()

		conf.WriteString(content)
		conf.Sync()

		removeFilesMatchWith("/tmp", "plain-rotate")
		defer func() {
			removeFilesMatchWith("/tmp", "plain-rotate")
		}()

		logger, err := NewLogger(conf.Name())
		So(err, ShouldBeNil)

		patches := gomonkey.ApplyFunc(os.Exit, func(code int) {
			panic(code)
		})
		defer patches.Reset()

		Convey("When log message larger than 20MB", func() {
			msg := ""
			for i := 0; i <= 1024; i++ {
				msg = msg + strconv.Itoa(i%10)
			}

			logger = logger.Named("mod_name")
			for i := 0; i <= 20*1024; i++ {
				logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
					With(StringField("TransactionID", "tid_1"), IntField("InstanceID", 1)).
					With(StringField("ObjectID", "obj_1"), StringField("ObjectType", "obj_t1")).
					Info(msg)
			}
			logger.Sync()

			Convey("Then the log file number is 1", func() {
				files, err := listFilesMatchWith("/tmp", "plain-rotate")
				So(err, ShouldBeNil)
				So(files, ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given a custom format config and use options", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
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

		logger, err := NewLogger(conf.Name(),
			CustomFormat("${T}\t${L}\t${H}\t${N}\t${M} ${E}\t[${C_F}]:[${C_L}]\n"),
			TimeEncoder("2006-01-02 15:04:05.000"),
			AddCallerSkip(1), AddCallerSkip(-1),
			CapitalLevelEncoder(map[string]string{
				"Warn": "Warning",
			}))
		So(err, ShouldBeNil)

		patches := gomonkey.ApplyFunc(os.Exit, func(code int) {
			panic(code)
		})
		defer patches.Reset()

		Convey("When log messages", func() {
			logger = logger.Named("mod_name")
			logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				With(StringField("TransactionID", "tid_1"), IntField("InstanceID", 1)).
				With(StringField("ObjectID", "obj_1"), StringField("ObjectType", "obj_t1")).
				Info("info level: a", "b", 1, 2, 3)
			logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				With(StringField("TransactionID", "tid_1"), IntField("InstanceID", 1)).
				With(StringField("ObjectID", "obj_1"), StringField("ObjectType", "obj_t1")).
				Infof("info level f: %s %s %d %d %d", "a", "b", 1, 2, 3)

			logger.Warn("warning level msg")

			logger.Sync()

			Convey("Then the log content should be expect(discard lower level log)", func() {
				hostName, err := os.Hostname()
				if err != nil {
					hostName = "unknown"
				}
				timeFormat := "\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\.\\d{3}"

				expect := timeFormat + "\tINFO\t" + hostName + "\tmod_name\t" +
					"info level: ab1 2 3 {\"a\": \"b\", \"c\": 1, \"d\": \"1h0m0s\", " +
					"\"TransactionID\": \"tid_1\", \"InstanceID\": 1, \"ObjectID\": \"obj_1\", \"ObjectType\": \"obj_t1\"}" +
					"\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tINFO\t" + hostName + "\tmod_name\t" +
					"info level f: a b 1 2 3 {\"a\": \"b\", \"c\": 1, \"d\": \"1h0m0s\", " +
					"\"TransactionID\": \"tid_1\", \"InstanceID\": 1, \"ObjectID\": \"obj_1\", \"ObjectType\": \"obj_t1\"}" +
					"\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n" +

					timeFormat + "\tWARNING\t" + hostName + "\tmod_name\t" +
					"warning level msg " +
					"\t\\[custom_format_logger_test\\.go\\]:\\[\\d+\\]\n$"

				for _, filePath := range []string{"/tmp/mars-json.log", "/tmp/mars-another-json.log"} {
					content, err := ioutil.ReadFile(filePath)
					So(err, ShouldBeNil)
					So(string(content), test.ShouldMatchWithRegex, expect)
				}
			})
		})
	})
}
