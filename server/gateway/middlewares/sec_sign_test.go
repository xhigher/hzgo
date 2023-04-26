package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/utils"
)

func TestSign(t *testing.T) {
	baseParams := &defines.BaseParams{
		Ap:   "com.xhigher.hzgo",
		Av:   "1.2.3",
		Dt:   1,
		Did:  "12131313131313",
		Bd:   "HUAWEI",
		Md:   "mate 40",
		Os:   "11.11",
		Nt:   1,
		Ch:   "weixin",
		Ip:   "",
		Loc:  "",
		Imei: "",
		Oaid: "",
		Idfa: "",
		Ds:   "",
		Sign: "",
		Ts:   utils.NowTime(),
	}

	secret := "ysntsFNlD4nwuolcc8evZtTPsToizNtA"
	signature := utils.NewDefaultSignature(secret)
	signature.UseReflect(true)
	sign, err := signature.Sign(baseParams)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	baseParams.Sign = sign
	fmt.Println("sign", sign)

	dataStr := utils.JSONString(baseParams)

	fmt.Println("dataStr", dataStr)
	headerData := base64.StdEncoding.EncodeToString([]byte(dataStr))
	fmt.Println("headerData", headerData)

	bytes, err := base64.StdEncoding.DecodeString(headerData)
	if err != nil {
		fmt.Println("error", err)
		return
	}

	var params2 *defines.BaseParams
	err = json.Unmarshal(bytes, &params2)
	if err != nil {
		fmt.Println("error", err)
		return
	}

	ok, err := signature.Verify(params2)

	fmt.Println("error", err, "ok", ok)
}
