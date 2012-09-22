package mdiff

import (
	"reflect"
)

type Difference struct {
	Updated map[string]interface{} `json:"updated"`
	Removed []string               `json:"removed"`
}

func Mapify(obj interface{}) (m map[string]interface{}) {
	v := reflect.ValueOf(obj)
	t := v.Type()

	m = make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		field := v.Field(i)
		m[name] = field.Interface()
	}

	return m
}

func DeMapify(dest *interface{}, m map[string]interface{}) {
	v := reflect.ValueOf(dest).Elem()
	t := v.Type()

	var otherField string

	other := make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("mdiff")

		switch {
		case tag == "other":
			otherField = field.Name
		}
	}

	for key, value := range m {
		field := v.FieldByName(key)
		if !field.IsValid() {
			other[key] = value
			continue
		}

		field.Set(reflect.ValueOf(value))
	}

	if otherField != "" {
		v.FieldByName(otherField).Set(reflect.ValueOf(other))
	}
}

func Copy(oldMap map[string]interface{}) (newMap map[string]interface{}) {
	newMap = make(map[string]interface{})

	for name, value := range oldMap {
		newMap[name] = value
	}

	return newMap
}

func Diff(oldMap map[string]interface{}, newMap map[string]interface{}) (diff *Difference) {
	diff = &Difference{
		Updated: make(map[string]interface{}),
	}

	for key, newValue := range newMap {
		oldValue, ok := oldMap[key]

		if !ok || newValue != oldValue {
			diff.Updated[key] = newValue
		}
	}

	for key, _ := range oldMap {
		_, ok := newMap[key]

		if !ok {
			diff.Removed = append(diff.Removed, key)
		}
	}

	return diff
}

func DiffObj(oldObj interface{}, newObj interface{}) (diff *Difference) {
	return Diff(Mapify(oldObj), Mapify(newObj))
}

func Apply(m map[string]interface{}, diff *Difference) {
	for key, value := range diff.Updated {
		m[key] = value
	}

	for _, key := range diff.Removed {
		delete(m, key)
	}
}

func ApplyCopy(oldMap map[string]interface{}, diff *Difference) (newMap map[string]interface{}) {
	newMap = Copy(oldMap)
	Apply(newMap, diff)
	return newMap
}

func ApplyObj(dest *interface{}, diff *Difference) {
	m := Mapify(*dest)
	Apply(m, diff)
	DeMapify(dest, m)
}
