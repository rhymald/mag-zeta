package api

import (
	"rhymald/mag-zeta/play"
	"rhymald/mag-zeta/connect"
	"rhymald/mag-zeta/base"
	"github.com/jackc/pgx/v5"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"sync"
	"math"
)

type Location struct {
	ID int
	ByID map[string]*play.State
	Writer *pgx.Conn
	// PosCache connection.to.table
	sync.Mutex
}

var (
	world = newWorld()
	tracer = otel.Tracer("api")
	GridCache = make(chan map[string][][3]int)
)

func newWorld() *Location {
	base.Wait(4096)
	buffer := Location{ ByID: make(map[string]*play.State), ID: base.EpochNS(), Writer: connect.ConnectCacheDB()[0] }
	location := &buffer
	go func(){ location.GridWriter_ByPush(GridCache) }()
	return location
}

func getAll(c *gin.Context) { 
	ctx, span := tracer.Start((*c).Request.Context(), "pull-all-objects")
	defer span.End()
	var buffer []play.Simplified
	_, spanPlayers := tracer.Start(ctx, "players")
	countOfPlayers, countOfFoes := 0, 0
	world.Lock() ; objLimit := len((*world).ByID)
	takenID := c.Param("myplayerid")
	myPlayer := &play.State{} 
	if _, ok := (*world).ByID[takenID] ; ok { myPlayer = (*world).ByID[takenID] } else { myPlayer = nil }
	plimit, flimit := base.Round(math.Log2( float64(objLimit) )) + 4, 16 + base.Round(math.Sqrt( float64(objLimit) ))
	first := [5][2]int{} ; if myPlayer != nil { 
		first = myPlayer.Path() 
		buffer = append(buffer, (*myPlayer).Current.Simplify(first, first[1]))
	}
	for id, each := range (*world).ByID { 
		distance := math.Sqrt( math.Pow(float64(each.Path()[1][0] - first[1][0]), 2) + math.Pow(float64(each.Path()[1][1] - first[1][1]), 2) )
		if (*each).Current.IsNPC() == false { 
			if countOfPlayers < plimit && distance < math.Sqrt2*1000 && id != takenID { buffer = append(buffer, (*each).Current.Simplify(each.Path(), first[1])) }
			countOfPlayers++ 
		} else { 
			if countOfFoes < flimit && distance < math.Sqrt2*1000 { buffer = append(buffer, (*each).Current.Simplify(each.Path(), first[1])) }
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
	c.IndentedJSON(200, buffer) 
}
