package utils

import (
	"bytes"
	"encoding/binary"
	"math"
)

// FloatSlice64To32 convert float64 slice to float32 slice
func FloatSlice64To32(f64 []float64) []float32 {
	f32 := make([]float32, len(f64))
	for i, v := range f64 {
		f32[i] = float32(v)
	}
	return f32
}

// Float64ToBytes convert float64 to bytes
func Float64ToByte(f float64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

func Float32SliceToBytes(s []float32) []byte {
	buf := new(bytes.Buffer)
	for _, f := range s {
		bs := Float64ToByte(float64(f))
		buf.Write(bs)
	}
	return buf.Bytes()
}
