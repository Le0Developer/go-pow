package pow

type AlgorithmName string

type Algorithm interface {
	Name() AlgorithmName
	Fulfil(req Request, data []byte) ([]byte, error)
	Verify(req Request, proof []byte, data []byte) bool
}

var registeredAlgorithms = make(map[AlgorithmName]Algorithm)

// This function is not thread-safe, so it should be called only once
// during the initialization of the application.
// Preferably call it in the init() function of the package
// that implements the Algorithm interface.
// See the [sha2bday](./sha2bday) package for an example.
func RegisterAlgorithm(alg Algorithm) {
	if alg == nil {
		return
	}
	name := alg.Name()
	if _, exists := registeredAlgorithms[name]; exists {
		return
	}
	registeredAlgorithms[name] = alg
}
