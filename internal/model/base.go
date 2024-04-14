package model

var _ ContractModel = (*customContractModel)(nil)

type (
	ContractModel interface {
		EncodeArgs() ([]byte, error)
	}

	customContractModel struct {
	}
)

func (c *customContractModel) EncodeArgs() ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func NewBaseContractModel() *SwapContractModel {
	return &SwapContractModel{
		signature: "",
		returns:   "",
	}
}
