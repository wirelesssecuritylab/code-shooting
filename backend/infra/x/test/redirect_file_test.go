package test

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"testing"
)

func TestRedirectFile(t *testing.T) {
	Convey("Given a redirect target file for stdout", t, func() {
		f, err := ioutil.TempFile(".", "des-*.log")
		So(err, ShouldBeNil)
		defer func() {
			f.Close()
			os.Remove(f.Name())
		}()

		Convey("When output \"redirect stdout to file\" to the console(stdout)", func() {
			DoActionsInFileRedirectedContext(os.Stdout, f, func() {
				fmt.Println("redirect stdout to file")
			})

			Convey("Then file content should be \"redirect stdout to file\"", func() {
				expect := "redirect stdout to file\n"
				content, _ := ioutil.ReadFile(f.Name())
				So(string(content), ShouldEqual, expect)
			})
		})
	})
}
