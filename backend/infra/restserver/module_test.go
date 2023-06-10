package restserver

import (
	"context"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"code-shooting/infra/restserver/internal"

	"go.uber.org/fx"

	. "github.com/agiledragon/gomonkey"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	. "github.com/smartystreets/goconvey/convey"

	marsconfig "code-shooting/infra/config"
	"code-shooting/infra/logger"
	"code-shooting/infra/x/test"
)

func waitRestServerSteady(addr string) {
	for i := 0; i < 10; i++ {
		res, err := http.Get("http://" + addr + "/")

		if err == nil && res.StatusCode == http.StatusOK {
			res.Body.Close()
			return
		}
		time.Sleep(500 * time.Microsecond)
	}
}

type user struct {
	Id   int    `json:"id" xml:"id" form:"id" query:"id" param:"id"`
	Name string `json:"name" xml:"name" form:"name" query:"name" param:"name"`
}

var userlist map[string]user

func hello(c Context) error {
	userlist = make(map[string]user)
	return c.String(http.StatusOK, "Hello, World!")
}

func saveUser(c Context) error {
	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, user{1, "Lily"})
	}
	name := c.FormValue("name")
	userlist[c.FormValue("id")] = user{Id: id, Name: name}
	return c.JSON(http.StatusOK, user{1, "Lily"})
}

func getUser(c Context) error {
	return c.JSON(http.StatusOK, userlist[c.Param("id")])
}

func updateUser(c Context) error {
	username := c.QueryParam("name")
	userid := c.QueryParam("id")
	if old, ok := userlist[userid]; ok {
		old.Name = username
	}
	return c.JSON(http.StatusOK, userlist[userid])
}

func deleteUser(c Context) error {
	delete(userlist, c.Param("id"))
	return c.NoContent(http.StatusOK)
}

func noNameServer(mRestServer *MultRestServer) error {
	server, err := mRestServer.GetRestServerByName("")
	if err != nil {
		return err
	}
	server.GET("/", hello)
	return nil
}

func noExistedServer(mRestServer *MultRestServer) error {
	server, err := mRestServer.GetRestServerByName("noserver")
	if err != nil {
		return err
	}
	server.GET("/", hello)
	return nil
}

func myServerRegister(mRestServer *MultRestServer) error {
	server, err := mRestServer.GetRestServerByName("myserver")
	if err != nil {
		return err
	}
	if server.RootGroupBox.RootGroup == nil {
		return nil
	}
	server.RootGroupBox.RootGroup.GET("/xxx", hello)
	server.RootGroupBox.RootGroup.Any("/maxconns", maxConnsHandler)
	server.GET("/", hello)
	server.POST("/users", saveUser)
	server.GET("/users/:id", getUser)
	server.PUT("/users/:id", updateUser)
	server.DELETE("/users/:id", deleteUser)
	return nil
}

func urServerRegister(mRestServer *MultRestServer) error {
	server, err := mRestServer.GetRestServerByName("urserver")
	if err != nil {
		return err
	}
	server.GET("/", hello)
	server.RootGroupBox.RootGroup.GET("/xxx", hello)
	return nil
}

