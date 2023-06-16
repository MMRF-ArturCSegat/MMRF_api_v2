package fiber_optic_components


type FiberCable struct {
    Id                  uint        `json:"id"`
    Attenuation         float64     `json:"attenuation"`
    Cost                float64     `json:"cost"`
}


func (fb FiberCable) GetId() uint {
    return fb.Id
}
