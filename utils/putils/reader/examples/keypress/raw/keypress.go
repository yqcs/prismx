package main

import (
	"log"
	"sync"
	"time"

	"prismx_cli/utils/putils/reader"
	stringsutil "prismx_cli/utils/putils/strings"
)

func main() {
	stdr := reader.KeyPressReader{
		Timeout: time.Duration(5 * time.Second),
		Once:    &sync.Once{},
		Raw:     true,
	}

	stdr.Start()
	defer stdr.Stop()

	for {
		data := make([]byte, 1)
		n, err := stdr.Read(data)
		if stringsutil.IsPrintable(string(data)) {
			log.Println(n, err)
		}

		if stringsutil.IsCTRLC(string(data)) {
			break
		}
	}
}
