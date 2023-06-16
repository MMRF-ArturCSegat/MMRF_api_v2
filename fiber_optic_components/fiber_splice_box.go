package fiber_optic_components


type FiberSpliceBox struct {
    Id                      uint        `json:"id"`
    Attenuation             float64     `json:"attenuation"`
    Cost                    float64     `json:"cost"`
}


func (fsb FiberSpliceBox) GetId() uint {
    return fsb.Id
}