func TestModuleStartFailed(t *testing.T) {
	Convey("Given a one-rest-server config file \n", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths:
    - stdout
  rest-servers:
  - name: myserver
    addr: $restserverip:$restserverport`

		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		os.Setenv("restserverip", "127.0.0.1")
		os.Setenv("restserverport", "9099")
		defer os.Unsetenv("restserverip")
		defer os.Unsetenv("restserverport")

		e := &echo.Echo{}
		patches := ApplyMethod(reflect.TypeOf(e), "StartServer", func(_ *echo.Echo, _ *http.Server) error {
			return errors.New("start server failed")
		})
		patches.ApplyFunc(buildListener, func(_ *internal.ListenerConf) (net.Listener, error) {
			return nil, nil
		})
		defer patches.Reset()

		old := logger.GetLogger()
		defer func() {
			logger.SetLogger(old)
		}()

		l, err := logger.NewLogger(confFile.Name())
		So(err, ShouldBeNil)
		logger.SetLogger(l)

		app := fx.New(
			fx.Logger(logger.GetLogger().CreateStdLogger()),
			marsconfig.NewModule(confFile.Name()),
			NewModule(),
			fx.Invoke(myServerRegister),
		)
		Convey("When start fx app and echo start server failed \n", func() {
			err := test.StartFxApp(app)
			Convey("Then start failed \n", func() {
				So(err.Error(), ShouldEqual, "context deadline exceeded")
			})
		})

		Convey("When start fx app failed and context-canceled \n", func() {
			ctx, cancel := context.WithTimeout(context.Background(), app.StartTimeout())
			err := app.Start(ctx)
			time.Sleep(time.Second)
			cancel()
			Convey("Then start failed \n", func() {
				So(err.Error(), ShouldEqual, "context deadline exceeded")
			})
		})
	})
}

func TestRestServerModule(t *testing.T) {
	Convey("Given a one-rest-server config file and a rest-server fx app", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths:
      - stdout
  rest-servers:
  - name: myserver
    addr: $restserverip:$restserverport
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384
    rootpath: /mypath
    middlewares:
    - name: bodylimit
      limit: 2M
    - name: stats
    - name: pprof
    - name: expvar`

		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		os.Setenv("restserverip", "127.0.0.1")
		os.Setenv("restserverport", "9080")
		defer os.Unsetenv("restserverip")
		defer os.Unsetenv("restserverport")

		old := logger.GetLogger()
		defer func() {
			logger.SetLogger(old)
		}()

		l, err := logger.NewLogger(confFile.Name())
		So(err, ShouldBeNil)
		logger.SetLogger(l)

		app := fx.New(
			fx.Logger(logger.GetLogger().CreateStdLogger()),
			marsconfig.NewModule(confFile.Name()),
			NewModule(),
			fx.Invoke(myServerRegister),
		)

		So(test.StartFxApp(app), ShouldBeNil)
		defer test.StopFxApp(app)
		waitRestServerSteady("127.0.0.1:9080")

		Convey("When access a URL of /mypath/xxx", func() {
			resp, err := http.Get("http://127.0.0.1:9080/mypath/xxx")
			So(err, ShouldBeNil)
			defer resp.Body.Close()

			Convey("Then status with 200 and response with 'Hello, World!'  should be received", func() {
				body, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(string(body), ShouldEqual, "Hello, World!")
			})
		})

		Convey("When access a URL of /mypath/stats", func() {
			resp, err := http.Get("http://127.0.0.1:9080/stats")
			So(err, ShouldBeNil)
			defer resp.Body.Close()

			Convey("Then response should contain 'total_status_code_count' ", func() {
				body, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(string(body), ShouldContainSubstring, "total_status_code_count")
			})
		})
	})
	Convey("Given a two-rest-servers config file and mult-rest-server fx app", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths:
    - stdout
  rest-servers:
  - name: myserver
    addr: 127.0.0.1:9081
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384
    rootpath: /mypath
  - name: urserver
    addr: 127.0.0.1:9082
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384
    rootpath: /urpath`

		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		old := logger.GetLogger()
		defer func() {
			logger.SetLogger(old)
		}()

		l, err := logger.NewLogger(confFile.Name())
		So(err, ShouldBeNil)
		logger.SetLogger(l)

		app := fx.New(
			fx.Logger(logger.GetLogger().CreateStdLogger()),
			marsconfig.NewModule(confFile.Name()),
			NewModule(),
			fx.Invoke(myServerRegister),
			fx.Invoke(urServerRegister),
		)

		So(test.StartFxApp(app), ShouldBeNil)
		defer test.StopFxApp(app)
		waitRestServerSteady("127.0.0.1:9081")
		waitRestServerSteady("127.0.0.1:9082")

		Convey("When access a URL of /mypath/xxx", func() {
			resp, err := http.Get("http://127.0.0.1:9081/mypath/xxx")
			So(err, ShouldBeNil)
			defer resp.Body.Close()

			Convey("Then status with 200 and response with 'Hello, World!' should be received", func() {
				body, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(string(body), ShouldEqual, "Hello, World!")
			})
		})

		Convey("When access a URL of /urpath/xxx", func() {
			resp, err := http.Get("http://127.0.0.1:9082/urpath/xxx")
			So(err, ShouldBeNil)
			defer resp.Body.Close()

			Convey("Then status with 200 and response with 'Hello, World!' should be received", func() {
				body, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(string(body), ShouldEqual, "Hello, World!")
			})
		})
	})

	Convey("Given a one-rest-server config file ", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths:
    - stdout
  rest-servers:
  - name: myserver
    addr: 127.0.0.1:9083
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384`

		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		Convey("When register a not existed rest-server and start fx app", func() {
			app := fx.New(
				fx.Logger(logger.GetLogger().CreateStdLogger()),
				marsconfig.NewModule(confFile.Name()),
				NewModule(),
				fx.Invoke(noExistedServer),
			)
			err := test.StartFxApp(app)
			defer test.StopFxApp(app)

			Convey("Then err should contain no rest server with name ", func() {
				So(err.Error(), ShouldContainSubstring, "no rest server with name")

			})
		})
	})
	Convey("Given a one-rest-server config file ", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths:
    - stdout
  rest-servers:
  - name: myserver
    addr: 127.0.0.1:9084
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384`

		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()
		Convey("When register a no-name rest-server and start fx app", func() {
			app := fx.New(
				fx.Logger(logger.GetLogger().CreateStdLogger()),
				marsconfig.NewModule(confFile.Name()),
				NewModule(),
				fx.Invoke(noNameServer),
			)
			err := test.StartFxApp(app)
			defer test.StopFxApp(app)

			Convey("Then err should contain rest server name is nil ", func() {
				So(err.Error(), ShouldContainSubstring, "rest server name is nil")
			})
		})
	})
}
func TestRestServerIpInvalid(t *testing.T) {

	Convey("Given a invalid-ip-server config file ", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths:
    - stdout
  rest-servers:
  - name: myserver
    addr: 127.0.:8080
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384
    rootpath: `

		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		old := logger.GetLogger()
		defer func() {
			logger.SetLogger(old)
		}()

		l, err := logger.NewLogger(confFile.Name())
		So(err, ShouldBeNil)
		logger.SetLogger(l)

		Convey("When start fx app", func() {
			app := fx.New(
				fx.Logger(logger.GetLogger().CreateStdLogger()),
				marsconfig.NewModule(confFile.Name()),
				NewModule(),
				fx.Invoke(myServerRegister),
			)
			err := test.StartFxApp(app)
			defer test.StopFxApp(app)

			Convey("Then err should contain 'not a valid textual representation of an IP address'", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "not a valid textual representation of an IP address")

			})
		})
	})
}
func TestRestServerPortUsed(t *testing.T) {
	Convey("Given a config file and a http-server with the same port ", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths:
    - stdout
  rest-servers:
  - name: myserver
    addr: 127.0.0.1:9081
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384
    rootpath: /mypath`

		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		httpserver := &http.Server{Addr: ":9081"}
		go httpserver.ListenAndServe()
		defer httpserver.Close()

		app := fx.New(
			fx.Logger(logger.GetLogger().CreateStdLogger()),
			marsconfig.NewModule(confFile.Name()),
			NewModule(),
			fx.Invoke(myServerRegister),
		)

		Convey("When start mult-rest-server fx app", func() {

			err := test.StartFxApp(app)

			Convey("Then start failed with message 'address already in use'", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "address already in use")
			})
		})

	})
}

var LocalhostCert = []byte(`-----BEGIN CERTIFICATE-----
MIICEzCCAXygAwIBAgIQMIMChMLGrR+QvmQvpwAU6zANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMCAXDTcwMDEwMTAwMDAwMFoYDzIwODQwMTI5MTYw
MDAwWjASMRAwDgYDVQQKEwdBY21lIENvMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCB
iQKBgQDuLnQAI3mDgey3VBzWnB2L39JUU4txjeVE6myuDqkM/uGlfjb9SjY1bIw4
iA5sBBZzHi3z0h1YV8QPuxEbi4nW91IJm2gsvvZhIrCHS3l6afab4pZBl2+XsDul
rKBxKKtD1rGxlG4LjncdabFn9gvLZad2bSysqz/qTAUStTvqJQIDAQABo2gwZjAO
BgNVHQ8BAf8EBAMCAqQwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUw
AwEB/zAuBgNVHREEJzAlggtleGFtcGxlLmNvbYcEfwAAAYcQAAAAAAAAAAAAAAAA
AAAAATANBgkqhkiG9w0BAQsFAAOBgQCEcetwO59EWk7WiJsG4x8SY+UIAA+flUI9
tyC4lNhbcF2Idq9greZwbYCqTTTr2XiRNSMLCOjKyI7ukPoPjo16ocHj+P3vZGfs
h1fIw3cSS2OolhloGw/XM6RWPWtPAlGykKLciQrBru5NAPvCMsb/I1DAceTiotQM
fblo6RBxUQ==
-----END CERTIFICATE-----`)

// LocalhostKey is the private key for localhostCert.
var LocalhostKey = []byte(testingKey(`-----BEGIN RSA TESTING KEY-----
MIICXgIBAAKBgQDuLnQAI3mDgey3VBzWnB2L39JUU4txjeVE6myuDqkM/uGlfjb9
SjY1bIw4iA5sBBZzHi3z0h1YV8QPuxEbi4nW91IJm2gsvvZhIrCHS3l6afab4pZB
l2+XsDulrKBxKKtD1rGxlG4LjncdabFn9gvLZad2bSysqz/qTAUStTvqJQIDAQAB
AoGAGRzwwir7XvBOAy5tM/uV6e+Zf6anZzus1s1Y1ClbjbE6HXbnWWF/wbZGOpet
3Zm4vD6MXc7jpTLryzTQIvVdfQbRc6+MUVeLKwZatTXtdZrhu+Jk7hx0nTPy8Jcb
uJqFk541aEw+mMogY/xEcfbWd6IOkp+4xqjlFLBEDytgbIECQQDvH/E6nk+hgN4H
qzzVtxxr397vWrjrIgPbJpQvBsafG7b0dA4AFjwVbFLmQcj2PprIMmPcQrooz8vp
jy4SHEg1AkEA/v13/5M47K9vCxmb8QeD/asydfsgS5TeuNi8DoUBEmiSJwma7FXY
fFUtxuvL7XvjwjN5B30pNEbc6Iuyt7y4MQJBAIt21su4b3sjXNueLKH85Q+phy2U
fQtuUE9txblTu14q3N7gHRZB4ZMhFYyDy8CKrN2cPg/Fvyt0Xlp/DoCzjA0CQQDU
y2ptGsuSmgUtWj3NM9xuwYPm+Z/F84K6+ARYiZ6PYj013sovGKUFfYAqVXVlxtIX
qyUBnu3X9ps8ZfjLZO7BAkEAlT4R5Yl6cGhaJQYZHOde3JEMhNRcVFMO8dJDaFeo
f9Oeos0UUothgiDktdQHxdNEwLjQf7lJJBzV+5OtwswCWA==
-----END RSA TESTING KEY-----`))

func testingKey(s string) string { return strings.ReplaceAll(s, "TESTING KEY", "PRIVATE KEY") }

func httpsServerRegister(restservers *MultRestServer) error {

	rsHttps, err := restservers.GetRestServerByName("https-server")
	if err != nil {
		return err
	}
	rsHttps.Pre(middleware.HTTPSRedirect())
	rsHttps.GET("/redirect", hello)
	return nil
}

func httpsRestServerOptions() ([]Option, error) {

	optFunc := func(s *http.Server) error {

		certs, err := tls.X509KeyPair(LocalhostCert, LocalhostKey)
		if err == nil {
			s.TLSConfig = &tls.Config{Certificates: []tls.Certificate{certs}}
		}
		return err
	}
	return []Option{WithOption("https-server", optFunc), WithOption("no-server", optFunc)}, nil
}

func TestHttpsRestServer(t *testing.T) {
	Convey("Given a https-rest-server config file and a rest-server fx app", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths:
    - stdout
  rest-servers:
  - name: https-server
    protocol: https
    addr: $restserverip:$restserverport
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384`

		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		os.Setenv("restserverip", "127.0.0.1")
		os.Setenv("restserverport", "9080")
		defer os.Unsetenv("restserverip")
		defer os.Unsetenv("restserverport")

		old := logger.GetLogger()
		defer func() {
			logger.SetLogger(old)
		}()

		l, err := logger.NewLogger(confFile.Name())
		So(err, ShouldBeNil)
		logger.SetLogger(l)

		app := fx.New(
			fx.Logger(logger.GetLogger().CreateStdLogger()),
			marsconfig.NewModule(confFile.Name()),
			NewModule(),
			fx.Provide(httpsRestServerOptions),
			fx.Invoke(httpsServerRegister),
		)

		So(test.StartFxApp(app), ShouldBeNil)
		defer test.StopFxApp(app)
		waitRestServerSteady("127.0.0.1:9080")

		Convey("When access a URL of /redirect", func() {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{Transport: tr}
			resp, err := client.Get("https://127.0.0.1:9080/redirect")
			So(err, ShouldBeNil)
			defer resp.Body.Close()

			Convey("Then status with 200 and response with 'Hello, World!' should be received", func() {
				body, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(string(body), ShouldEqual, "Hello, World!")
			})
		})
	})

	Convey("Given a https-rest-server config file with cert and key, and a rest-server fx app", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths:
    - stdout
  rest-servers:
  - name: https-server
    protocol: https
    certfile: ./crt/server.crt
    keyfile: ./crt/server.key
    addr: $restserverip:$restserverport
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384`

		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		os.Setenv("restserverip", "127.0.0.1")
		os.Setenv("restserverport", "9080")
		defer os.Unsetenv("restserverip")
		defer os.Unsetenv("restserverport")

		patches := ApplyFunc(ioutil.ReadFile, func(filename string) ([]byte, error) {
			if strings.Contains(filename, ".crt") {
				return LocalhostCert, nil
			} else if strings.Contains(filename, ".key") {
				return LocalhostKey, nil
			} else if strings.Contains(filename, "config") {
				return []byte(content), nil
			}
			return []byte{}, errors.New("no such file or directory")
		})
		defer patches.Reset()

		old := logger.GetLogger()
		defer func() {
			logger.SetLogger(old)
		}()
		l, err := logger.NewLogger(confFile.Name())
		So(err, ShouldBeNil)
		logger.SetLogger(l)

		app := fx.New(
			fx.Logger(logger.GetLogger().CreateStdLogger()),
			marsconfig.NewModule(confFile.Name()),
			NewModule(),
			fx.Provide(httpsRestServerOptions),
			fx.Invoke(httpsServerRegister),
		)

		So(test.StartFxApp(app), ShouldBeNil)
		defer test.StopFxApp(app)
		waitRestServerSteady("127.0.0.1:9080")

		Convey("When access a URL of /redirect", func() {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{Transport: tr}
			resp, err := client.Get("https://127.0.0.1:9080/redirect")
			So(err, ShouldBeNil)
			defer resp.Body.Close()

			Convey("Then status with 200 and response with 'Hello, World!' should be received", func() {
				body, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(string(body), ShouldEqual, "Hello, World!")
			})
		})
	})

	Convey("Given a https-rest-server config file with error-cert\n", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths:
    - stdout
  rest-servers:
  - name: https-server
    protocol: https
    certfile: ./crt/server.crt
    keyfile: ./crt/server.key
    addr: $restserverip:$restserverport
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384`

		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		os.Setenv("restserverip", "127.0.0.1")
		os.Setenv("restserverport", "9080")
		defer os.Unsetenv("restserverip")
		defer os.Unsetenv("restserverport")

		patches := ApplyFunc(ioutil.ReadFile, func(filename string) ([]byte, error) {
			if strings.Contains(filename, ".crt") {
				return []byte{}, errors.New("cann't open")
			} else if strings.Contains(filename, ".key") {
				return LocalhostKey, nil
			} else if strings.Contains(filename, "config") {
				return []byte(content), nil
			}
			return []byte{}, errors.New("no such file or directory")
		})
		defer patches.Reset()

		old := logger.GetLogger()
		defer func() {
			logger.SetLogger(old)
		}()

		l, err := logger.NewLogger(confFile.Name())
		So(err, ShouldBeNil)
		logger.SetLogger(l)

		app := fx.New(
			fx.Logger(logger.GetLogger().CreateStdLogger()),
			marsconfig.NewModule(confFile.Name()),
			NewModule(),
			fx.Provide(httpsRestServerOptions),
			fx.Invoke(httpsServerRegister),
		)

		Convey("When start a rest-server fx app", func() {
			err := test.StartFxApp(app)

			Convey("Then an error with 'read TLS cert file' should be received", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "read TLS cert file")
			})
		})
	})

	Convey("Given a https-rest-server config file with error-key\n", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths:
    - stdout
  rest-servers:
  - name: https-server
    protocol: https
    certfile: ./crt/server.crt
    keyfile: ./crt/server.key
    addr: $restserverip:$restserverport
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384`

		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		os.Setenv("restserverip", "127.0.0.1")
		os.Setenv("restserverport", "9080")
		defer os.Unsetenv("restserverip")
		defer os.Unsetenv("restserverport")

		patches := ApplyFunc(ioutil.ReadFile, func(filename string) ([]byte, error) {
			if strings.Contains(filename, ".crt") {
				return LocalhostCert, nil
			} else if strings.Contains(filename, ".key") {
				return []byte{}, errors.New("cann't open")
			} else if strings.Contains(filename, "config") {
				return []byte(content), nil
			}
			return []byte{}, errors.New("no such file or directory")
		})
		defer patches.Reset()

		old := logger.GetLogger()
		defer func() {
			logger.SetLogger(old)
		}()

		l, err := logger.NewLogger(confFile.Name())
		So(err, ShouldBeNil)
		logger.SetLogger(l)

		app := fx.New(
			fx.Logger(logger.GetLogger().CreateStdLogger()),
			marsconfig.NewModule(confFile.Name()),
			NewModule(),
			fx.Provide(httpsRestServerOptions),
			fx.Invoke(httpsServerRegister),
		)

		Convey("When start a rest-server fx app", func() {
			err := test.StartFxApp(app)

			Convey("Then an error with 'read TLS key file' should be received", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "read TLS key file")
			})
		})
	})

	Convey("Given a https-rest-server config file with empty-key\n", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths:
    - stdout
  rest-servers:
  - name: https-server
    protocol: https
    certfile: ./crt/server.crt
    keyfile: ./crt/server.key
    addr: $restserverip:$restserverport
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384`

		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		os.Setenv("restserverip", "127.0.0.1")
		os.Setenv("restserverport", "9080")
		defer os.Unsetenv("restserverip")
		defer os.Unsetenv("restserverport")

		patches := ApplyFunc(ioutil.ReadFile, func(filename string) ([]byte, error) {
			if strings.Contains(filename, ".crt") {
				return LocalhostCert, nil
			} else if strings.Contains(filename, ".key") {
				return []byte{}, nil
			} else if strings.Contains(filename, "config") {
				return []byte(content), nil
			}
			return []byte{}, errors.New("no such file or directory")
		})
		defer patches.Reset()

		old := logger.GetLogger()
		defer func() {
			logger.SetLogger(old)
		}()

		l, err := logger.NewLogger(confFile.Name())
		So(err, ShouldBeNil)
		logger.SetLogger(l)

		app := fx.New(
			fx.Logger(logger.GetLogger().CreateStdLogger()),
			marsconfig.NewModule(confFile.Name()),
			NewModule(),
			fx.Provide(httpsRestServerOptions),
			fx.Invoke(httpsServerRegister),
		)

		Convey("When start a rest-server fx app", func() {
			err := test.StartFxApp(app)

			Convey("Then an error with 'TLS X509KeyPair' should be received", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "TLS X509KeyPair")
			})
		})
	})
}

