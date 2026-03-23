package main

type KVStore struct {
	data map[string]string
}

func NewKVStore() *KVStore {
	return &KVStore{
		data: make(map[string]string),
	}
}

func (kv *KVStore) Set(key string, value string) {
	kv.data[key] = value
}

func (kv *KVStore) Get(key string) (string, bool) {
	value, exists := kv.data[key]
	return value, exists
}

func (kv *KVStore) Delete(key string) {
	delete(kv.data, key)
}