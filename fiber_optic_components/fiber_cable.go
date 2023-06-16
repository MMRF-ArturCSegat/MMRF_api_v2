package fiber_optic_components

import "fmt"


type FiberCable struct {
    Id                  uint        `json:"id"`
    Attenuation         float64     `json:"attenuation"`
    Cost                float64     `json:"cost"`
}


func (fb FiberCable) GetId() uint {
    return fb.Id
}


func (fb FiberCable) String() string{
    str := fmt.Sprintf("ID\n%v\n", fb.Id)
    str += fmt.Sprintf("Attenuation\n%v\n", fb.Attenuation)
    str += fmt.Sprintf("Cost\n%v\n", fb.Cost)
    return str
}
