package uuid

func GenerateC24() (C24, error) {
	return globalNode.Generate()
}
