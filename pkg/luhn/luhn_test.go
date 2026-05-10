package luhn

import "testing"

func TestValid(t *testing.T) {
	if !Valid("79927398713") {
		t.Fatal("expected valid luhn number")
	}
	if Valid("79927398710") {
		t.Fatal("expected invalid luhn number")
	}
	if Valid("abc") {
		t.Fatal("expected invalid for non-digit")
	}
}
