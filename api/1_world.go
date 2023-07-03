package api

import (
	"rhymald/mag-zeta/play"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"sync"
)

type Location struct {
	ByID map[int]*play.State
	sync.Mutex
}

var (
	world = &Location{ ByID: make(map[int]*play.State) }
	tracer = otel.Tracer("api")
)

func getAll(c *gin.Context) { 
	ctx, span := tracer.Start((*c).Request.Context(), "pull-all-objects")
	defer span.End()

	var buffer []play.Simplified
	_, spanPlayers := tracer.Start(ctx, "players")
	countOfPlayers := 0
	world.Lock()
	for _, each := range (*world).ByID { buffer = append(buffer, (*each).Current.Simplify()) ; if (*each).Current.IsNPC() == false { countOfPlayers++ }}
	world.Unlock()
	span.SetAttributes(attribute.Int("Players", countOfPlayers))
	span.SetAttributes(attribute.Int("NPCs", len(buffer)-countOfPlayers))
	spanPlayers.End()

	_, spanResponse := tracer.Start(ctx, "responding")
	defer spanResponse.End()
	c.IndentedJSON(200, buffer) 
}
