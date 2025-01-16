package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	path := "/root/.kube/config"
	d := filepath.Dir(path)
	fmt.Println(d)

}
