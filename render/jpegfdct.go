// Forward DCT for JPEG encoder.
//
// Forked from Go's standard library image/jpeg package.
// Original source: https://go.googlesource.com/go/+/refs/heads/master/src/image/jpeg/fdct.go
// Copyright 2011 The Go Authors. All rights reserved.
// Licensed under the BSD 3-Clause License: https://go.dev/LICENSE
//
// Based on jfdctint.c from the Independent JPEG Group (http://www.ijg.org).
// See the original source for the IJG license terms.
package render

const (
	jpegFix_0_298631336 = 2446
	jpegFix_0_390180644 = 3196
	jpegFix_0_541196100 = 4433
	jpegFix_0_765366865 = 6270
	jpegFix_0_899976223 = 7373
	jpegFix_1_175875602 = 9633
	jpegFix_1_501321110 = 12299
	jpegFix_1_847759065 = 15137
	jpegFix_1_961570560 = 16069
	jpegFix_2_053119869 = 16819
	jpegFix_2_562915447 = 20995
	jpegFix_3_072711026 = 25172
)

const (
	jpegConstBits     = 13
	jpegPass1Bits     = 2
	jpegCenterJSample = 128
)

func jpegFDCT(b *jpegBlock) {
	for y := 0; y < 8; y++ {
		y8 := y * 8
		s := b[y8 : y8+8 : y8+8]
		x0 := s[0]
		x1 := s[1]
		x2 := s[2]
		x3 := s[3]
		x4 := s[4]
		x5 := s[5]
		x6 := s[6]
		x7 := s[7]

		tmp0 := x0 + x7
		tmp1 := x1 + x6
		tmp2 := x2 + x5
		tmp3 := x3 + x4

		tmp10 := tmp0 + tmp3
		tmp12 := tmp0 - tmp3
		tmp11 := tmp1 + tmp2
		tmp13 := tmp1 - tmp2

		tmp0 = x0 - x7
		tmp1 = x1 - x6
		tmp2 = x2 - x5
		tmp3 = x3 - x4

		s[0] = (tmp10 + tmp11 - 8*jpegCenterJSample) << jpegPass1Bits
		s[4] = (tmp10 - tmp11) << jpegPass1Bits
		z1 := (tmp12 + tmp13) * jpegFix_0_541196100
		z1 += 1 << (jpegConstBits - jpegPass1Bits - 1)
		s[2] = (z1 + tmp12*jpegFix_0_765366865) >> (jpegConstBits - jpegPass1Bits)
		s[6] = (z1 - tmp13*jpegFix_1_847759065) >> (jpegConstBits - jpegPass1Bits)

		tmp10 = tmp0 + tmp3
		tmp11 = tmp1 + tmp2
		tmp12 = tmp0 + tmp2
		tmp13 = tmp1 + tmp3
		z1 = (tmp12 + tmp13) * jpegFix_1_175875602
		z1 += 1 << (jpegConstBits - jpegPass1Bits - 1)
		tmp0 *= jpegFix_1_501321110
		tmp1 *= jpegFix_3_072711026
		tmp2 *= jpegFix_2_053119869
		tmp3 *= jpegFix_0_298631336
		tmp10 *= -jpegFix_0_899976223
		tmp11 *= -jpegFix_2_562915447
		tmp12 *= -jpegFix_0_390180644
		tmp13 *= -jpegFix_1_961570560

		tmp12 += z1
		tmp13 += z1
		s[1] = (tmp0 + tmp10 + tmp12) >> (jpegConstBits - jpegPass1Bits)
		s[3] = (tmp1 + tmp11 + tmp13) >> (jpegConstBits - jpegPass1Bits)
		s[5] = (tmp2 + tmp11 + tmp12) >> (jpegConstBits - jpegPass1Bits)
		s[7] = (tmp3 + tmp10 + tmp13) >> (jpegConstBits - jpegPass1Bits)
	}
	for x := 0; x < 8; x++ {
		tmp0 := b[0*8+x] + b[7*8+x]
		tmp1 := b[1*8+x] + b[6*8+x]
		tmp2 := b[2*8+x] + b[5*8+x]
		tmp3 := b[3*8+x] + b[4*8+x]

		tmp10 := tmp0 + tmp3 + 1<<(jpegPass1Bits-1)
		tmp12 := tmp0 - tmp3
		tmp11 := tmp1 + tmp2
		tmp13 := tmp1 - tmp2

		tmp0 = b[0*8+x] - b[7*8+x]
		tmp1 = b[1*8+x] - b[6*8+x]
		tmp2 = b[2*8+x] - b[5*8+x]
		tmp3 = b[3*8+x] - b[4*8+x]

		b[0*8+x] = (tmp10 + tmp11) >> jpegPass1Bits
		b[4*8+x] = (tmp10 - tmp11) >> jpegPass1Bits

		z1 := (tmp12 + tmp13) * jpegFix_0_541196100
		z1 += 1 << (jpegConstBits + jpegPass1Bits - 1)
		b[2*8+x] = (z1 + tmp12*jpegFix_0_765366865) >> (jpegConstBits + jpegPass1Bits)
		b[6*8+x] = (z1 - tmp13*jpegFix_1_847759065) >> (jpegConstBits + jpegPass1Bits)

		tmp10 = tmp0 + tmp3
		tmp11 = tmp1 + tmp2
		tmp12 = tmp0 + tmp2
		tmp13 = tmp1 + tmp3
		z1 = (tmp12 + tmp13) * jpegFix_1_175875602
		z1 += 1 << (jpegConstBits + jpegPass1Bits - 1)
		tmp0 *= jpegFix_1_501321110
		tmp1 *= jpegFix_3_072711026
		tmp2 *= jpegFix_2_053119869
		tmp3 *= jpegFix_0_298631336
		tmp10 *= -jpegFix_0_899976223
		tmp11 *= -jpegFix_2_562915447
		tmp12 *= -jpegFix_0_390180644
		tmp13 *= -jpegFix_1_961570560

		tmp12 += z1
		tmp13 += z1
		b[1*8+x] = (tmp0 + tmp10 + tmp12) >> (jpegConstBits + jpegPass1Bits)
		b[3*8+x] = (tmp1 + tmp11 + tmp13) >> (jpegConstBits + jpegPass1Bits)
		b[5*8+x] = (tmp2 + tmp11 + tmp12) >> (jpegConstBits + jpegPass1Bits)
		b[7*8+x] = (tmp3 + tmp10 + tmp13) >> (jpegConstBits + jpegPass1Bits)
	}
}
