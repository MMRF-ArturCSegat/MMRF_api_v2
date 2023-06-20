package fiber_optic_components

import "fmt"

type FiberSpliceBox struct {
    Id                      uint32        `json:"id"`
    Attenuation             float64     `json:"attenuation"`
    Cost                    float64     `json:"cost"`
}


func (fsb FiberSpliceBox) GetId() uint32 {
    return fsb.Id
}


func (fb FiberSpliceBox) String() string{
    str := fmt.Sprintf("ID\n%v\n", fb.Id)
    str += fmt.Sprintf("Attenuation\n%v\n", fb.Attenuation)
    str += fmt.Sprintf("Cost\n%v\n", fb.Cost)
    return str
}

