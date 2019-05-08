package util

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// MapFillStruct 用map填充结构体，遇到未知字段或数值转换错误返回具体错误
func MapFillStruct(data map[string]interface{}, obj interface{}) error {
	for k, v := range data {
		if v == nil {
			continue
		}
		err := setField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// MapFillStructMust 用map填充结构体，忽略所有不能填充的类型
func MapFillStructMust(data map[string]interface{}, obj interface{}) {
	for k, v := range data {
		if v == nil {
			continue
		}
		setField(obj, k, v)
	}
}

// StructToMap 将结构体的字段填充到map中
func StructToMap(v interface{}, data map[string]interface{}) {
	if v == nil || data == nil {
		return
	}
	structToMap(reflect.ValueOf(v), data)
}

func structToMap(v reflect.Value, data map[string]interface{}) {
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	// Only struct are supported
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		if !t.Field(i).Anonymous {
			name := t.Field(i).Tag.Get("map")
			if name == "-" {
				continue
			}
			if name == "" {
				name = t.Field(i).Name
			}
			data[name] = v.Field(i).Interface()
		} else {
			structToMap(v.Field(i).Addr(), data)
		}
	}
}

// StructToStruct 结构体COPY
func StructToStruct(src, dst interface{}) {
	if src == nil || dst == nil {
		return
	}
	structToStruct(reflect.ValueOf(src), reflect.ValueOf(dst))
}

func structToStruct(src, dst reflect.Value) {
	st := src.Type()
	dt := dst.Type()
	if st.Kind() == reflect.Ptr {
		src = src.Elem()
		st = st.Elem()
	}
	if dt.Kind() == reflect.Ptr {
		dst = dst.Elem()
		dt = dt.Elem()
	}
	// Only struct are supported
	if st.Kind() != reflect.Struct || dt.Kind() != reflect.Struct {
		return
	}
	var field reflect.Value
	for i := 0; i < st.NumField(); i++ {
		if !st.Field(i).Anonymous {
			field = dst.FieldByName(st.Field(i).Name)
			if field.IsValid() && field.CanSet() {
				field.Set(src.Field(i))
			}
		} else {
			structToStruct(src.Field(i).Addr(), dst)
		}
	}
}

// SortByID 对结构体Slice的field字段参照sortIDs顺序排序
func SortByID(sortIDs []int, in interface{}, field string) interface{} {
	var tmpMap = make(map[int]reflect.Value)
	v := reflect.ValueOf(in)
	if v.Type().Kind() != reflect.Slice {
		return in
	}
	for i := 0; i < v.Len(); i++ {
		sortID := v.Index(i).Elem().FieldByName(field).Int()
		tmpMap[int(sortID)] = v.Index(i)
	}

	out := reflect.MakeSlice(v.Type(), 0, 0)
	for _, id := range sortIDs {
		if v, ok := tmpMap[id]; ok {
			out = reflect.Append(out, v)
		}
	}

	return out.Interface()
}

// 用map的值替换结构的值
func setField(obj interface{}, name string, value interface{}) error {
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
		if !rv.Elem().IsValid() && rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return fmt.Errorf("Cannot set %s on nil", name)
	}

	// 普通命名
	field := rv.FieldByName(name)
	if !field.IsValid() {
		// 驼峰命名
		field = rv.FieldByName(CamelCase(name))
		if !field.IsValid() {
			// 严格驼峰命名，强制首字母缩写命名规范
			field = rv.FieldByName(CamelCaseInitialism(name))
			if !field.IsValid() {
				return fmt.Errorf("No such field: %s in obj", name)
			}
		}
	}

	if !field.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	fieldType := field.Type()     //结构体的类型
	val := reflect.ValueOf(value) //map值的反射值

	var err error
	if fieldType != val.Type() {
		val, err = typeConversion(fmt.Sprintf("%v", value), field.Type().Name()) //类型转换
		if err != nil {
			return err
		}
	}

	field.Set(val)
	return nil
}

// 类型转换
func typeConversion(value string, ntype string) (reflect.Value, error) {
	switch ntype {
	case "string":
		return reflect.ValueOf(value), nil
	case "time.Time", "Time", "time":
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	case "bool":
		b, err := strconv.ParseBool(value)
		return reflect.ValueOf(b), err
	case "int":
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	case "int8":
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	case "int16":
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int16(i)), err
	case "int32":
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int32(i)), err
	case "int64":
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	case "uint":
		i, err := strconv.ParseUint(value, 10, 64)
		return reflect.ValueOf(uint(i)), err
	case "uint8":
		i, err := strconv.ParseUint(value, 10, 64)
		return reflect.ValueOf(uint8(i)), err
	case "uint16":
		i, err := strconv.ParseUint(value, 10, 64)
		return reflect.ValueOf(uint16(i)), err
	case "uint32":
		i, err := strconv.ParseUint(value, 10, 64)
		return reflect.ValueOf(uint32(i)), err
	case "uint64":
		i, err := strconv.ParseUint(value, 10, 64)
		return reflect.ValueOf(i), err
	case "float32":
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	case "float64":
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//.......增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype + " " + value)
}
