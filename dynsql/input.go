package dynsql

import (
	"encoding/json"
	"net/url"
)

type Input struct {
	Filters []*filter `json:"filters" query:"filters" validate:"dive"`
	Group   string    `json:"group"   query:"group"`
	Order   string    `json:"order"   query:"order"`
	Desc    bool      `json:"desc"    query:"desc"`
}

func (in Input) empty() bool {
	return len(in.Filters) == 0 && in.Group == "" && in.Order == ""
}

type filter struct {
	Col string `json:"key"      validate:"required,lte=50"`
	Op  string `json:"operator" validate:"omitempty,oneof=eq ne gt lt gte lte in notin like notlike"`
	Val string `json:"value"    validate:"lte=100"`
}

func (f *filter) UnmarshalBind(raw string) error {
	data, err := url.QueryUnescape(raw)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), f)
}
