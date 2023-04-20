package api

// get foe
// get player
// new player

import (
	"rhymald/mag-zeta/play"
	"github.com/gin-gonic/gin"
)

var (
	foes []*play.Character
	players []*play.Character
)

func RunAPI() {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/spawn", newFoe)
	router.GET("/around", getAll)
	router.GET("/login", newPlayer)
	router.Run("0.0.0.0:4917")
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
