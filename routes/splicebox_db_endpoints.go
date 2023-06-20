package routes

import (
	foc "github.com/UFSM-Routelib/routelib_api/fiber_optic_components"
	"github.com/gin-gonic/gin"
    "net/http"
    "strconv"
)


func add_box(c * gin.Context){
    var box foc.FiberSpliceBox

	if err := c.ShouldBindJSON(&box); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error in binding": err.Error(),})
		return
	}

    err := foc.AddObj(&box)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error in database:": err.Error(),})
		return
	}

	c.JSON(http.StatusOK, gin.H{"box added": box})
}


func get_box(c * gin.Context){
    id_string := c.Param("id")
    id, _ := strconv.ParseUint(id_string, 10, 32)

    var obj foc.FiberSpliceBox

    box, err := foc.GetOne(uint32(id), obj)

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"box": box})
}


func all_boxes(c * gin.Context){
    var objs []foc.FiberSpliceBox

    boxes, err := foc.GetAll(objs)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"boxes": boxes,})
}
