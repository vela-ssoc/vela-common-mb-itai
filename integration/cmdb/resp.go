package cmdb

import "github.com/vela-ssoc/vela-common-mb-itai/dal/model"

type reply struct {
	Result model.Cmdbs `json:"result"`
}

type Host struct {
	ID   int64
	Inet string
}

type Hosts []*Host

func (hs Hosts) inets() ([]string, map[string]int64) {
	sz := len(hs)
	ips := make([]string, 0, sz)
	ipm := make(map[string]int64, sz)

	for _, h := range hs {
		id, inet := h.ID, h.Inet
		if id <= 0 || inet == "" {
			continue
		}

		if _, ok := ipm[inet]; !ok {
			ipm[inet] = h.ID
			ips = append(ips, inet)
		}
	}

	return ips, ipm
}
