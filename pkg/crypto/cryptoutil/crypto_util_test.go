package cryptoutil

import (
	"os"
	"testing"
)

var ciphkey string

func TestMain(m *testing.M) {
	ciphkey = "65dTxbqk7rE3IFly1hnI1234"
	os.Exit(m.Run())
}

func TestEncrypt(t *testing.T) {
	want := "Q0jhx4gItMsD"
	got := Encrypt([]byte(ciphkey), "radiation")
	if got != want {
		t.Errorf("Encrypt got %v, want %v", got, want)
	}
}

func TestDecrypt(t *testing.T) {
	want := "radiation"
	got := Decrypt([]byte(ciphkey), "Q0jhx4gItMsD")
	if got != want {
		t.Errorf("Decrypt got %v, want %v", got, want)
	}
}

func TestSumHash(t *testing.T) {
	hash := SumHash("12345678")
	if len(hash) == 0 {
		t.Errorf("SumHash got an empty hash")
	}
}
