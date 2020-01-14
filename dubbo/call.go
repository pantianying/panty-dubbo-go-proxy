package dubbo

import (
	"errors"
	"fmt"
	"github.com/pantianying/dubbo-go-proxy/common/config"
	"github.com/pantianying/dubbo-go-proxy/common/logger"
	"reflect"
	"strings"
)

type InvokeData struct {
	InterfaceName  string
	Version        string
	Group          string
	Method         string
	ParameterTypes []string
	ReqData        []interface{}
}

func (d *GenericClientPool) Call(inData InvokeData) (resp interface{}, err error) {
	c := d.Get(inData.InterfaceName, inData.Version, inData.Group)
	resp, err = c.Invoke([]interface{}{inData.Method, inData.ParameterTypes, inData.ReqData})
	if err != nil {
		logger.Errorf("GenericClient call get err:%v, InvokeData:%+v", err, inData)
		return
	}
	resp, err = dealResp(resp, config.Config.ResultFiledHumpToLine)
	if err != nil {
		logger.Errorf("deal resp err:%v", err)
		return
	}
	return
}

func dealResp(in interface{}, HumpToLine bool) (interface{}, error) {
	if in == nil {
		return in, nil
	}
	switch reflect.TypeOf(in).Kind() {
	case reflect.Map:
		if _, ok := in.(map[interface{}]interface{}); ok {
			m := mapIItoMapSI(in)
			if HumpToLine {
				m = map2x_y(m)
			}
			return m, nil
		} else if inm, ok := in.(map[string]interface{}); ok {
			if HumpToLine {
				m := map2x_y(in)
				return m, nil
			}
			return inm, nil
		}
	case reflect.Slice:
		value := reflect.ValueOf(in)
		newTemps := make([]interface{}, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			if value.Index(i).CanInterface() {
				newTemp, e := dealResp(value.Index(i).Interface(), HumpToLine)
				if e != nil {
					return nil, e
				}
				newTemps = append(newTemps, newTemp)
			} else {
				return nil, errors.New(fmt.Sprintf("unexpect err,value:%+v", value))
			}
		}
		return newTemps, nil
	default:
		return in, nil
	}
	return in, nil
}
func mapIItoMapSI(in interface{}) interface{} {
	var inMap = make(map[interface{}]interface{})
	if v, ok := in.(map[interface{}]interface{}); !ok {
		return in
	} else {
		inMap = v
	}
	outMap := make(map[string]interface{}, len(inMap))

	for k, v := range inMap {
		if v == nil {
			continue
		}
		s := fmt.Sprint(k)
		if s == "class" {
			//ignore the "class" field
			continue
		}
		vt := reflect.TypeOf(v)
		switch vt.Kind() {
		case reflect.Map:
			if _, ok := v.(map[interface{}]interface{}); ok {
				v = mapIItoMapSI(v)
			}
		case reflect.Slice:
			vl := reflect.ValueOf(v)
			os := make([]interface{}, 0, vl.Len())
			for i := 0; i < vl.Len(); i++ {
				if vl.Index(i).CanInterface() {
					osv := mapIItoMapSI(vl.Index(i).Interface())
					os = append(os, osv)
				}
			}
			v = os
		}
		outMap[s] = v

	}
	return outMap
}
func map2x_y(in interface{}) interface{} {

	var m map[string]interface{}
	if v, ok := in.(map[string]interface{}); ok {
		m = v
	} else {
		return in
	}

	var out = make(map[string]interface{}, len(m))
	for k1, v1 := range m {
		x := humpToUnderline(k1)

		if v1 == nil {
			out[x] = v1
		} else if reflect.TypeOf(v1).Kind() == reflect.Struct {
			out[x] = map2x_y(struct2Map(v1))
		} else if reflect.TypeOf(v1).Kind() == reflect.Slice {
			value := reflect.ValueOf(v1)
			var newTemps = make([]interface{}, 0, value.Len())
			for i := 0; i < value.Len(); i++ {
				newTemp := map2x_y(value.Index(i).Interface())
				newTemps = append(newTemps, newTemp)
			}
			out[x] = newTemps
		} else if reflect.TypeOf(v1).Kind() == reflect.Map {
			out[x] = map2x_y(v1)
		} else {
			out[x] = v1
		}
	}
	return out
}
func humpToUnderline(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}
func struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
