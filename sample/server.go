package main

import (
	"fmt"
	"github.com/smartwalle/pks"
)

func main() {
	var s = pks.New()

	s.Handle("p", func(req *pks.Request, rsp *pks.Response) error {
		fmt.Println(req.Header)
		return nil
	})

	s.RunWithName("ss")
}
