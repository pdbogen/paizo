package main

import (
	"fmt"
	"github.com/pdbogen/paizo"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		panic("usage: " + os.Args[0] + " <username> <password>")
	}
	sess, err := paizo.NewSession()
	if err != nil {
		panic(err)
	}
	if err := sess.Authenticate(os.Args[0], os.Args[1]); err != nil {
		panic(err)
	}
	fmt.Println(sess)
}
