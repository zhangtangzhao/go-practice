package user

import (
	"context"
	"fmt"
	"github.com/go-practice/api/user"
	"google.golang.org/grpc"
)

var grpcClient user.UserClient

var GrpcConn grpc.ClientConn

func init(){
	GrpcConn,err := grpc.Dial("127.0.0.1:8800",grpc.WithInsecure())
	if err != nil{
		fmt.Printf("grpc dial err: ",err)
		return
	}

	grpcClient = user.NewUserClient(GrpcConn)


}

func GetClient(id int64) (*user.UserEntry, error){
	var user user.Request
	user.Id = id
	reps ,err := grpcClient.GetUser(context.Background(),&user)
	if err != nil{

	}
	return reps.UserEntry,err;
}