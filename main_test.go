package lockfreehash

import (
	"testing"
)

// func TeststringKey(t *testing.T) {
// 	key1 := StringKey("test1")
// 	key2 := StringKey("test1")
// 	if !key1.Equal(key2) {
// 		t.Fatal("StringKey equal error")
// 	}
// 	t.Log(key1.GetHash())
// 	t.Log(key2.GetHash())
// }

/*
func TestPut(t *testing.T) {
	hash := New()
	hash.Put(StringKey("testing"), 200)
	hash.Put(StringKey("textmate"), 300)
	hash.Put(StringKey("sublime text2"), 400)
	hash.Put(StringKey("vi"), 500)
	hash.Put(StringKey("emacs"), 600)
	val, ok := hash.Get(StringKey("vi"))
	if ok == false || val != 600 {
		t.Fatal("Put error!")
	}
}
*/

func TestBitReverse(t *testing.T) {
	var x uint32 = 0x000000ff
	y := BitReverse(x)
	if y != 0xff000000 {
		t.Fatal("BitReverse wrong!")
	}
	x = 0x00000050f
	y = BitReverse(x)
	if y != 0xf0a00000 {
		t.Fatal("BitReverse wrong 0x0000050f")
	}
	x = uint32(186)
	t.Log(BitReverse(x))
}
