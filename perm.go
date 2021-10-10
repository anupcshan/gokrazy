package gokrazy

import (
	"runtime"

	"github.com/gokrazy/internal/rootdev"
)

func permPartition() string {
	if runtime.GOARCH == "arm" {
		return "/dev/sda1"
	}
	return rootdev.Partition(rootdev.Perm)
}
