package two

import (
	"fmt"
	"go_mod_demo/first"
)

func init() {
	fmt.Println("two init")
}

func Two() {
	first.First()
}