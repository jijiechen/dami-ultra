package helpers

import (
	"hash/fnv"
	"os"
)

func GetHostID() int64 {
	hostname, err := os.Hostname()
	if err != nil {
		return 0
	}

	h := fnv.New64a()
	h.Write([]byte(hostname))
	return int64(h.Sum64() % 1024)
}
