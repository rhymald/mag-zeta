package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"encoding/json"
	"fmt"
)

type LogEvent struct {
	Status  int    `json:"Status"`
	Time    string `json:"Time"`
	Latency string `json:"Latency"`
	Source  string `json:"Source"`
	Method  string `json:"Method"`
	Path    string `json:"Path"`
}

func RunAPI() {
	// router := gin.Default()
	router := gin.New()
	router.Use(gin.Recovery()) // to recover gin automatically
	router.Use(jsonLoggerMiddleware())
	router.GET("/", hiThere)
	router.GET("/around", getAll)
	router.GET("/around/:myplayerid", getAll)
	router.GET("/login", newPlayer)
	router.GET("/test", showGrid)
	router.GET("/test/spawn", newFoe)
	router.GET("/test/:id", showState)
	metrics := gin.New()
	metrics.Use(gin.Recovery()) // to recover gin automatically
	metrics.Use(jsonLoggerMiddleware())
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
			if len(params.Path) >= 7 {
				if params.StatusCode == 200 && params.Path[:7] == "/around" { return "" }
			}
			log := LogEvent{
				Status:  params.StatusCode,
				Latency: fmt.Sprintf("%0.3fms", float64(params.Latency.Microseconds())/1000),
				Method:  params.Method,
				Path:    params.Path,
				Time:    params.TimeStamp.Format("2006/01/02 15:04:05.999"),
				Source:  params.ClientIP,
			}
 			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}
