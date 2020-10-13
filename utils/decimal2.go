package utils

import (
	"encoding/json"
)

type Decimal2 int

func (d Decimal2) MarshalJSON() ([]byte, error) {
	value := float32(d) / 100.0
	return json.Marshal(value)
}

func (d *Decimal2) UnmarshalJSON(data []byte) error {
	var f float32
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	result := int(f * 100)
	*d = Decimal2(result)
	return nil
}
