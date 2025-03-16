package main

import "fmt"

// responsible for instantiation and
// start up of our go application
func Run() error {
	return nil
}

func main() {
	fmt.Println("Go rest api course")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
