package file

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"code-shooting/infra/config/model"

	. "github.com/smartystreets/goconvey/convey"
)

func writeToFile(fileName string, content string) error {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	n, _ := f.Seek(0, os.SEEK_END)
	_, err = f.WriteAt([]byte(content), n)
	defer f.Close()
	return nil
}

func TestGetConfigFromFileSource(t *testing.T) {

	Convey("Given yaml file config source . \n", t, func() {
		file, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			file.Close()
			os.Remove(file.Name())
		}()
		file.WriteString(`
        code-shooting:
          apiversion: v1
          obj:
            name: RoundRobin
            type: Recreate
          slice:
          - name: a
            age: 1
          - name: b 
            age: 2`)
		file.Sync()
		source, err := NewFileSource(DefaultYamlSource, ParseYamlFile, file.Name())
		So(err, ShouldBeNil)
		stop := make(chan struct{})
		defer close(stop)
		source.Start(stop)
		time.Sleep(time.Second)

		Convey("When get configItem. \n", func() {
			configItem := source.Get("code-shooting.apiversion")
			Convey("Then get the configItem success. \n", func() {
				So(configItem, ShouldNotBeNil)
				So(configItem.Value, ShouldEqual, "v1")
			})
		})
	})
}

func TestFileSourceCreateEvents(t *testing.T) {

	Convey("Given yaml file config source . \n", t, func() {
		file, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			file.Close()
			os.Remove(file.Name())
		}()
		file.WriteString(`
        code-shooting:
          apiversion: v1
          obj:
            name: RoundRobin
            type: Recreate
          slice:
          - name: a
            age: 1
          - name: b 
            age: 2`)
		file.Sync()
		source, err := NewFileSource(DefaultYamlSource, ParseYamlFile, file.Name())
		So(err, ShouldBeNil)
		stop := make(chan struct{})
		defer close(stop)
		source.Start(stop)
		time.Sleep(time.Second)
		Convey("When create new config item. \n", func() {
			file.Close()
			content := `
        code-shooting:
          apiversion: v1
          obj:
            name: RoundRobin
            type: Recreate
            id: 1`
			err := writeToFile(file.Name(), content)
			So(err, ShouldBeNil)

			events := make([]*model.Event, 0)
			source.RegisterEventHandler("code-shooting.obj", func(es []*model.Event) {
				events = es
			})
			Convey("Then got create event. \n", func() {
				time.Sleep(2 * time.Second)
				So(len(events), ShouldEqual, 1)
				So(events[0].EventType, ShouldEqual, model.Create)
				So(events[0].Value, ShouldEqual, 1)

			})
		})
	})
}

func TestFileSourceUpdateEvents(t *testing.T) {

	Convey("Given yaml file config source . \n", t, func() {
		file, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			file.Close()
			os.Remove(file.Name())
		}()
		file.WriteString(`
        code-shooting:
          apiversion: v1
          obj:
            name: RoundRobin
            type: Recreate
          slice:
          - name: a
            age: 1
          - name: b 
            age: 2`)
		file.Sync()
		source, err := NewFileSource(DefaultYamlSource, ParseYamlFile, file.Name())
		So(err, ShouldBeNil)
		stop := make(chan struct{})
		defer close(stop)
		source.Start(stop)
		time.Sleep(time.Second)
		Convey("When  update code-shooting.obj.type \n", func() {
			file.Close()
			content := `
        code-shooting:
          apiversion: v1
          obj:
            name: RoundRobin
            type: Recreate1`
			err := writeToFile(file.Name(), content)
			So(err, ShouldBeNil)

			events := make([]*model.Event, 0)
			source.RegisterEventHandler("code-shooting.obj", func(es []*model.Event) {
				events = es
			})

			Convey("Then got update events. \n", func() {
				time.Sleep(2 * time.Second)
				So(len(events), ShouldEqual, 1)
				So(events[0].EventType, ShouldEqual, model.Update)
				So(events[0].Value, ShouldEqual, "Recreate1")

			})
		})
	})
}

