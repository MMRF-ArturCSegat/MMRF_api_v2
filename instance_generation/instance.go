package instance_generation

import (
    gm"github.com/UFSM-Routelib/routelib_api/graph_model"
    "github.com/UFSM-Routelib/routelib_api/util"
	foc "github.com/UFSM-Routelib/routelib_api/fiber_optic_components"
    "os"
    "fmt"
    "errors"
)

// in theory, Clients should be aligned to Paths 
// so the 0th client is the root that generated the 0th path 
type Instance struct {
    Paths               [][]gm.GraphPath     `json:"paths"`
    Clients             []util.Coord         `json:"clients"`
    OLT                 util.Coord           `json:"olt"`
    Cable_id            uint32               `json:"Cable_id"`
    Splicebox_id        uint32               `json:"Splicebox_id"`
    Uspliters_id        []uint32             `json:"Uspliters_id"`
    Bspliters_id        []uint32             `json:"Bspliters_id"`
}

func (i Instance) GetUspliters() []foc.FiberUnbalancedSpliter {
    uspliters := make([]foc.FiberUnbalancedSpliter, len(i.Uspliters_id))
    for idx, id := range i.Uspliters_id {
        err := foc.GetOne(id, &uspliters[idx])
        if err != nil{
            continue;
        }
    }
    return uspliters
}
func (i Instance) GetBspliters() []foc.FiberBalancedSpliter {
    bspliters := make([]foc.FiberBalancedSpliter, len(i.Bspliters_id))
    for idx, id := range i.Bspliters_id {
        err := foc.GetOne(id, &bspliters[idx])
        if err != nil{
            continue;
        }
    }
    return bspliters
}
func (i Instance) GetCable () *foc.FiberCable {
    cable := &foc.FiberCable{}
    foc.GetOne(i.Cable_id, cable)
    return cable
}
func (i Instance) GetSpliceBox () *foc.FiberSpliceBox {
    box := &foc.FiberSpliceBox{}
    foc.GetOne(i.Splicebox_id, box)
    return box
}

// all that nodes the need to be connected to the OLT
// in order for it to have connections with all sub-graphs in the CSV_Graph
func (i Instance) OltNecessaryColetions(csvg *gm.CSV_Graph) []uint32 {
    connections := make([]uint32 ,0)
    invalid_nodes := make([]*gm.GraphNode, len(csvg.Nodes)/2)
    validate_node := func (node *gm.GraphNode, reference util.Coord, dist float32) bool {
        if node.GetCoord().DistanceToInMeters(reference) < dist && !util.In(node, invalid_nodes){
            return true
        }
        return false
    }

    for {
        short, _, err := csvg.ClosestNodeFunc(i.OLT, validate_node)
        if err != nil {
            // err will be not nill when the file "validate_node" is strict enough that no nodes in csvg pass it
            break
        }
        connections = append(connections, short.ID)
        invalid_nodes = append(invalid_nodes, csvg.AllVisitableNodesFrom(short, make([]*gm.GraphNode, 0))...)
    }
        
    return connections
}

func (i Instance) GenerateSubGraphOptimizationFile(csvg * gm.CSV_Graph) (*os.File, error){
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
    
    if len(i.Clients) != len(i.Paths){
        return nil, errors.New("invalid sub_graphs or invalid clients")
    }

    for index, client := range i.Clients { // clients will be added as nodes
        nodes_content += fmt.Sprintf("0%v\t%v\t%v\n", index, client.Lat, client.Lng)
        nodes_count++
    }

    // pairs all nodes in a Path array to its respective client
    for index, client_paths := range i.Paths{
        read_nodes := make([]uint32, 0)
        for _, path := range client_paths{
            for _, node := range path.Nodes {
                if util.In(node.ID, read_nodes) {
                    continue
                }
                edges_content += fmt.Sprintf("0%v\t%v\tvirtual\n", index, node.ID)
                edges_count++
                read_nodes = append(read_nodes, node.ID)
            }
        }
    }

    file_content := fmt.Sprintf("Clients %v\n", len(i.Clients))
    file_content += fmt.Sprintf("Nodes %v\n", nodes_count + 1) // + 1 necesseary for OLT
    file_content += fmt.Sprintf("OLT \t%v\t%v\n", i.OLT.Lat, i.OLT.Lng)
    file_content += nodes_content
    file_content += fmt.Sprintf("Edges %v\n", edges_count)
    for _, e := range i.OltNecessaryColetions(csvg) {
        file_content += fmt.Sprintf("OLT\t%v\n", e)
    }
    file_content += edges_content

    // adding all fiber components
    file_content += "UnbalancedSpliter " + fmt.Sprint(len(i.Uspliters_id)) + "\n"
    for _, uspliter := range i.GetUspliters() { 
        file_content += uspliter.String()
    }

    file_content += "BalancedSpliter " + fmt.Sprint(len(i.Bspliters_id)) + "\n"
    for _, bspliter := range i.GetBspliters() { 
        file_content += bspliter.String()
    }

    file_content += "Cable " + "\n"
    file_content += i.GetCable().String()

    file_content += "SpliceBox " + "\n"
    file_content += i.GetSpliceBox().String()

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
