package fiber_optic_components

import "errors"

// interface for all fibre optic components to implement as
// they have very similar querys
type db_obj interface {
    GetId() uint32
    String() string
}

func GetAll[T db_obj](objs []T) ([]T, error){
    res := db.Find(&objs)

    if res.Error != nil{
        return nil, errors.New("Failed to retrieve objects")
    }
    return objs, nil
}


func GetOne(id uint32, obj db_obj) (db_obj, error){
    res := db.Find(&obj, id)

	if res.RowsAffected == 0 || res.Error != nil {
		return nil, errors.New("Not in database")
	} 
    return obj, nil
}


func AddObj[T db_obj](obj *T) error {
    res := db.Create(&obj)

    if res.RowsAffected == 0 {
        return errors.New("could not add to database")
    }
    return nil
}
