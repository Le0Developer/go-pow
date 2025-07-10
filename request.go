package pow

import (
	"encoding/base64"
	"errors"
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

// UnmarshalRequest takes a string representation of a proof-of-work request
// and returns a Request struct. It returns an error if the string is not
// formatted correctly.
func UnmarshalRequest(request string) (Request, error) {
	var req Request
	err := req.UnmarshalText([]byte(request))
	return req, err
}

// String returns a string representation of the request with the following parts separated by dashes:
//   - Algorithm name (e.g., "sha2bday")
//   - Difficulty (as a decimal number)
//   - Nonce (base64 encoded)
//
// This can be sent to the client to solve.
func (req Request) String() string {
	return string(req.Alg) + "-" + strconv.FormatUint(uint64(req.Difficulty), 10) + "-" +
		string(base64.RawStdEncoding.EncodeToString(req.Nonce))
}

// MarshalText implements the encoding.TextMarshaler interface for Request.
// It uses the String method to produce a string representation of the request.
func (req Request) MarshalText() ([]byte, error) {
	return []byte(req.String()), nil
}

var ErrInvalidRequest = errors.New("invalid request format")

// UnmarshalText implements the encoding.TextUnmarshaler interface for Request.
func (req *Request) UnmarshalText(buf []byte) error {
	bits := strings.SplitN(string(buf), "-", 3)
	if len(bits) != 3 {
		return ErrInvalidRequest
	}
	alg := AlgorithmName(bits[0])
	req.Alg = alg
	diff, err := strconv.ParseUint(bits[1], 10, 32)
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

var ErrNoSuchAlgorithm = errors.New("no such algorithm registered")

// Fulfil the proof-of-work request.
func (req *Request) Fulfil(data []byte) (Proof, error) {
	alg, ok := registeredAlgorithms[req.Alg]
	if !ok {
		return nil, ErrNoSuchAlgorithm
	}
	res, err := alg.Fulfil(*req, data)
	if err != nil {
		return nil, err
	}
	return Proof(res), nil
}
