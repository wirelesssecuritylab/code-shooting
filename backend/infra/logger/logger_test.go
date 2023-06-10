package logger

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"code-shooting/infra/common"
	"code-shooting/infra/config"
	"code-shooting/infra/config/model"
	"code-shooting/infra/x/test"

	"github.com/agiledragon/gomonkey"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func listFilesMatchWith(dir, fileNamePrefix string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filePaths []string
	for _, file := range files {
		if !file.IsDir() {
			matched, _ := regexp.MatchString("^"+fileNamePrefix, file.Name())
			if matched {
				filePaths = append(filePaths, filepath.Join(dir, file.Name()))
			}
		}
	}

	return filePaths, nil
}

func removeFilesMatchWith(dir, fileNamePrefix string) {
	files, _ := listFilesMatchWith(dir, fileNamePrefix)
	for _, file := range files {
		os.Remove(file)
	}
}

type tmpStringer struct {
	str string
}

func (s *tmpStringer) String() string {
	return s.str
}

type tmpMarshaler []int

func (s tmpMarshaler) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for _, v := range s {
		arr.AppendInt(v)
	}
	return nil
}

type tmpObjectMarshaler struct {
	ID    string
	Index uint
}

func (s tmpObjectMarshaler) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("ID", s.ID)
	enc.AddUint("Index", s.Index)
	return errors.New("marshal log object of tmp")
}

func newInvalidField() zap.Field {
	field := zap.String("invalid", "")
	field.Type = zapcore.FieldType(255)
	return field
}

