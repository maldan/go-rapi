package rapi_rest

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func ApplyString(field *reflect.Value, v interface{}) {
	field.SetString(fmt.Sprintf("%v", v))
}

func ApplyInt(field *reflect.Value, v interface{}) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Uint64:
	case reflect.Uint32:
	case reflect.Uint16:
	case reflect.Uint8:
	case reflect.Uint:
	case reflect.Int64:
	case reflect.Int32:
	case reflect.Int16:
	case reflect.Int8:
	case reflect.Int:
		field.SetInt(v.(int64))
	case reflect.Float32:
	case reflect.Float64:
		field.SetInt(int64(v.(float64)))
	case reflect.String:
		i, _ := strconv.ParseInt(v.(string), 10, 64)
		field.SetInt(i)
	}
}

func ApplyFloat(field *reflect.Value, v interface{}) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint,
		reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		field.SetFloat(float64(v.(int64)))
	case reflect.Float32, reflect.Float64:
		field.SetFloat(v.(float64))
	case reflect.String:
		i, _ := strconv.ParseFloat(v.(string), 64)
		field.SetFloat(i)
	}
}

func ApplyBool(field *reflect.Value, v interface{}) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Bool:
		field.SetBool(v.(bool))
	case reflect.String:
		if v.(string) == "true" {
			field.SetBool(true)
		} else {
			field.SetBool(false)
		}
	}
}

func ApplySlice(field *reflect.Value, v interface{}) {
	length := reflect.ValueOf(v).Len()
	kind := reflect.ValueOf(field)

	// Raw message
	if fmt.Sprintf("%v", kind) == "<json.RawMessage Value>" {
		stringJson, err := json.Marshal(v.(map[string]interface{}))
		if err == nil {
			field.Set(reflect.ValueOf(json.RawMessage(stringJson)))
		}
		return
	}

	// Bytes
	if fmt.Sprintf("%v", kind) == "<[]uint8 Value>" {
		slice := make([]byte, 0)
		for i := 0; i < length; i++ {
			slice = append(slice, reflect.ValueOf(v).Index(i).Interface().(byte))
		}
		field.Set(reflect.ValueOf(slice))
		return
	}

	// Bytes
	if fmt.Sprintf("%v", kind) == "<[][]uint8 Value>" {
		slice := make([][]byte, 0)
		for i := 0; i < length; i++ {
			slice = append(slice, reflect.ValueOf(v).Index(i).Interface().([]byte))
		}
		field.Set(reflect.ValueOf(slice))
		return
	}

	// Fill other slice
	slice := reflect.MakeSlice(field.Type(), length, length)

	for i := 0; i < length; i++ {
		index := slice.Index(i)
		// fmt.Printf("%v\n", )

		// Each element of slice contains map
		if reflect.ValueOf(v).Index(i).Elem().Type().Kind() == reflect.Map {
			//elem := reflect.ValueOf(reflect.ValueOf(v).Index(i).Elem())
			val := reflect.ValueOf(v).Index(i).Interface().(map[string]interface{})

			//fmt.Printf("%v\n", index.Type())
			//fmt.Printf("%v\n", &index)
			//fmt.Printf("%v\n", val)

			FillFieldList(
				&index,
				index.Type(),
				val,
			)
		} else {
			index.Set(reflect.ValueOf(v).Index(i).Elem().Convert(index.Type()))
		}
	}
	field.Set(slice)
}

func ApplyTime(field *reflect.Value, v interface{}) {
	t := time.Now()

	// RFC3339Nano without T
	t1, err := time.Parse("2006-01-02 15:04:05.999999999Z07:00", v.(string))
	if err == nil {
		t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), t1.Nanosecond(), t1.Location())
		field.Set(reflect.ValueOf(t1))
		return
	}

	// My time format
	t1, err = time.Parse("2006-01-02 15:04:05 -07:00", v.(string))
	if err == nil {
		fmt.Printf("MY TIME FORMAT %v\n", v.(string))
		t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), t1.Nanosecond(), t1.Location())
		field.Set(reflect.ValueOf(t1))
		return
	}

	// RFC3339
	t1, err = time.Parse("2006-01-02T15:04:05Z07:00", v.(string))
	if err == nil {
		fmt.Printf("RFC3339 %v\n", v.(string))
		t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), t1.Nanosecond(), t1.Location())
		field.Set(reflect.ValueOf(t1))
		return
	}

	// RFC3339Nano
	t1, err = time.Parse("2006-01-02T15:04:05.999999999Z07:00", v.(string))
	if err == nil {
		fmt.Printf("RFC3339Nano %v\n", v.(string))
		t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), t1.Nanosecond(), t1.Location())
		field.Set(reflect.ValueOf(t1))
		return
	}

	// JSON date time with nanoseconds
	t1, err = time.Parse("2006-01-02T15:04:05.999Z", v.(string))
	if err == nil {
		t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), t1.Nanosecond(), t.Location())
		field.Set(reflect.ValueOf(t1))
		return
	}

	// JSON date time
	t1, err = time.Parse("2006-01-02T15:04:05", v.(string))
	if err == nil {
		t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), 0, t.Location())
		field.Set(reflect.ValueOf(t1))
		return
	}

	// Date time
	t1, err = time.Parse("2006-01-02 15:04:05", v.(string))
	if err == nil {
		t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), 0, t.Location())
		field.Set(reflect.ValueOf(t1))
		return
	}

	// Only date
	t1, err = time.Parse("2006-01-02", v.(string))
	if err == nil {
		t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t.Location())
		field.Set(reflect.ValueOf(t1))
		return
	}
}

func ApplyPtr(field *reflect.Value, v interface{}) {
	if reflect.TypeOf(v).Kind() == reflect.Map {
		field.Set(reflect.New(field.Type().Elem()))
		x := field.Elem()
		FillFieldList(&x, field.Elem().Type(), v.(map[string]interface{}))
	}

	// If value is *float64
	if reflect.TypeOf(v).Kind() == reflect.Float64 {
		// If field is *float32
		if field.Type().Elem().Kind() == reflect.Float32 {
			pointerToData := reflect.New(reflect.TypeOf(float32(0)))
			pointerToData.Elem().SetFloat(v.(float64))
			field.Set(pointerToData)
		}
	}

	// If value is *string
	if reflect.TypeOf(v).Kind() == reflect.String {
		// If field is *string
		if field.Type().Elem().Kind() == reflect.String {
			pointerToData := reflect.New(reflect.TypeOf(""))
			pointerToData.Elem().SetString(v.(string))
			field.Set(pointerToData)
		}
	}
}

func ApplyMap(field *reflect.Value, v interface{}) {
	mm := reflect.MakeMap(field.Type())
	iter := reflect.ValueOf(v).MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value().Elem()
		mm.SetMapIndex(k, v)
	}
	field.Set(mm)
}
