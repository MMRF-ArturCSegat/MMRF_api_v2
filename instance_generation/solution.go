package instance_generation

import (
	"bufio"
	"bytes"
	"mime/multipart"
	"strconv"
	"strings"
	foc "github.com/UFSM-Routelib/routelib_api/fiber_optic_components"
	gm "github.com/UFSM-Routelib/routelib_api/graph_model"
	"github.com/UFSM-Routelib/routelib_api/util"
)

type solutionOrVirtNet interface {
    addPath([2]util.Coord)
} 

type Solution struct {
    Path                [][2]util.Coord
    BspliterMap         map[uint32]foc.FiberBalancedSpliter
    UspliterMap         map[uint32]foc.FiberUnbalancedSpliter
    CableMap            map[uint32]foc.FiberCable
    SpliceboxNodesId    []uint32
}
func (s * Solution) addPath(c [2]util.Coord){s.Path = append(s.Path, c)}

type VirtualNetwork struct {
    Path                [][2]util.Coord
    BspliterMap         map[uint32]foc.FiberBalancedSpliter
}
func (s * VirtualNetwork) addPath(c [2]util.Coord){s.Path = append(s.Path, c)}

func loadInfostructure(obj_path solutionOrVirtNet, info_buffer [][]byte, csvg  *gm.CSV_Graph) {
    for _, line := range info_buffer {
        items := strings.FieldsFunc(string(line), func(c rune)bool{return c == '#'})
        values := strings.FieldsFunc(items[len(items) - 1], func(c rune)bool{return c == ' '})
        items[len(items) - 1] = values[0]
        value_str := values[1]
        value, _ := strconv.ParseFloat(value_str, 32)
        if value <= 0 {
            continue
        }

        id1, _ := strconv.ParseUint(items[1], 10, 32) // ignores all error's trusting the file
        id2, _ := strconv.ParseUint(items[2], 10, 32)
        node1, _ := csvg.FindNode(uint32(id1))
        node2, _ := csvg.FindNode(uint32(id2))
        
        obj_path.addPath([2]util.Coord{node1.GetCoord(), node2.GetCoord()})
    }
}

