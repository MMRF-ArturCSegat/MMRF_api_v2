package fiber_optic_components


type FiberBalancedSpliter struct{
    Id                      uint        `json:"id"`
    Cost                    float64     `json:"cost"`
    Loss_ratio              float64     `json:"loss1"`
    Split_ratio             float64     `json:"loss2"`
}


func (fbs FiberBalancedSpliter) GetId() uint {
    return fbs.Id
}
