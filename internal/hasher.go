// Package internal
// Author: hyphen
// Copyright 2024 hyphen. All rights reserved.
// Create-time: 2024/3/18
package internal

import (
	"hash/maphash"
	"reflect"
	"unsafe"
)

type Hasher func(unsafe.Pointer, uintptr) uintptr

func GetRuntimeHasher[K comparable]() (h Hasher, seed uintptr) {
	switch reflect.TypeFor[K]().Kind() {
	case reflect.String:
		h = stringHasher(maphash.MakeSeed())
		seed = 0
		return
	default:
		a := any(make(map[K]struct{}))
		i := (*mapiface)(unsafe.Pointer(&a))
		h, seed = i.typ.hasher, uintptr(i.val.hash0)
		return
	}
}

func stringHasher(s maphash.Seed) Hasher {
	return func(pointer unsafe.Pointer, u uintptr) uintptr {
		ss := *(*string)(pointer)
		return uintptr(maphash.String(s, ss))
	}
}

//func stringStructOf(sp *string) *stringStruct {
//	return (*stringStruct)(unsafe.Pointer(sp))
//}

type stringStruct struct {
	str unsafe.Pointer
	len int
}

func Tophash(hash uintptr) uint8 {
	top := hash >> (PtrSize*8 - 24)
	return uint8(top)
}
func Lowhash(hash uintptr, B uint8) uintptr {
	low := hash & (1<<B - 1)
	return low
}

const PtrSize = 4 << (^uintptr(0) >> 63)

func OverLoadFactor(count int, B uint8) bool {
	return count > bucketCnt && uintptr(count) > loadFactorNum*(bucketShift(B)/loadFactorDen)
}

func bucketShift(b uint8) uintptr {
	// Masking the shift amount allows overflow checks to be elided.
	return uintptr(1) << (b & (PtrSize*8 - 1))
}

type mapiface struct {
	typ *maptype
	val *hmap
}

// this is a mirror of go1.22 src/internal/abi/type.go/MapType.
// no change, keep sync with go src
type maptype struct {
	typ    _type
	key    *_type
	elem   *_type
	bucket *_type
	// function for hashing keys (ptr to key, seed) -> hash
	hasher     func(unsafe.Pointer, uintptr) uintptr
	keysize    uint8
	elemsize   uint8
	bucketsize uint16
	flags      uint32
}

type _type struct {
	size_       uintptr
	ptrbytes    uintptr // number of (prefix) bytes in the type that can contain pointers
	hash        uint32  // hash of type; avoids computation in hash tables
	tflag       tflag   // extra type information flags
	align_      uint8   // alignment of variable with this type
	fieldalign_ uint8   // alignment of struct field with this type
	kind_       uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
	// GCData stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, GCData is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	gcdata    *byte
	str       nameoff // string form
	ptrtohhis typeoff // type for pointer to this type, may be zero
}

type tflag uint8

type nameoff int32

type typeoff int32

type hmap struct {
	// Note: the format of the hmap is also encoded in cmd/compile/internal/reflectdata/reflect.go.
	// Make sure this stays in sync with the compiler's definition.
	count     int // # live cells == size of map.  Must be first (used by len() builtin)
	flags     uint8
	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	hash0     uint32 // hash seed

	buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

	extra *mapextra // optional fields
}

type mapextra struct {
	// If both key and elem do not contain pointers and are inline, then we mark bucket
	// type as containing no pointers. This avoids scanning such maps.
	// However, bmap.overflow is a pointer. In order to keep overflow buckets
	// alive, we store pointers to all overflow buckets in hmap.extra.overflow and hmap.extra.oldoverflow.
	// overflow and oldoverflow are only used if key and elem do not contain pointers.
	// overflow contains overflow buckets for hmap.buckets.
	// oldoverflow contains overflow buckets for hmap.oldbuckets.
	// The indirection allows to store a pointer to the slice in hiter.
	overflow    *[]*bmap
	oldoverflow *[]*bmap

	// nextOverflow holds a pointer to a free overflow bucket.
	nextOverflow *bmap
}

type bmap struct {
	// tophash generally contains the top byte of the hash value
	// for each key in this bucket. If tophash[0] < minTopHash,
	// tophash[0] is a bucket evacuation state instead.
	tophash [bucketCnt]uint8
	// Followed by bucketCnt keys and then bucketCnt elems.
	// NOTE: packing all the keys together and then all the elems together makes the
	// code a bit more complicated than alternating key/elem/key/elem/... but it allows
	// us to eliminate padding which would be needed for, e.g., map[int64]int8.
	// Followed by an overflow pointer.
}

const (
	bucketCnt          = mapbucketcount
	mapbucketcountbits = 3 // log2 of number of elements in a bucket.
	mapbucketcount     = 1 << mapbucketcountbits
	loadFactorNum      = loadFactorDen * bucketCnt * 13 / 16
	loadFactorDen      = 2
)
