package logger

import (
	"go.uber.org/zap"
)

var _ Meta = (*meta)(nil)

type Meta interface {
	Key() string
	Value() any
	meta()
}

type meta struct {
	key   string
	value any
}

func NewMeta(key string, value any) Meta {
	return &meta{
		key:   key,
		value: value,
	}
}

func (m *meta) Key() string {
	return m.key
}

func (m *meta) Value() any {
	return m.value
}

func (m *meta) meta() {}

func WrapMeta(err error, metas ...Meta) (fields []zap.Field) {
	capacity := len(metas) + 1
	if err != nil {
		capacity++
	}

	fields = make([]zap.Field, 0, capacity)
	if err != nil {
		fields = append(fields, zap.Error(err))
	}

	fields = append(fields, zap.Namespace("meta"))
	for _, meta := range metas {
		fields = append(fields, zap.Any(meta.Key(), meta.Value()))
	}

	return
}