func maxConnsHandler(c Context) error {
	// time.Sleep(10 * time.Millisecond)
	time.Sleep(3 * time.Second)
	return c.String(http.StatusOK, "connection ok")
}

func TestMaxConnections(t *testing.T) {
	Convey("Given a config file and a http-server with the maxconnections \n", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths:
    - stdout
  rest-servers:
  - name: myserver
    addr: 127.0.0.1:9086
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384
    maxconnections: 2
    rootpath: /mypath`
		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		app := fx.New(
			fx.Logger(logger.GetLogger().CreateStdLogger()),
			marsconfig.NewModule(confFile.Name()),
			NewModule(),
			fx.Invoke(myServerRegister),
		)
		So(test.StartFxApp(app), ShouldBeNil)
		defer test.StopFxApp(app)
		waitRestServerSteady("127.0.0.1:9086")

		Convey("When receive one request", func() {
			res, err := http.Get("http://127.0.0.1:9086/mypath/maxconns")

			Convey("Then response of request should be StatusOK", func() {
				So(err, ShouldBeNil)
				defer res.Body.Close()
				So(res.StatusCode, ShouldEqual, http.StatusOK)
			})
		})

		Convey("When receive more than max requests\n", func() {
			var res *http.Response
			var err error
			var failed int32
			wg := sync.WaitGroup{}
			for i := 0; i < 4; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					c := http.Client{Timeout: 4 * time.Second}
					res, err = c.Get("http://127.0.0.1:9086/mypath/maxconns")
					if err != nil {
						atomic.AddInt32(&failed, 1)
						return
					}
					defer res.Body.Close()

				}()
			}
			wg.Wait()

			// We expect some Gets to fail as the kernel's accept queue is filled,
			// but most should succeed.
			Convey("Then failed should not be zero\n", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "context deadline exceeded")
				So(failed, ShouldEqual, 2)
			})
		})
	})
}

func TestRateLimit(t *testing.T) {
	Convey("Given a config file and a http-server with ratelimit \n", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths:
    - stdout
  rest-servers:
  - name: myserver
    addr: 127.0.0.1:9086
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384
    rootpath: /mypath
    middlewares:
    - name: ratelimit
      maxrequests: 3
      requestspersec: 1`
		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		app := fx.New(
			fx.Logger(logger.GetLogger().CreateStdLogger()),
			marsconfig.NewModule(confFile.Name()),
			NewModule(),
			fx.Invoke(myServerRegister),
		)
		So(test.StartFxApp(app), ShouldBeNil)
		defer test.StopFxApp(app)
		waitRestServerSteady("127.0.0.1:9086") // here use a token !!!

		Convey("When receive 7 requests\n", func() {
			var res *http.Response
			var err error
			var failedNum, refusedNum, acceptedNum int32
			var retryAfter string

			for i := 0; i < 7; i++ {
				c := http.Client{Timeout: 4 * time.Second}
				res, err = c.Get("http://127.0.0.1:9086/")
				if err != nil {
					atomic.AddInt32(&failedNum, 1)
					return
				} else if res.StatusCode == http.StatusTooManyRequests {
					retryAfter = res.Header.Get("Retry-After")
					atomic.AddInt32(&refusedNum, 1)
				} else if res.StatusCode == http.StatusOK {
					atomic.AddInt32(&acceptedNum, 1)
				}
				defer res.Body.Close()

				time.Sleep(time.Millisecond)
			}

			Convey("Then 2 requests should be accepted and 5 requests should be refused\n", func() {
				So(failedNum, ShouldEqual, 0)
				So(acceptedNum, ShouldEqual, 2)
				So(refusedNum, ShouldEqual, 5)
				So(retryAfter, ShouldEqual, strconv.Itoa(3/1)) // retry-after calculated by ceil(maxrequests / requestspersec)
			})
		})
	})
}

