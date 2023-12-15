package main

import (
	"One.2/api"
	"github.com/gin-gonic/gin"
)

func main() {
	api.InitRouter()
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})
}
