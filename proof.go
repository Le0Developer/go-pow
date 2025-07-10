package pow

import "encoding/base64"

// Represents a completed proof-of-work
type Proof []byte

// UnmarshalProof takes a string representation of a Proof
// and returns a Proof struct. It returns an error if the string is not
// formatted correctly.
func UnmarshalProof(proof string) (Proof, error) {
	var p Proof
	err := p.UnmarshalText([]byte(proof))
	return p, err
}

// MarshalText implements the encoding.TextMarshaler interface for Proof.
// It encodes the proof as a base64 string.
func (proof Proof) MarshalText() ([]byte, error) {
	return []byte(base64.RawStdEncoding.EncodeToString(proof)), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Proof.
// It decodes the proof from a base64 string.
func (proof *Proof) UnmarshalText(buf []byte) error {
	var err error
	buf, err = base64.RawStdEncoding.DecodeString(string(buf))
	if err != nil {
		return err
	}
	*proof = buf
	return nil
}

// Check whether the proof is ok
func (proof *Proof) Check(req Request, data []byte) bool {
	alg, ok := registeredAlgorithms[req.Alg]
	if !ok {
		return false
	}
	return alg.Verify(req, *proof, data)
}
