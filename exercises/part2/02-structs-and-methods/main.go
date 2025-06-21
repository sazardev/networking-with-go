package main

import "fmt"

type Server struct {
	Host string
	Port int
}

func (s Server) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func main() {
	srv := Server{Host: "localhost", Port: 8080}
	fmt.Println("Server address:", srv.Address())
}
