package utils

import (
	"encoding/json"
	"fmt"
	"github.com/xhigher/hzgo/logger"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)
type SignEncrypt int
const (
	SignMD5 SignEncrypt = 0
	SignHmacSha1 SignEncrypt = 1
)


type Signature struct {
	secret string        // 签名secret
	keySign string
	keyTimestamp string
	duration       time.Duration // 签名有效时间，标准是120秒，也可按照业务标准来设置
	reflect     bool          //是否使用reflect反射机制来转换请求参数到map，默认为false不使用，不使用就使用json来转换
	encrypt SignEncrypt
}



func NewDefaultSignature(secret string) *Signature {
	return &Signature{
		secret: secret,
		keySign: "sign",
		keyTimestamp: "",
		duration: 0,
		reflect:     false,
		encrypt: SignMD5,

	}
}

func NewSignature(secret string, keySign string) *Signature {
	return &Signature{
		secret: secret,
		keySign: keySign,
		keyTimestamp: "",
		duration: 0,
		reflect:     false,
		encrypt: SignMD5,
	}
}

func (s *Signature) SetSignEncrypt(encypt SignEncrypt) {
	s.encrypt = encypt
}

func (s *Signature) CheckTimestamp(keyName string, duration time.Duration) {
	s.keyTimestamp = keyName
	s.duration = duration
}

func (s *Signature) UseReflect(yes bool) {
	s.reflect = yes
}

/**
计算签名
参数key-value对，value为int类型或string，暂不支持struct结构，struct结构可转成json格式的字符串作为参数
参数key-value对可以是map，也可以是struct结构数据，签名时会统一转成map数据
*/
func (s *Signature) Sign(params interface{}) (string, error) {
	var paramsMap map[string]interface{}
	var err error
	if !s.reflect {
		paramsMap, err = transParamsToMap(params)
	} else {
		paramsMap, err = transParamsToMapByReflect(params)
	}
	if err != nil {
		return "", err
	}
	keys := make([]string, 0)
	for k := range paramsMap {
		// 忽略签名字段
		if k == s.keySign {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	signParams := make([]string, 0)
	for _, key := range keys {
		item := fmt.Sprintf("%s=%v", key, paramsMap[key])
		signParams = append(signParams, item)
	}
	signStr := strings.Join(signParams, "&")
	var sign string
	if s.encrypt == SignHmacSha1 {
		sign = EncryptHmacSha1(signStr, s.secret)
	}else{
		sign = MD5(fmt.Sprintf("%s:%s", signStr, s.secret))
	}
	return sign, nil
}

/**
校验签名
params包含签名字段
*/
func (s *Signature) Verify(params interface{}) (bool, error) {
	var paramsMap map[string]interface{}
	var err error
	if !s.reflect {
		paramsMap, err = transParamsToMap(params)
	} else {
		paramsMap, err = transParamsToMapByReflect(params)
	}
	if err != nil {
		return false, err
	}
	//检查时间戳是否过期
	if len(s.keyTimestamp) > 0 {
		if paramsMap[s.keyTimestamp] == nil {
			return false, fmt.Errorf("timestamp not exist")
		}
		tsStr, ok := paramsMap[s.keyTimestamp].(string)
		if !ok || tsStr == "" {
			return false, fmt.Errorf("timestamp wrong value")
		}
		//以秒计算有效期
		if len(tsStr) > 10 {
			tsStr = tsStr[:10]
		}
		timestamp, _ := strconv.ParseInt(tsStr, 10, 64)
		signTime := time.Unix(timestamp, 0)
		nowTime := time.Now()
		// 在签名端时间比验证端时间快时，会出现当前验证时间小于签名时间的情况，理论上签名不会通过，所以加个几秒的误差，允许签名端时间比验证端时间快几秒
		// 签名端时间比验证端时间慢的时候，不会出现上述问题
		if nowTime.Before(signTime.Add(-10*time.Second)) || nowTime.After(signTime.Add(s.duration)) {
			return false, fmt.Errorf("timestamp expired")
		}
	}

	//检查签名串是否正确
	if paramsMap[s.keySign] == nil {
		return false, fmt.Errorf("sign not exist")
	}
	verifySign, ok := paramsMap[s.keySign].(string)
	if !ok || verifySign == "" {
		return false, fmt.Errorf("sign param wrong")
	}
	sign, err := s.Sign(params)
	if err != nil {
		return false, err
	}
	if sign != verifySign {
		logger.Warnf("check sign error, sign:%s, verify_sign:%s", sign, verifySign)
		return false, fmt.Errorf("sign not equal")
	}
	return true, nil
}

/**
将map类型的params，或struct类型，或指针指向的struct类型，转换为map[string]interface{}
利用json进行转换
已有问题：params为struct时，转成json字符串，如果value为空字符或0，key可能会被忽略去掉，然后转成map就会丢失部分字段，造成签名校验失败
如果参数都没有为空或空字符串，这个问题就没影响
*/
func transParamsToMap(params interface{}) (map[string]interface{}, error) {
	paramsByte, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	paramsMap := make(map[string]interface{}, 0)
	err = json.Unmarshal(paramsByte, &paramsMap)
	if err != nil {
		return nil, err
	}
	return paramsMap, nil
}

/**
将map类型的params，或struct类型，或指针指向的struct类型，转换为map[string]interface{}
利用反射进行转换
1、map类型的，通过json转成map[string]interface{}返回
2、struct类型的，只提取详细json标签的公有字段到map[string]interface{}并返回，需要每个字段单独标注json-tag
3、其他类型的params，不进行转换，直接返回错误提示信息
*/
func transParamsToMapByReflect(params interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(params)
	//map类型的params，通过json转成map[string]interface{}返回
	if v.Kind() == reflect.Map {
		paramsByte, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		mapRes := make(map[string]interface{}, 0)
		err = json.Unmarshal(paramsByte, &mapRes)
		if err != nil {
			return nil, err
		}
		for key, val := range mapRes {
			// json转换时，整数会转成float64类型，大整数情况展示可能采用科学记数法展示，会有问题，就强制转成string
			if valFloat, ok := val.(float64); ok {
				mapRes[key] = fmt.Sprintf("%d", uint64(valFloat))
			}
		}
		return mapRes, nil
	}
	v = reflect.Indirect(v)
	//非struct类型的，不进行转换，直接返回错误提示信息
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("if params is not map, it should be struct or non-nil pointer which point to struct")
	}
	t := v.Type()
	//struct类型的，只提取详细json标签的公有字段到map[string]interface{}并返回
	mapRes := make(map[string]interface{}, 0)
	for i := 0; i < t.NumField(); i++ {
		// 忽略struct里面的私有字段，这些字段无法正常访问
		if !v.Field(i).CanInterface() {
			continue
		}
		// 忽略struct里面没有json标签，或json标签为空失效的字段，这些字段不在请求参数内
		jsontag := t.Field(i).Tag.Get("json")
		jlist := strings.Split(jsontag, ",")
		if len(jlist) == 0 || strings.TrimSpace(jlist[0]) == "" || strings.TrimSpace(jlist[0]) == "-" {
			continue
		}
		name := strings.TrimSpace(jlist[0])
		value := v.Field(i).Interface()
		mapRes[name] = value
	}
	return mapRes, nil
}
