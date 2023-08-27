package routes

import (
	foc "github.com/UFSM-Routelib/routelib_api/fiber_optic_components"
	"github.com/gin-gonic/gin"
    "net/http"
    "strconv"
)


func add_b_spliter(c * gin.Context){
    var fbs foc.FiberBalancedSpliter 

	if err := c.ShouldBindJSON(&fbs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error in binding": err.Error(),})
		return
	}

    err := foc.AddObj(&fbs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error in database:": err.Error(),})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fbs added": fbs})
}


func get_b_spliter(c * gin.Context){
    id_string := c.Param("id")
    id, _ := strconv.ParseUint(id_string, 10, 32)

    var bspliter foc.FiberBalancedSpliter

    err := foc.GetOne(uint32(id), &bspliter)

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"bspliter": bspliter})
}


func all_b_spliters(c * gin.Context){
    var objs []foc.FiberBalancedSpliter

    bspliters, err := foc.GetAll(objs)
    
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bspliters": bspliters,})
}
