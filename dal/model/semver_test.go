package model_test

import (
	"testing"

	"github.com/vela-ssoc/vela-common-mb-itai/dal/model"
)

func TestSemver_Int64(t *testing.T) {
	cases := map[model.Semver]int64{
		"1.1.0":             10000100000,     // 1.00001.00000
		"0.1.0":             100000,          // 1.00000
		"0.0.1":             1,               // 0.0.00001
		"1.2.3":             10000200003,     // 1.00002.00003
		"1.2.46":            10000200046,     // 1.00002.00046
		"12.344.5421":       120034405421,    // 12.00344.05421
		"1.1.0.2":           10000100000,     // 1.00001.00000
		"1.2.3-beta":        10000200003,     // 1.00002.00003
		"1.2":               10000200000,     // 1.00002.00000
		"devel":             0,               // 0.0.0
		"99999.99999.99999": 999999999999999, // 99999.99999.99999 maximum
	}
	for k, v := range cases {
		i := k.Int64()
		match := i == v

		if match {
			t.Logf("[OK] %s => %d", k, v)
		} else {
			t.Errorf("[ERROR] %s 期望: %d，但是计算结果却是：%d", k, v, i)
		}
	}
}
