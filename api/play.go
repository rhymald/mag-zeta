package api

import (
	"errors"
	"rhymald/mag-zeta/play"
	"github.com/gin-gonic/gin"
	// For metrics:
	"github.com/prometheus/client_golang/prometheus/promhttp"
	// Create span:
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	foes []*play.Character
	players []*play.Character
	tracer = otel.Tracer("api")
)

func RunAPI() {
	router := gin.Default()
	// router.Use(otelgin.Middleware("mag"))
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
	c.IndentedJSON(200, "Hello world!")
}


// MODIFY
func newFoe(c *gin.Context) { 
	ctx, span := tracer.Start((*c).Request.Context(), "spawn-foe")
	defer span.End()

	_, spanGenerate := tracer.Start(ctx, "generating-basic-stats")
	foe := play.MakeNPC()
	spanGenerate.AddEvent("Character generated with ID: 123456789-1234-1-1234567")
	span.SetAttributes(attribute.String("CharacterID","123456789-1234-1-1234567"))
	spanGenerate.End()

	_, spanCalculate := tracer.Start(ctx, "calculating-attributes-from-basic")
	err := foe.CalculateAttributes()
	if err != nil { spanCalculate.RecordError(err) }
	spanCalculate.End()

	_, spanResponse := tracer.Start(ctx, "responding")
	defer spanResponse.End()
	if err == nil {
		foes = append(foes, foe)
		c.IndentedJSON(200, "Successfully spawned")
	} else {
		c.AbortWithError(500, errors.New("Invalid foe character"))
	}
}

func newPlayer(c *gin.Context) { 
	ctx, span := tracer.Start((*c).Request.Context(), "login-player")
	defer span.End()

	_, spanGenerate := tracer.Start(ctx, "generating-basic-stats")
	player := play.MakePlayer()
	spanGenerate.AddEvent("Character generated with ID: 123456789-1234-1-1234567")
	span.SetAttributes(attribute.String("CharacterID","123456789-1234-1-1234567"))
	spanGenerate.End()

	_, spanCalculate := tracer.Start(ctx, "calculating-attributes-from-basic")
	err := player.CalculateAttributes()
	if err != nil { spanCalculate.RecordError(err) }
	spanCalculate.End()

	_, spanResponse := tracer.Start(ctx, "responding")
	defer spanResponse.End()
	if err == nil {
		players = append(players, player)
		c.IndentedJSON(200, "Successfully logged in")
	} else {
		c.AbortWithError(500, errors.New("Invalid player character"))
	}
}


// READ
func getAll(c *gin.Context) { 
	ctx, span := tracer.Start((*c).Request.Context(), "pull-all-objects")
	defer span.End()

	var buffer []play.Simplified
	_, spanPlayers := tracer.Start(ctx, "players")
	for _, each := range players { buffer = append(buffer, each.Simplify()) }
	countOfPlayers := len(buffer)
	span.SetAttributes(attribute.Int("Players", countOfPlayers))
	spanPlayers.End()

	_, spanNPC := tracer.Start(ctx, "npc")
	for _, each := range foes { buffer = append(buffer, each.Simplify()) }
	span.SetAttributes(attribute.Int("NPCs", len(buffer)-countOfPlayers))
	spanNPC.End()

	_, spanResponse := tracer.Start(ctx, "responding")
	defer spanResponse.End()
	c.IndentedJSON(200, buffer) 
}
