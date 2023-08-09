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
	takenID := "" // c.Request.Header["myplayerid"]
	if _, ok := c.Request.Header["myplayerid"] ; ok { 
		takenID = c.GetHeader("myplayerid") } else { takenID = c.Param("myplayerid") 
	}
	myPlayer := &play.State{} 
	if _, ok := (*world).ByID[takenID] ; ok { myPlayer = (*world).ByID[takenID] } else { myPlayer = nil }
	plimit, flimit := base.Round(math.Log2( float64(objLimit) )) + 4, 16 + base.Round(math.Sqrt( float64(objLimit) ))
	radius := math.Sqrt(3)*4000 ; first := [5][2]int{} ; if myPlayer != nil { 
		first = myPlayer.Path() 
		for i:=2; i<5; i++ { first[i][0] += -first[1][0] ; first[i][1] += -first[1][1] }
		// if first[0][0] > 0 { first[0][0] += -1000 } else { first[0][0] += 1000 }
		// angle := float64(first[0][0]) / 1000 * 180
		// first[1][0], first[1][1] = first[1][0] - base.Round(float64(radius/2)*math.Sin(angle)), first[1][1] - base.Round(float64(radius/2)*math.Cos(angle))
		buffer = append(buffer, (*myPlayer).Current.Simplify(first))
	}
	for id, each := range (*world).ByID { 
		if countOfFoes + countOfPlayers >= flimit { break } 
		path := each.Path()
		distance := math.Sqrt( math.Pow(float64(path[1][0] - first[1][0]), 2) + math.Pow(float64(path[1][1] - first[1][1]), 2) )
		if (*each).Current.IsNPC() == false { 
			if countOfPlayers < plimit && distance < float64(radius) && id != takenID { 
				for i:=1; i<5; i++ { path[i][0] += -first[1][0] ; path[i][1] += -first[1][1] }
				buffer = append(buffer, (*each).Current.Simplify(path)) 
				countOfPlayers++ 
			}
		} else { 
			if distance < float64(radius) { 
				for i:=1; i<5; i++ { path[i][0] += -first[1][0] ; path[i][1] += -first[1][1] }
				buffer = append(buffer, (*each).Current.Simplify(path)) 
				countOfFoes++ 
			}
		}
	} // ; if countOfPlayers < 10 {break}}
	world.Unlock()
	span.SetAttributes(attribute.Int("Players", countOfPlayers))
	span.SetAttributes(attribute.Int("NPCs", countOfFoes))
	spanPlayers.End()

	_, spanResponse := tracer.Start(ctx, "responding")
	defer spanResponse.End()
	c.JSON(200, buffer) 
}
