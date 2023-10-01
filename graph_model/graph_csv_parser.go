package graph_model

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/UFSM-Routelib/routelib_api/util"
	"github.com/icholy/utm"
)


func New_csvg(csv_file multipart.File, coord_limiter util.Square) (*CSV_Graph, error){
    csv_graph := CSV_Graph{Nodes: make(map[uint32]*GraphNode)}

    reader := csv.NewReader(csv_file)
    reader.FieldsPerRecord = -1 
    reader.LazyQuotes = true
    reader.Comma = '\t'

    // for some reason, reading the whole file at once runs faster then reading the lines individually
    // bench marked
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
            // onnly 1 is invalid o add 2
            csv_graph.AddNode(&node2)
            continue
        }
        
        if !node2.GetCoord().IsInSquare(coord_limiter) && node1.GetCoord().IsInSquare(coord_limiter) {
            csv_graph.AddNode(&node1)
            continue
        }
        
        if !node2.GetCoord().IsInSquare(coord_limiter) && !node1.GetCoord().IsInSquare(coord_limiter) {
            continue
        }

        csv_graph.AddEdge(&node1, &node2)
    }
    
    print("\nbefore crash\n")
    clean_bad_nodes(&csv_graph)
    print("\nafter crash\n")
    csv_graph.Print()
    print("\nafter crash\n")
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
