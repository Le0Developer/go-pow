package pow

type AlgorithmName string

type Algorithm interface {
	Name() AlgorithmName
	Fulfil(req Request, data []byte) ([]byte, error)
	Verify(req Request, proof []byte, data []byte) bool
}

var registeredAlgorithms = make(map[AlgorithmName]Algorithm)

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
