package utils

func ArrayMap[T any, U any](arr []T, fn func(T) U) []U {
    ret := []U{}
    
    for _, item := range arr {
        ret = append(ret, fn(item))
    }
    
    return ret
}

func ArrayFilter[T any](arr []T, fn func(T) bool) []T {
    ret := []T{}
    
    for _, item := range arr {
        if fn(item) {
            ret = append(ret, item)
        }
    }
    
    return ret
}

func ArrayFind[T any](arr []T, fn func(T) bool) *T {
    for _, item := range arr {
        if fn(item) {
            return &item
        }
    }
    
    return nil
}

func ArrayHas[T any](arr []T, fn func(T) bool) bool {
    for _, item := range arr {
        if fn(item) {
            return true
        }
    }
    
    return false
}