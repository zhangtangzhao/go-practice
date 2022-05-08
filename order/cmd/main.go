package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-practice/order/internal/metrics"
	userClinet "github.com/go-practice/order/service/user"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func index(w http.ResponseWriter, r *http.Request){
	io.WriteString(w,"hello world")
}

func startServer(srv *http.Server) error{
	fmt.Println("http server start....")
	http.HandleFunc("/",index)
	http.HandleFunc("/getOrder",orderHandle)
	http.Handle("/metrics",promhttp.Handler())
	error := srv.ListenAndServe()
	return error
}

func orderHandle(writer http.ResponseWriter, request *http.Request) {
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
	query := request.URL.Query()
	id := query.Get("id")
	id64,err := strconv.ParseInt(id,10,64)
	if err != nil{
		io.WriteString(writer,"id not int")
		return
	}
	userEntity, err := userClinet.GetClient(id64)
	if err != nil{
		io.WriteString(writer,"user data err")
		return
	}
	data,err := json.Marshal(userEntity)
	io.WriteString(writer,string(data))

}

func main(){
	rootctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(rootctx)
	srv := &http.Server{Addr: ":8080"}
	g.Go(func() error {
		return startServer(srv)
	})

	g.Go(func() error {
		<- ctx.Done()
		fmt.Println("http server down...")
		return srv.Shutdown(ctx)
	})

	chanel := make(chan os.Signal)
	signal.Notify(chanel,syscall.SIGINT,syscall.SIGTERM,syscall.SIGKILL)

	g.Go(func() error {
		for  {
			select {
			case <- ctx.Done():
				userClinet.GrpcConn.Close()
				return ctx.Err()
			case s:= <- chanel:
				switch s{
				case syscall.SIGINT,syscall.SIGTERM,syscall.SIGKILL:
					userClinet.GrpcConn.Close()
					cancel()
				default:
					fmt.Println("unsignal syscall...")
				}
			}
		}
		return nil
	})

	if err := g.Wait(); err != nil && err != context.Canceled {
		fmt.Printf("errgroup err : %+v\n",err.Error())
	}

	fmt.Println("httpserver shutdown .....")
}
