package zcnsc

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"0chain.net/smartcontract/provider"
	"github.com/tinylib/msgp/msgp"
)

// MarshalMsg implements msgp.Marshaler
func (z *AuthorizerConfig) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Fee"
	o = append(o, 0x81, 0xa3, 0x46, 0x65, 0x65)
	o, err = z.Fee.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Fee")
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *AuthorizerConfig) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "Fee":
			bts, err = z.Fee.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Fee")
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
func (z *AuthorizerConfig) Msgsize() (s int) {
	s = 1 + 4 + z.Fee.Msgsize()
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *AuthorizerNode) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "Provider"
	o = append(o, 0x84, 0xa8, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72)
	if z.Provider == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.Provider.MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, "Provider")
			return
		}
	}
	// string "PublicKey"
	o = append(o, 0xa9, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79)
	o = msgp.AppendString(o, z.PublicKey)
	// string "URL"
	o = append(o, 0xa3, 0x55, 0x52, 0x4c)
	o = msgp.AppendString(o, z.URL)
	// string "Config"
	o = append(o, 0xa6, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67)
	if z.Config == nil {
		o = msgp.AppendNil(o)
	} else {
		// map header, size 1
		// string "Fee"
		o = append(o, 0x81, 0xa3, 0x46, 0x65, 0x65)
		o, err = z.Config.Fee.MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, "Config", "Fee")
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *AuthorizerNode) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Provider = nil
			} else {
				if z.Provider == nil {
					z.Provider = new(provider.Provider)
				}
				bts, err = z.Provider.UnmarshalMsg(bts)
				if err != nil {
					err = msgp.WrapError(err, "Provider")
					return
				}
			}
		case "PublicKey":
			z.PublicKey, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "PublicKey")
				return
			}
		case "URL":
			z.URL, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "URL")
				return
			}
		case "Config":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.Config = nil
			} else {
				if z.Config == nil {
					z.Config = new(AuthorizerConfig)
				}
				var zb0002 uint32
				zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Config")
					return
				}
				for zb0002 > 0 {
					zb0002--
					field, bts, err = msgp.ReadMapKeyZC(bts)
					if err != nil {
						err = msgp.WrapError(err, "Config")
						return
					}
					switch msgp.UnsafeString(field) {
					case "Fee":
						bts, err = z.Config.Fee.UnmarshalMsg(bts)
						if err != nil {
							err = msgp.WrapError(err, "Config", "Fee")
							return
						}
					default:
						bts, err = msgp.Skip(bts)
						if err != nil {
							err = msgp.WrapError(err, "Config")
							return
						}
					}
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
func (z *AuthorizerNode) Msgsize() (s int) {
	s = 1 + 9
	if z.Provider == nil {
		s += msgp.NilSize
	} else {
		s += z.Provider.Msgsize()
	}
	s += 10 + msgp.StringPrefixSize + len(z.PublicKey) + 4 + msgp.StringPrefixSize + len(z.URL) + 7
	if z.Config == nil {
		s += msgp.NilSize
	} else {
		s += 1 + 4 + z.Config.Fee.Msgsize()
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *GlobalNode) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "ZCNSConfig"
	o = append(o, 0x82, 0xaa, 0x5a, 0x43, 0x4e, 0x53, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67)
	if z.ZCNSConfig == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.ZCNSConfig.MarshalMsg(o)
		if err != nil {
			err = msgp.WrapError(err, "ZCNSConfig")
			return
		}
	}
	// string "ID"
	o = append(o, 0xa2, 0x49, 0x44)
	o = msgp.AppendString(o, z.ID)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *GlobalNode) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "ZCNSConfig":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.ZCNSConfig = nil
			} else {
				if z.ZCNSConfig == nil {
					z.ZCNSConfig = new(ZCNSConfig)
				}
				bts, err = z.ZCNSConfig.UnmarshalMsg(bts)
				if err != nil {
					err = msgp.WrapError(err, "ZCNSConfig")
					return
				}
			}
		case "ID":
			z.ID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ID")
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
func (z *GlobalNode) Msgsize() (s int) {
	s = 1 + 11
	if z.ZCNSConfig == nil {
		s += msgp.NilSize
	} else {
		s += z.ZCNSConfig.Msgsize()
	}
	s += 3 + msgp.StringPrefixSize + len(z.ID)
	return
}

