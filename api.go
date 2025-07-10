// Create and fulfill proof of work requests.
package pow

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

// Represents a proof-of-work request.
type Request struct {
	// The requested algorithm
	Alg AlgorithmName

	// The requested difficulty
	Difficulty uint32

	// Nonce to diversify the request
	Nonce []byte
}

// Represents a completed proof-of-work
type Proof []byte

// Convenience function to create a new sha2bday proof-of-work request
// as a string
func NewRequest(difficulty uint32, nonce []byte, name AlgorithmName) string {
	req := Request{
		Difficulty: difficulty,
		Nonce:      nonce,
		Alg:        name,
	}
	s, _ := req.MarshalText()
	return string(s)
}

func (proof Proof) MarshalText() ([]byte, error) {
	return []byte(base64.RawStdEncoding.EncodeToString(proof)), nil
}

func (proof *Proof) UnmarshalText(buf []byte) error {
	var err error
	buf, err = base64.RawStdEncoding.DecodeString(string(buf))
	if err != nil {
		return err
	}
	*proof = buf
	return nil
}

func (req Request) String() string {
	return fmt.Sprintf("%s-%d-%s",
		req.Alg,
		req.Difficulty,
		string(base64.RawStdEncoding.EncodeToString(req.Nonce)))
}

func (req Request) MarshalText() ([]byte, error) {
	return []byte(req.String()), nil
}

func (req *Request) UnmarshalText(buf []byte) error {
	bits := strings.SplitN(string(buf), "-", 3)
	if len(bits) != 3 {
		return fmt.Errorf("There should be two dashes in a PoW request")
	}
	alg := AlgorithmName(bits[0])
	req.Alg = alg
	diff, err := strconv.Atoi(bits[1])
	if err != nil {
		return err
	}
	req.Difficulty = uint32(diff)
	req.Nonce, err = base64.RawStdEncoding.DecodeString(bits[2])
	if err != nil {
		return err
	}
	return nil
}

// Convenience function to check whether a proof of work is fulfilled
func Check(request, proof string, data []byte) (bool, error) {
	var req Request
	var prf Proof
	err := req.UnmarshalText([]byte(request))
	if err != nil {
		return false, err
	}
	err = prf.UnmarshalText([]byte(proof))
	if err != nil {
		return false, err
	}
	return prf.Check(req, data), nil
}

// Fulfil the proof-of-work request.
func (req *Request) Fulfil(data []byte) (Proof, error) {
	alg, ok := registeredAlgorithms[req.Alg]
	if !ok {
		return nil, fmt.Errorf("No such algorithm: %s", req.Alg)
	}
	res, err := alg.Fulfil(*req, data)
	if err != nil {
		return nil, err
	}
	return Proof(res), nil
}

// Convenience function to fulfil the proof of work request
func Fulfil(request string, data []byte) (string, error) {
	var req Request
	err := req.UnmarshalText([]byte(request))
	if err != nil {
		return "", err
	}
	proof, err := req.Fulfil(data)
	if err != nil {
		return "", err
	}
	s, _ := proof.MarshalText()
	return string(s), nil
}

// Check whether the proof is ok
func (proof *Proof) Check(req Request, data []byte) bool {
	alg, ok := registeredAlgorithms[req.Alg]
	if !ok {
		return false
	}
	return alg.Verify(req, *proof, data)
}
