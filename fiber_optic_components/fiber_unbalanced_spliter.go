package fiber_optic_components


type FiberUnbalancedSpliter struct{
    Id                      uint        `json:"id"`
    Cost                    float64     `json:"cost"`
    Loss_ratio1             float64     `json:"loss1"`
    Loss_ratio2             float64     `json:"loss2"`
}


func (fus FiberUnbalancedSpliter) GetId() uint {
    return fus.Id
}
