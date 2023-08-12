package http

import (
	"fmt"
	"midscale/midscale/app/controller/http/handler/group"
	"midscale/midscale/app/data/mgr/transport/key"
	"midscale/midscale/app/data/model/define"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
}

func init() {
	key.SetUp()
}

func NewHttpServer() *HttpServer {
	return &HttpServer{}
}
func (httpServer *HttpServer) Start() {

	go httpServer.StartServerForRemote(define.DefaultListenAddrForRemote)
	go httpServer.StartServerForLocal(define.DefaultListenAddrForLocal)

	waitSignal()
}

func waitSignal() {
	quitChan := make(chan os.Signal, 64)
	signal.Notify(quitChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for s := range quitChan {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Printf("signal %v\n", s)
			gracefullExit()
		default:
			// fmt.Println("other signal", s)
		}
	}
}

func gracefullExit() {
	os.Exit(0)
}

func (httpServer *HttpServer) StartServerForRemote(listenAddr string) error {
	router := gin.Default()

	ms := router.Group("/ms")
	{
		v1 := ms.Group("/v1")
		{
			api := v1.Group("/api")
			{
				api.GET("/ping", func(c *gin.Context) {
					c.JSON(200, gin.H{
						"message": "pong",
					})
				})

				grp := api.Group("/group")
				{
					grp.POST("/apply", group.Apply)
				}

			}
		}

		return router.Run(listenAddr)
	}
}

func (httpServer *HttpServer) StartServerForLocal(listenAddr string) error {
	router := gin.Default()

	ms := router.Group("/ms")
	{
		v1 := ms.Group("/v1")
		{
			api := v1.Group("/api")
			{
				api.GET("/ping", func(c *gin.Context) {
					c.JSON(200, gin.H{
						"message": "pong",
					})
				})

				grp := api.Group("/group")
				{
					grp.POST("/create", group.Create)
					grp.POST("/join", group.Join)
					grp.GET("/info/:GroupName", group.GetInfo)
					grp.POST("/tunnel", group.Tunnel)
				}

			}
		}

		return router.Run(listenAddr)
	}
}
