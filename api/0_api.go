package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"encoding/json"
)

func RunAPI() {
	// router := gin.Default()
	router := gin.New()
	router.Use(gin.Recovery()) // to recover gin automatically
	router.Use(jsonLoggerMiddleware()) // we'll define it later
	router.GET("/", hiThere)
	router.GET("/around", getAll)
	router.GET("/around/:myplayerid", getAll)
	router.GET("/login", newPlayer)
	router.GET("/test", showGrid)
	router.GET("/test/spawn", newFoe)
	router.GET("/test/:id", showState)
	metrics := gin.New()
	metrics.Use(gin.Recovery()) // to recover gin automatically
	metrics.Use(jsonLoggerMiddleware()) // we'll define it later
	metrics.GET("/metrics", gin.WrapH(promhttp.Handler()))
	go func(){ router.Run(":4917") }()
	go func(){ metrics.Run(":9093") }()
	select {}
}

// TEST
func hiThere(c *gin.Context) { 
	c.IndentedJSON(200, "Hello world!")
}

func jsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})
			log["status_code"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			log["start_time"] = params.TimeStamp.Format("2006/01/02 - 15:04:05")
			log["remote_addr"] = params.ClientIP
			log["response_time"] = params.Latency.String()
 			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}
