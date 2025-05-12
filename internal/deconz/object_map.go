package deconz

import "math"

type StateObject interface {
	Has(key string) bool
	ValueToBool(key string) bool
	ValueToInt(key string) int
	ValueToPercent(key string) int
	ValueToString(key string) string
}
type SimpleStateMap map[string]interface{}

func (o SimpleStateMap) Has(key string) bool {
	return o[key] != nil
}

func (o SimpleStateMap) ValueToBool(key string) bool {
	return o[key].(bool)
}

func (o SimpleStateMap) ValueToInt(key string) int {
	return int(o[key].(float64))
}

func (o SimpleStateMap) ValueToString(key string) string {
	return o[key].(string)
}

func (o SimpleStateMap) ValueToPercent(key string) int {
	value := o[key].(float64)
	return int(math.Round(value * 100.0 / 255.0))
}

type ExtendedStateMap map[string]*struct {
	LastUpdated string      `json:"lastupdated"`
	Value       interface{} `json:"value"`
}

func (o ExtendedStateMap) Has(key string) bool {
	return o[key] != nil
}

func (o ExtendedStateMap) ValueToBool(key string) bool {
	return o[key].Value.(bool)
}

func (o ExtendedStateMap) ValueToInt(key string) int {
	return int(o[key].Value.(float64))
}

func (o ExtendedStateMap) ValueToString(key string) string {
	return o[key].Value.(string)
}

func (o ExtendedStateMap) ValueToPercent(key string) int {
	value := o[key].Value.(float64)
	return int(math.Round(value * 100.0 / 255.0))
}
