package api

import (
	"rhymald/mag-zeta/play"
	"rhymald/mag-zeta/base"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"sync"
	"math"
)

type Location struct {
	ByID map[string]*play.State
	Grid struct {
		X map[string]int // axis
		Y map[string]int // axis
		T map[string]int // time
		Vec map[int]string // direction from 0, 0
		Rad map[int]string // distance from 0, 0
		Zero [2]int // 0, 0
	}
	// PosCache connection.to.table
	sync.Mutex
}

var (
	world = newWorld()
	tracer = otel.Tracer("api")
)

func newWorld() *Location {
	buffer := Location{ ByID: make(map[string]*play.State) }
	buffer.Grid.Zero = [2]int{0, 0}
	buffer.Grid.Y = make(map[string]int)
	buffer.Grid.X = make(map[string]int)
	buffer.Grid.T = make(map[string]int)
	buffer.Grid.Vec = make(map[int]string)
	buffer.Grid.Rad = make(map[int]string)
	return &buffer
}

func getAll(c *gin.Context) { 
	ctx, span := tracer.Start((*c).Request.Context(), "pull-all-objects")
	defer span.End()

	var buffer []play.Simplified
	_, spanPlayers := tracer.Start(ctx, "players")
	countOfPlayers, countOfFoes := 0, 0
	world.Lock()
	plimit, flimit := base.Round(math.Log2(float64(len((*world).ByID)))+1), base.Round(math.Sqrt(float64(len((*world).ByID))))
	for _, each := range (*world).ByID { 
		if (*each).Current.IsNPC() == false { 
			countOfPlayers++ 
			if countOfPlayers < plimit { buffer = append(buffer, (*each).Current.Simplify(each.Path())) }
		} else { 
			countOfFoes++ 
			if countOfFoes < flimit { buffer = append(buffer, (*each).Current.Simplify(each.Path())) }
		}
		if countOfFoes + countOfPlayers >= plimit + flimit { break } 
	} // ; if countOfPlayers < 10 {break}}
	world.Unlock()
	span.SetAttributes(attribute.Int("Players", countOfPlayers))
	span.SetAttributes(attribute.Int("NPCs", len(buffer)-countOfPlayers))
	spanPlayers.End()

	_, spanResponse := tracer.Start(ctx, "responding")
	defer spanResponse.End()
	c.JSON(200, buffer) 
}
