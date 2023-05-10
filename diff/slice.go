package diff

type Slice[T any] struct {
    OnCompare func(a, b T) bool
    OnDelete  func(T)
    OnInsert  func(T)
    OnUpdate  func(prev, curr T)
}

func (o Slice[T]) Compare(prev []T, curr []T) {
    for _, s := range prev {
        if !o.find(curr, s, true) {
            if o.OnDelete != nil {
                o.OnDelete(s)
            }
        }
    }
    for _, m := range curr {
        if !o.find(prev, m, false) {
            if o.OnInsert != nil {
                o.OnInsert(m)
            }
        }
    }
}

func (o Slice[T]) find(s []T, what T, byCurr bool) bool {
    for _, v := range s {
        if o.OnCompare(v, what) {
            if byCurr && o.OnUpdate != nil {
                o.OnUpdate(what, v)
            }
            return true
        }
    }
    return false
}
