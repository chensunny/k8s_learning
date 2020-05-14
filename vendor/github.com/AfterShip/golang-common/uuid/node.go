package uuid

import (
	"bytes"
	"errors"
	"strconv"
	"sync"
	"time"
)

var maxStep, _ = strconv.ParseInt("zzz", 36, 64)
var maxRandom, _ = strconv.ParseInt("zzzzz", 36, 64)

var ErrInvalidTime = errors.New("timestep err when create c24 ")

//**具体ID格式**
//>[g-z random] + random(12) + date(8) + increase(3 random)
//> random(12) = random(5) +  processUnique(7)

type C24 [24]byte

// NilC24 is the zero value for NilC24.
var NilC24 C24

var processUnique = processUniqueBytes()

type node struct {
	mu            sync.Mutex
	lastTimestamp int64
	step          int64
	random        int64
}

var globalNode = &node{}

func (n *node) Generate() (C24, error) {
	//方便测试做抽离注入
	return n.GenerateFromTimestamp(timeGen())
}

func (n *node) GenerateFromTimestamp(timestamp int64) (C24, error) {
	n.mu.Lock()

	//防止时间回拨，同时对预支是的时间做一个补偿
	if timestamp < n.lastTimestamp {
		offset := n.lastTimestamp - timestamp
		if offset <= 5 {
			time.Sleep(time.Duration(offset<<1) * time.Millisecond)
			timestamp = timeGen()
			if timestamp < n.lastTimestamp {
				return NilC24, ErrInvalidTime
			}
		} else {
			return NilC24, ErrInvalidTime
		}
	}
	//如果当前毫秒的 step 分配完了，尝试预支下面几毫秒的
	if timestamp == n.lastTimestamp {
		n.step++
		//seq 为 maxStep 的时候表示是下一毫秒时间开始对seq做随机
		if n.step >= maxStep {
			n.step = randInt63n(100)
			timestamp = tillNextMillis(n.lastTimestamp)
		}
	} else {
		n.step = randInt63n(100)
	}

	n.lastTimestamp = timestamp
	n.random = randInt63n(maxRandom)

	var b [24]byte

	b[0] = firstByte()

	// compare performance RandBytesByMathLib(5)
	copy36base(b[1:6], []byte(strconv.FormatInt(n.random, 36)))
	//随机数
	//copy(b[1:6], RandBytesByMathLib(5))

	//进程标示
	copy(b[6:13], processUnique[:])

	//时间戳
	copy(b[13:21], []byte(strconv.FormatInt(timestamp, 36)))

	//step
	copy36base(b[21:24], []byte(strconv.FormatInt(n.step, 36)))

	n.mu.Unlock()

	return b, nil
}

func (b C24) Random() int64 {
	p, _ := strconv.ParseInt(string(b[1:6]), 36, 64)
	return p
}

func (b C24) ProcessMark() int64 {
	r, _ := strconv.ParseInt(string(b[6:13]), 36, 64)
	return r
}

func (b C24) Time() int64 {
	time, _ := strconv.ParseInt(string(b[13:21]), 36, 64)
	return time
}

func (b C24) Step() int64 {
	step, _ := strconv.ParseInt(string(b[21:24]), 36, 64)
	return step
}

func (b C24) String() string {
	return string(b[:])
}

func (b C24) IsZero() bool {
	return bytes.Equal(b[:], NilC24[:])
}
