// Copyright (c) 2020 Reto Brunner. All rights reserved.
// Copyright (c) 2015 Andy Balholm. All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package pgpmail

import "golang.org/x/text/transform"

// CRLFTransformer takes LF line endings in src, and converts them to CRLF in dst.
type CRLFTransformer struct {
	transform.NopResetter
}

func (CRLFTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for nDst < len(dst) && nSrc < len(src) {
		if c := src[nSrc]; c == '\n' {
			if nDst+1 == len(dst) {
				break
			}
			dst[nDst] = '\r'
			dst[nDst+1] = '\n'
			nSrc++
			nDst += 2
		} else {
			dst[nDst] = c
			nSrc++
			nDst++
		}
	}
	if nSrc < len(src) {
		err = transform.ErrShortDst
	}
	return
}

// LFTransformer takes CRLF, CR, or LF line endings in src, and converts them
// to LF in dst.
type LFTransformer struct {
	prev byte
}

func (lt *LFTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for nDst < len(dst) && nSrc < len(src) {
		c := src[nSrc]
		switch c {
		case '\r':
			dst[nDst] = '\n'
		case '\n':
			if lt.prev == '\r' {
				nSrc++
				lt.prev = c
				continue
			}
			dst[nDst] = '\n'
		default:
			dst[nDst] = c
		}
		lt.prev = c
		nDst++
		nSrc++
	}
	if nSrc < len(src) {
		err = transform.ErrShortDst
	}
	return
}

func (lt *LFTransformer) Reset() {
	lt.prev = 0
}
