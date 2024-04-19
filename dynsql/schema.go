package dynsql

type Schema struct {
	Filters columnSchemas `json:"conditions"`
	Groups  nameSchemas   `json:"groups"`
	Orders  nameSchemas   `json:"orders"`
}

type nameSchema struct {
	Col  string `json:"col"`
	Name string `json:"name"`
}

type nameSchemas []*nameSchema

type columnSchema struct {
	Col       string            `json:"key"`       // `json:"col"`
	Name      string            `json:"desc"`      // `json:"name"`
	Type      string            `json:"type"`      // `json:"type"`
	Operators []*operatorSchema `json:"operators"` // `json:"operators"`
	Enum      bool              `json:"enum"`      // `json:"enums"`
	Enums     []*enumItem       `json:"enums"`     // `json:"enums"`
}

type columnSchemas []*columnSchema
