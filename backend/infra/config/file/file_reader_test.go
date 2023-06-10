package file

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReadYamlFile(t *testing.T) {

	Convey("Given a config file ", t, func() {
		content := `code-shooting:
  apiversion: v1
  namespace: ${OPENPALETTE_NAMESPACE:director}
  struct:
    name: RoundRobin
    type: Recreate
  slice:
    - name: a
      age: 1
      addr:
        phone:
          - 123
          - 456
        street: s
    - name: b 
      age: 2`

		confFile, _ := ioutil.TempFile(".", "config-*.yml")
		defer func() {
			confFile.Close()
			os.Remove(confFile.Name())
		}()

		confFile.WriteString(content)
		confFile.Sync()

		Convey("When read yaml conf file. \n", func() {
			conf, err := ParseYamlFile([]string{confFile.Name()})
			Convey("Then return the configItems. \n", func() {
				So(err, ShouldBeNil)
				So(len(conf), ShouldEqual, 5)
				So(conf["code-shooting#namespace"], ShouldEqual, "director")
			})
		})
	})
}
