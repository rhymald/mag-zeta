package api

import (
	"rhymald/mag-zeta/play"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"sync"
)

type Location struct {
	ByID map[int]*play.Character
	sync.Mutex
}

var (
	world = &Location{ ByID: make(map[int]*play.Character) }
	tracer = otel.Tracer("api")
)

func getAll(c *gin.Context) { 
	ctx, span := tracer.Start((*c).Request.Context(), "pull-all-objects")
	defer span.End()

	var buffer []play.Simplified
	_, spanPlayers := tracer.Start(ctx, "players")
	countOfPlayers := 0
	world.Lock()
	for _, each := range (*world).ByID { buffer = append(buffer, each.Simplify()) ; if each.IsNPC() == false { countOfPlayers++ }}
	world.Unlock()
	span.SetAttributes(attribute.Int("Players", countOfPlayers))
	span.SetAttributes(attribute.Int("NPCs", len(buffer)-countOfPlayers))
	spanPlayers.End()

	// _, spanNPC := tracer.Start(ctx, "npc")
	// world.Lock()
	// for _, each := range (*world).NPCs { buffer = append(buffer, each.Simplify()) }
	// world.Unlock()
	// span.SetAttributes(attribute.Int("NPCs", len(buffer)-countOfPlayers))
	// spanNPC.End()

	_, spanResponse := tracer.Start(ctx, "responding")
	defer spanResponse.End()
	c.IndentedJSON(200, buffer) 
}
