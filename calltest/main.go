package main

import (
	"context"
	"fmt"
)

type Calltest struct{}

func (m *Calltest) CheckTest(ctx context.Context) {
	fmt.Println("I'm within calltest")
}
