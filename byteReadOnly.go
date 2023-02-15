package minicache

// 用于缓存的字节数组
type ByteReadOnly struct {
	b []byte
}

// 获取字节数组的长度
func (v ByteReadOnly) Len() int {
	return len(v.b)
}

// 获取字节数组内容的string类型
func (v ByteReadOnly) String() string {
	return string(v.b)
}

// 获取字节数组内容的拷贝
func (v ByteReadOnly) ByteSliceCopy() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
