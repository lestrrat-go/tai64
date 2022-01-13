// Package tai64 contains the implementation for TAI64N timestamps

package tai64

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

// N represents a timestamp in tai64n format.
type N struct {
	secs  uint64
	nsecs uint32
}

// Unused as of yet
// type NA struct {
// 	N
// 	asecs uint32
// }

const (
	epochBase = uint64(0x4000000000000000)

	// Size is the size of tai64 objects
	Size = 8

	// NSize is the size of tai64n objects
	NSize = 12

	// NLabelSize represents the number of characters in a tai64n label
	NLabelSize = NSize * 2

	// Unused as of yet
	// NASize = 16
	// NALabelSize = NASize * 2
)

var zeroN = N{}

// ParseN parses a sequence of raw data representing a TAI64N label.
// The input must be exactly 12 bytes.
func ParseN(in []byte) (N, error) {
	if len(in) != NSize {
		return zeroN, fmt.Errorf(`expected input to be %d bytes long, got %d bytes`, NSize, len(in))
	}

	var n N
	n.secs = binary.BigEndian.Uint64(in[:Size]) - epochBase
	n.nsecs = binary.BigEndian.Uint32(in[Size:NSize])
	return n, nil
}

// ParseNLabel parses a tai64n label.
// The input must be either 24 or 25 bytes. If the input size is 25 bytes,
// it should contain a '@' at the beginning
func ParseNLabel(in []byte) (N, error) {
	switch l := len(in); l {
	case NLabelSize + 1:
		if in[0] != '@' {
			return zeroN, fmt.Errorf(`expected label of size %d to start with '@'`, l)
		}
		in = in[1:]
	case NLabelSize:
		// nothing to do
	default:
		return zeroN, fmt.Errorf(`label must be of size %d (or %d for labels starting with '@')`, NLabelSize, NLabelSize+1)
	}

	var tmp [NSize]byte
	w, err := hex.Decode(tmp[:], in)
	if err != nil {
		return zeroN, fmt.Errorf(`failed to decode hex string: %w`, err)
	}
	if w != len(tmp) {
		return zeroN, fmt.Errorf(`failed to decode hex string (expected %d bytes, only decoded %d bytes)`, len(tmp), w)
	}
	return ParseN(tmp[:])
}

// Time returns the time.Time representation of tai64n object.
func (n N) Time() time.Time {
	return time.Unix(int64(n.secs), int64(n.nsecs))
}

// Format writes the hex-encoded representation of the TAI64N object.
// The resulting string does not contain a '@' in the beginning.
//
// Use `Write()` if you would like to write the raw bytes.
func (n N) Format(dst []byte) error {
	var tmp [NSize]byte
	n.Write(tmp[:])
	hex.Encode(dst, tmp[:])
	return nil
}

// Write writes the raw bytes of the TAI64N object.
// The destination must be more than 12 bytes.
func (n N) Write(dst []byte) (int, error) {
	binary.BigEndian.PutUint64(dst[:Size], n.secs+epochBase)
	binary.BigEndian.PutUint32(dst[Size:NSize], n.nsecs)
	return NLabelSize, nil
}

// MarshalText is exactly the same as Format, but it exists to
// fulfill the TextMarshaler interface
func (n N) MarshalText() ([]byte, error) {
	dst := make([]byte, NLabelSize)
	if err := n.Format(dst); err != nil {
		return nil, fmt.Errorf(`failed to write tai64n value into byte array: %w`, err)
	}
	return dst, nil
}
