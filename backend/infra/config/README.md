# code-shooting: config 

Convenient, injection-friendly YAML configuration.

## Installation

In your code:
```golang
importcode-shooting/config
```

If you have a working Go 1.6/1.7 environment:
```shell
go getcode-shooting
```

## Quick Start
If you have a config file: /app/config/base.yaml
```yaml
module:
  parameter: foo
````

```golang
import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"go.uber.org/fx"

	"code-shooting/config"
)

// Model your application's configuration using a Go struct.
type Module struct {
	Parameter string
}

var conf config.Config
app := fx.New(
	config.NewModule("/app/config/base.yaml"),
	fx.Invoke(func(c config.Config) {
		conf = c
	}),
)

startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()
if err := app.Start(startCtx); err != nil {
	fmt.Println(err)
}

mod := Module{}
if err:= conf.Get("module", &mod); err != nil{
	return err // handle error
}
fmt.Printf("%+v\n", mod)
// Output:
// {Parameter:foo}

quit := make(chan os.Signal)
signal.Notify(quit, os.Interrupt)
<-quit
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()
if err := app.Stop(ctx); err != nil {
	fmt.Println(err)
}

```

If you have a config direction: /app/config, there are two config files in it: base.yaml and override.yml:
```yaml
// base.yaml
module:
  parameter: foo
````
```yaml
// override.yml
module:
  parameter: bar
````

```golang
import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"go.uber.org/fx"

	"code-shooting/config"
)

// Model your application's configuration using a Go struct.
type Module struct {
	Parameter string
}

var conf config.Config
app := fx.New(
	config.NewModule("/app/config"),
	fx.Invoke(func(c config.Config) {
		conf = c
	}),
)

startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()
if err := app.Start(startCtx); err != nil {
	fmt.Println(err)
}

mod := Module{}
if err := conf.Get("module", &mod); err != nil{
	return err // handle error
}

// Merge the two config files into a Provider. Subsequent loaded config are higher-priority.
// See the top-level package documentation for details on the merging logic.
fmt.Printf("%+v\n", mod)
// Output:
// {Parameter:bar}

quit := make(chan os.Signal)
signal.Notify(quit, os.Interrupt)
<-quit
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()
if err := app.Stop(ctx); err != nil {
	fmt.Println(err)
}

```
When Get return an error, you could use IsNotExist to know whether the error report that a key is not exist:
```
if err := conf.Get("module", &mod); err != nil {
	if IsNotExist(err):
		panic("the key is not exist in config file)
	else:
		panic("the value type with the key is not consistent with in config file)
}
```