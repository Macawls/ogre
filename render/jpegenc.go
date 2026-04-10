// JPEG 4:4:4 encoder.
//
// Forked from Go's standard library image/jpeg package.
// Original source: https://go.googlesource.com/go/+/refs/heads/master/src/image/jpeg/writer.go
// Copyright 2011 The Go Authors. All rights reserved.
// Licensed under the BSD 3-Clause License: https://go.dev/LICENSE
//
// The stdlib encoder hardcodes 4:2:0 chroma subsampling (see Go issue #13614)
// which causes visible color shifts on dark uniform backgrounds. This fork
// changes the sampling factors to 4:4:4 (no subsampling) for accurate color
// reproduction. The only structural change is in writeSOF0 (sampling factors)
// and writeSOS (8x8 MCU iteration instead of 16x16 with scale()).
package render

import (
	"bufio"
	"errors"
	"image"
	"image/color"
	"io"
)

const jpegBlockSize = 64

type jpegBlock [jpegBlockSize]int32

func jpegDiv(a, b int32) int32 {
	if a >= 0 {
		return (a + (b >> 1)) / b
	}
	return -((-a + (b >> 1)) / b)
}

var jpegBitCount = [256]byte{
	0, 1, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
	8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8,
}

type jpegQuantIndex int

const (
	jpegQuantLuminance jpegQuantIndex = iota
	jpegQuantChrominance
	jpegNQuantIndex
)

var jpegUnscaledQuant = [jpegNQuantIndex][jpegBlockSize]byte{
	{
		16, 11, 12, 14, 12, 10, 16, 14,
		13, 14, 18, 17, 16, 19, 24, 40,
		26, 24, 22, 22, 24, 49, 35, 37,
		29, 40, 58, 51, 61, 60, 57, 51,
		56, 55, 64, 72, 92, 78, 64, 68,
		87, 69, 55, 56, 80, 109, 81, 87,
		95, 98, 103, 104, 103, 62, 77, 113,
		121, 112, 100, 120, 92, 101, 103, 99,
	},
	{
		17, 18, 18, 24, 21, 24, 47, 26,
		26, 47, 99, 66, 56, 66, 99, 99,
		99, 99, 99, 99, 99, 99, 99, 99,
		99, 99, 99, 99, 99, 99, 99, 99,
		99, 99, 99, 99, 99, 99, 99, 99,
		99, 99, 99, 99, 99, 99, 99, 99,
		99, 99, 99, 99, 99, 99, 99, 99,
		99, 99, 99, 99, 99, 99, 99, 99,
	},
}

type jpegHuffIndex int

const (
	jpegHuffLumDC jpegHuffIndex = iota
	jpegHuffLumAC
	jpegHuffChrDC
	jpegHuffChrAC
	jpegNHuffIndex
)

type jpegHuffSpec struct {
	count [16]byte
	value []byte
}

