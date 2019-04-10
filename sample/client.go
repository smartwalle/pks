package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/pks"
)

func main() {
	var c = pks.New()

	fmt.Println(c.Request(context.Background(), "ss", "p", nil, nil))
}
