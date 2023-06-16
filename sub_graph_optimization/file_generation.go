package sub_graph_optimization

import (
	"fmt"
	"os"

	foc "github.com/UFSM-Routelib/routelib_api/fiber_optic_components"
	gm "github.com/UFSM-Routelib/routelib_api/graph_model"
)


func GenerateSubGraphOptimizationFile(sb *gm.CSV_Graph, cable_IDs, splicebox_IDs, uspliter_IDs, bspliter_IDs []uint) (*os.File, error){
    file_content := sb.Build_txt_file_content() + "\n"
    file_content += "Cable: " + fmt.Sprint(len(cable_IDs)) + "\n"
    for _, id := range cable_IDs{
        cable, err := foc.GetOne(id, &foc.FiberCable{})
        if err != nil{
            continue;
        }
        file_content += cable.String()
    }

    file_content += "SpliceBox: " + fmt.Sprint(len(splicebox_IDs)) + "\n"
    for _, id := range splicebox_IDs{
        box, err := foc.GetOne(id, &foc.FiberSpliceBox{})
        if err != nil{
            continue;
        }
        file_content += box.String()
    }

    file_content += "UnbalancedSpliter: " + fmt.Sprint(len(uspliter_IDs)) + "\n"
    for _, id := range uspliter_IDs{
        uspliter, err := foc.GetOne(id, &foc.FiberUnbalancedSpliter{})
        if err != nil{
            continue;
        }
        file_content += uspliter.String()
    }

    file_content += "BalancedSpliter: " + fmt.Sprint(len(bspliter_IDs)) + "\n"
    for _, id := range bspliter_IDs{
        bspliter, err := foc.GetOne(id, &foc.FiberBalancedSpliter{})
        if err != nil{
            continue;
        }
        file_content += bspliter.String()
    }

    file, err := os.Create("sub_graph.txt")
    if err != nil {
        return file, err
    }
    _, err = file.WriteString(file_content)
    if err != nil {
        return file, err
    }
    return file, nil
}

