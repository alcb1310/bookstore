package main

import "github.com/alcb1310/bookstore/internal/router"

func main() {
	var port uint16 = 42069

	s := router.New(port)
	if err := s.Router(); err != nil {
		panic(err)
	}
}
