// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lookup3

import (
	"bytes"
	"hash"
	"testing"
)

type golden struct {
	sum  []byte
	text string
}

var golden32 = []golden{
    {[]byte{0xde, 0xad, 0xbe, 0xef}, ""},
    {[]byte{0x58, 0xd6, 0x87, 0x8}, "a"},
    {[]byte{0xfb, 0xb3, 0xa8, 0xdf}, "ab"},
    {[]byte{0xe, 0x39, 0x76, 0x31}, "abc"},
    {[]byte{0x49, 0x6f, 0x81, 0x9a}, "sK2iisPVTchSsRXIBTPUSCWswVsWVB0s9Qsve"},
    {[]byte{0xb0, 0x83, 0xcc, 0xe3}, "J4b58TgCOdroAvWzHN1HFZQQ"},
    {[]byte{0xce, 0x4f, 0xb, 0xb5}, "LbycZRyoRYqYtw9dzyBOuvQQByaOUcY"},
    {[]byte{0xde, 0xad, 0xbe, 0xef}, ""},
    {[]byte{0x93, 0x10, 0x42, 0x49}, "Sy2ZzcNt6avMfdQo4e2pQTjGs4hfAi7rQo"},
    {[]byte{0x31, 0x1d, 0x95, 0xff}, "DmsUiSW65STrO9MYz9UZEiHoA9W"},
    {[]byte{0xdc, 0xf0, 0xab, 0xfa}, "RL5"},
    {[]byte{0x33, 0xea, 0xc7, 0xd8}, "67"},
    {[]byte{0x8, 0x5d, 0x88, 0x1a}, "NJz02SkBRGkGn9d0nztLLSL9g8YW3p4d7xAfgB"},
    {[]byte{0xaa, 0xbe, 0x2, 0xcb}, "C"},
    {[]byte{0xba, 0xdd, 0xeb, 0xe8}, "6V0CbxRjAuvTgONMsMM4f"},
    {[]byte{0x7d, 0xe6, 0x5c, 0xf8}, "UpAx2XrCe23Dupo4aePyuUFyIJMQTg"},
    {[]byte{0xa6, 0xed, 0x36, 0x3d}, "yDVs38VovVv7qUbzOSzvSbbIwdeW4er"},
    {[]byte{0x65, 0x72, 0xe6, 0x3c}, "qfFZyU"},
    {[]byte{0x2b, 0xd4, 0x56, 0x51}, "5Em2SulDbzArw6j"},
    {[]byte{0x3b, 0x8e, 0x14, 0x3e}, "kkscBEhp"},
    {[]byte{0x49, 0x2f, 0xc4, 0x2}, "Zg5yRgd6dsnz02zPeSi6a4PjaRzD8Qdgo"},
    {[]byte{0x4c, 0xd5, 0x61, 0x37}, "uSHMkV6Fvhcaald2j2RdYU96ctq"},
    {[]byte{0xa2, 0x57, 0x4d, 0xb}, "7BakZCTxLR"},
    {[]byte{0xb6, 0xbd, 0x8e, 0x51}, "URvZqDQRaPZMy3Fpi5nz"},
    {[]byte{0x34, 0x55, 0x69, 0x3c}, "M5qKA3vUAmOJ17wIeBa0c4U6iwuAaxRF8L"},
    {[]byte{0x93, 0x6d, 0x8f, 0x36}, "vvGwVWK2QDZRePcPhbEAZeNm6AB3oP0TCb"},
    {[]byte{0x54, 0x95, 0xf1, 0x16}, "IJZpZ2tJ4SaqNEz25oV6ceBxSCX4lqF8ElmwXfw"},
    {[]byte{0x98, 0x47, 0x30, 0x79}, "7LdX6NWCjkTeQbTYS4S2rzMrbNFPleXbGWeSQt"},
    {[]byte{0x4b, 0x75, 0x27, 0x17}, "9rbidzBNqzuqazhmkQENPnWrhJrxHiUQP"},
    {[]byte{0x5f, 0x9d, 0xfd, 0xe8}, "qjXQN28P42FmdNaHl6iQLFcKT"},
    {[]byte{0x46, 0x7a, 0xe8, 0x8e}, "Q47POeCdVhRZjTX0"},
    {[]byte{0xa4, 0x12, 0x6b, 0xd}, "UbyJd5VDvCaoKBJzdz7yE824h1dsAT4MpdZ"},
    {[]byte{0xd9, 0xc5, 0xbe, 0xf7}, "qE4P"},
    {[]byte{0x47, 0x23, 0xcc, 0xec}, "t2062iRqkiOEc65V7GMtIbAHt"},
    {[]byte{0xde, 0xad, 0xbe, 0xef}, ""},
    {[]byte{0xb5, 0xda, 0xa2, 0x1a}, "7bTduA"},
}

func TestOnce(t *testing.T) {
	h := HashLittle()
	h.Write([]byte("foo"))
	h1 := h.Sum32()
	if h1 != 0xe18f6896 {
		t.Errorf("hashlittle(%q) = 0x%x want 0x%x", "foo", h1, 0xe18f6896)
	}
	done, error := h.Write([]byte("bar"))
	if error == nil || done != 0 {
		t.Errorf("expect an error when hashlittle is written to twice")
	}
	h.Reset()
	done, error = h.Write([]byte("bar"))
	if error != nil {
		t.Fatalf("write error: %s", error)
	}
	if done != 3 {
		t.Fatalf("wrote only %d out of 3 bytes", done)
	}
	h2 := h.Sum32()
	if h2 != 0xa6c2ceb8 {
		t.Errorf("hashlittle(%q) = 0x%x want 0x%x", "foo", h2, 0xa6c2ceb8)
	}
}

func TestGolden32(t *testing.T) {
	testGolden(t, HashLittle(), golden32)
}

func testGolden(t *testing.T, hash hash.Hash, gold []golden) {
	for _, g := range gold {
		hash.Reset()
		done, error := hash.Write([]byte(g.text))
		if error != nil {
			t.Fatalf("write error: %s", error)
		}
		if done != len(g.text) {
			t.Fatalf("wrote only %d out of %d bytes", done, len(g.text))
		}
		if actual := hash.Sum(nil); !bytes.Equal(g.sum, actual) {
			t.Errorf("hashlittle(%q) = 0x%x want 0x%x", g.text, actual, g.sum)
		}
	}
}

func BenchmarkHashLittleKB(b *testing.B) {
	benchmarkKB(b, HashLittle())
}

func benchmarkKB(b *testing.B, h hash.Hash) {
	b.SetBytes(1024)
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i)
	}
	in := make([]byte, 0, h.Size())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Reset()
		h.Write(data)
		h.Sum(in)
	}
}
