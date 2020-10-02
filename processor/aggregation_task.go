package processor

import "math"

type AggregatedConfig struct {
	DomainName         string   `json:"domain_name"`
	AggregatedID       string   `json:"aggregated_id"`
	GroupByKeyID       []string `json:"group_by_key_id"`
	AggregatedFunction []struct {
		PropertyName string `json:"property_name"`
		FieldName    string `json:"field_name"`
		Function     string `json:"function"`
	} `json:"aggregated_function"`
}

type AggregatedTask struct {
	ID           string
	GroupByKeyID []string
	Functions    []AggregatedFunction
}

type AggregatedFunction struct {
	PropertyName string
	FieldName    string
	Function     Function
}

type Function interface {
	Apply(old interface{}, new interface{}) interface{}
}

type SumFunction struct {
}

func (f *SumFunction) Apply(old interface{}, new interface{}) interface{} {
	if old == nil {
		old = 0.0
	}

	oldNumber, oldOk := old.(float64)
	newNumber, newOk := new.(float64)

	if oldOk && newOk {
		return oldNumber + newNumber
	}

	return old
}

type MinFunction struct {
}

func (f *MinFunction) Apply(old interface{}, new interface{}) interface{} {
	if old == nil {
		old = 0.0
	}

	oldNumber, oldOk := old.(float64)
	newNumber, newOk := new.(float64)

	if oldOk && newOk {
		return math.Min(oldNumber, newNumber)
	}

	return old
}

type MaxFunction struct {
}

func (f *MaxFunction) Apply(old interface{}, new interface{}) interface{} {
	if old == nil {
		old = 0.0
	}

	oldNumber, oldOk := old.(float64)
	newNumber, newOk := new.(float64)

	if oldOk && newOk {
		return math.Max(oldNumber, newNumber)
	}

	return old
}

type CountFunction struct {
}

func (f *CountFunction) Apply(old interface{}, new interface{}) interface{} {
	if old == nil {
		old = 0.0
	}

	oldNumber, oldOk := old.(float64)

	if oldOk && new != nil {
		return oldNumber + 1
	}

	return old
}

type LastFunction struct {
}

func (f *LastFunction) Apply(old interface{}, new interface{}) interface{} {
	return new
}
