package goid

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"

	"github.com/en-v/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const DEF_SIZE = 8

type BBID struct {
	size int
	data []byte
}

func New() BBID {
	return NewCustom(DEF_SIZE)
}

func NewCustom(size int) BBID {
	new := BBID{
		data: make([]byte, size),
	}
	_, err := rand.Read(new.data)
	if err != nil {
		panic(err)
	}
	return new
}

func Empty() BBID {
	return EmptyCustom(DEF_SIZE)
}

func EmptyCustom(size int) BBID {
	return BBID{
		data: make([]byte, size),
	}
}

func (self *BBID) IsEmpty() bool {
	for i := range self.data {
		if self.data[i] != 0 {
			return false
		}
	}
	return true
}

func (self *BBID) String() string {
	return hex.EncodeToString(self.data)
}

func Parse(str string) BBID {
	arr, err := hex.DecodeString(str)
	if err != nil {
		println(err)
		return Empty()
	}
	new := EmptyCustom(len(arr))
	new.data = arr
	return new
}

func (self *BBID) Len() int {
	return len(self.data)
}

func (self *BBID) UInt64() uint64 {
	return binary.LittleEndian.Uint64(self.data)
}

func JustString() string {
	b := NewCustom(DEF_SIZE)
	return b.String()
}

func JustCustomString(size int) string {
	b := NewCustom(size)
	return b.String()
}

// Custom JSON marshallind and unmarshalling

func (self *BBID) MarshalJSON() ([]byte, error) {
	return []byte("\"" + self.String() + "\""), nil
}

func (self *BBID) UnmarshalJSON(data []byte) (err error) {

	self.data, err = hex.DecodeString(string(data[1 : len(data)-1]))
	if err != nil {
		self.data = make([]byte, DEF_SIZE)
		self.size = DEF_SIZE
		return err
	}
	self.size = len(self.data)
	return nil
}

// Custom BSON marshallind and unmarshalling

func (self *BBID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(self.String())
}

func (self *BBID) UnmarshalBSONValue(t bsontype.Type, data []byte) (err error) {
	self.data, err = hex.DecodeString(string(data[4 : len(data)-1]))
	if err != nil {
		log.Error(err)
		self.data = make([]byte, DEF_SIZE)
		self.size = DEF_SIZE
		return err
	}

	self.size = len(self.data)
	return nil
}
