package api

import (
	"rhymald/mag-zeta/play"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	foes []*play.Character
	players []*play.Character
	tracer = otel.Tracer("api")
)

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
