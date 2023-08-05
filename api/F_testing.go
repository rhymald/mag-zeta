package api

import (
	"github.com/gin-gonic/gin"
)

func showGrid(c *gin.Context) { 
	c.IndentedJSON(200, *world) 
}

func showState(c *gin.Context) { 
	id := c.Param("id")
	c.IndentedJSON(200, (*world).ByID[ id ]) 
}
