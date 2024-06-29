package vector_io

import (
	"eigen_db/cfg"
	t "eigen_db/types"
	"fmt"
)

type Vector struct {
	Id         t.VectorId         `json:"id"`
	Components t.VectorComponents `json:"components"`
}

func (v *Vector) Insert() {
	vectorStoreInstance.writeVector(v)
}

type VectorFactory struct{}

func (factory *VectorFactory) NewVector(components t.VectorComponents) (IVector, error) {
	dimensions := cfg.GetConfig().GetHNSWParamsDimensions()
	if len(components) == int(dimensions) {
		v := &Vector{}
		v.Components = components
		return v, nil
	}
	return nil, fmt.Errorf("provided a %d-dimensional vector while the vector space is %d-dimensional", len(components), dimensions)
}
