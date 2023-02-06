package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestIntToBase36(t *testing.T) {

	micro := NowTimeMicro() - 1040000000000000

	code := IntToBase36(micro)
	fmt.Println("TestIntToBase36: ", micro, code)
}

type testParam struct {
	Version string `json:"version"`//是	版本信息，固定值v2
	Nonce string `json:"nonce"`//是	用于防重放。最长32个字符
	Timestamp string `json:"timestamp"`//是	发起此查询操作的时间。UNIX时间戳，单位：毫秒
	SecretId string `json:"secretId"`//是	产品密钥id，由易盾验证码服务分配
	Signature string `json:"signature"`//是	此次请求的签名，用来验证请求的合法性。具体算法见 接口鉴权
	CaptchaId string `json:"captchaId"`
	Validate string `json:"validate"`
}

func TestSignature(t *testing.T) {
	secret := "YRf34Ubaid3SSZdEqNwq2Jo6bny7A2aF"
	signature := NewSignature(secret, "signature")
	signature.CheckTimestamp("timestamp", 2*time.Minute)

	params := testParam{
		Version: "1.2.3",
		Nonce: RandString(20),
		Timestamp: fmt.Sprintf("%d", NowTimeMillis()),
		SecretId: "12311231312",
		CaptchaId: "abcdefg",
		Validate: "12313213213213213213213213",
	}
	sign, err := signature.Sign(params)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("sign:", sign)
	params.Signature = sign

	time.Sleep(1*time.Second)
	ok, err := signature.Verify(params)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("verify:", ok)

	time.Sleep(3*time.Minute)
	ok, err = signature.Verify(params)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("verify:", ok)
}