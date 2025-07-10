// Create and fulfill proof of work requests.
package pow

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
