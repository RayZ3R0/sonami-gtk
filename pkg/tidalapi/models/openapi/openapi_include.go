package openapi

import (
	"encoding/json"
	"slices"
)

type ObjectType string

type IncludedObjects []IncludedObject

type IncludedObject struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Raw  json.RawMessage
}

func (i *IncludedObject) UnmarshalJSON(b []byte) error {
	i.Raw = append([]byte(nil), b...)
	type Alias IncludedObject
	return json.Unmarshal(b, (*Alias)(i))
}

func (i IncludedObjects) FromRelationships(relationships []Relationship, t string) IncludedObjects {
	var ids []string
	for _, rel := range relationships {
		ids = append(ids, rel.ID)
	}

	var seenIds = []string{}
	var included IncludedObjects
	for _, obj := range i {
		if slices.Contains(ids, obj.ID) && (obj.Type == t || t == "") && !slices.Contains(seenIds, obj.ID) {
			included = append(included, obj)
			seenIds = append(seenIds, obj.ID)
		}
	}
	return included
}

func (i IncludedObjects) FromType(t string) IncludedObjects {
	var included IncludedObjects
	for _, obj := range i {
		if obj.Type == t {
			included = append(included, obj)
		}
	}
	return included
}
