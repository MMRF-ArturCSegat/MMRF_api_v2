package instance_generation

import (
	"fmt"
	"os"
	foc "github.com/UFSM-Routelib/routelib_api/fiber_optic_components"
	gm "github.com/UFSM-Routelib/routelib_api/graph_model"
	"github.com/UFSM-Routelib/routelib_api/util"
)


func GenerateSubGraphOptimizationFile(sub_graphs []*gm.CSV_Graph, clients []util.Coord, entrys, cable_IDs, splicebox_IDs, uspliter_IDs, bspliter_IDs []uint32) (*os.File, error){
    // in theory, entys should be aligned to sub_graphs
    // so the 0th entry is the root that generated the 0th sub_graph

    file_content := fmt.Sprintf("Clients %v\n", len(clients))
        
    // adding all nodes and edges
    read_nodes := make(map[uint32]bool)
    read_nodes_edges := make(map[uint32]bool)
    nodes_content := ""
    edges_content := ""
    edges_count := 0
    for _, sb := range sub_graphs {
        graph_content, _ := sb.String_from_all_nodes(read_nodes)
        nodes_content += graph_content

        graph_content_edges, count := sb.String_from_all_edges(read_nodes_edges)
        edges_count += count
        edges_content += graph_content_edges
    }
    file_content += fmt.Sprintf("Nodes %v\n", (len(read_nodes) + len(clients) + 1)) // weird adding:
    file_content += "OLT\tOLT\tOLT\n" // temporary                                  // the nodes category includes the nodes + clients + the OLT
    file_content += nodes_content
    for index, client := range clients { // clients will be added as nodes
        file_content += fmt.Sprintf("0%v\t%v\t%v\n", index, client.Lat, client.Lng)
    }
    file_content += fmt.Sprintf("Edges %v\n", edges_count)
    file_content += edges_content

    // adding virtual edges, optimization of the combinatory algo
    // TODO: add proper reference 

    for index, sb := range sub_graphs{
        for _, node := range sb.AllNodes(){
            file_content += fmt.Sprintf("%v\t%v\tvirtual\n", entrys[index], node.ID)
        }
    }

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

