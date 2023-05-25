package api

import (
	"rhymald/mag-zeta/play"
	"github.com/gin-gonic/gin"
	// For metrics:
	"github.com/prometheus/client_golang/prometheus/promhttp"
	// Create span:
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	// "go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/trace"
	// "go.opentelemetry.io/otel/attribute"
	// "go.opentelemetry.io/otel/exporters/prometheus"
	// "go.opentelemetry.io/otel/metric"
	// "go.opentelemetry.io/otel/metric/instrument"
	// sdk "go.opentelemetry.io/otel/sdk/metric"
	"errors"
)

var (
	foes []*play.Character
	players []*play.Character
	// tracer = otel.Tracer("main")
)

func RunAPI() {
	// ctx := context.Background()
	router := gin.Default()
	router.Use(otelgin.Middleware("mag"))
	router.GET("/", hiThere)
	router.GET("/spawn", newFoe)
	router.GET("/around", getAll)
	router.GET("/login", newPlayer)
	metrics := gin.Default()
	metrics.GET("/metrics", gin.WrapH(promhttp.Handler()))
	go func(){ router.Run(":4917") }()
	go func(){ metrics.Run(":9093") }()
	select {}
}

// TEST
func hiThere(c *gin.Context) { 
	// ctx := (*c).Request.Context()
	// span := trace.SpanFromContext(ctx)
	// span.SetAttributes(attribute.String("ProjectsID","4917"))
	c.IndentedJSON(200, "Hello world!")
}


// MODIFY
func newFoe(c *gin.Context) { 
	foe := play.MakeNPC()
	err := foe.CalculateAttributes()
	if err == nil {
		foes = append(foes, foe)
		c.IndentedJSON(200, "Successfully spawned")
	} else {
		c.AbortWithError(500, errors.New("Invalid foe character"))
	}
}

func newPlayer(c *gin.Context) { 
	player := play.MakePlayer()
	err := player.CalculateAttributes()
	if err == nil {
		players = append(players, player)
		c.IndentedJSON(200, "Successfully logged in")
	} else {
		c.AbortWithError(500, errors.New("Invalid player character"))
	}
}


// READ
func getAll(c *gin.Context) { 
	var buffer []play.Simplified
	for _, each := range players { buffer = append(buffer, each.Simplify()) }
	for _, each := range foes { buffer = append(buffer, each.Simplify()) }
	c.IndentedJSON(200, buffer) 
}
