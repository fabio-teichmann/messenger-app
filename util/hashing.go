package util

import (
	"hash/fnv"
	"time"
)

func CreateHash(text []byte) uint32 {
	hash := fnv.New32a()
	t, _ := time.Now().MarshalBinary()
	hash.Write(append(text, t...))
	return hash.Sum32()
}