var jpegHuffSpecs = [jpegNHuffIndex]jpegHuffSpec{
	// Luminance DC.
	{[16]byte{0, 1, 5, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0}, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}},
	// Luminance AC.
	{[16]byte{0, 2, 1, 3, 3, 2, 4, 3, 5, 5, 4, 4, 0, 0, 1, 125}, []byte{
		0x01, 0x02, 0x03, 0x00, 0x04, 0x11, 0x05, 0x12, 0x21, 0x31, 0x41, 0x06, 0x13, 0x51, 0x61, 0x07,
		0x22, 0x71, 0x14, 0x32, 0x81, 0x91, 0xa1, 0x08, 0x23, 0x42, 0xb1, 0xc1, 0x15, 0x52, 0xd1, 0xf0,
		0x24, 0x33, 0x62, 0x72, 0x82, 0x09, 0x0a, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x25, 0x26, 0x27, 0x28,
		0x29, 0x2a, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49,
		0x4a, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69,
		0x6a, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89,
		0x8a, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7,
		0xa8, 0xa9, 0xaa, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xc2, 0xc3, 0xc4, 0xc5,
		0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xe1, 0xe2,
		0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8,
		0xf9, 0xfa,
	}},
	// Chrominance DC.
	{[16]byte{0, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0}, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}},
	// Chrominance AC.
	{[16]byte{0, 2, 1, 2, 4, 4, 3, 4, 7, 5, 4, 4, 0, 1, 2, 119}, []byte{
		0x00, 0x01, 0x02, 0x03, 0x11, 0x04, 0x05, 0x21, 0x31, 0x06, 0x12, 0x41, 0x51, 0x07, 0x61, 0x71,
		0x13, 0x22, 0x32, 0x81, 0x08, 0x14, 0x42, 0x91, 0xa1, 0xb1, 0xc1, 0x09, 0x23, 0x33, 0x52, 0xf0,
		0x15, 0x62, 0x72, 0xd1, 0x0a, 0x16, 0x24, 0x34, 0xe1, 0x25, 0xf1, 0x17, 0x18, 0x19, 0x1a, 0x26,
		0x27, 0x28, 0x29, 0x2a, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48,
		0x49, 0x4a, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68,
		0x69, 0x6a, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87,
		0x88, 0x89, 0x8a, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0xa2, 0xa3, 0xa4, 0xa5,
		0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xc2, 0xc3,
		0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda,
		0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8,
		0xf9, 0xfa,
	}},
}

type jpegHuffLUT []uint32

func (h *jpegHuffLUT) init(s jpegHuffSpec) {
	maxValue := 0
	for _, v := range s.value {
		if int(v) > maxValue {
			maxValue = int(v)
		}
	}
	*h = make([]uint32, maxValue+1)
	code, k := uint32(0), 0
	for i := 0; i < len(s.count); i++ {
		nBits := uint32(i+1) << 24
		for j := uint8(0); j < s.count[i]; j++ {
			(*h)[s.value[k]] = nBits | code
			code++
			k++
		}
		code <<= 1
	}
}

var jpegHuffLUTs [4]jpegHuffLUT

func init() {
	for i, s := range jpegHuffSpecs {
		jpegHuffLUTs[i].init(s)
	}
}

var jpegUnzig = [jpegBlockSize]int{
	0, 1, 8, 16, 9, 2, 3, 10,
	17, 24, 32, 25, 18, 11, 4, 5,
	12, 19, 26, 33, 40, 48, 41, 34,
	27, 20, 13, 6, 7, 14, 21, 28,
	35, 42, 49, 56, 57, 50, 43, 36,
	29, 22, 15, 23, 30, 37, 44, 51,
	58, 59, 52, 45, 38, 31, 39, 46,
	53, 60, 61, 54, 47, 55, 62, 63,
}

const (
	jpegSOF0 = 0xc0
	jpegDHT  = 0xc4
	jpegSOI  = 0xd8
	jpegEOI  = 0xd9
	jpegSOS  = 0xda
	jpegDQT  = 0xdb
)

type jpegWriter interface {
	Flush() error
	io.Writer
	io.ByteWriter
}

type jpegEncoder struct {
	w          jpegWriter
	err        error
	buf        [16]byte
	bits       uint32
	nBits      uint32
	quant      [jpegNQuantIndex][jpegBlockSize]byte
}

func (e *jpegEncoder) flush() {
	if e.err != nil {
		return
	}
	e.err = e.w.Flush()
}

func (e *jpegEncoder) write(p []byte) {
	if e.err != nil {
		return
	}
	_, e.err = e.w.Write(p)
}

func (e *jpegEncoder) writeByte(b byte) {
	if e.err != nil {
		return
	}
	e.err = e.w.WriteByte(b)
}

func (e *jpegEncoder) emit(bits, nBits uint32) {
	nBits += e.nBits
	bits <<= 32 - nBits
	bits |= e.bits
	for nBits >= 8 {
		b := uint8(bits >> 24)
		e.writeByte(b)
		if b == 0xff {
			e.writeByte(0x00)
		}
		bits <<= 8
		nBits -= 8
	}
	e.bits, e.nBits = bits, nBits
}

