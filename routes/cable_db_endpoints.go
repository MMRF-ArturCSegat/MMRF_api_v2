package routes

import (
	foc "github.com/UFSM-Routelib/routelib_api/fiber_optic_components"
	"github.com/gin-gonic/gin"
    "net/http"
    "strconv"
)


func add_cable(c * gin.Context){
    var cable foc.FiberCable 

	if err := c.ShouldBindJSON(&cable); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error in binding": err.Error(),})
		return
	}

    err := foc.AddObj(&cable)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error in database:": err.Error(),})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cable added": cable})
}


func get_cable(c * gin.Context){
    id_string := c.Param("id")
    id, _ := strconv.ParseUint(id_string, 10, 32)

    var obj foc.FiberCable

    cable, err := foc.GetOne(uint32(id), obj)

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"cable": cable})
}


func all_cables(c * gin.Context){
    var objs []foc.FiberCable

    cables, err := foc.GetAll(objs)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cables": cables,})
}
