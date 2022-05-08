package main

import (
	"fmt"
	"github.com/go-practice/api/user"
	userService "github.com/go-practice/user/internal/service/user"
	"google.golang.org/grpc"
	"net"
)


func main(){
	grpcServer := grpc.NewServer()
	user.RegisterUserServer(grpcServer,new(userService.UserServiceImpl))

	listener, err := net.Listen("tcp","127.0.0.1:8800")

	if err != nil{
		fmt.Printf("Listen err : ",err)
		return
	}
	defer listener.Close()

	grpcServer.Serve(listener)
}