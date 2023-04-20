package api

import (
	"rhymald/mag-zeta/play"
	"github.com/gin-gonic/gin"
	// For metrics:
	"github.com/prometheus/client_golang/prometheus/promhttp"
	// Create span:
	// "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/attribute"
)

var (
	foes []*play.Character
	players []*play.Character
	tracer = otel.Tracer("main")
)

func RunAPI() {
	router := gin.Default()
	router.GET("/", hiThere)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/spawn", newFoe)
	router.GET("/around", getAll)
	router.GET("/login", newPlayer)
	router.Run(":4917")
}

// TEST
func hiThere(c *gin.Context) { 
	ctx := (*c).Request.Context()
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ProjectsID","4917"))
	c.IndentedJSON(200, "Hello world!")
}


// MODIFY
func newFoe(c *gin.Context) { 
	foe := play.MakeNPC()
	err := foe.CalculateAttributes()
	if err == nil {
		foes = append(foes, foe)
		c.IndentedJSON(201, "Successfully spawned")
	} else {
		c.IndentedJSON(403, "Invalid foe character")
	}
}

func newPlayer(c *gin.Context) { 
	player := play.MakePlayer()
	err := player.CalculateAttributes()
	if err == nil {
		players = append(players, player)
		c.IndentedJSON(201, "Successfully logged in")
	} else {
		c.IndentedJSON(403, "Invalid player character")
	}
}


// READ
func getAll(c *gin.Context) { 
	var buffer []play.Simplified
	for _, each := range players { buffer = append(buffer, each.Simplify()) }
	for _, each := range foes { buffer = append(buffer, each.Simplify()) }
	c.IndentedJSON(200, buffer) 
}
