### Jenkins hashlittle for Go

This package implements the 32-bit `hashlittle()` hash function from Bob Jenkins' [lookup3.c](http://burtleburtle.net/bob/c/lookup3.c) library. See also the [Wikipedia article](https://en.wikipedia.org/wiki/Jenkins_hash_function).

It is written for portability and compatibility, not speed, and implements the Go standard library [hash.Hash32](https://golang.org/pkg/hash/#Hash32) interface. Hashlittle is meant to hash short keys, however, and does not work incrementally as the hash itself uses the key length. Hashing `"foo"` then `"bar"` would produce different results than hashing `"foobar"`, for that reason attempting to call Write() more than once on an object returned by `lookup3.HashLittle()` without resetting in-between will raise an error.

This code is hereby placed in the public domain.
