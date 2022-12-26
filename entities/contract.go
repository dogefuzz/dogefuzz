package entities

type Contract struct {
	Id            string
	Address       string
	Source        string
	CompiledCode  string
	AbiDefinition string
	Name          string
	Args          string
	ConstructorId string
	CFG           string
	DistanceMap   string
}