// MarshalMsg implements msgp.Marshaler
func (z UserNode) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "ID"
	o = append(o, 0x82, 0xa2, 0x49, 0x44)
	o = msgp.AppendString(o, z.ID)
	// string "BurnNonce"
	o = append(o, 0xa9, 0x42, 0x75, 0x72, 0x6e, 0x4e, 0x6f, 0x6e, 0x63, 0x65)
	o = msgp.AppendInt64(o, z.BurnNonce)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *UserNode) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "ID":
			z.ID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ID")
				return
			}
		case "BurnNonce":
			z.BurnNonce, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "BurnNonce")
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
func (z UserNode) Msgsize() (s int) {
	s = 1 + 3 + msgp.StringPrefixSize + len(z.ID) + 10 + msgp.Int64Size
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ZCNSConfig) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 11
	// string "MinMintAmount"
	o = append(o, 0x8b, 0xad, 0x4d, 0x69, 0x6e, 0x4d, 0x69, 0x6e, 0x74, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74)
	o, err = z.MinMintAmount.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "MinMintAmount")
		return
	}
	// string "MinBurnAmount"
	o = append(o, 0xad, 0x4d, 0x69, 0x6e, 0x42, 0x75, 0x72, 0x6e, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74)
	o, err = z.MinBurnAmount.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "MinBurnAmount")
		return
	}
	// string "MinStakeAmount"
	o = append(o, 0xae, 0x4d, 0x69, 0x6e, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74)
	o, err = z.MinStakeAmount.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "MinStakeAmount")
		return
	}
	// string "MinLockAmount"
	o = append(o, 0xad, 0x4d, 0x69, 0x6e, 0x4c, 0x6f, 0x63, 0x6b, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74)
	o, err = z.MinLockAmount.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "MinLockAmount")
		return
	}
	// string "MinAuthorizers"
	o = append(o, 0xae, 0x4d, 0x69, 0x6e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x72, 0x73)
	o = msgp.AppendInt64(o, z.MinAuthorizers)
	// string "PercentAuthorizers"
	o = append(o, 0xb2, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x72, 0x73)
	o = msgp.AppendFloat64(o, z.PercentAuthorizers)
	// string "MaxFee"
	o = append(o, 0xa6, 0x4d, 0x61, 0x78, 0x46, 0x65, 0x65)
	o, err = z.MaxFee.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "MaxFee")
		return
	}
	// string "BurnAddress"
	o = append(o, 0xab, 0x42, 0x75, 0x72, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73)
	o = msgp.AppendString(o, z.BurnAddress)
	// string "OwnerId"
	o = append(o, 0xa7, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x49, 0x64)
	o = msgp.AppendString(o, z.OwnerId)
	// string "Cost"
	o = append(o, 0xa4, 0x43, 0x6f, 0x73, 0x74)
	o = msgp.AppendMapHeader(o, uint32(len(z.Cost)))
	keys_za0001 := make([]string, 0, len(z.Cost))
	for k := range z.Cost {
		keys_za0001 = append(keys_za0001, k)
	}
	msgp.Sort(keys_za0001)
	for _, k := range keys_za0001 {
		za0002 := z.Cost[k]
		o = msgp.AppendString(o, k)
		o = msgp.AppendInt(o, za0002)
	}
	// string "MaxDelegates"
	o = append(o, 0xac, 0x4d, 0x61, 0x78, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x73)
	o = msgp.AppendInt(o, z.MaxDelegates)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ZCNSConfig) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
		case "MinMintAmount":
			bts, err = z.MinMintAmount.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "MinMintAmount")
				return
			}
		case "MinBurnAmount":
			bts, err = z.MinBurnAmount.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "MinBurnAmount")
				return
			}
		case "MinStakeAmount":
			bts, err = z.MinStakeAmount.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "MinStakeAmount")
				return
			}
		case "MinLockAmount":
			bts, err = z.MinLockAmount.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "MinLockAmount")
				return
			}
		case "MinAuthorizers":
			z.MinAuthorizers, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MinAuthorizers")
				return
			}
		case "PercentAuthorizers":
			z.PercentAuthorizers, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "PercentAuthorizers")
				return
			}
		case "MaxFee":
			bts, err = z.MaxFee.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "MaxFee")
				return
			}
		case "BurnAddress":
			z.BurnAddress, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "BurnAddress")
				return
			}
		case "OwnerId":
			z.OwnerId, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "OwnerId")
				return
			}
		case "Cost":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Cost")
				return
			}
			if z.Cost == nil {
				z.Cost = make(map[string]int, zb0002)
			} else if len(z.Cost) > 0 {
				for key := range z.Cost {
					delete(z.Cost, key)
				}
			}
			for zb0002 > 0 {
				var za0001 string
				var za0002 int
				zb0002--
				za0001, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Cost")
					return
				}
				za0002, bts, err = msgp.ReadIntBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Cost", za0001)
					return
				}
				z.Cost[za0001] = za0002
			}
		case "MaxDelegates":
			z.MaxDelegates, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MaxDelegates")
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
func (z *ZCNSConfig) Msgsize() (s int) {
	s = 1 + 14 + z.MinMintAmount.Msgsize() + 14 + z.MinBurnAmount.Msgsize() + 15 + z.MinStakeAmount.Msgsize() + 14 + z.MinLockAmount.Msgsize() + 15 + msgp.Int64Size + 19 + msgp.Float64Size + 7 + z.MaxFee.Msgsize() + 12 + msgp.StringPrefixSize + len(z.BurnAddress) + 8 + msgp.StringPrefixSize + len(z.OwnerId) + 5 + msgp.MapHeaderSize
	if z.Cost != nil {
		for za0001, za0002 := range z.Cost {
			_ = za0002
			s += msgp.StringPrefixSize + len(za0001) + msgp.IntSize
		}
	}
	s += 13 + msgp.IntSize
	return
}
