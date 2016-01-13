// Copyright (c) 2015 Apsalar Inc. All rights reserved.

/*
Package lookup3 implements Bob Jenkins' lookup3 hashlittle
non-cryptographic hash function.

See:

https://en.wikipedia.org/wiki/Jenkins_hash_function

http://burtleburtle.net/bob/c/lookup3.c

While this package uses the Go standard library hash package interface
hashlittle calls cannot be chained, i.e. all the bytes of the key need to be
written to the hash object in a single call to Write(), as the hash is
dependent on the key length. Making multiple calls to Write() would not
yield the same result, and for this reason throws an error

This code is hereby placed in the public domain
*/
package lookup3

import (
	"hash"
)

type (
	sum32  uint32
)

// Hashlittle returns a new hash.Hash32
// to calculate the 32-bit lookup3 hashlittle
//
// Its Sum() method will lay the value out in big-endian byte order.
// and Sum32() will return the 32-bit unsigned value
func HashLittle() hash.Hash32 {
	var s sum32 = 0
	return &s
}
func (s *sum32) Reset()  { *s = 0 }

func (s *sum32) Sum32() uint32  { return uint32(*s) }

type Error string
func (e Error) Error() string {
	return string(e)
}

// this follows the third (slower) code path in lookup3.c
// we care more about correctness than speed and hackiness
// http://commandcenter.blogspot.de/2012/04/byte-order-fallacy.html
func (s *sum32) Write(data []byte) (int, error) {
	var a, b, c uint32

	if *s != 0 {
		return 0, Error("lookup3.HashLittle can only be written to once")
	}
	k := data[:]
	a = 0xdeadbeef + uint32(len(k)) + uint32(*s)
	b = 0xdeadbeef + uint32(len(k)) + uint32(*s)
	c = 0xdeadbeef + uint32(len(k)) + uint32(*s)

    /*--------------- all but the last block: affect some 32 bits of (a,b,c) */
    for len(k) > 12 {
		a += uint32(k[0])
		a += uint32(k[1]) << 8
		a += uint32(k[2]) << 16
		a += uint32(k[3]) << 24
		b += uint32(k[4])
		b += uint32(k[5]) << 8
		b += uint32(k[6]) << 16
		b += uint32(k[7]) << 24
		c += uint32(k[8])
		c += uint32(k[9]) << 8
		c += uint32(k[10]) << 16
		c += uint32(k[11]) << 24
		// mix(a,b,c)
		a -= c
		// a ^= rot(c, 4)
		// #define rot(x,k) (((x) << (k)) | ((x) >> (32-(k))))
		a ^= ((c << 4) | (c >> 28))
		c += b;
		b -= a
		// b ^= rot(a, 6)
		b ^= ((a << 6) | (a >> 26))
		a += c;
		c -= b
		// c ^= rot(b, 8)
		c ^= ((b << 8) | (b >> 24))
		b += a;
		a -= c
		// a ^= rot(c,16)
		a ^= ((c << 16) | (c >> 16))
		c += b;
		b -= a
		// b ^= rot(a,19)
		b ^= ((a << 19) | (a >> 13))
		a += c;
		c -= b
		// c ^= rot(b, 4)
		c ^= ((b << 4) | (b >> 28))
		b += a;

		k = k[12:]
    }

    /*-------------------------------- last block: affect all 32 bits of (c) */
	switch(len(k)) {                 /* all the case statements fall through */
	case 12:
		c += uint32(k[11])<<24
		fallthrough
    case 11:
		c += uint32(k[10])<<16
		fallthrough
    case 10:
		c += uint32(k[9])<<8
		fallthrough
    case 9:
		c += uint32(k[8])
		fallthrough
    case 8:
		b += uint32(k[7])<<24
		fallthrough
    case 7:
		b += uint32(k[6])<<16
		fallthrough
    case 6:
		b += uint32(k[5])<<8
		fallthrough
    case 5:
		b += uint32(k[4])
		fallthrough
    case 4:
		a += uint32(k[3])<<24
		fallthrough
    case 3:
		a += uint32(k[2])<<16
		fallthrough
    case 2:
		a += uint32(k[1])<<8
		fallthrough
    case 1:
		a += uint32(k[0])
    case 0:
		*s = sum32(c);
		return len(data), nil
	}
    
	//final(a,b,c);
	c ^= b
	// c -= rot(b,14)
	c -= ((b << 14) | (b >> 18))
	a ^= c
	// a -= rot(c,11)
	a -= ((c << 11) | (c >> 21))
	b ^= a
	// b -= rot(a,25)
	b -= ((a << 25) | (a >> 7))
	c ^= b
	// c -= rot(b,16)
	c -= ((b << 16) | (b >> 16))
	a ^= c
	// a -= rot(c,4)
	a -= ((c << 4) | (c >> 28))
	b ^= a
	// b -= rot(a,14)
	b -= ((a << 14) | (a >> 18))
	c ^= b
	// c -= rot(b,24)
	c -= ((b << 24) | (b >> 8))
		
	*s = sum32(c);
	return len(data), nil
}

func (s *sum32) Size() int  { return 4 }

func (s *sum32) BlockSize() int  { return 1 }

func (s *sum32) Sum(in []byte) []byte {
	v := uint32(*s)
	return append(in, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}
