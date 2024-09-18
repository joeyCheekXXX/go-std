package arrayx

import (
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"math/rand"
	"time"
)

var Util = arrayUtil{}

// arrayUtil 数组工具类
type arrayUtil struct{}

// ListToTree 字典列表转树形结构
func (au arrayUtil) ListToTree(arr []map[string]interface{}, id string, pid string, child string) (mapList []interface{}) {
	mapList = []interface{}{}
	// 遍历以id_为key生成map
	idValMap := make(map[uint]interface{})
	for _, m := range arr {
		if idVal, ok := m[id]; ok {
			idValMap[idVal.(uint)] = m
		}
	}
	// 遍历
	for _, m := range arr {
		// 获取父节点
		if pidVal, ok := m[pid]; ok {
			if pNode, pok := idValMap[pidVal.(uint)]; pok {
				// 有父节点则添加到父节点子集
				if cVal, cok := pNode.(map[string]interface{})[child]; cok {
					if cVal == nil {
						cVal = []interface{}{m}
					} else {
						cVal = append(cVal.([]interface{}), m)
					}
					pNode.(map[string]interface{})[child] = cVal
					continue
				}
			}
		}
		mapList = append(mapList, m)
	}
	return
}

// StringArrayToIntArray 将string array转为int arrayx
func (au arrayUtil) StringArrayToIntArray(strArray []string) (intArray []int) {
	for _, str := range strArray {
		intArray = append(intArray, cast.ToInt(str))
	}
	return
}

// FilterItem 过滤数组 返回第一个符合条件的元素
func FilterItem[T any](s []T, cond func(t T) bool) (res T) {
	for _, v := range s {
		if cond(v) {
			res = v
			return
		}
	}
	return
}

// Filter 过滤数组 返回新数组
func Filter[T any](s []T, cond func(t T) bool) []T {
	res := []T{}
	for _, v := range s {
		if cond(v) {
			res = append(res, v)
		}
	}
	return res
}

// FilterCallback 过滤数组 执行回调
func FilterCallback[T any](s []T, cond func(t T) bool, callback func(t T)) {
	for _, v := range s {
		if cond(v) {
			callback(v)
		}
	}
}

// FilterRevise 过滤数组 并修改
func FilterRevise[T any](s []T, cond func(t T) bool, callback func(t T) (res T)) {
	for i, v := range s {
		if cond(v) {
			s[i] = callback(v)
		}
	}
}

func FindMinKeyWithCondition(treeMap *treemap.Map, condition func(int) bool) string {
	minValue := int(^uint(0) >> 1) // 初始化为最大的整数值
	var minKey string

	treeMap.Each(func(key interface{}, value interface{}) {
		if intValue, ok := value.(int); ok {
			if intValue < minValue && condition(intValue) {
				minValue = intValue
				minKey = key.(string)
			}
		}
	})
	return minKey
}

func FindKeysByValue(treeMap *treemap.Map, targetValue int) []string {
	var keys []string

	treeMap.Each(func(key interface{}, value interface{}) {
		if intValue, ok := value.(int); ok {
			if intValue == targetValue {
				keys = append(keys, key.(string))
			}
		}
	})
	return keys
}

func RandStringValue(values []string) string {
	rand.NewSource(time.Now().UnixNano())
	if len(values) == 0 {
		return ""
	}
	randIndex := rand.Intn(len(values))
	return values[randIndex]
}

func FindMinKeysSatisfyingCondition(treeMap *treemap.Map, condition func(key string, value decimal.Decimal) bool) []string {
	minKeys := []string{}
	minValue := decimal.NewFromFloat(0)

	treeMap.Each(func(key, value interface{}) {
		if condition(key.(string), value.(decimal.Decimal)) {
			if value.(decimal.Decimal).LessThan(minValue) {
				minKeys = []string{key.(string)}
				minValue = value.(decimal.Decimal)
			} else if value.(decimal.Decimal).Equals(minValue) {
				minKeys = append(minKeys, key.(string))
			}
		}
	})
	return minKeys
}

// RemoveRepeatedElement 数组去重
func RemoveRepeatedElement[T any](s []T) []T {
	result := make([]T, 0)
	m := make(map[any]bool) //map的值不重要
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}

func Foreach[T any](arr []T, callback func(key int, value T)) {
	for i, v := range arr {
		callback(i, v)
	}
}
