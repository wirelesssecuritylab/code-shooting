package test

import "os"

func DoActionsInFileRedirectedContext(src *os.File, des *os.File, action func()) {
	tmp := *src
	*src = *des
	*des = tmp

	defer func() {
		*des = *src
		*src = tmp
	}()

	action()
}
