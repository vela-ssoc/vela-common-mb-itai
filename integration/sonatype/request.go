package sonatype

import "github.com/vela-ssoc/vela-common-mb-itai/dal/model"

// sonatypeRequest sonatype HTTP 请求体
type sonatypeRequest struct {
	Coordinates []string `json:"coordinates"` // PURL
}

// sonatypeComponent 漏洞组件信息
type sonatypeComponent struct {
	Coordinates     string                   `json:"coordinates"`     // PURL
	Description     string                   `json:"description"`     // 漏洞简短描述
	Reference       string                   `json:"reference"`       // 组件详情链接（组件维度）
	Vulnerabilities []*sonatypeVulnerability `json:"vulnerabilities"` // 漏洞详细信息
}

// sonatypeVulnerability 漏洞详情
type sonatypeVulnerability struct {
	ID                 string          `json:"id"`                 // 漏洞编号
	DisplayName        string          `json:"displayName"`        // 通常和漏洞编号一样
	Title              string          `json:"title"`              // 漏洞标题
	Description        string          `json:"description"`        // 漏洞简述
	CVSSScore          model.CVSSScore `json:"cvssScore"`          // CVSS 评分
	CVSSVector         string          `json:"cvssVector"`         // CVSS Vector
	CWE                string          `json:"cwe"`                // CWE
	CVE                string          `json:"cve"`                // CVE
	Reference          string          `json:"reference"`          // 漏洞详情链接（漏洞维度）
	ExternalReferences []string        `json:"externalReferences"` // 其他的漏洞描述链接
}

// sonatypeComponents sonatype HTTP 响应体
type sonatypeComponents []*sonatypeComponent

// Vulns 将数据转为业务数据库格式的数据
func (scs sonatypeComponents) Vulns() []*model.SBOMVuln {
	ret := make([]*model.SBOMVuln, 0, len(scs))
	for _, sc := range scs {
		purl := sc.Coordinates
		for _, vn := range sc.Vulnerabilities {
			score := vn.CVSSScore
			vul := &model.SBOMVuln{
				VulnID:      vn.ID,
				PURL:        purl,
				Title:       vn.Title,
				Description: vn.Description,
				Score:       score,
				Level:       score.Level(),
				Vector:      vn.CVSSVector,
				CVE:         vn.CVE,
				CWE:         vn.CWE,
				Reference:   sc.Reference,
				References:  vn.ExternalReferences,
			}
			ret = append(ret, vul)
		}
	}

	return ret
}
