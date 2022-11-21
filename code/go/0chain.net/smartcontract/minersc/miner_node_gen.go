package minersc

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"0chain.net/smartcontract/stakepool"
	"github.com/tinylib/msgp/msgp"
)

// MarshalMsg implements msgp.Marshaler
func (z *MinerNode) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "SimpleNode"
	o = append(o, 0x82, 0xaa, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x4e, 0x6f, 0x64, 0x65)
	if z.SimpleNode == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.SimpleNode.MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, "SimpleNode")
			return
		}
	}
	// string "Provider"
	o = append(o, 0xa9, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x50, 0x6f, 0x6f, 0x6c)
	if z.StakePool == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.StakePool.MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, "Provider")
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *MinerNode) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "SimpleNode":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.SimpleNode = nil
			} else {
				if z.SimpleNode == nil {
					z.SimpleNode = new(SimpleNode)
				}
				bts, err = z.SimpleNode.UnmarshalMsg(bts)
				if err != nil {
					err = msgp.WrapError(err, "SimpleNode")
					return
				}
			}
		case "Provider":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.StakePool = nil
			} else {
				if z.StakePool == nil {
					z.StakePool = new(stakepool.StakePool)
				}
				bts, err = z.StakePool.UnmarshalMsg(bts)
				if err != nil {
					err = msgp.WrapError(err, "Provider")
					return
				}
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
func (z *MinerNode) Msgsize() (s int) {
	s = 1 + 11
	if z.SimpleNode == nil {
		s += msgp.NilSize
	} else {
		s += z.SimpleNode.Msgsize()
	}
	s += 10
	if z.StakePool == nil {
		s += msgp.NilSize
	} else {
		s += z.StakePool.Msgsize()
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *NodePool) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "PoolID"
	o = append(o, 0x82, 0xa6, 0x50, 0x6f, 0x6f, 0x6c, 0x49, 0x44)
	o = msgp.AppendString(o, z.PoolID)
	// string "DelegatePool"
	o = append(o, 0xac, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x6f, 0x6c)
	if z.DelegatePool == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.DelegatePool.MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, "DelegatePool")
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *NodePool) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "PoolID":
			z.PoolID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "PoolID")
				return
			}
		case "DelegatePool":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.DelegatePool = nil
			} else {
				if z.DelegatePool == nil {
					z.DelegatePool = new(stakepool.DelegatePool)
				}
				bts, err = z.DelegatePool.UnmarshalMsg(bts)
				if err != nil {
					err = msgp.WrapError(err, "DelegatePool")
					return
				}
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
func (z *NodePool) Msgsize() (s int) {
	s = 1 + 7 + msgp.StringPrefixSize + len(z.PoolID) + 13
	if z.DelegatePool == nil {
		s += msgp.NilSize
	} else {
		s += z.DelegatePool.Msgsize()
	}
	return
}
