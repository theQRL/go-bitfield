package bitfield

import (
	"encoding/binary"
	"math/bits"
)

var _ = Bitfield(Bitvector16{})

// Bitvector16 is a bitfield with a fixed defined size of 16. There is no length bit
// present in the underlying byte array.
type Bitvector16 []byte

const bitvector16ByteSize = 2
const bitvector16BitSize = bitvector16ByteSize * 8

// NewBitvector16 creates a new bitvector of size 16.
func NewBitvector16() Bitvector16 {
	byteArray := [bitvector16ByteSize]byte{}
	return byteArray[:]
}

// BitAt returns the bit value at the given index. If the index requested
// exceeds the number of bits in the bitvector, then this method returns false.
func (b Bitvector16) BitAt(idx uint64) bool {
	// Out of bounds, must be false.
	if idx >= b.Len() || len(b) != bitvector16ByteSize {
		return false
	}

	i := uint8(1 << (idx % 8))
	return b[idx/8]&i == i
}

// SetBitAt will set the bit at the given index to the given value. If the index
// requested exceeds the number of bits in the bitvector, then this method returns
// false.
func (b Bitvector16) SetBitAt(idx uint64, val bool) {
	// Out of bounds, do nothing.
	if idx >= b.Len() || len(b) != bitvector16ByteSize {
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
func (b Bitvector16) Len() uint64 {
	return bitvector16BitSize
}

// Count returns the number of 1s in the bitvector.
func (b Bitvector16) Count() uint64 {
	if len(b) == 0 {
		return 0
	}
	c := 0
	for i, bt := range b {
		if i >= bitvector16ByteSize {
			break
		}
		c += bits.OnesCount8(bt)
	}
	return uint64(c)
}

// Bytes returns the bytes data representing the bitvector16.
func (b Bitvector16) Bytes() []byte {
	if len(b) == 0 {
		return []byte{}
	}
	ln := min(len(b), bitvector16ByteSize)
	ret := make([]byte, ln)
	copy(ret, b[:ln])
	return ret[:]
}

// Shift bitvector by i. If i >= 0, perform left shift, otherwise right shift.
func (b Bitvector16) Shift(i int) {
	if len(b) == 0 {
		return
	}

	// Shifting greater than 16 bits is pointless and can have unexpected behavior.
	if i > bitvector16BitSize {
		i = bitvector16BitSize
	} else if i < -bitvector16BitSize {
		i = -bitvector16BitSize
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

// BitIndices returns the list of indices which are set to 1.
func (b Bitvector16) BitIndices() []int {
	indices := make([]int, 0, bitvector16BitSize)
	for i, bt := range b {
		if i >= bitvector16ByteSize {
			break
		}
		for j := 0; j < 8; j++ {
			bit := byte(1 << uint(j))
			if bt&bit == bit {
				indices = append(indices, i*8+j)
			}
		}
	}

	return indices
}

// Contains returns true if the bitlist contains all of the bits from the provided argument
// bitlist. This method will return an error if bitlists are not the same length.
func (b Bitvector16) Contains(c Bitvector16) (bool, error) {
	if b.Len() != c.Len() {
		return false, ErrBitvectorDifferentLength
	}

	// To ensure all of the bits in c are present in b, we iterate over every byte, combine
	// the byte from b and c, then XOR them against b. If the result of this is non-zero, then we
	// are assured that a byte in c had bits not present in b.
	for i := 0; i < len(b); i++ {
		if b[i]^(b[i]|c[i]) != 0 {
			return false, nil
		}
	}

	return true, nil
}

// Overlaps returns true if the bitlist contains one of the bits from the provided argument
// bitlist. This method will return an error if bitlists are not the same length.
func (b Bitvector16) Overlaps(c Bitvector16) (bool, error) {
	lenB, lenC := b.Len(), c.Len()
	if b.Len() != c.Len() {
		return false, ErrBitvectorDifferentLength
	}

	if lenB == 0 || lenC == 0 {
		return false, nil
	}

	// To ensure all of the bits in c are not overlapped in b, we iterate over every byte, invert b
	// and xor the byte from b and c, then and it against c. If the result is non-zero, then
	// we can be assured that byte in c had bits not overlapped in b.
	for i := 0; i < len(b); i++ {
		// If this byte is the last byte in the array, mask the length bit.
		mask := uint8(0xFF)

		if (^b[i]^c[i])&c[i]&mask != 0 {
			return true, nil
		}
	}
	return false, nil
}

// Or returns the OR result of the two bitfields. This method will return an error if the bitlists are not the same length.
func (b Bitvector16) Or(c Bitvector16) (Bitvector16, error) {
	if b.Len() != c.Len() {
		return nil, ErrBitvectorDifferentLength
	}

	ret := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		ret[i] = b[i] | c[i]
	}

	return ret, nil
}
