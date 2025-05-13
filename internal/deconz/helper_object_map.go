package deconz

import "math"

type MapObject interface {
	Has(key string) bool
	ValueToBool(key string) bool
	ValueToInt(key string) int
	ValueToPercent(key string) int
	ValueToString(key string) string
}
type ObjectMap map[string]interface{}

func (obj ObjectMap) Has(key string) bool {
	return obj[key] != nil
}

func (obj ObjectMap) ValueToBool(key string) bool {
	return obj[key].(bool)
}

func (obj ObjectMap) ValueToInt(key string) int {
	return int(obj[key].(float64))
}

func (obj ObjectMap) ValueToString(key string) string {
	return obj[key].(string)
}

func (obj ObjectMap) ValueToPercent(key string) int {
	value := obj[key].(float64)
	return int(math.Round(value * 100.0 / 255.0))
}

type ExtendedObjectMap map[string]*struct {
	LastUpdated string      `json:"lastupdated"`
	Value       interface{} `json:"value"`
}

func (obj ExtendedObjectMap) Has(key string) bool {
	return obj[key] != nil
}

func (obj ExtendedObjectMap) ValueToBool(key string) bool {
	return obj[key].Value.(bool)
}

func (obj ExtendedObjectMap) ValueToInt(key string) int {
	return int(obj[key].Value.(float64))
}

func (obj ExtendedObjectMap) ValueToString(key string) string {
	return obj[key].Value.(string)
}

func (obj ExtendedObjectMap) ValueToPercent(key string) int {
	value := obj[key].Value.(float64)
	return int(math.Round(value * 100.0 / 255.0))
}
