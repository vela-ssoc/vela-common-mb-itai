package model

import "database/sql/driver"

const (
	RLvlLow RiskLevel = iota
	RLvlCritical
	RLvlHigh
	RLvlMiddle
)

var riskLvlSI = map[string]RiskLevel{
	"低危": RLvlLow,
	"中危": RLvlMiddle,
	"高危": RLvlHigh,
	"紧急": RLvlCritical,
}

var riskLvlIS = map[RiskLevel]string{
	RLvlLow:      "低危",
	RLvlMiddle:   "中危",
	RLvlHigh:     "高危",
	RLvlCritical: "紧急",
}

// RiskLevel 风险级别，支持直接比较
type RiskLevel int

func (rl RiskLevel) Value() (driver.Value, error) {
	str := riskLvlIS[rl]
	return str, nil
}

func (rl *RiskLevel) Scan(src any) error {
	if str, ok := src.(string); ok {
		lv := riskLvlSI[str]
		*rl = lv
	}
	return nil
}

func (rl *RiskLevel) UnmarshalText(raw []byte) error {
	n := riskLvlSI[string(raw)]
	*rl = n
	return nil
}

func (rl RiskLevel) MarshalText() ([]byte, error) {
	str := riskLvlIS[rl]
	return []byte(str), nil
}

func (rl RiskLevel) String() string {
	return riskLvlIS[rl]
}
