package easycache

// 缓存值的抽象与封装,不可篡改，也就是只读的。
//b将会存储真实的缓存值。选择 byte 类型是为了能够支持任意的数据类型的存储，例如字符串、图片等。
type ByteView struct {
	b []byte
}

// Len returns the view's length
//返回byteview的字节长度
//我们在 lru.Cache 的实现中，要求被缓存对象必须实现 Value 接口，即 Len() int 方法，返回其所占的内存大小。
func (bv ByteView) Len() int {
	return len(bv.b)
}

//因为要保证b是只读的，这个方法会返回b的一个拷贝，防止缓存值被外部程序修改。
func (bv ByteView) ByteSlice() []byte {
	return cloneBytes(bv.b)
}

//将byteview的数据转换成字符串形式返回
func (bv ByteView) String() string {
	return string(bv.b)
}

//做了一个拷贝返回
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