func (e *jpegEncoder) emitHuff(h jpegHuffIndex, value int32) {
	x := jpegHuffLUTs[h][value]
	e.emit(x&(1<<24-1), x>>24)
}

func (e *jpegEncoder) emitHuffRLE(h jpegHuffIndex, runLength, value int32) {
	a, b := value, value
	if a < 0 {
		a, b = -value, value-1
	}
	var nBits uint32
	if a < 0x100 {
		nBits = uint32(jpegBitCount[a])
	} else {
		nBits = 8 + uint32(jpegBitCount[a>>8])
	}
	e.emitHuff(h, runLength<<4|int32(nBits))
	if nBits > 0 {
		e.emit(uint32(b)&(1<<nBits-1), nBits)
	}
}

func (e *jpegEncoder) writeMarkerHeader(marker uint8, markerlen int) {
	e.buf[0] = 0xff
	e.buf[1] = marker
	e.buf[2] = uint8(markerlen >> 8)
	e.buf[3] = uint8(markerlen & 0xff)
	e.write(e.buf[:4])
}

func (e *jpegEncoder) writeDQT() {
	const markerlen = 2 + int(jpegNQuantIndex)*(1+jpegBlockSize)
	e.writeMarkerHeader(jpegDQT, markerlen)
	for i := range e.quant {
		e.writeByte(uint8(i))
		e.write(e.quant[i][:])
	}
}

// writeSOF0 writes the Start Of Frame marker with 4:4:4 sampling.
func (e *jpegEncoder) writeSOF0(size image.Point) {
	const nComponent = 3
	markerlen := 8 + 3*nComponent
	e.writeMarkerHeader(jpegSOF0, markerlen)
	e.buf[0] = 8
	e.buf[1] = uint8(size.Y >> 8)
	e.buf[2] = uint8(size.Y & 0xff)
	e.buf[3] = uint8(size.X >> 8)
	e.buf[4] = uint8(size.X & 0xff)
	e.buf[5] = nComponent
	for i := 0; i < nComponent; i++ {
		e.buf[3*i+6] = uint8(i + 1)
		// 4:4:4 — all components use 1x1 sampling.
		e.buf[3*i+7] = 0x11
		e.buf[3*i+8] = "\x00\x01\x01"[i]
	}
	e.write(e.buf[:3*(nComponent-1)+9])
}

func (e *jpegEncoder) writeDHT() {
	const nComponent = 3
	markerlen := 2
	for _, s := range jpegHuffSpecs {
		markerlen += 1 + 16 + len(s.value)
	}
	e.writeMarkerHeader(jpegDHT, markerlen)
	for i, s := range jpegHuffSpecs {
		e.writeByte("\x00\x10\x01\x11"[i])
		e.write(s.count[:])
		e.write(s.value)
	}
}

func (e *jpegEncoder) writeBlock(b *jpegBlock, q jpegQuantIndex, prevDC int32) int32 {
	jpegFDCT(b)
	dc := jpegDiv(b[0], 8*int32(e.quant[q][0]))
	e.emitHuffRLE(jpegHuffIndex(2*q+0), 0, dc-prevDC)
	h, runLength := jpegHuffIndex(2*q+1), int32(0)
	for zig := 1; zig < jpegBlockSize; zig++ {
		ac := jpegDiv(b[jpegUnzig[zig]], 8*int32(e.quant[q][zig]))
		if ac == 0 {
			runLength++
		} else {
			for runLength > 15 {
				e.emitHuff(h, 0xf0)
				runLength -= 16
			}
			e.emitHuffRLE(h, runLength, ac)
			runLength = 0
		}
	}
	if runLength > 0 {
		e.emitHuff(h, 0x00)
	}
	return dc
}

