package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunAPI() {
	router := gin.Default()
	router.GET("/", hiThere)
	router.GET("/spawn", newFoe)
	router.GET("/around", getAll)
	router.GET("/login", newPlayer)
	router.GET("/test", showGrid)
	router.GET("/test/:id", showState)
	metrics := gin.Default()
	metrics.GET("/metrics", gin.WrapH(promhttp.Handler()))
	go func(){ router.Run(":4917") }()
	go func(){ metrics.Run(":9093") }()
	select {}
}

// TEST
func hiThere(c *gin.Context) { 
	c.IndentedJSON(200, "Hello world!")
}
