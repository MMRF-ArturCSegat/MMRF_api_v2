package graph_model

import (
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
    "slices"
	"github.com/UFSM-Routelib/routelib_api/util"
	"github.com/icholy/utm"
)


func New_csvg(csv_file multipart.File, olt util.Coord, coord_limiter util.Square) (*CSV_Graph, error){
    csv_graph := CSV_Graph{
        Nodes: make(map[uint32]*GraphNode),
        Limiter: coord_limiter,
        Olt: olt,
    }

    reader := csv.NewReader(csv_file)
    reader.FieldsPerRecord = -1 
    reader.LazyQuotes = true
    reader.Comma = '\t'

    // for some reason, reading the whole file at once runs faster then reading the lines individually
    // benchmarked
    lines, err := reader.ReadAll()
    if err != nil{
        return nil, err
    }

    // define the column name indexes  
    // SUBESTACAO	ALIMENTADOR	FISICO_FONTE	FONTEX	FONTEY	FISICO_NO	NOX	NOY	NIVEL	TIPO_UNIDADE	I_ADM	USEC	TELECOMANDADA	MANOBRA_ANEL	ABERTURA_CARGA	LOAD_BUSTER	PROTECAO	FUSE	UREG	ENDERECO	U_REF	UR	UX	UCAP	Q_NOM	UTRD	PROPRIETARIO	S_NOM	UTRF	ESTADO	ESTADO_NORMAL	FASES	QTDFASES	BIT_FAS	MAT_FAS	COMPRIMENTO	BLOCO	R1	X1	R0	X0
    start_id, err1 := util.IndexOf("FISICO_FONTE", lines[0])
    start_lat, err2 := util.IndexOf("FONTEX", lines[0])
    start_lng, err3 := util.IndexOf("FONTEY", lines[0])
    end_id, err4 := util.IndexOf("FISICO_NO", lines[0])   
    end_lat, err5 := util.IndexOf("NOX", lines[0])
    end_lng, err6 := util.IndexOf("NOY", lines[0])

    errs := []error{err1, err2, err3, err4, err5, err6}
    for _, e := range errs {
        if e != nil {
            return nil, e
        }
    }

    zone := utm.LatLonZone(coord_limiter.Bot.Lat, coord_limiter.Top.Lng)

    for _, line := range lines{
        n1_id, id_err1 :=  strconv.ParseUint(line[start_id], 10, 32)
        n1_x := fix_float(line[start_lat])
        n1_y := fix_float(line[start_lng])

        n2_id, id_err2 :=  strconv.ParseUint(line[end_id], 10, 32)
        n2_x := fix_float(line[end_lat])
        n2_y := fix_float(line[end_lng])

        if id_err1 != nil || id_err2 != nil {
            // invalid id, probably bad node
            continue
        }

        n1_lat, n1_lng := n1_x, n1_y
        n2_lat, n2_lng := n2_x, n2_y

        if n1_lat > 999 || n1_lng > 999 || n2_lat > 999 || n2_lng >999 {
            n1_lat, n1_lng = zone.ToLatLon(n1_x, n1_y)
            n2_lat, n2_lng = zone.ToLatLon(n2_x, n2_y)
        }

        node1 := GraphNode{ID: uint32(n1_id), Lat: n1_lat, Lng: n1_lng}
        node2 := GraphNode{ID: uint32(n2_id), Lat: n2_lat, Lng: n2_lng}
    
        if !node1.GetCoord().IsInSquare(coord_limiter) && node2.GetCoord().IsInSquare(coord_limiter) {
            continue
        }
        
        if !node2.GetCoord().IsInSquare(coord_limiter) && node1.GetCoord().IsInSquare(coord_limiter) {
            continue
        }
        
        if !node2.GetCoord().IsInSquare(coord_limiter) && !node1.GetCoord().IsInSquare(coord_limiter) {
            continue
        }

        csv_graph.AddEdge(&node1, &node2)
    }
    
    old_n_conns := csv_graph.OltNecessaryConnections()
    old_len := len(old_n_conns)
    
    if old_len != 1{
        new_limiter := coord_limiter.Expand1km()

        aux_graph := CSV_Graph{
            Nodes: make(map[uint32]*GraphNode, len(csv_graph.Nodes)),
            Limiter: new_limiter,
            Olt: olt,
        }

        for _, line := range lines{
            n1_id, id_err1 :=  strconv.ParseUint(line[start_id], 10, 32)
            n1_x := fix_float(line[start_lat])
            n1_y := fix_float(line[start_lng])

            n2_id, id_err2 :=  strconv.ParseUint(line[end_id], 10, 32)
            n2_x := fix_float(line[end_lat])
            n2_y := fix_float(line[end_lng])

            if id_err1 != nil || id_err2 != nil {
                // invalid id, probably bad node
                continue
            }

            n1_lat, n1_lng := n1_x, n1_y
            n2_lat, n2_lng := n2_x, n2_y

            if n1_lat > 999 || n1_lng > 999 || n2_lat > 999 || n2_lng >999 {
                n1_lat, n1_lng = zone.ToLatLon(n1_x, n1_y)
                n2_lat, n2_lng = zone.ToLatLon(n2_x, n2_y)
            }

            node1 := GraphNode{ID: uint32(n1_id), Lat: n1_lat, Lng: n1_lng}
            node2 := GraphNode{ID: uint32(n2_id), Lat: n2_lat, Lng: n2_lng}

            if !node1.GetCoord().IsInSquare(new_limiter) && node2.GetCoord().IsInSquare(new_limiter) {
                continue
            }

            if !node2.GetCoord().IsInSquare(new_limiter) && node1.GetCoord().IsInSquare(new_limiter) {
                continue
            }

            if !node2.GetCoord().IsInSquare(new_limiter) && !node1.GetCoord().IsInSquare(new_limiter) {
                continue
            }

            aux_graph.AddEdge(&node1, &node2)
        }
        
        must_reach := make([]*GraphNode, 0, old_len)
        for _, id := range old_n_conns {
            node, err := aux_graph.FindNode(id)
            if err != nil {
                return nil, errors.New("Exapanded Graph differs too much from original: key connection node not found")
            }
            must_reach = append(must_reach, node)
        }
        
        close_sm, _ := csv_graph.ClosestNode(olt)
        close_bg, err := aux_graph.FindNode(close_sm.ID);

        if err != nil {
            return nil, errors.New("Exapanded Graph differs too much from original: entry node not found")
        }

        full_paths := make([]GraphPath, 0, old_len)
        
        for _, node := range must_reach {
            full_paths = append(full_paths, aux_graph.DepthFirstSearch(node, close_bg))
        }
        
        for _, path := range full_paths{
            for index := 1; index < len(path.Nodes) - 1; index++ {
                csv_graph.AddEdge(path.Nodes[index], path.Nodes[index - 1])
                csv_graph.AddEdge(path.Nodes[index], path.Nodes[index + 1])
            }
        }

        if len(csv_graph.OltNecessaryConnections()) != 1 {
            return nil, errors.New("Graph expansion was not suficient to connect sub-graphs")
        }
    }
    
    clean_bad_nodes(&csv_graph)
    csv_graph.Print()
    return &csv_graph, nil 
}