func jpegRGBAToYCbCr(m *image.RGBA, p image.Point, yBlock, cbBlock, crBlock *jpegBlock) {
	b := m.Bounds()
	xmax := b.Max.X - 1
	ymax := b.Max.Y - 1
	for j := 0; j < 8; j++ {
		sj := p.Y + j
		if sj > ymax {
			sj = ymax
		}
		offset := (sj-b.Min.Y)*m.Stride - b.Min.X*4
		for i := 0; i < 8; i++ {
			sx := p.X + i
			if sx > xmax {
				sx = xmax
			}
			pix := m.Pix[offset+sx*4:]
			yy, cb, cr := color.RGBToYCbCr(pix[0], pix[1], pix[2])
			yBlock[8*j+i] = int32(yy)
			cbBlock[8*j+i] = int32(cb)
			crBlock[8*j+i] = int32(cr)
		}
	}
}

func jpegToYCbCr(m image.Image, p image.Point, yBlock, cbBlock, crBlock *jpegBlock) {
	b := m.Bounds()
	xmax := b.Max.X - 1
	ymax := b.Max.Y - 1
	for j := 0; j < 8; j++ {
		for i := 0; i < 8; i++ {
			r, g, bb, _ := m.At(min(p.X+i, xmax), min(p.Y+j, ymax)).RGBA()
			yy, cb, cr := color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(bb>>8))
			yBlock[8*j+i] = int32(yy)
			cbBlock[8*j+i] = int32(cb)
			crBlock[8*j+i] = int32(cr)
		}
	}
}

// SOS header for YCbCr.
var jpegSOSHeaderYCbCr = []byte{
	0xff, 0xda, 0x00, 0x0c, 0x03, 0x01, 0x00, 0x02,
	0x11, 0x03, 0x11, 0x00, 0x3f, 0x00,
}

// writeSOS writes scan data using 4:4:4 — one Y, one Cb, one Cr block per 8x8 MCU.
func (e *jpegEncoder) writeSOS(m image.Image) {
	e.write(jpegSOSHeaderYCbCr)
	var (
		yb, cbb, crb                   jpegBlock
		prevDCY, prevDCCb, prevDCCr int32
	)
	bounds := m.Bounds()
	rgba, _ := m.(*image.RGBA)
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 8 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 8 {
			p := image.Pt(x, y)
			if rgba != nil {
				jpegRGBAToYCbCr(rgba, p, &yb, &cbb, &crb)
			} else {
				jpegToYCbCr(m, p, &yb, &cbb, &crb)
			}
			prevDCY = e.writeBlock(&yb, jpegQuantLuminance, prevDCY)
			prevDCCb = e.writeBlock(&cbb, jpegQuantChrominance, prevDCCb)
			prevDCCr = e.writeBlock(&crb, jpegQuantChrominance, prevDCCr)
		}
	}
	e.emit(0x7f, 7)
}

// encodeJPEG444 encodes m as JPEG with 4:4:4 chroma subsampling (no subsampling).
func encodeJPEG444(w io.Writer, m image.Image, quality int) error {
	b := m.Bounds()
	if b.Dx() >= 1<<16 || b.Dy() >= 1<<16 {
		return errors.New("jpeg: image is too large to encode")
	}
	var e jpegEncoder
	if ww, ok := w.(jpegWriter); ok {
		e.w = ww
	} else {
		e.w = bufio.NewWriter(w)
	}
	if quality < 1 {
		quality = 1
	} else if quality > 100 {
		quality = 100
	}
	var scale int
	if quality < 50 {
		scale = 5000 / quality
	} else {
		scale = 200 - quality*2
	}
	for i := range e.quant {
		for j := range e.quant[i] {
			x := int(jpegUnscaledQuant[i][j])
			x = (x*scale + 50) / 100
			if x < 1 {
				x = 1
			} else if x > 255 {
				x = 255
			}
			e.quant[i][j] = uint8(x)
		}
	}
	// SOI
	e.buf[0] = 0xff
	e.buf[1] = jpegSOI
	e.write(e.buf[:2])
	e.writeDQT()
	e.writeSOF0(b.Size())
	e.writeDHT()
	e.writeSOS(m)
	// EOI
	e.buf[0] = 0xff
	e.buf[1] = jpegEOI
	e.write(e.buf[:2])
	e.flush()
	return e.err
}
