package utils

import (
	"fmt"
	"sync"
)

type Name struct {
	First string `json:"first"` // this notation is to map the json key to the struct field
	Last  string `json:"last"`
}

type Record struct {
	Id string `json:"id"`
	Name Name   `json:"name"`
	Score int 	`json:"score"`
}

type RecordStore struct {
	mut sync.Mutex
	records []Record
}

func NewRecordStore() *RecordStore {
	return &RecordStore{
	}
}

func (rs *RecordStore) List() []Record {
	rs.mut.Lock()
	defer rs.mut.Unlock()
	return rs.records
}

func (rs *RecordStore) Get(n int) (Record, error) {
	rs.mut.Lock()
	defer rs.mut.Unlock()
	if n < 0 || n >= len(rs.records) {
		return Record{}, fmt.Errorf("index out of range")
	}
	return rs.records[n], nil
}

func (rs *RecordStore) Put(r Record) (int, error) {
	if !isValid(r) {
		return -1, fmt.Errorf("invalid record")
	}
	rs.mut.Lock()
	defer rs.mut.Unlock()
	rs.records = append(rs.records, r)
	return len(rs.records) - 1, nil
}