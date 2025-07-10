package pow_test

import (
	"testing"

	"github.com/le0developer/go-pow"
	"github.com/le0developer/go-pow/sha2bday"
)

func TestSha2BDay(t *testing.T) {
	nonce := []byte{1, 2, 3, 4, 5}
	data := []byte{2, 2, 3, 4, 5}
	r := pow.NewRequest(5, nonce, sha2bday.Sha2BDay)
	proof, err := pow.Fulfil(r, data)
	if err != nil {
		t.Fatalf("Fulfil: %v", err)
	}
	ok, err := pow.Check(r, proof, data)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !ok {
		t.Fatalf("Proof of work should be ok")
	}
	ok, err = pow.Check(r, proof, nonce)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if ok {
		t.Fatalf("Proof of work should not be ok")
	}
}

func BenchmarkCheck5(b *testing.B)  { benchmarkCheck(5, b) }
func BenchmarkCheck10(b *testing.B) { benchmarkCheck(10, b) }
func BenchmarkCheck15(b *testing.B) { benchmarkCheck(15, b) }
func BenchmarkCheck20(b *testing.B) { benchmarkCheck(20, b) }

func benchmarkCheck(diff uint32, b *testing.B) {
	req := pow.NewRequest(diff, []byte{1, 2, 3, 4, 5}, sha2bday.Sha2BDay)
	prf, _ := pow.Fulfil(req, []byte{6, 7, 8, 9})
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		pow.Check(req, prf, []byte{6, 7, 8, 9})
	}
}

func BenchmarkFulfil5(b *testing.B)  { benchmarkFulfil(5, b) }
func BenchmarkFulfil10(b *testing.B) { benchmarkFulfil(10, b) }
func BenchmarkFulfil15(b *testing.B) { benchmarkFulfil(15, b) }
func BenchmarkFulfil20(b *testing.B) { benchmarkFulfil(20, b) }

func benchmarkFulfil(diff uint32, b *testing.B) {
	req := pow.NewRequest(diff, []byte{1, 2, 3, 4, 5}, sha2bday.Sha2BDay)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		pow.Fulfil(req, []byte{6, 7, 8, 9})
	}
}
