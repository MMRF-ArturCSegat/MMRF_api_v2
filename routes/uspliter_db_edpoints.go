package routes

import (
	foc "github.com/UFSM-Routelib/routelib_api/fiber_optic_components"
	"github.com/gin-gonic/gin"
    "net/http"
    "strconv"
)

func add_u_spliter(c * gin.Context){
    var fus foc.FiberUnbalancedSpliter 

	if err := c.ShouldBindJSON(&fus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error in binding": err.Error(),})
		return
	}

    err := foc.AddObj(&fus)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error in database:": err.Error(),})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fus added": fus})
}


func get_u_spliter(c * gin.Context){
    id_string := c.Param("id")
    id, _ := strconv.ParseUint(id_string, 10, 32)

    var uspliter foc.FiberUnbalancedSpliter

    err := foc.GetOne(uint32(id), &uspliter)

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"uspliter": uspliter})
}


func all_u_spliters(c * gin.Context){
    var objs []foc.FiberUnbalancedSpliter

    uspliters, err := foc.GetAll(objs)
    
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"uspliters": uspliters,})
}
