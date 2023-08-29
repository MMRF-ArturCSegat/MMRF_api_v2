package fiber_optic_components

import "errors"

// interface for all fibre optic components to implement as
// they have very similar querys
type FiberComponent interface {
    GetId() uint32
    String() string
}

func GetAll[T FiberComponent](objs []T) ([]T, error){
    res := db.Find(&objs)

    if res.Error != nil{
        return nil, errors.New("Failed to retrieve objects")
    }
    return objs, nil
}


func GetOne[T FiberComponent](id uint32, obj *T) (error){
    res := db.Find(obj, id)

	if res.RowsAffected == 0 || res.Error != nil {
		return errors.New("Not in database")
	} 
    return nil
}


func AddObj[T FiberComponent](obj *T) error {
    res := db.Create(&obj)

    if res.RowsAffected == 0 {
        return errors.New("could not add to database")
    }
    return nil
}
