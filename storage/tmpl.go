package storage

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"text/template"
)

type storeRender interface {
	Valuer
	Rend(ctx context.Context, v any) *bytes.Buffer
}

type storeTemplate struct {
	filter []func([]byte) []byte
	value  Valuer
	mutex  sync.RWMutex
	done   bool
	tmpl   *template.Template
	err    error
}

func (st *storeTemplate) Value(ctx context.Context) ([]byte, error) {
	return st.value.Value(ctx)
}

func (st *storeTemplate) ID() string {
	return st.value.ID()
}

func (st *storeTemplate) Shared() bool {
	return st.value.Shared()
}

func (st *storeTemplate) Reset() {
	st.mutex.Lock()
	st.done = false
	st.value.Reset()
	st.mutex.Unlock()
}

func (st *storeTemplate) Invalid(dat []byte) bool {
	return st.value.Invalid(dat)
}

func (st *storeTemplate) Rend(ctx context.Context, v any) *bytes.Buffer {
	id := st.value.ID()
	buf := new(bytes.Buffer)
	tmpl, err := st.load(ctx)
	if err != nil {
		msg := fmt.Sprintf("加载 %s 模板出错：%s", id, err)
		buf.WriteString(msg)
		return buf
	}

	if err = tmpl.Execute(buf, v); err != nil {
		buf.Reset()
		msg := fmt.Sprintf("渲染 %s 模板出错：%s", id, err)
		buf.WriteString(msg)
	}

	return buf
}

func (st *storeTemplate) load(ctx context.Context) (*template.Template, error) {
	st.mutex.RLock()
	done, tmpl, err := st.done, st.tmpl, st.err
	st.mutex.RUnlock()
	if done {
		return tmpl, err
	}
	return st.slowLoad(ctx)
}

func (st *storeTemplate) slowLoad(ctx context.Context) (*template.Template, error) {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	if st.done {
		return st.tmpl, st.err
	}

	dat, err := st.value.Value(ctx)
	if err != nil {
		st.err = err
		st.done = true
		return nil, err
	}

	for _, fn := range st.filter {
		dat = fn(dat)
	}

	id := st.value.ID()
	tmpl, exx := template.New(id).Parse(string(dat))
	st.tmpl, st.err = tmpl, exx
	st.done = true

	return tmpl, exx
}
