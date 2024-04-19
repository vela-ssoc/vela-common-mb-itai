package alarm

import "github.com/vela-ssoc/vela-common-mb-itai/dal/model"

type Center interface {
	Risk(rsk *model.Risk) error
	Event(evt *model.Event) error
}

type name struct{}
