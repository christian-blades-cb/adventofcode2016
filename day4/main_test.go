package main

import (
	"testing"
)

// TestSectorOf tests the sectorOf method
func TestSectorOf(t *testing.T) {
	input := "aabbbcccddddde123[dbcae]"
	sector, ok := sectorOf(input)
	if sector != 123 {
		t.Logf("sector %d != 123", sector)
		t.Fail()
	}
	if !ok {
		t.Log("expected checksum to pass")
		t.Fail()
	}
}
