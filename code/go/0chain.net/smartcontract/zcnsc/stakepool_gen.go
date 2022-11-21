package zcnsc

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// MarshalMsg implements msgp.Marshaler
func (z *StakePool) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Provider"
	o = append(o, 0x81, 0xa9, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x50, 0x6f, 0x6f, 0x6c)
	o, err = z.StakePool.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Provider")
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *StakePool) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Provider":
			bts, err = z.StakePool.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Provider")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *StakePool) Msgsize() (s int) {
	s = 1 + 10 + z.StakePool.Msgsize()
	return
}
