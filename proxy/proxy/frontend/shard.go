package frontend

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	"git.2dfire.net/zerodb/proxy/pkg/errcode"
)

/*由分片ID找到分片，可用文件中的函数*/
type KeyError string

func NewKeyError(format string, args ...interface{}) KeyError {
	return KeyError(fmt.Sprintf(format, args...))
}

func (ke KeyError) Error() string {
	return string(ke)
}

// Uint64Key is a uint64 that can be converted into a KeyspaceId.
type Uint64Key uint64

func (i Uint64Key) String() string {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint64(i))
	return buf.String()
}

func HashValue(value interface{}, rule string) (int64, error) {
	switch rule {
	case "int":
		intVal, ok := value.(int64)
		if !ok {
			//panic(NewKeyError("Unexpected key variable type %T", value))
			strVal, ok := value.(string)
			if !ok {
				return -1, errcode.BuildError(errcode.ShardingKeyTypeInvalid, value)
			}
			i, err := strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				return -1, errcode.BuildError(errcode.ShardingKeyTypeInvalid, value)
			}
			intVal = i
		}

		return intVal, nil
	case "string":
		val, ok := value.(string)
		if !ok {
			return -1, errcode.BuildError(errcode.ShardingKeyTypeInvalid, value)
		}
		return stringHash(val), nil
	default:
		return -1, errcode.BuildError(errcode.ShardingKeyTypeInvalid, value)
	}
}

func stringHash(s string) int64 {
	start := 0
	var end int
	if len(s) > 32 {
		end = 32
	} else {
		end = len(s)
	}
	var hash int64 = 0
	for i := start; i < end; i++ {
		hash = (hash << 5) - hash + int64(s[i])
	}
	return hash
}

// for sharding extension
type Shard interface {
	FindForKey(key interface{}) (int, error)
}

type IntHashShard struct {
	SlotNum     int
	ShardLength int
}

func (s *IntHashShard) FindForKey(key interface{}) (int, error) {
	h, err := HashValue(key, "int")
	if err != nil {
		return int(h), err
	}

	var partitionLength int64
	partitionLength = int64(s.SlotNum)
	var andVal int64
	andVal = partitionLength - 1
	return int(h & andVal / int64(s.ShardLength)), nil
}

type StringHashShard struct {
	SlotNum     int
	ShardLength int
}

func (s *StringHashShard) FindForKey(key interface{}) (int, error) {
	h, err := HashValue(key, "string")
	if err != nil {
		return int(h), err
	}

	var partitionLength int64
	partitionLength = int64(s.SlotNum)
	var andVal int64
	andVal = partitionLength - 1
	return int(h & andVal / int64(s.ShardLength)), nil
}

type DefaultShard struct {
}

func (s *DefaultShard) FindForKey(key interface{}) (int, error) {
	return 0, nil
}
