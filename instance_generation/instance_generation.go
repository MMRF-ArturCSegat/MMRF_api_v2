package instance_generation

import (
	"errors"
	"fmt"
	"os"

	foc "github.com/UFSM-Routelib/routelib_api/fiber_optic_components"
	gm "github.com/UFSM-Routelib/routelib_api/graph_model"
	"github.com/UFSM-Routelib/routelib_api/util"
)


func GenerateSubGraphOptimizationFile(csvg *gm.CSV_Graph,sub_graphs []*gm.CSV_Graph, OLT util.Coord, clients []util.Coord, cable_IDs, splicebox_IDs, uspliter_IDs, bspliter_IDs []uint32) (*os.File, error){
    // in theory, entys should be aligned to sub_graphs
    // so the 0th entry is the root that generated the 0th sub_graph
    nodes_content := ""
    edges_content := ""
    edges_count := 0
    nodes_count := 0
    nodes := csvg.AllNodes()
    for _, node := range nodes{
        nodes_content += node.String() + "\n"
        nodes_count++
        for _, neighbour_id := range node.NeighboursID{
            edges_content += fmt.Sprintf("%v\t%v\n", node.ID, neighbour_id)
            edges_count++
        }
    }
    
    if len(clients) != len(sub_graphs){
        return nil, errors.New("invalid sub_graphs or invalid clients")
    }

    for index, client := range clients { // clients will be added as nodes
        nodes_content += fmt.Sprintf("0%v\t%v\t%v\n", index, client.Lat, client.Lng)
        nodes_count++
    }

    for index, sb := range sub_graphs{
        for _, node := range sb.AllNodes(){
            edges_content += fmt.Sprintf("0%v\t%v\tvirtual\n", index, node.ID)
            edges_count++
        }
    }

    file_content := fmt.Sprintf("Clients %v\n", len(clients))
    file_content += fmt.Sprintf("Nodes %v\n", nodes_count + 1) // + 1 necesseary for OLT
    file_content += fmt.Sprintf("OLT \t%v\t%v\n", OLT.Lat, OLT.Lng)
    file_content += nodes_content
    file_content += fmt.Sprintf("Edges %v\n", edges_count)
    file_content += edges_content

    // adding all fiber components
    file_content += "Cable " + fmt.Sprint(len(cable_IDs)) + "\n"
        for _, id := range cable_IDs{
        cable, err := foc.GetOne(id, &foc.FiberCable{})
        if err != nil{
            continue;
        }
        file_content += cable.String()
    }

    file_content += "SpliceBox " + fmt.Sprint(len(splicebox_IDs)) + "\n"
    for _, id := range splicebox_IDs{
        box, err := foc.GetOne(id, &foc.FiberSpliceBox{})
        if err != nil{
            continue;
        }
        file_content += box.String()
    }

    file_content += "UnbalancedSpliter " + fmt.Sprint(len(uspliter_IDs)) + "\n"
    for _, id := range uspliter_IDs{
        uspliter, err := foc.GetOne(id, &foc.FiberUnbalancedSpliter{})
        if err != nil{
            continue;
        }
        file_content += uspliter.String()
    }

    file_content += "BalancedSpliter " + fmt.Sprint(len(bspliter_IDs)) + "\n"
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

