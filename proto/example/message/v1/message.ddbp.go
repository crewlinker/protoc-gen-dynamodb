package messagev1

import (
	"strconv"
)

type ListBasic struct{ V string }

func (p ListBasic) At(i int) string {
	return p.V + "[" + strconv.Itoa(i) + "]"
}

type ListP[T interface{ Set(v string) T }] struct{ V string }

func (p ListP[T]) At(i int) T {
	var v T
	return v.Set(p.V + "[" + strconv.Itoa(i) + "]")
}

// EngineP provides path construction for the Engine message
type EngineP struct{ V string }

func (p EngineP) Set(v string) EngineP {
	p.V = v
	return p
}

func (p EngineP) Brand() string {
	return p.V + ".5"
}

// KitchenP provides path construction for the Kitchen message
type KitchenP struct{ V string }

func (p KitchenP) Brand() string {
	return p.V + ".1"
}

func (p KitchenP) AnotherKitchen() KitchenP {
	return KitchenP{V: p.V + ".2"}
}

func (p KitchenP) ApplianceBrands() ListBasic {
	return ListBasic{V: p.V + ".4"}
}

func (p KitchenP) ApplianceEngines() ListP[EngineP] {
	return ListP[EngineP]{V: p.V + ".3"}
}
