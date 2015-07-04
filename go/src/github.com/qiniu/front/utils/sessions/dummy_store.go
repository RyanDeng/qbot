package sessions

import (
	"hd.qiniu.com/utils/object"
)

type dummyStore struct {
	sid    string
	values map[string]interface{}

	changed bool // flag when values has changed
	closed  bool // flag when data destroy
}

var _ SessionStore = new(dummyStore)

// get session id
func (m *dummyStore) Sid() string {
	return m.sid
}

func (m *dummyStore) Set(key string, value interface{}) {
	m.changed = true
	m.values[key] = value
}

func (m *dummyStore) Get(key string) *object.Value {
	if v, ok := m.values[key]; ok {
		return object.ValueTo(v)
	}
	return object.ValueTo(nil)
}

func (m *dummyStore) Delete(key string) {
	m.changed = true
	delete(m.values, key)
}

func (m *dummyStore) Has(key string) bool {
	_, ok := m.values[key]
	return ok
}

func (m *dummyStore) Values() map[string]interface{} {
	values := make(map[string]interface{}, len(m.values))
	for key, value := range m.values {
		values[key] = value
	}
	return values
}

func (m *dummyStore) Clean() {
	m.changed = true
	m.values = make(map[string]interface{})
}

func (m *dummyStore) Flush() error {
	// has destroy
	if m.closed {
		return nil
	}

	// no changes
	if !m.changed {
		return nil
	}

	m.changed = false
	return nil
}

func (m *dummyStore) Destroy() error {
	m.closed = true
	m.values = make(map[string]interface{})
	return nil
}

func NewDummySessionStore(sid string, values map[string]interface{}) SessionStore {
	sess := &dummyStore{
		sid:    sid,
		values: values,
	}

	if sess.values == nil {
		sess.values = make(map[string]interface{})
		sess.changed = true
	}
	return sess
}
