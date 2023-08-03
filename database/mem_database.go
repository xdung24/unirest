package database

import (
	"fmt"
	"sync"
)

type MemDatabase struct {
	mu         sync.Mutex
	namespaces map[string]namespace
}

type namespace struct {
	data map[string][]byte
}

func newNamespace() namespace {
	return namespace{
		data: make(map[string][]byte),
	}
}

func (m *MemDatabase) Init() {
	m.namespaces = make(map[string]namespace)
}

func (m *MemDatabase) Disconnect() {
	// Do nothing
}

func (m *MemDatabase) Upsert(namespace string, key string, value []byte, allowOverWrite bool) *DbError {
	m.mu.Lock()
	defer m.mu.Unlock()

	ns, ok := m.namespaces[namespace]
	if !ok {
		ns = newNamespace()
		m.namespaces[namespace] = ns
	}
	ns.data[key] = value
	return nil
}

func (m *MemDatabase) Get(namespace string, key string) ([]byte, *DbError) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ns, ok := m.namespaces[namespace]
	if !ok {
		return nil, &DbError{
			ErrorCode: NAMESPACE_NOT_FOUND,
			Message:   fmt.Sprintf("namespace '%v' does not exist.", namespace),
		}
	}
	val, ok := ns.data[key]
	if !ok {
		return nil, &DbError{
			ErrorCode: ID_NOT_FOUND,
			Message:   fmt.Sprintf("value not found in namespace '%v' for key '%v'", namespace, key),
		}
	}
	return val, nil
}

func (m *MemDatabase) GetAll(namespace string) (map[string][]byte, *DbError) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ns, ok := m.namespaces[namespace]
	if !ok {
		return nil, &DbError{
			ErrorCode: NAMESPACE_NOT_FOUND,
			Message:   fmt.Sprintf("namespace '%v' does not exist.", namespace),
		}
	}
	return ns.data, nil
}

func (m *MemDatabase) Delete(namespace string, key string) *DbError {
	m.mu.Lock()
	defer m.mu.Unlock()

	ns, ok := m.namespaces[namespace]
	if !ok {
		return &DbError{
			ErrorCode: NAMESPACE_NOT_FOUND,
			Message:   fmt.Sprintf("namespace '%v' does not exist.", namespace),
		}
	}
	_, ok = ns.data[key]
	if !ok {
		return &DbError{
			ErrorCode: ID_NOT_FOUND,
			Message:   fmt.Sprintf("value not found in namespace '%v' for key '%v'", namespace, key),
		}
	}

	delete(ns.data, key)
	return nil
}

func (m *MemDatabase) DeleteAll(namespace string) *DbError {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.namespaces[namespace]
	if !ok {
		return &DbError{
			ErrorCode: NAMESPACE_NOT_FOUND,
			Message:   fmt.Sprintf("namespace '%v' does not exist.", namespace),
		}
	}
	delete(m.namespaces, namespace)
	return nil
}

func (m *MemDatabase) GetNamespaces() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	ret := make([]string, 0)
	for k := range m.namespaces {
		ret = append(ret, k)
	}
	return ret
}
