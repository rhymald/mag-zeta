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
	// id := c.Param("id")

	var buffer []play.Simplified
	_, spanPlayers := tracer.Start(ctx, "players")
	countOfPlayers, countOfFoes := 0, 0
	world.Lock() ; objLimit := len((*world).ByID)
	plimit, flimit := base.Round(math.Log2(float64( objLimit ))+1), base.Round(math.Sqrt(float64( objLimit )))
	first := [2]int{}// (*world).ByID[id].Path()[1]
	for _, each := range (*world).ByID { 
		distance := 0.0
		if countOfFoes + countOfPlayers == 0 {
			first = each.Path()[1]
		} else {
			distance = math.Sqrt( math.Pow(float64(each.Path()[1][0] - first[0]), 2) + math.Pow(float64(each.Path()[1][1] - first[1]), 2) )
		}
		if (*each).Current.IsNPC() == false { 
			if countOfPlayers < plimit && distance < 500 { buffer = append(buffer, (*each).Current.Simplify(each.Path())) }
			countOfPlayers++ 
		} else { 
			if countOfFoes < flimit && distance < 500 { buffer = append(buffer, (*each).Current.Simplify(each.Path())) }
			countOfFoes++ 
		}
		if countOfFoes + countOfPlayers >= plimit + flimit { break } 
	} // ; if countOfPlayers < 10 {break}}
	world.Unlock()
	span.SetAttributes(attribute.Int("Players", countOfPlayers))
	span.SetAttributes(attribute.Int("NPCs", countOfFoes))
	spanPlayers.End()

	_, spanResponse := tracer.Start(ctx, "responding")
	defer spanResponse.End()
	c.JSON(200, buffer) 
}