func TestFileSourceDeletEvents(t *testing.T) {

	Convey("Given yaml file config source . \n", t, func() {
		file, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			file.Close()
			os.Remove(file.Name())
		}()
		file.WriteString(`
        code-shooting:
          apiversion: v1
          obj:
            name: RoundRobin
            type: Recreate
          slice:
          - name: a
            age: 1
          - name: b 
            age: 2`)
		file.Sync()
		source, err := NewFileSource(DefaultYamlSource, ParseYamlFile, file.Name())
		So(err, ShouldBeNil)
		stop := make(chan struct{})
		defer close(stop)
		source.Start(stop)
		time.Sleep(time.Second)
		Convey("When delete code-shooting.slice \n", func() {
			file.Close()
			content := `
        code-shooting:
          apiversion: v1
          obj:
            name: RoundRobin
            type: Recreate1`
			err := writeToFile(file.Name(), content)
			So(err, ShouldBeNil)

			events := make([]*model.Event, 0)
			source.RegisterEventHandler("code-shooting.slice", func(es []*model.Event) {
				events = es
			})

			Convey("Then got delete events. \n", func() {
				time.Sleep(2 * time.Second)
				So(len(events), ShouldEqual, 1)
				So(events[0].EventType, ShouldEqual, model.Delete)
			})
		})
	})
}

func TestFileSourceEvents(t *testing.T) {

	Convey("Given yaml file config source . \n", t, func() {
		file, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			file.Close()
			os.Remove(file.Name())
		}()
		file.WriteString(`
        code-shooting:
          apiversion: v1
          obj:
            name: RoundRobin
            type: Recreate
          slice:
          - name: a
            age: 1
          - name: b 
            age: 2`)
		file.Sync()
		source, err := NewFileSource(DefaultYamlSource, ParseYamlFile, file.Name())
		So(err, ShouldBeNil)
		stop := make(chan struct{})
		defer close(stop)
		source.Start(stop)
		time.Sleep(time.Second)
		Convey("When config file change. \n", func() {
			file.Close()
			content := `
        code-shooting:
          apiversion: v1
          obj:
            type: Recreate2
            id: 1
          slice:
          - name: a
            age: 1`
			err := writeToFile(file.Name(), content)
			So(err, ShouldBeNil)

			events := make([]*model.Event, 0)
			source.RegisterEventHandler("code-shooting", func(es []*model.Event) {
				events = es
			})

			Convey("Then get the configItem success. \n", func() {
				time.Sleep(2 * time.Second)
				configItem := source.Get("code-shooting.apiversion")
				So(configItem, ShouldNotBeNil)
				So(len(events), ShouldEqual, 4)
			})
		})
	})
}

func TestFileSourceForCreateAfterDeleteEvents(t *testing.T) {
	Convey("Given yaml file config source . \n", t, func() {
		file, _ := ioutil.TempFile("/tmp", "config-*.yaml")
		defer func() {
			file.Close()
			os.Remove(file.Name())
		}()
		file.WriteString(`
        code-shooting:
          apiversion: v1
          obj:
            name: RoundRobin
            type: Recreate
          slice:
          - name: a
            age: 1
          - name: b 
            age: 2`)
		file.Sync()
		file.Close()
		source, err := NewFileSource(DefaultYamlSource, ParseYamlFile, file.Name())
		So(err, ShouldBeNil)
		stop := make(chan struct{})
		defer close(stop)
		source.Start(stop)
		time.Sleep(time.Second)

		events := make([]*model.Event, 0)
		source.RegisterEventHandler("code-shooting.apiversion", func(es []*model.Event) {
			events = es
		})

		Convey("When delete file and create a new file which name is as same as the deleted, and modify apiversion two times", func() {
			os.Remove(file.Name())
			content := `
        code-shooting:
          apiversion: v2
          obj:
            name: RoundRobin
            type: Recreate
          slice:
          - name: a
            age: 1
          - name: b 
            age: 2`

			err := ioutil.WriteFile(file.Name(), []byte(content), 0777)
			So(err, ShouldBeNil)

			time.Sleep(1 * time.Second)

			content2 := `
        code-shooting:
          apiversion: v3
          obj:
            name: RoundRobin
            type: Recreate
          slice:
          - name: a
            age: 1
          - name: b 
            age: 2`

			err = ioutil.WriteFile(file.Name(), []byte(content2), 0777)
			So(err, ShouldBeNil)

			Convey("Then get the configItem success. \n", func() {
				time.Sleep(2 * time.Second)
				configItem := source.Get("code-shooting.apiversion")
				So(configItem, ShouldNotBeNil)
				So(configItem.Value, ShouldEqual, "v3")
				So(len(events), ShouldEqual, 1)
			})
		})
	})
}
