package util

import "errors"

func In[T comparable](e T, list []T) bool{
    for _, item := range list{
        if item == e{
            return true
        }
    }
    return false
}

func SliceComp[T comparable](s1, s2[]T) bool{
    for i, e := range s1{
        if s2[i] != e{
            return false
        }
    }
    return true
}

func SliceInMatrix[T comparable](m [][]T, s[]T) bool{
    for _, e := range m{
        if SliceComp(e, s) == true{
            return true
        }
    }
    return false
}


func IndexOf[T comparable](thing T, array []T) (int, error){
    for index, value := range array{
        if value == thing{
            return index, nil
        }
    }
    return -1, errors.New("item not in array")
}