// weird function to ignore errors less verbousely
func fix_float(bad_float_string string) float64{
    good_float_string := strings.Replace(bad_float_string, ",", ".", 1)
    good_float, _ := strconv.ParseFloat(good_float_string, 64)
    return good_float
}

func clean_bad_nodes(csvg *CSV_Graph) {
    fmt.Printf("\ntotal was: %v\n", len(csvg.AllNodes()))
    deleted := 0
    for _, node := range csvg.AllNodes() {
        if len(node.NeighboursID) == 0 {
            delete(csvg.Nodes, node.ID)
            fmt.Printf("\ndeleted: %v\n", node.ID)
            deleted += 1
        }
    }
    fmt.Printf("\ndeleted was: %v\n", deleted)
}

// all that nodes the need to be connected to the OLT
// in order for it to have connections with all sub-graphs in the CSV_Graph
func (csvg * CSV_Graph) OltNecessaryConnections() []uint32 {
    connections := make([]uint32 ,0)
    invalid_nodes := make([]*GraphNode, 0)
    validate_node := func (node *GraphNode, reference util.Coord, dist float32) bool {
        if node.GetCoord().DistanceToInMeters(reference) < dist && len(node.NeighboursID) != 0 {
            return !(slices.ContainsFunc(invalid_nodes, func(n *GraphNode)bool{return n.ID == node.ID}))
        }
        return false
    }

    short, _ := csvg.ClosestNode(csvg.Olt)
    connections = append(connections, short.ID)
    invalid_nodes = csvg.AllVisitableNodesFrom(short, invalid_nodes)

    for {
        short, _, err := csvg.ClosestNodeFunc(csvg.Olt, validate_node)
        if err != nil {
            // err will be not nill once there are no sub-networks left to visit
            break
        }
        connections = append(connections, short.ID)
        invalid_nodes = csvg.AllVisitableNodesFrom(short, invalid_nodes)
    }
        
    return connections
}
