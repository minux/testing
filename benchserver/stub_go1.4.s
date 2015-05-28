// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.4,!go1.5

#include "textflag.h"

#ifdef GOARCH_amd64
#define MOV MOVQ
#define REG AX
#endif

#ifdef GOARCH_386
#define MOV MOVL
#define REG AX
#endif

#ifdef GOARCH_amd64p32
#define MOV MOVL
#define REG AX
#endif

#ifdef GOARCH_arm
#define MOV MOVW
#define REG R0
#endif

TEXT ·init14(SB),NOSPLIT,$0
	MOV	$main·benchmarks(SB), REG
	MOV	REG, ·benchmarks(SB)
	MOV	$main·tests(SB), REG
	MOV	REG, ·tests(SB)
	RET
