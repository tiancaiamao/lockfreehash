package lockfreehash

type Uint32Key uint32

func (key Uint32Key) GetHash() uint32 {
	return uint32(key)
}
func (key Uint32Key) Equal(other Key) bool {
	if v, ok := other.(Uint32Key); !ok {
		panic("type  not match")
	} else {
		return key == v
	}
}

type StringKey string

func (str StringKey) GetHash() uint32 {
	var ret uint32
	for v := range str {
		ret += uint32(v<<5 - v)
	}
	return ret
}
func (str StringKey) Equal(other Key) bool {
	if v, ok := other.(StringKey); !ok {
		panic("type not match")
	} else {
		return string(v) == string(str)
	}
}
