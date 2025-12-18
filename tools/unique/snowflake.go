package unique

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const (
	epoch         = 1609459200000 // 起始时间:2021-01-01 00:00:00 UTC
	timestampBits = 41            // 时间戳位数
	machineIDBits = 10            // 机器ID位数
	sequenceBits  = 12            // 序列号位数

	maxMachineID = -1 ^ (-1 << machineIDBits) // 最大机器ID
	maxSequence  = -1 ^ (-1 << sequenceBits)  // 最大序列号

	timestampShift = machineIDBits + sequenceBits // 时间戳位移
	machineIDShift = sequenceBits                 // 机器ID位移
)

var defaultSnowflake *Snowflake

// Snowflake 结构体
type Snowflake struct {
	mutex     sync.Mutex
	lastStamp int64
	machineID int64
	sequence  int64
}

func init() {
	podIPStr := os.Getenv("POD_IP")
	podIP := net.ParseIP(podIPStr)
	if podIP == nil {
		podIP = []byte{127, 0, 0, 1}
	}
	machineID, err := IPToInt(podIP)
	if err != nil {
		machineID = 1
	}
	machineID = machineID % maxMachineID
	defaultSnowflake, _ = NewSnowflake(int64(machineID))
}

func NewSnowflakeID() uint64 {
	return uint64(defaultSnowflake.NextID())
}

// NewSnowflake 创建一个 Snowflake 实例
func NewSnowflake(machineID int64) (*Snowflake, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, errors.New("machine ID out of range")
	}
	return &Snowflake{
		lastStamp: 0,
		machineID: machineID,
		sequence:  0,
	}, nil
}

// NextID 生成下一个唯一ID
func (s *Snowflake) NextID() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 获取当前时间戳
	currentStamp := time.Now().UnixMilli()

	// 检查时钟回拨
	if currentStamp < s.lastStamp {
		panic("clock moved backwards")
	}

	// 处理同一毫秒内的序列号
	if currentStamp == s.lastStamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// 序列号溢出，等待下一毫秒
			for currentStamp <= s.lastStamp {
				currentStamp = time.Now().UnixMilli()
			}
		}
	} else {
		// 新的一毫秒，重置序列号
		s.sequence = 0
	}

	// 更新时间戳
	s.lastStamp = currentStamp

	// 生成ID
	id := ((currentStamp - epoch) << timestampShift) |
		(s.machineID << machineIDShift) |
		s.sequence
	return id
}

// IPToInt 将IPv4地址转换为32位整数
func IPToInt(ip net.IP) (uint32, error) {
	if ip.To4() == nil {
		return 0, fmt.Errorf("not an IPv4 address")
	}
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3]), nil
}
