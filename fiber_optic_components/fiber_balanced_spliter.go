package fiber_optic_components

import "fmt"


type FiberBalancedSpliter struct{
    Id                      uint32        `json:"id"`
    Cost                    float64     `json:"cost"`
    Loss_ratio              float64     `json:"loss1"`
    Split_ratio             float64     `json:"loss2"`
}


func (fbs FiberBalancedSpliter) GetId() uint32 {
    return fbs.Id
}


func (fb FiberBalancedSpliter) String() string{
    str := fmt.Sprintf("ID\n%v\n", fb.Id)
    str += fmt.Sprintf("Cost\n%v\n", fb.Cost)
    str += fmt.Sprintf("LossRatio\n%v\n", fb.Loss_ratio)
    str += fmt.Sprintf("SplitRation\n%v\n", fb.Split_ratio)
    return str
}

