package dynsql

type enumItem struct {
	Val  string `json:"key"`  // `json:"val"`
	Name string `json:"desc"` // `json:"name"`
}

type Enums struct {
	items []*enumItem
}

func (es Enums) toMap() map[string]struct{} {
	size := len(es.items)
	if size == 0 {
		return nil
	}

	ret := make(map[string]struct{}, size)
	for _, e := range es.items {
		ret[e.Val] = struct{}{}
	}

	return ret
}
