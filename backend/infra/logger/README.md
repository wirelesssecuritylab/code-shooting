# code-shooting: logger 

Structured, customizable, leveled logging in Go.

## Installation

In your code:

```golang
import code-shooting/infra/logger
```

If you have a working Go 1.6/1.7 environment:

```shell
go get code-shooting/infra
```

## Quick Start

If you have a configure file: ```/app/config/base.yaml```

``` yaml
code-shooting:
  log:
    level: info # <debug|info|warn|error|panic|fatal>
    encoder: json # <json|plain>
    outputPaths: # Log output files
    - /tmp/app.log
    rotateConfig: 
      maxSize: 1 # The maximum size in megabytes of the log file before it gets rotated. If maxSize less than or equal to 0, then set to default(10MB).
      maxBackups: 2 # The maximum number of old log files to retain.  If maxBackups less than or equal to 0, this means to retain all old log files(though maxAge may still cause them to get deleted).
      maxAge: 7 # The maximum number of days to retain old log files based on the timestamp encoded in their filename.  Note that a day is defined as 24 hours and may not exactly correspond to calendar days due to daylight savings, leap seconds, etc. If maxAge less than or equal to 0, this means not to remove old log files based on age.
      compress: false # Determines if the rotated log files should be compressed using gzip. The default is not to perform compression.
```



**Notes**:

> If maxBackups and maxAge both are less than or equal to 0, then set maxBackups to 2.



``` go
import (
	"code-shooting/infra/config"
    "go.uber.org/fx"
    "time"
    "log"
    "os"
    "os/signal"
)

...

l, err := logger.NewLogger("/app/config/")
if err != nil {
    log.Fatal("failed create code-shooting logger", err)
    return
}

logger.SetLogger(l)

app := fx.New(
    fx.Logger(logger.GetLogger().CreateStdLogger()),
    config.NewModule("/app/config/"),
    fx.Invoke(...),
    ...
)

startCtx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
defer cancel()
if err := app.Start(startCtx); err != nil {
    log.Fatal(err)
}

logger.Info("xxx")

...

quit := make(chan os.Signal)
signal.Notify(quit, os.Interrupt)
<-quit

stopCtx, cancel := context.WithTimeout(context.Backgroup(), 15 * time.Second)
def cancel()
if err := app.Stop(stopCtx); err != nil {
    log.Fatal(err)
}
```



The output when the encoder is ```json```:

``` text
{"L":"INFO","T":"2020-11-25T09:40:39.935+0800","N":"mod_name","C":"logger/logger_test.go:142","M":"msg","key":"value"}
```

The output when the encoder is ```plain```:

``` text
2020-11-25T09:53:15.069+08:00	INFO	hostName	mod_name	msg	{"key": "value"}	[logger_test.go][1166]
```



---



If you want to customize the output format to ```plain```, then:

```yaml
code-shooting:
  log:
    level: info  
    encoder: plain
    format: "$${T}\t$${L}\t$${H}\t$${N}\t$${M} $${E}\t[$${C_F}]:[$${C_L}]\n" 
    outputPaths: 
    - /tmp/app.log
    rotateConfig: 
      maxSize: 1 
      maxBackups: 2
      maxAge: 7 
      compress: false 
```

> ${...} is placeholder for specified field, other is the plain character:
>
> * ${T} -> Time
> * ${L} -> Level
> * ${H} -> Host name
> * ${N} -> Mod name
> * ${M} -> Message
> * ${C} -> Caller
> * ${C_F} -> Caller file
> * ${C_L} -> Caller line number
> * ${XX:xx} -> XX filed, and the xx is the default value
> * ${E} -> Extends fields(all the unspecified fields)



The output is:

```text
2020-11-25T09:53:15.069+08:00	INFO	hostName	mod_name	msg	XX=value YY=null	[logger_test.go][1166]
```



```API``` interface:

```go
func NewLogger(configPath string, options ...Option) (Logger, error)
```



```logger``` interface:

``` go
type Logger interface {
        Debug(args ...interface{})
        Info(args ...interface{})
        Warn(args ...interface{})
        Error(args ...interface{})
        Panic(args ...interface{})
        Fatal(args ...interface{})
        Debugf(format string, args ...interface{})
        Infof(format string, args ...interface{})
        Warnf(format string, args ...interface{})
        Errorf(format string, args ...interface{})
        Panicf(format string, args ...interface{})
        Fatalf(format string, args ...interface{})

	    CreateStdLogger() *log.Logger
    
	    AddCallerSkip(skip int) Logger
        Named(name string) Logger // set logger name
        With(fields ...Field) Logger // add extended fields info
        SetLevel(level string) // dynamic adjustment level
	    GetLevel() string 
        Sync() error 
}
```



current Fields:

```go
func StringField(k, v string) Field
func IntField(k string, v int64) Field
func DurationField(k string, v time.Duration) Field
```

