package graph_model

import (
	"encoding/csv"
	"errors"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/UFSM-Routelib/routelib_api/util"
)


func New_csvg(csv_file multipart.File) (*CSV_Graph, error){
    csv_graph := CSV_Graph{Nodes: make(map[int64]*GraphNode)}

    lines, err := csv.NewReader(csv_file).ReadAll()
    if err != nil{
        return nil, errors.New("Failed to read bad csv file")
    }

    // define the column name indexes  
    start_id, err1 := util.IndexOf("FISICO_FONTE", lines[0])
    start_lat, err2 := util.IndexOf("FONTEX", lines[0])
    start_lng, err3 := util.IndexOf("FONTEY", lines[0])
    end_id, err4 := util.IndexOf("FISICO_NO", lines[0])   
    end_lat, err5 := util.IndexOf("NOX", lines[0])
    end_lng, err6 := util.IndexOf("NOY", lines[0])

    if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil{
        return nil, errors.New("csv files has bad column names")
    }

    for _, line := range lines{
        n1_id, id_err1 :=  strconv.ParseInt(line[start_id], 10, 64)
        n1_lat := fix_float(line[start_lat])
        n1_lng := fix_float(line[start_lng])

        n2_id, id_err2 :=  strconv.ParseInt(line[end_id], 10, 64)
        n2_lat := fix_float(line[end_lat])
        n2_lng := fix_float(line[end_lng])

        if id_err1 != nil || id_err2 != nil {
            // invalid id, probably bad node
            continue
        }

        node1 := GraphNode{ID: n1_id, Lat: n1_lat, Lng: n1_lng}
        node2 := GraphNode{ID: n2_id, Lat: n2_lat, Lng: n2_lng}

        csv_graph.AddEdge(&node1, &node2)
    }
    return &csv_graph, err
}

// weird function to ignore errors less verbousely
func fix_float(bad_float_string string) float64{
    good_float_string := strings.Replace(bad_float_string, ",", ".", 1)
    good_float, _ := strconv.ParseFloat(good_float_string, 64)
    return good_float
}
