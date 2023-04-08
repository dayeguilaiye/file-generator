package utils

import "io"

func HandleClose(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(err)
	}
}
