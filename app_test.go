// Copyright (C) 2019 MizukiSonoko. All rights reserved.

package main

import (
	"testing"
)

func BenchmarkGetIssueId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = getIssueId("github.com/MizukiSonoko/habakari-honoka/pull/13579")
	}
}

func BenchmarkToMP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = toMP(1234567890)
	}
}