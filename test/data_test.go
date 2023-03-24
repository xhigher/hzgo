package test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"math/rand"
	"testing"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AES加密,CBC
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func CommonAesEncrypt(key, data string) string {
	if len(data) > 0 {
		encrypted, err := AesEncrypt([]byte(data), []byte(key))
		if err == nil {
			return base64.StdEncoding.EncodeToString(encrypted)
		}
	}
	return data
}

func CommonAesDecrypt(key, data string) string {
	if len(data) > 0 {
		encrypted, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return data
		}
		tempBytes, err := AesDecrypt(encrypted, []byte(key))
		if err == nil {
			return string(tempBytes)
		}
	}
	return data
}

func TestA(t *testing.T) {

	key:= "dPpEvwepqi$dF9HE#tHVWly2a&vpUR2@"
	text := "[{\"domain\":\".amazon.com\",\"expirationDate\":1986017706,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"_delighted_web\",\"path\":\"/\",\"sameSite\":\"lax\",\"secure\":false,\"session\":false,\"storeId\":\"0\",\"value\":\"{%22L66xj3IYPSF6JZGA%22:{%22_delighted_fst%22:{%22t%22:%221657870503078%22}%2C%22_delighted_lst%22:{%22t%22:%221655283442000%22%2C%22m%22:{%22token%22:%221DpSjJ8Js38unbLfwqtVsmbD%22}}}}\",\"id\":1},{\"domain\":\".amazon.com\",\"expirationDate\":1704702069,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"AMCV_7742037254C95E840A4C98A6%40AdobeOrg\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":false,\"session\":false,\"storeId\":\"0\",\"value\":\"MCMID|35176176545934014784333302565956035563\",\"id\":2},{\"domain\":\".amazon.com\",\"expirationDate\":1708565953,\"hostOnly\":false,\"httpOnly\":true,\"name\":\"at-main\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":true,\"session\":false,\"storeId\":\"0\",\"value\":\"Atza|IwEBIB9akUMz6cGWHk6cllR58sejf_cc1k0n9Ted0v5_L6vRzn7jTPiNdpLl2N49ie7GN7IU6ylN-GdDK3RAwaWF961exevM6jGI70UgmYXC6Z9YH9oc1of3on286V2Y3X00sBv8dMwxug_6QB3qPsyc91kq0dZoQZItR0T6qAdRU9JU6eENHzJZJ1Zxq_sdjaZV050gMk7rWOwHjIPdgSpKhte-wZn-34JABA_pMhmLmaEcvQ\",\"id\":3},{\"domain\":\".amazon.com\",\"expirationDate\":1679571952,\"hostOnly\":false,\"httpOnly\":true,\"name\":\"av-profile\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":true,\"session\":false,\"storeId\":\"0\",\"value\":\"cGlkPWFtem4xLmFjdG9yLnBlcnNvbi5vaWQuQTJYNFFKRjNIREk3WkcmdGltZXN0YW1wPTE2NzY5ODAwNzkyNTkmdmVyc2lvbj12MQ.g_BPqyL_qJn5aEbpcZi76pHIkWLJhzcgWb-1e5MiRUsOAAAAAQAAAABj9K9vcmF3AAAAAPgWC9WfHH8iB-olH_E9xQ\",\"id\":4},{\"domain\":\".amazon.com\",\"expirationDate\":1701327931,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"av-timezone\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":false,\"session\":false,\"storeId\":\"0\",\"value\":\"Asia/Hong_Kong\",\"id\":5},{\"domain\":\".amazon.com\",\"expirationDate\":1708568725,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"i18n-prefs\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":false,\"session\":false,\"storeId\":\"0\",\"value\":\"USD\",\"id\":6},{\"domain\":\".amazon.com\",\"expirationDate\":1705134069,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"kndctr_7742037254C95E840A4C98A6_AdobeOrg_identity\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":false,\"session\":false,\"storeId\":\"0\",\"value\":\"CiYzNTE3NjE3NjU0NTkzNDAxNDc4NDMzMzMwMjU2NTk1NjAzNTU2M1IPCJqY5f3QMBgBKgRKUE4z8AGamOX90DA=\",\"id\":7},{\"domain\":\".amazon.com\",\"expirationDate\":2082787202.382692,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"lc-main\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":true,\"session\":false,\"storeId\":\"0\",\"value\":\"en_US\",\"id\":8},{\"domain\":\".amazon.com\",\"expirationDate\":1769743150,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"s_dslv\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":false,\"session\":false,\"storeId\":\"0\",\"value\":\"1675135150443\",\"id\":9},{\"domain\":\".amazon.com\",\"expirationDate\":2107135150,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"s_nr\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":false,\"session\":false,\"storeId\":\"0\",\"value\":\"1675135150441-New\",\"id\":10},{\"domain\":\".amazon.com\",\"expirationDate\":1711589944,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"s_pers\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":false,\"session\":false,\"storeId\":\"0\",\"value\":\"%20s_fid%3D13AF2F354D22B692-2A60C10D2AEF10AB%7C1834796344427%3B%20s_dl%3D1%7C1677031744429%3B%20s_ev15%3D%255B%255B%2527SCUSWPDirect%2527%252C%25271675647392334%2527%255D%252C%255B%2527SCSOAlogin%2527%252C%25271675683579328%2527%255D%252C%255B%2527SCUSWPDirect%2527%252C%25271677029944434%2527%255D%255D%7C1834796344434%3B\",\"id\":11},{\"domain\":\".amazon.com\",\"expirationDate\":2107135117,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"s_vnum\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":false,\"session\":false,\"storeId\":\"0\",\"value\":\"2107135117375%26vn%3D1\",\"id\":12},{\"domain\":\".amazon.com\",\"expirationDate\":1708565954,\"hostOnly\":false,\"httpOnly\":true,\"name\":\"sess-at-main\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":true,\"session\":false,\"storeId\":\"0\",\"value\":\"\\\"YdbH8V8FJ+DiZJA04d/L5QqZwd5/S/enS5le37B2EWs=\\\"\",\"id\":13},{\"domain\":\".amazon.com\",\"expirationDate\":2082787201.795,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"session-id\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":false,\"session\":false,\"storeId\":\"0\",\"value\":\"144-2276449-9885648\",\"id\":14},{\"domain\":\".amazon.com\",\"expirationDate\":2082787201.795108,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"session-id-time\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":false,\"session\":false,\"storeId\":\"0\",\"value\":\"2082787201l\",\"id\":15},{\"domain\":\".amazon.com\",\"expirationDate\":1708587040.621836,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"session-token\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":true,\"session\":false,\"storeId\":\"0\",\"value\":\"wb4owH0rifBWX1HCJw7/YHYzllYBl4uvotyzV3mDi1wX408cB55SmgXL6nhZDQlCRuiNV2aCDMzEsmESBsPbsyIq6Lg8oFrYzJeaE1+SshFcB2TIDsYY127fmFLQzODhI3lzfEemY33g6pJUvEGX6GPa9MuFz9yktPL6Vc+EMlkAPKLL7CsF+HImPVrhvTc4BUyt3+JHQGRZZXS22LB0qd1d3vP/sQnjbmUrM2kqKikzeKVlAMbThA==\",\"id\":16},{\"domain\":\".amazon.com\",\"expirationDate\":1707619668,\"hostOnly\":false,\"httpOnly\":true,\"name\":\"sid\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":true,\"session\":false,\"storeId\":\"0\",\"value\":\"\\\"06KOVd4DBQL9z9fNJYWI8Q==|ZldA/h+NDQoPSREb4+ANBrUMAbxtcUcU17rZk5V2NEA=\\\"\",\"id\":17},{\"domain\":\".amazon.com\",\"expirationDate\":1708565953,\"hostOnly\":false,\"httpOnly\":true,\"name\":\"sst-main\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":true,\"session\":false,\"storeId\":\"0\",\"value\":\"Sst1|PQFv19cC7fm636uXxWXdLbmlCQDSyooeGrMK9i3FMGWcwsmAJ7862nmWsUswyomMNRpZdiMXal3mNdWFLZeNbJEWtUrdTnBOB-UHl18XCni32wTr0ic5cf4Izuo1QfdQM9Wf8oJVIBc101yHzCdfDMYB6UbsmvhSxKUO2RQGiYJ2XlekiY-U-4_wlrN6oRvSFaO-uKVmwMi_7210bcGrHGquetx4KJ5QflUFQ13zmHxIY5iVHSMypmPfOqUGzbDJXuyCxET89qdp3VEcUy1arOeGqi89ERv_jNGuLSyeTmNFLa0\",\"id\":18},{\"domain\":\".amazon.com\",\"expirationDate\":2082787201,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"ubid-acbus\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":true,\"session\":false,\"storeId\":\"0\",\"value\":\"130-6544934-9772331\",\"id\":19},{\"domain\":\".amazon.com\",\"expirationDate\":1708568725,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"ubid-main\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":true,\"session\":false,\"storeId\":\"0\",\"value\":\"130-9680917-1849440\",\"id\":20},{\"domain\":\".amazon.com\",\"expirationDate\":1708568725,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"x-main\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":true,\"session\":false,\"storeId\":\"0\",\"value\":\"\\\"6CbRQRnfUzq1F0mDMVlK?TwNYoL3n?zRkRkmnNXM?w1d@S@9uB1BjdLYwkcdRyJE\\\"\",\"id\":21},{\"domain\":\".sellercentral.amazon.com\",\"expirationDate\":1679580181,\"hostOnly\":false,\"httpOnly\":false,\"name\":\"cwr_u\",\"path\":\"/\",\"sameSite\":\"strict\",\"secure\":true,\"session\":false,\"storeId\":\"0\",\"value\":\"2aa4ea9e-ff12-4993-8644-f2e217dd53c4\",\"id\":22},{\"domain\":\"sellercentral.amazon.com\",\"expirationDate\":1708582737.308849,\"hostOnly\":true,\"httpOnly\":true,\"name\":\"__Host_mlang\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":true,\"session\":false,\"storeId\":\"0\",\"value\":\"zh_CN\",\"id\":23},{\"domain\":\"sellercentral.amazon.com\",\"expirationDate\":1708575664,\"hostOnly\":true,\"httpOnly\":true,\"name\":\"__Host-mselc\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":true,\"session\":false,\"storeId\":\"0\",\"value\":\"H4sIAAAAAAAA/6tWSs5MUbJSSsytyjPUS0xOzi/NK9HLT85M0XP0cTX1tLA0DPIw8fVW0lHKRVKYm1qUnJEIUompLBtZXQFIRUhYgIu3p3eEgYtrkFItAHN1LddzAAAA\",\"id\":24},{\"domain\":\"sellercentral.amazon.com\",\"expirationDate\":1707269953,\"hostOnly\":true,\"httpOnly\":false,\"name\":\"csm-hit\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":false,\"session\":false,\"storeId\":\"0\",\"value\":\"tb:37SR6771ZCF26DR5CQG5+s-CRFR4Y01B9TZ63RG1553|1677029953531&t:1677029953531&adb:adblk_no\",\"id\":25},{\"domain\":\"sellercentral.amazon.com\",\"expirationDate\":1678354698,\"hostOnly\":true,\"httpOnly\":false,\"name\":\"ld\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":false,\"session\":false,\"storeId\":\"0\",\"value\":\"SCUSWPDirect\",\"id\":26},{\"domain\":\"sellercentral.amazon.com\",\"expirationDate\":1677655846.210327,\"hostOnly\":true,\"httpOnly\":true,\"name\":\"stck\",\"path\":\"/\",\"sameSite\":\"unspecified\",\"secure\":true,\"session\":false,\"storeId\":\"0\",\"value\":\"NA\",\"id\":27}]"
	text = CommonAesEncrypt(key, text)
	fmt.Println(text)

	text2 := text
	text2 = CommonAesDecrypt(key, text2)
	fmt.Println(text2)

}

func TestMap(t *testing.T) {
	playerNum := 4
	id := (((playerNum-1) / 2)+1) * 200 + rand.Intn(4)
	fmt.Println("map.id", id)
}
