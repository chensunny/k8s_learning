package uuid

import (
	"time"
)

const (
	letterBytes   = "0123456789abcdefghijklmnopqrstuvwxyz"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits
)

// nolint
func randBytesByMathLib(n int) []byte {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return b
}

func randBytesByCryptoLib(n int) []byte {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, readRandomUint64(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = readRandomUint64(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return b
}

func tillNextMillis(lastTimestamp int64) int64 {
	now := timeGen()
	for now <= lastTimestamp {
		now = timeGen()
	}
	return now
}

func timeGen() int64 {
	return time.Now().UnixNano() / 1e6
}

// 从 g～w 一共17个,xyz 做预留
func firstByte() byte {
	first := byte(randInt63n(17))
	return first + byte('g')
}

func processUniqueBytes() [7]byte {
	var b [7]byte
	processUnique := randBytesByCryptoLib(7)
	copy(b[:], processUnique[:])
	return b
}

func copy36base(target []byte, base []byte) {
	tLen := len(target)
	bLen := len(base)
	if bLen > tLen {
		return
	}
	diff := tLen - bLen
	for i, b := range base {
		target[i+diff] = b
	}
	for i := 0; i < diff; i++ {
		target[i] = byte('0')
	}
}
