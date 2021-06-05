package bitfield

import (
	"encoding/binary"
	"math/bits"
)

var _ = Bitfield(Bitvector256{})

// Bitvector256 is a bitfield with a fixed defined size of 256. There is no length bit
// present in the underlying byte array.
type Bitvector256 []byte

// NewBitvector256 creates a new bitvector of size 256.
func NewBitvector256() Bitvector256 {
	byteArray := [32]byte{}
	return byteArray[:]
}

// BitAt returns the bit value at the given index. If the index requested
// exceeds the number of bits in the bitvector, then this method returns false.
func (b Bitvector256) BitAt(idx uint64) bool {
	// Out of bounds, must be false.
	if idx >= b.Len() {
		return false
	}

	i := uint8(1 << (idx % 8))
	return b[idx/8]&i == i
}

// SetBitAt will set the bit at the given index to the given value. If the index
// requested exceeds the number of bits in the bitvector, then this method returns
// false.
func (b Bitvector256) SetBitAt(idx uint64, val bool) {
	// Out of bounds, do nothing.
	if idx >= b.Len() {
		return
	}

	bit := uint8(1 << (idx % 8))
	if val {
		b[idx/8] |= bit
	} else {
		b[idx/8] &^= bit
	}
}

// Len returns the number of bits in the bitvector.
func (b Bitvector256) Len() uint64 {
	return 256
}

// Count returns the number of 1s in the bitvector.
func (b Bitvector256) Count() uint64 {
	if len(b) == 0 {
		return 0
	}
	c := 0
	for _, bt := range b {
		c += bits.OnesCount8(bt)
	}
	return uint64(c)
}

// Bytes returns the bytes data representing the Bitvector256. This method
// bitmasks the underlying data to ensure that it is an accurate representation.
func (b Bitvector256) Bytes() []byte {
	if len(b) == 0 {
		return []byte{}
	}
	ret := make([]byte, len(b))
	copy(ret, b[:])
	return ret[:]
}

// Shift bitvector by i. If i >= 0, perform left shift, otherwise right shift.
func (b Bitvector256) Shift(i int) {
	if len(b) == 0 {
		return
	}

	// Shifting greater than 256 bits is pointless and can have unexpected behavior.
	if i > 256 {
		i = 256
	} else if i < -256 {
		i = -256
	}
	if i >= 0 {
		num := binary.BigEndian.Uint64(b)
		num <<= uint8(i)
		binary.BigEndian.PutUint64(b, num)
	} else {
		num := binary.BigEndian.Uint64(b)
		num >>= uint8(i * -1)
		binary.BigEndian.PutUint64(b, num)
	}
}

// BitIndices returns the list of indices that are set to 1.
func (b Bitvector256) BitIndices() []int {
	indices := make([]int, 0, 256)
	for i, bt := range b {
		for j := 0; j < 8; j++ {
			bit := byte(1 << uint(j))
			if bt&bit == bit {
				indices = append(indices, i*8+j)
			}
		}
	}

	return indices
}
