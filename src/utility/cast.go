package utility

import (
	"encoding/binary"
	"github.com/google/uuid"
)

func BytesToInt(data []byte) int {
	return int(binary.BigEndian.Uint32(data))
}

func BytesToString(data []byte) string {
	return string(data)
}

func BytesToUid(data []byte) string {
	uid := uuid.NullUUID{}
	_ = uid.UnmarshalBinary(data)
	val, _ := uid.MarshalText()
	return string(val)
}