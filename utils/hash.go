package utils

import (
	"hash/fnv"
	"strconv"
)

func Hash(s string) string {
	return HashBytes([]byte(s))
}

func HashBytes(b []byte) string {
	h := fnv.New32a()
	h.Write(b)
	hashNum := h.Sum32()
	return strconv.FormatUint(uint64(hashNum), 16)
}
