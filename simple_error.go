package main 

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("/doesnt.exist")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(f.Name(), "opened successfully")
}