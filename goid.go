package goid

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"

	"github.com/en-v/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const DEFAULT_SIZE = 8

type GoId struct {
	size int
	data []byte
}

func New() GoId {
	return NewCustom(DEFAULT_SIZE)
}

func NewCustom(size int) GoId {
	new := GoId{
		data: make([]byte, size),
	}
	_, err := rand.Read(new.data)
	if err != nil {
		panic(err)
	}
	return new
}

func Empty() GoId {
	return EmptyCustom(DEFAULT_SIZE)
}

func EmptyCustom(size int) GoId {
	return GoId{
		data: make([]byte, size),
	}
}

func (self *GoId) IsEmpty() bool {
	for i := range self.data {
		if self.data[i] != 0 {
			return false
		}
	}
	return true
}

func (self *GoId) String() string {
	return hex.EncodeToString(self.data)
}

func Parse(str string) GoId {
	arr, err := hex.DecodeString(str)
	if err != nil {
		println(err)
		return Empty()
	}
	new := EmptyCustom(len(arr))
	new.data = arr
	return new
}

func (self *GoId) Len() int {
	return len(self.data)
}

func (self *GoId) UInt64() uint64 {
	return binary.LittleEndian.Uint64(self.data)
}

func JustString() string {
	b := NewCustom(DEFAULT_SIZE)
	return b.String()
}

func JustCustomString(size int) string {
	b := NewCustom(size)
	return b.String()
}

// Custom JSON marshallind and unmarshalling

func (self *GoId) MarshalJSON() ([]byte, error) {
	return []byte("\"" + self.String() + "\""), nil
}

func (self *GoId) UnmarshalJSON(data []byte) (err error) {

	self.data, err = hex.DecodeString(string(data[1 : len(data)-1]))
	if err != nil {
		self.data = make([]byte, DEFAULT_SIZE)
		self.size = DEFAULT_SIZE
		return err
	}
	self.size = len(self.data)
	return nil
}

// Custom BSON marshallind and unmarshalling

func (self *GoId) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(self.String())
}

func (self *GoId) UnmarshalBSONValue(t bsontype.Type, data []byte) (err error) {
	self.data, err = hex.DecodeString(string(data[4 : len(data)-1]))
	if err != nil {
		log.Error(err)
		self.data = make([]byte, DEFAULT_SIZE)
		self.size = DEFAULT_SIZE
		return err
	}

	self.size = len(self.data)
	return nil
}
