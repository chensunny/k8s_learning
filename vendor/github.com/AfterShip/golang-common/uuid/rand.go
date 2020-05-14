package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	mrand "math/rand"
)

// 原来可以使用纳秒，但是担心不同语言随机算法不同，
// 容易受到 Seed 的影响，这里统一使用readRandomUint64
var src = mrand.NewSource(int64(readRandomUint64()))

// randInt63n returns an int64, a non-negative pseudo-random number in [0,n).
// It panics if n <= 0.
// 这里没使用 "math/rand" 下默认的的 rand 是减少锁的影响
// 原来包下默认对象有个全局的锁
func randInt63n(n int64) int64 {
	if n <= 0 {
		panic("invalid argument to Int63n")
	}
	if n&(n-1) == 0 { // n is power of two, can mask
		return src.Int63() & (n - 1)
	}
	max := int64((1 << 63) - 1 - (1<<63)%uint64(n))
	v := src.Int63()
	for v > max {
		v = src.Int63()
	}
	return v % n
}

// nolint
func readRandomUint32() uint32 {
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("cannot initialize objectid package with crypto.rand.Reader: %v", err))
	}
	return binary.BigEndian.Uint32(b[:])
}

func readRandomUint64() uint64 {
	var b [8]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("cannot initialize objectid package with crypto.rand.Reader: %v", err))
	}

	return binary.BigEndian.Uint64(b[:])
}