func TestRateLimitWithSetRateLimit(t *testing.T) {
	Convey("Given a config file and a http-server with ratelimit \n", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths:
    - stdout
  rest-servers:
  - name: myserver
    addr: 127.0.0.1:19086
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384
    rootpath: /mypath
    middlewares:
    - name: ratelimit
      maxrequests: 3
      requestspersec: 1`
		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		var multRestServer *MultRestServer

		app := fx.New(
			fx.Logger(logger.GetLogger().CreateStdLogger()),
			marsconfig.NewModule(confFile.Name()),
			NewModule(),
			fx.Invoke(myServerRegister),
			fx.Invoke(func(mRestServer *MultRestServer) {
				multRestServer = mRestServer
			}),
		)
		So(test.StartFxApp(app), ShouldBeNil)
		defer test.StopFxApp(app)
		//waitRestServerSteady("127.0.0.1:19086") // here use a token !!!

		Convey("When update the ratelimit config and receive 7 requests\n", func() {
			var res *http.Response
			var err error
			var failedNum, refusedNum, acceptedNum int32
			var retryAfter string
			var server *RestServer

			server, err = multRestServer.GetRestServerByName("myserver55")
			So(server, ShouldBeNil)
			So(err, ShouldNotBeNil)

			server, err = multRestServer.GetRestServerByName("myserver")
			So(server, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(server.SetRateLimitConfig(5, 4), ShouldBeNil)

			time.Sleep(1 * time.Second)

			for i := 0; i < 7; i++ {
				c := http.Client{Timeout: 4 * time.Second}
				res, err = c.Get("http://127.0.0.1:19086/")
				if err != nil {
					atomic.AddInt32(&failedNum, 1)
					return
				} else if res.StatusCode == http.StatusTooManyRequests {
					retryAfter = res.Header.Get("Retry-After")
					atomic.AddInt32(&refusedNum, 1)
				} else if res.StatusCode == http.StatusOK {
					atomic.AddInt32(&acceptedNum, 1)
				}
				defer res.Body.Close()

				time.Sleep(time.Millisecond)
			}

			Convey("Then 5 requests should be accepted and 2 requests should be refused\n", func() {
				So(failedNum, ShouldEqual, 0)
				So(acceptedNum, ShouldEqual, 5)
				So(refusedNum, ShouldEqual, 2)
				So(retryAfter, ShouldEqual, strconv.Itoa(5/4))
				// retry-after calculated by ceil(maxrequests / requestspersec)
			})
		})
	})
}

func TestMethodOverrideWithSetMethodOverride(t *testing.T) {
	Convey("Given a config file and a http-server with methodOverride \n", t, func() {
		content := `code-shooting:
  log:
    level: info
    encoder: plain
    outputPaths:
    - stdout
  rest-servers:
  - name: myserver
    addr: 127.0.0.1:19086
    readtimeout: 180s
    writetimeout: 60s
    maxheaderbytes: 16384
    rootpath: /mypath
    middlewares:
    - name: methodOverride`
		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		var multRestServer *MultRestServer

		app := fx.New(
			fx.Logger(logger.GetLogger().CreateStdLogger()),
			marsconfig.NewModule(confFile.Name()),
			NewModule(),
			fx.Invoke(myServerRegister),
			fx.Invoke(func(mRestServer *MultRestServer) {
				multRestServer = mRestServer
			}),
		)
		So(test.StartFxApp(app), ShouldBeNil)
		defer test.StopFxApp(app)
		//waitRestServerSteady("127.0.0.1:19086") // here use a token !!!

		Convey("When send a request with head contains X_HTTP_METHOD_OVERRIDE\n", func() {
			var res *http.Response
			var err error
			var server *RestServer

			server, err = multRestServer.GetRestServerByName("myserver")
			So(server, ShouldNotBeNil)
			So(err, ShouldBeNil)

			time.Sleep(1 * time.Second)

			c := http.Client{Timeout: 4 * time.Second}
			req, _ := http.NewRequest("GET", "http://127.0.0.1:19086/", nil)
			req.Header.Set("X-HTTP-Method-Override", "POST")
			res, err = c.Do(req)

			Convey("Then the response statusCode should be 405\n", func() {
				So(err, ShouldEqual, nil)
				So(res.StatusCode, ShouldEqual, 405)
			})
			defer res.Body.Close()
		})
	})
}