func TestGoMarsNLogger(t *testing.T) {
	Convey("Given a logger with json encoder config", t, func() {
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
				Info("info level: a", "b", 1, 2, 3)
			logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				Infof("info level f: %s %s %d %d %d", "a", "b", 1, 2, 3)
			logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				Warn("warn level: a", "b", 1, 2, 3)
			logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				Warnf("warn level f: %s %s %d %d %d", "a", "b", 1, 2, 3)
			logger.With(StringField("s", "a")).With(IntField("b", 1)).
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
				timeFormat := "\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{3}(([+-]\\d{4})|Z)"
				expect := "{\"L\":\"INFO\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"info level: ab1 2 3\",\"a\":\"b\",\"c\":1,\"d\":\"1h0m0s\"}\n" +

					"{\"L\":\"INFO\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"info level f: a b 1 2 3\",\"a\":\"b\",\"c\":1,\"d\":\"1h0m0s\"}\n" +

					"{\"L\":\"WARN\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"warn level: ab1 2 3\",\"a\":\"b\",\"c\":1,\"d\":\"1h0m0s\"}\n" +

					"{\"L\":\"WARN\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"warn level f: a b 1 2 3\",\"a\":\"b\",\"c\":1,\"d\":\"1h0m0s\"}\n" +

					"{\"L\":\"ERROR\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"error level: a1 2\",\"s\":\"a\",\"b\":1}\n" +

					"{\"L\":\"ERROR\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"error level f: a 1 2\",\"s\":\"a\",\"d\":\"1m0s\"}\n" +

					"{\"L\":\"PANIC\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"panic level: ab1 2\",\"a\":\"b\"}\n" +

					"{\"L\":\"PANIC\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"panic level f: a b 1 2\",\"a\":\"b\"}\n" +

					"{\"L\":\"FATAL\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"fatal level: ab1 2\",\"f\":\"b\"}\n" +

					"{\"L\":\"FATAL\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"fatal level f: a b 1 2\",\"f\":\"b\"}\n" +

					"{\"L\":\"WARN\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"warn level 1: ab1 2\",\"l\":\"w\"}\n" +

					"{\"L\":\"WARN\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"warn level f 1: a b 1 2\",\"l\":\"w\"}\n" +

					"{\"L\":\"ERROR\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"error level 2: ab1 2\",\"l\":\"e\"}\n" +

					"{\"L\":\"ERROR\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"error level f 2: a b 1 2\",\"l\":\"e\"}\n" +

					"{\"L\":\"PANIC\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"panic level 3: ab1 2\",\"l\":\"p\"}\n" +

					"{\"L\":\"PANIC\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"panic level f 3: a b 1 2\",\"l\":\"p\"}\n" +

					"{\"L\":\"FATAL\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"fatal level 4: ab1 2\",\"l\":\"f\"}\n" +

					"{\"L\":\"FATAL\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"fatal level f 4: a b 1 2\",\"l\":\"f\"}\n" +

					"{\"L\":\"DEBUG\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"debug level 5: ab1 2\",\"l\":\"d\"}\n" +

					"{\"L\":\"DEBUG\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
					"\"C\":\"logger/logger_test\\.go:\\d+\"," +
					"\"M\":\"debug level f 5: a b 1 2\",\"l\":\"d\"}\n$"

				for _, filePath := range []string{"/tmp/mars-json.log", "/tmp/mars-another-json.log"} {
					content, err := ioutil.ReadFile(filePath)
					So(err, ShouldBeNil)
					So(string(content), test.ShouldMatchWithRegex, expect)
				}
			})
		})
	})

	Convey("Given a logger with plain encoder config", t, func() {
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
				Info("info level: a", "b", 1, 2, 3)
			logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				Infof("info level f: %s %s %d %d %d", "a", "b", 1, 2, 3)
			logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				Warn("warn level: a", "b", 1, 2, 3)
			logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
				Warnf("warn level f: %s %s %d %d %d", "a", "b", 1, 2, 3)
			logger.With(StringField("s", "a")).With(IntField("b", 1)).
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
					"info level: ab1 2 3\t{\"a\": \"b\", \"c\": 1, \"d\": \"1h0m0s\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tINFO\t" + hostName + "\tmod_name\t" +
					"info level f: a b 1 2 3\t{\"a\": \"b\", \"c\": 1, \"d\": \"1h0m0s\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tWARN\t" + hostName + "\tmod_name\t" +
					"warn level: ab1 2 3\t{\"a\": \"b\", \"c\": 1, \"d\": \"1h0m0s\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tWARN\t" + hostName + "\tmod_name\t" +
					"warn level f: a b 1 2 3\t{\"a\": \"b\", \"c\": 1, \"d\": \"1h0m0s\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tERROR\t" + hostName + "\tmod_name\t" +
					"error level: a1 2\t{\"s\": \"a\", \"b\": 1}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tERROR\t" + hostName + "\tmod_name\t" +
					"error level f: a 1 2\t{\"s\": \"a\", \"d\": \"1m0s\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tPANIC\t" + hostName + "\tmod_name\t" +
					"panic level: ab1 2\t{\"a\": \"b\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tPANIC\t" + hostName + "\tmod_name\t" +
					"panic level f: a b 1 2\t{\"a\": \"b\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tFATAL\t" + hostName + "\tmod_name\t" +
					"fatal level: ab1 2\t{\"f\": \"b\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tFATAL\t" + hostName + "\tmod_name\t" +
					"fatal level f: a b 1 2\t{\"f\": \"b\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tWARN\t" + hostName + "\tmod_name\t" +
					"warn level 1: ab1 2\t{\"l\": \"w\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tWARN\t" + hostName + "\tmod_name\t" +
					"warn level f 1: a b 1 2\t{\"l\": \"w\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tERROR\t" + hostName + "\tmod_name\t" +
					"error level 2: ab1 2\t{\"l\": \"e\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tERROR\t" + hostName + "\tmod_name\t" +
					"error level f 2: a b 1 2\t{\"l\": \"e\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tPANIC\t" + hostName + "\tmod_name\t" +
					"panic level 3: ab1 2\t{\"l\": \"p\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tPANIC\t" + hostName + "\tmod_name\t" +
					"panic level f 3: a b 1 2\t{\"l\": \"p\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tFATAL\t" + hostName + "\tmod_name\t" +
					"fatal level 4: ab1 2\t{\"l\": \"f\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tFATAL\t" + hostName + "\tmod_name\t" +
					"fatal level f 4: a b 1 2\t{\"l\": \"f\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tDEBUG\t" + hostName + "\tmod_name\t" +
					"debug level 5: ab1 2\t{\"l\": \"d\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tDEBUG\t" + hostName + "\tmod_name\t" +
					"debug level f 5: a b 1 2\t{\"l\": \"d\"}\t\\[logger_test.go\\]\\[\\d+\\]\n$"

				for _, filePath := range []string{"/tmp/mars-plain.log", "/tmp/mars-another-plain.log"} {
					content, err := ioutil.ReadFile(filePath)
					So(err, ShouldBeNil)
					So(string(content), test.ShouldMatchWithRegex, expect)
				}
			})
		})
	})

	Convey("Given a logger with plain encoder(output to stdout and stderr) config", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths: 
    - stdout
    - stderr
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
		f1, _ := ioutil.TempFile(".", "stdout-*.log")
		defer func() {
			f1.Close()
			os.Remove(f1.Name())
		}()

		f2, _ := ioutil.TempFile(".", "stderr-*.log")
		defer func() {
			f2.Close()
			os.Remove(f2.Name())
		}()

		logger, err := NewLogger(conf.Name())
		So(err, ShouldBeNil)

		patches := gomonkey.ApplyFunc(os.Exit, func(code int) {
			panic(code)
		})
		defer patches.Reset()

		Convey("When log messages with dynamic adjust level", func() {
			test.DoActionsInFileRedirectedContext(os.Stdout, f1, func() {
				test.DoActionsInFileRedirectedContext(os.Stderr, f2, func() {
					logger = logger.Named("mod_name")
					logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
						Debug("debug level: a", "b", 1, 2, 3)
					logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
						Debugf("debug level f: %s %s %d %d %d", "a", "b", 1, 2, 3)
					logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
						Info("info level: a", "b", 1, 2, 3)
					logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
						Infof("info level f: %s %s %d %d %d", "a", "b", 1, 2, 3)
					logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
						Warn("warn level: a", "b", 1, 2, 3)
					logger.With(StringField("a", "b"), IntField("c", 1), DurationField("d", time.Duration(1*time.Hour))).
						Warnf("warn level f: %s %s %d %d %d", "a", "b", 1, 2, 3)
					logger.With(StringField("s", "a")).With(IntField("b", 1)).
						Error("error level: a", 1, 2)
					logger.With(StringField("s", "a")).With(DurationField("d", time.Duration(1*time.Minute))).
						Errorf("error level f: %s %d %d", "a", 1, 2)

					logger.SetLevel("Debug")
					logger.With(StringField("l", "d")).
						Debug("debug level 5: a", "b", 1, 2)
					logger.With(StringField("l", "d")).
						Debugf("debug level f 5: %s %s %d %d", "a", "b", 1, 2)

					logger.Sync()
				})
			})

			Convey("Then the log content should be expect(discard lower level log)", func() {
				hostName, err := os.Hostname()
				if err != nil {
					hostName = "unknown"
				}
				timeFormat := "\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{3}(([+-]\\d{2}:\\d{2})|Z)"

				expect := timeFormat + "\tINFO\t" + hostName + "\tmod_name\t" +
					"info level: ab1 2 3\t{\"a\": \"b\", \"c\": 1, \"d\": \"1h0m0s\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tINFO\t" + hostName + "\tmod_name\t" +
					"info level f: a b 1 2 3\t{\"a\": \"b\", \"c\": 1, \"d\": \"1h0m0s\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tWARN\t" + hostName + "\tmod_name\t" +
					"warn level: ab1 2 3\t{\"a\": \"b\", \"c\": 1, \"d\": \"1h0m0s\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tWARN\t" + hostName + "\tmod_name\t" +
					"warn level f: a b 1 2 3\t{\"a\": \"b\", \"c\": 1, \"d\": \"1h0m0s\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tERROR\t" + hostName + "\tmod_name\t" +
					"error level: a1 2\t{\"s\": \"a\", \"b\": 1}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tERROR\t" + hostName + "\tmod_name\t" +
					"error level f: a 1 2\t{\"s\": \"a\", \"d\": \"1m0s\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tDEBUG\t" + hostName + "\tmod_name\t" +
					"debug level 5: ab1 2\t{\"l\": \"d\"}\t\\[logger_test.go\\]\\[\\d+\\]\n" +

					timeFormat + "\tDEBUG\t" + hostName + "\tmod_name\t" +
					"debug level f 5: a b 1 2\t{\"l\": \"d\"}\t\\[logger_test.go\\]\\[\\d+\\]\n$"

				for _, filePath := range []string{f1.Name()} {
					content, err := ioutil.ReadFile(filePath)
					So(err, ShouldBeNil)
					So(string(content), test.ShouldMatchWithRegex, expect)
				}
			})
		})
	})

	Convey("Given a logger config with unsupported encoder", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: console
    outputPaths: 
    - /tmp/mars-json.log
    - /tmp/mars-another-json.log
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

		Convey("When create logger by above config", func() {
			logger, err := NewLogger(conf.Name())

			Convey("Then the result is abnormal and tips: unsupported encoder", func() {
				So(logger, ShouldBeNil)
				So(err.Error(), ShouldContainSubstring, "unsupported encoder")
			})
		})
	})
}
func TestLogConfigChange(t *testing.T) {
	Convey("Given a logger with json encoder config", t, func() {
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
      compress: false`

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

		patches := gomonkey.ApplyFunc(os.Exit, func(code int) {
			panic(code)
		})
		defer patches.Reset()

		Convey("When log messages with dynamic adjust level", func() {
			logger = logger.Named("mod_name")
			conf.WriteAt([]byte(`code-shooting:
			log:
			  level: info
			  encoder: json
			  outputPaths: 
			  - /tmp/mars-json.log
			  - /tmp/mars-another-json.log
			  rotateConfig:
				maxSize: 1 # MB
				maxBackups: 3 # 
				maxAge: 7 # days
				compress: false
                fileMode: "0644"`), 0)
			config.ProcessConfigEvent([]*model.Event{
				{
					ConfigItem: model.ConfigItem{Key: common.ROOT + "." + common.LOG_LEVEL, Value: "a"},
					EventType:  model.Update,
				},
				{
					ConfigItem: model.ConfigItem{Key: common.ROOT + "." + common.LOG_ROTATE_CONFIG_MAXSIZE, Value: 2},
					EventType:  model.Update,
				},
			})

			config.ProcessConfigEvent([]*model.Event{
				{
					ConfigItem: model.ConfigItem{Key: common.ROOT + "." + common.LOG_LEVEL, Value: "debug"},
					EventType:  model.Update,
				},
				{
					ConfigItem: model.ConfigItem{Key: common.ROOT + "." + common.LOG_ROTATE_CONFIG_MAXSIZE, Value: 2},
					EventType:  model.Update,
				},
				{
					ConfigItem: model.ConfigItem{Key: common.ROOT + "." + common.LOG_ROTATE_CONFIG_MAXBACKUPS, Value: 2},
					EventType:  model.Update,
				},
				{
					ConfigItem: model.ConfigItem{Key: common.ROOT + "." + common.LOG_ROTATE_CONFIG_MAXAGE, Value: 5},
					EventType:  model.Update,
				},
				{
					ConfigItem: model.ConfigItem{Key: common.ROOT + "." + common.LOG_ROTATE_CONFIG_COMPRESS, Value: 1},
					EventType:  model.Update,
				},
				{
					ConfigItem: model.ConfigItem{Key: common.ROOT + "." + common.LOG_ROTATE_CONFIG_COMPRESS, Value: true},
					EventType:  model.Update,
				},
			})

			logger.With(StringField("l", "d")).
				Debug("debug level 5: a", "b", 1, 2)
			logger.With(StringField("l", "d")).
				Debugf("debug level f 5: %s %s %d %d", "a", "b", 1, 2)

			logger.Sync()

			Convey("Then the log content should be expect(discard lower level log)", func() {
				timeFormat := "\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{3}(([+-]\\d{4})|Z)"
				expect :=
					"{\"L\":\"DEBUG\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
						"\"C\":\"logger/logger_test\\.go:\\d+\"," +
						"\"M\":\"debug level 5: ab1 2\",\"l\":\"d\"}\n" +

						"{\"L\":\"DEBUG\",\"T\":\"" + timeFormat + "\",\"N\":\"mod_name\"," +
						"\"C\":\"logger/logger_test\\.go:\\d+\"," +
						"\"M\":\"debug level f 5: a b 1 2\",\"l\":\"d\"}\n$"

				content, err := ioutil.ReadFile("/tmp/mars-json.log")
				So(err, ShouldBeNil)
				So(string(content), test.ShouldMatchWithRegex, expect)

			})
		})
	})
}
