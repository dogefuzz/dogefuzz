package mapper

import (
	"encoding/json"

	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/dto"
)

type contractMapper struct{}

func NewContractMapper() *contractMapper {
	return &contractMapper{}
}

func (m *contractMapper) MapNewDTOToEntity(d *dto.NewContractDTO) *entities.Contract {
	return &entities.Contract{
		TaskId:        d.TaskId,
		Source:        d.Source,
		CompiledCode:  d.CompiledCode,
		AbiDefinition: d.AbiDefinition,
		Name:          d.Name,
	}
}

func (m *contractMapper) MapDTOToEntity(d *dto.ContractDTO) *entities.Contract {
	cfg, _ := json.Marshal(d.CFG)
	distanceMap, _ := json.Marshal(d.DistanceMap)

	return &entities.Contract{
		Id:            d.Id,
		TaskId:        d.TaskId,
		Address:       d.Address,
		Source:        d.Source,
		CompiledCode:  d.CompiledCode,
		AbiDefinition: d.AbiDefinition,
		Name:          d.Name,
		CFG:           string(cfg),
		DistanceMap:   string(distanceMap),
	}
}

func (m *contractMapper) MapEntityToDTO(c *entities.Contract) *dto.ContractDTO {
	var cfg common.CFG
	_ = json.Unmarshal([]byte(c.CFG), &cfg)
	var distanceMap common.DistanceMap
	_ = json.Unmarshal([]byte(c.DistanceMap), &distanceMap)

	return &dto.ContractDTO{
		Id:            c.Id,
		TaskId:        c.TaskId,
		Address:       c.Address,
		Source:        c.Source,
		CompiledCode:  c.CompiledCode,
		AbiDefinition: c.AbiDefinition,
		Name:          c.Name,
		CFG:           cfg,
		DistanceMap:   distanceMap,
	}
}

func (m *contractMapper) MapDTOToCommon(c *dto.ContractDTO) *common.Contract {
	return &common.Contract{
		Address:       c.Address,
		AbiDefinition: c.AbiDefinition,
		CompiledCode:  c.CompiledCode,
		Name:          c.Name,
	}
}