func ParseSolutionFile(inst Instance, csvg *gm.CSV_Graph, sol_file multipart.File) (*Solution, *VirtualNetwork) {
    scanner := bufio.NewScanner(sol_file)

    // solution buffers
    y_buffer := make([][]byte, 0) // arcs
    a_buffer := make([][]byte, 0) // unbalanced spliters
    b_buffer := make([][]byte, 0) // balanced spliters
    u_buffer := make([][]byte, 0) // spliceboxes
    
    // virt net buffers
    x2_buffer := make([][]byte, 0) // arcs of the virtual network
    b2_buffer := make([][]byte, 0) // balanced spliters of the virutal network

    for scanner.Scan() {
        line := scanner.Bytes()
        
        if len(line) == 0 {
            continue
        }
            
        // the # cheking is to separte lines like "a#1#2" from lines like "another bla bla bla"
        if bytes.Compare(line[0:2], []byte("x2#")) == 0 {
            x2_buffer = append(x2_buffer, line)
        } else if bytes.Compare(line[0:2], []byte("b2#")) == 0 {
            b2_buffer = append(b2_buffer, line)
        } else if line[0] == 'y' && line[1] == '#' {
            y_buffer = append(y_buffer, line)
        } else if line[0] == 'a' && line[1] == '#' {
            a_buffer = append(a_buffer, line)
        } else if line[0] == 'b' && line[1] == '#' {
            b_buffer = append(b_buffer, line)
        } else if line[0] == 'u' && line[1] == '#' {
            u_buffer = append(u_buffer, line)
        }
    }
    
    solution := Solution{
        Path: make([][2]util.Coord, 0, len(y_buffer)),
        BspliterMap: make(map[uint32]foc.FiberBalancedSpliter, len(b_buffer)),
        UspliterMap: make(map[uint32]foc.FiberUnbalancedSpliter, len(a_buffer)),
        CableMap: make(map[uint32]foc.FiberCable),
        SpliceboxNodesId: make([]uint32, 0, len(u_buffer)),
    }
    virtual_net := VirtualNetwork {
        Path: make([][2]util.Coord, 0, len(x2_buffer)),
        BspliterMap: make(map[uint32]foc.FiberBalancedSpliter, len(b2_buffer)),
    }

    all_bspliters := inst.GetBspliters()
    all_uspliters := inst.GetUspliters()

    // a "y line" would look like:
    // y#1#2
    // where the edge between 1 and 2 should be drawn as is used in the solution
    loadInfostructure(&solution, y_buffer, csvg)
    // a "x2 line" is similar to the "y lin" but represents an edge in a virtual network
    loadInfostructure(&virtual_net, x2_buffer, csvg)
    
    // a "b line" looks like
    // b#1#2
    // where the balanced spliter 1 is applied to the node with ID 2
    for _, line := range b_buffer {
        items := strings.FieldsFunc(string(line), func(c rune)bool{return c == '#'})
        values := strings.FieldsFunc(items[len(items) - 1], func(c rune)bool{return c == ' '})
        items[len(items) - 1] = values[0]
        value_str := values[1]
        value, _ := strconv.ParseFloat(value_str, 32)
        if value <= 0 {
            continue
        }

        nodeID, _ := strconv.ParseUint(items[1], 10, 32) // ignores all error's trusting the file
        spliterID, _ := strconv.ParseUint(items[2], 10, 32)
    
        for _, spl := range all_bspliters {
            if spl.Id == uint32(spliterID) {
                solution.BspliterMap[uint32(nodeID)] = spl
                break
            }
        }
    }
    // does the same thing but with the b2 buffer for the virtual network
    for _, line := range b2_buffer {
        items := strings.FieldsFunc(string(line), func(c rune)bool{return c == '#'})
        values := strings.FieldsFunc(items[len(items) - 1], func(c rune)bool{return c == ' '})
        items[len(items) - 1] = values[0]
        value_str := values[1]
        value, _ := strconv.ParseFloat(value_str, 32)
        if value <= 0 {
            continue
        }

        nodeID, _ := strconv.ParseUint(items[1], 10, 32) // ignores all error's trusting the file
        spliterID, _ := strconv.ParseUint(items[2], 10, 32)
    
        for _, spl := range all_bspliters {
            if spl.Id == uint32(spliterID) {
                virtual_net.BspliterMap[uint32(nodeID)] = spl
                break
            }
        }
    }

    // a "a line" looks like
    // a#1#2#3
    // where the unbalanced spliter 1 is applied to the node with ID 3
    // the 2 is not important to the scope of this function
    for _, line := range a_buffer {
        items := strings.FieldsFunc(string(line), func(c rune)bool{return c == '#'})
        values := strings.FieldsFunc(items[len(items) - 1], func(c rune)bool{return c == ' '})
        items[len(items) - 1] = values[0]
        value_str := values[1]
        value, _ := strconv.ParseFloat(value_str, 32)
        if value <= 0 {
            continue
        }

        nodeID, _ := strconv.ParseUint(items[1], 10, 32) // ignores all error's trusting the file
        spliterID, _ := strconv.ParseUint(items[3], 10, 32)
    
        for _, spl := range all_uspliters {
            if spl.Id == uint32(spliterID) {
                solution.UspliterMap[uint32(nodeID)] = spl
                break
            }
        }
    }
    // a "u line" looks like
    // u#1
    // where the splicebox is applied to the node with ID 1
    // there is only ever 1 splicebox per instance
    for _, line := range u_buffer {
        items := strings.FieldsFunc(string(line), func(c rune)bool{return c == '#'})
        values := strings.FieldsFunc(string(items[len(items) - 1]), func(c rune)bool{return c == ' '})
        items[len(items) - 1] = values[0]
        value_str := values[1]
        value, _ := strconv.ParseFloat(value_str, 32)
        if value <= 0 {
            continue
        }

        nodeID, _ := strconv.ParseUint(items[1], 10, 32) // ignores all error's trusting the file
        solution.SpliceboxNodesId = append(solution.SpliceboxNodesId, uint32(nodeID))
    }

    return &solution, &virtual_net
}
