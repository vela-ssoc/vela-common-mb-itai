package dbms

import (
	"context"
	"net"
	"reflect"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	workerBits  uint8 = 10
	numberBits  uint8 = 12
	workerMax   int64 = -1 ^ (-1 << workerBits)
	numberMax   int64 = -1 ^ (-1 << numberBits)
	epoch       int64 = 1643644800000 // 2022 Chinese New Year
	timeShift         = workerBits + numberBits
	workerShift       = numberBits
)

func newSnow() *snow {
	mid := machineID()
	return &snow{machineID: mid}
}

type snow struct {
	mutex     sync.Mutex
	timestamp int64
	machineID int64
	sequence  int64
}

func (s *snow) Int64() int64 {
	return s.nextID()
}

func (s *snow) String() string {
	id := s.nextID()
	return strconv.FormatInt(id, 10)
}

// nextID 生成下一个ID
func (s *snow) nextID() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now().UnixNano() / 1e6
	if s.timestamp == now {
		s.sequence++
		if s.sequence > numberMax {
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.timestamp, s.sequence = now, 0
	}

	return (now-epoch)<<timeShift | (s.machineID << workerShift) | s.sequence
}

func (s *snow) autoID(db *gorm.DB) {
	stmt := db.Statement
	shm := stmt.Schema
	if stmt == nil || shm == nil {
		return
	}

	pk := shm.PrioritizedPrimaryField
	rv := stmt.ReflectValue
	ctx := stmt.Context
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			s.setID(ctx, pk, rv.Index(i))
		}
	case reflect.Struct:
		s.setID(ctx, pk, rv)
	}
}

func (s *snow) setID(ctx context.Context, field *schema.Field, rv reflect.Value) {
	if field == nil {
		return
	}
	val, zero := field.ValueOf(ctx, rv)
	if !zero {
		return
	}

	switch val.(type) {
	case int64:
		_ = field.Set(ctx, rv, s.Int64())
	case string:
		_ = field.Set(ctx, rv, s.String())
	}
}

// machineID 根据 mac 计算本机机器码
// FIXME: 这种算法还是有相当大重复概率的. 但是实现简单, 实际上线后节点应该不会太多
func machineID() int64 {
	var mac net.HardwareAddr
	faces, _ := net.Interfaces()
	for _, face := range faces {
		if len(face.HardwareAddr) > 0 {
			mac = face.HardwareAddr
			break
		}
	}

	var ret int64
	for _, b := range mac {
		ret += int64(b)
	}

	return ret % workerMax
}
