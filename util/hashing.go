package util

import "hash/fnv"

func CreateHash(text []byte) uint32 {
	hash := fnv.New32a()
	hash.Write(text)
	return hash.Sum32()
}
