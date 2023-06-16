package fiber_optic_components

import "fmt"

type FiberUnbalancedSpliter struct{
    Id                      uint        `json:"id"`
    Cost                    float64     `json:"cost"`
    Loss_ratio1             float64     `json:"loss1"`
    Loss_ratio2             float64     `json:"loss2"`
}


func (fus FiberUnbalancedSpliter) GetId() uint {
    return fus.Id
}


func (fb FiberUnbalancedSpliter) String() string{
    str := fmt.Sprintf("ID\n%v\n", fb.Id)
    str += fmt.Sprintf("Cost\n%v\n", fb.Cost)
    str += fmt.Sprintf("LossRatio1\n%v\n", fb.Loss_ratio1)
    str += fmt.Sprintf("LossRatio2\n%v\n", fb.Loss_ratio2)
    return str
}

