/*
 * paperback: resilient paper backups for the very paranoid
 * Copyright (C) 2018 Aleksa Sarai <cyphar@cyphar.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package shamir

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"testing"
	"time"
)

// rng is the global random number generator used for all non-important RNG
// operations in our tests.
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// mustRandomBytes returns a slice of random bytes of the given size.
func mustRandomBytes(size uint) []byte {
	bytes := make([]byte, size)
	if _, err := io.ReadFull(rng, bytes); err != nil {
		panic(err)
	}
	return bytes
}

// shuffleShares takes a slice of shares and in-place scrambles it.
func shuffleShares(shares []Share) {
	for i := 0; i < len(shares); i++ {
		j := rng.Intn(i + 1)
		shares[i], shares[j] = shares[j], shares[i]
	}
}

// extendBytes takes a series of byte slices and combines them into a single
// slice, by appending one after the other (using minimal allocations).
func extendBytes(slices ...[]byte) []byte {
	var totalLen int
	for _, slice := range slices {
		totalLen += len(slice)
	}
	result := make([]byte, totalLen)
	var start int
	for _, slice := range slices {
		start += copy(result[start:], slice)
	}
	return result
}

// secretVectors is the set of vectors which are used to test a variety of
// edge-cases.
var secretVectors = [][]byte{
	// Some very simple test vectors.
	[]byte("Hello, world!"),
	[]byte("A slightly longer test string, which spans multiple parts."),
	[]byte("The quick brown fox jumps over the lazy dog."),
	[]byte("Some numeric values: π=3.14156926 e=2.71828. Punctuation!?@#%69*&#!@(&%(!@)#)!)%*(!@#{}{}:|:,"),
	// Make sure that leading zero are not dropped.
	extendBytes([]byte{0x00}, mustRandomBytes(DefaultBlockSize)),
	// Make sure that zeros on a block boundary are not dropped.
	extendBytes(mustRandomBytes(DefaultBlockSize), []byte{0x00}, mustRandomBytes(DefaultBlockSize)),
	// Make sure that zeros in the final block are not dropped.
	extendBytes(mustRandomBytes(DefaultBlockSize), []byte{0x00, 0x01}),
	// Random vectors. These are used to ensure that we can handle arbitrary binary data.
	mustRandomBytes(DefaultBlockSize / 2),
	mustRandomBytes(DefaultBlockSize - 1),
	mustRandomBytes(DefaultBlockSize*2 + 1),
	mustRandomBytes(DefaultBlockSize*16 - 2),
	mustRandomBytes(DefaultBlockSize*34 + 23),
	mustRandomBytes(DefaultBlockSize*74 - 19),
	mustRandomBytes(DefaultBlockSize*143 - 18),
}

// testSchemeHelper is a helper which just provides the ability to run test
// functions on the set of secretVectors, and with multiple schemes.
func testSchemeHelper(t *testing.T, fn func(t *testing.T, k, n uint, secret []byte)) {
	// Test only a few (k,n) combinations.
	schemes := []struct{ k, n uint }{
		{2, 2},
		{2, 3},
		{2, 4},
		{4, 4},
		{4, 7},
		{7, 7},
		{7, 8},
		{7, 12},
		{8, 9},
		{14, 19},
	}
	for _, scheme := range schemes {
		k, n := scheme.k, scheme.n
		tn := fmt.Sprintf("Scheme:k=%d,n=%d", k, n)
		t.Run(tn, func(t *testing.T) {
			for _, secret := range secretVectors {
				fn(t, k, n, secret)
			}
		})
	}
}

func copyShare(s Share) Share {
	payload, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	var newShare Share
	if err := json.Unmarshal(payload, &newShare); err != nil {
		panic(err)
	}
	return newShare
}

func copyShares(shares []Share) []Share {
	var newShares []Share
	for _, share := range shares {
		newShares = append(newShares, copyShare(share))
	}
	return newShares
}
