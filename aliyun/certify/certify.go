package certify

import (
	"encoding/json"
	"errors"
	"fmt"
	cloudauth "github.com/alibabacloud-go/cloudauth-20190307/client"
	tearpc "github.com/alibabacloud-go/tea-rpc/client"
	teautils "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/logger"
	"strconv"
)

type CertifyConfig struct {
	Endpoints       []string
	AccessKeyId     string
	AccessKeySecret string
	SceneId         int64
}

type CertifyInitResult struct {
	Code      string
	Message   string
	CertifyId string
}

func (result *CertifyInitResult) Success() bool {
	return result.Code == "200"
}

type CertifyInitData struct {
	Orderid         string
	Type            int32
	SceneId         int64
	Userid          int64
	IdCardCode      string
	IdCardName      string
	Avatar          string
	CertifiedAvatar string
}

type CertifyCompareData struct {
	Orderid         string
	SceneId         int64
	SourceImage     string
	TargetImage     string
	SourceCertifyId string
	TargetCertifyId string
}

type CertifyQueryResult struct {
	Code         string  `json:"code"`
	Message      string  `json:"message"`
	CertifyId    string  `json:"certify_id"`
	Passed       bool    `json:"passed"`
	Reason       string  `json:"reason"`
	QualityScore float32 `json:"quality_score"`
	VerifyScore  float32 `json:"verify_score"`
	PictureUrl   string  `json:"picture_url"`
}

func (result *CertifyQueryResult) Success() bool {
	return result.Code == "200"
}

type ResultMaterialPicture struct {
	FaceAttackScore          float32 `json:"faceAttackScore"`
	QualityScore             float32 `json:"qualityScore"`
	VerifyScore              float32 `json:"verifyScore"`
	PictureUrl               string  `json:"pictureUrl"`
	OssBucketName            string  `json:"ossBucketName"`
	OssObjectName            string  `json:"ossObjectName"`
	OssVerifyVideoObjectName string  `json:"ossVerifyVideoObjectName"`
	VerifyVideoUrl           string  `json:"verifyVideoUrl"`
}

type ResultMaterialInfo struct {
	FaceAttack         string                `json:"faceAttack"`
	FaceOcclusion      string                `json:"faceOcclusion"`
	FacialPictureFront ResultMaterialPicture `json:"facialPictureFront"`
}

var (
	certConfig  *CertifyConfig
	certClients []*cloudauth.Client

	resultSubCodes = map[string]string{
		"200": "认证通过",
		"201": "姓名和身份证不一致",
		"202": "查询不到身份信息",
		"203": "查询不到照片或照片不可用",
		"204": "人脸比对不一致，请更换头像",
		"205": "活体检测存在风险",
		"206": "业务策略限制",
		"207": "身份证人脸比对不一致",
		"209": "权威比对源异常",
		"210": "认证通过",
	}
	initRespCodes = map[string]string{
		"401": "姓名或身份证号码长度不符合国家标准，请检查",
		"414": "当前移动设备不支持刷脸认证，请更换设备后操作",
		"415": "当前认证SDK版本不支持刷脸认证，请升级SDK后操作",
		"416": "前操作系统版本不支持刷脸认证，请升级系统或更换设备操作",
		"417": "无法使用刷脸服务，当前身份信息无可用比对源",
		"418": "当天刷脸认证次数过多，请明天再试",
		"419": "传入头像不可用，请更换头像",
		"421": "传入头像不符合认证要求，请更换头像",
		"500": "系统错误",
	}
)

func Init(conf *CertifyConfig) (err error) {
	if conf == nil {
		err = errors.New("error config nil")
		return
	}
	if len(conf.AccessKeyId) == 0 {
		err = errors.New("error AccessKeyId nil")
		return
	}
	if len(conf.AccessKeySecret) == 0 {
		err = errors.New("error AccessKeySecret nil")
		return
	}
	if conf.SceneId == 0 {
		err = errors.New("error SceneId nil")
		return
	}
	if len(conf.Endpoints) == 0 {
		conf.Endpoints = []string{
			"cloudauth.cn-shanghai.aliyuncs.com",
			"cloudauth.cn-beijing.aliyuncs.com",
		}
	}
	certConfig = conf
	certClients = make([]*cloudauth.Client, len(conf.Endpoints))
	for i, endpoint := range certConfig.Endpoints {
		config := &tearpc.Config{
			// 您的AccessKey ID。
			AccessKeyId: tea.String(certConfig.AccessKeyId),
			// 您的AccessKey Secret。
			AccessKeySecret: tea.String(certConfig.AccessKeySecret),
			Endpoint:        &endpoint,
			// 设置HTTP代理。
			// HttpProxy: tea.String("http://xx.xx.xx.xx:xxxx"),
			// 设置HTTPS代理。
			// HttpsProxy: tea.String("https://xx.xx.xx.xx:xxxx"),
		}
		certClients[i], err = cloudauth.NewClient(config)
		if err != nil {
			err = errors.New("error client new failed")
			return
		}
	}
	return
}

func CertifyInit(data *CertifyInitData, metaInfo string) (result *CertifyInitResult, err error) {
	request := &cloudauth.InitFaceVerifyRequest{
		SceneId:      tea.Int64(certConfig.SceneId), //tea.Int64(data.SceneId),
		OuterOrderNo: tea.String(data.Orderid),
		Model:        tea.String("LIVENESS"),
		MetaInfo:     tea.String(metaInfo),
	}

	if data.Type == consts.UserCertifyRealName {
		request.ProductCode = tea.String("ID_PRO")
		request.CertType = tea.String("IDENTITY_CARD")
		request.CertName = tea.String(data.IdCardName)
		request.CertNo = tea.String(data.IdCardCode)
	} else {
		outUserid := strconv.FormatInt(data.Userid+75321968, 10)
		request.ProductCode = tea.String("PV_FV")
		request.UserId = tea.String(outUserid)
		request.FaceContrastPictureUrl = tea.String(data.Avatar)
	}

	runtime := &teautils.RuntimeOptions{
		ReadTimeout:    tea.Int(10000),
		ConnectTimeout: tea.Int(10000),
	}
	var response *cloudauth.InitFaceVerifyResponse
	for _, client := range certClients {
		response, err = client.InitFaceVerify(request, runtime)
		if err != nil {
			logger.Errorf("CertifyInit error: %v", err)
			return
		}
		if response == nil {
			logger.Errorf("CertifyInit response nil")
			continue
		}

		respStr := teautils.ToJSONString(response)
		if tea.BoolValue(teautils.EqualString(tea.String("500"), response.Code)) {
			logger.Errorf("CertifyInit response: %v", *respStr)
			continue
		}
		result = &CertifyInitResult{
			Code:    *response.Code,
			Message: *response.Message,
		}
		if *response.Code == "200" {
			result.CertifyId = *response.ResultObject.CertifyId
		} else {
			logger.Errorf("CertifyInit response: %v", *respStr)
		}
		return
	}
	return
}

func CertifyQuery(orderid, certifyId string) (result *CertifyQueryResult, err error) {
	request := &cloudauth.DescribeFaceVerifyRequest{
		SceneId:   tea.Int64(certConfig.SceneId), //tea.Int64(sceneId),
		CertifyId: tea.String(certifyId),
	}
	runtime := &teautils.RuntimeOptions{
		ReadTimeout:    tea.Int(10000),
		ConnectTimeout: tea.Int(10000),
	}
	var response *cloudauth.DescribeFaceVerifyResponse
	for _, client := range certClients {
		response, err = client.DescribeFaceVerify(request, runtime)
		if err != nil {
			logger.Errorf("CertifyQuery error: %v", err)
			return
		}
		if response == nil {
			logger.Errorf("CertifyQuery response nil")
			continue
		}

		respStr := teautils.ToJSONString(response)
		logger.Infof("CertifyQuery response: %v, %v, %v", orderid, certifyId, *respStr)
		if tea.BoolValue(teautils.EqualString(tea.String("500"), response.Code)) {
			continue
		}
		result = &CertifyQueryResult{
			Code:    *response.Code,
			Message: *response.Message,
		}
		if *response.Code == "200" {
			passed := false
			if *response.ResultObject.Passed == "T" {
				passed = true
			}
			result.Passed = passed

			materialInfo := &ResultMaterialInfo{}
			json.Unmarshal([]byte(*response.ResultObject.MaterialInfo), &materialInfo)
			result.Reason = *response.ResultObject.SubCode
			result.VerifyScore = materialInfo.FacialPictureFront.VerifyScore
			result.QualityScore = materialInfo.FacialPictureFront.QualityScore
			result.PictureUrl = materialInfo.FacialPictureFront.PictureUrl

		} else {
			logger.Errorf("CertifyQuery response: %v", *respStr)
		}
		return
	}
	return
}

func GetInitRespCodeDesc(code string) string {
	if desc, ok := initRespCodes[code]; ok {
		return desc
	}
	return fmt.Sprintf("系统内部错误：%s", code)
}

func GetResultReasonDesc(subCode string) string {
	if desc, ok := resultSubCodes[subCode]; ok {
		return desc
	}
	return "认证失败"
}

func CertifyCompare(data *CertifyCompareData) (result *CertifyQueryResult, err error) {
	request := &cloudauth.CompareFaceVerifyRequest{
		// 请输入场景ID。
		SceneId:      tea.Int64(certConfig.SceneId), //tea.Int64(data.SceneId),
		ProductCode:  tea.String("PV_FC"),
		OuterOrderNo: tea.String(data.Orderid),
	}

	if len(data.SourceImage) > 0 {
		request.SourceFaceContrastPictureUrl = tea.String(data.SourceImage)
	}
	if len(data.SourceCertifyId) > 0 {
		request.SourceCertifyId = tea.String(data.SourceCertifyId)
	}
	if len(data.TargetImage) > 0 {
		request.TargetFaceContrastPictureUrl = tea.String(data.TargetImage)
	}
	if len(data.TargetCertifyId) > 0 {
		request.TargetCertifyId = tea.String(data.TargetCertifyId)
	}

	runtime := &teautils.RuntimeOptions{
		ReadTimeout:    tea.Int(10000),
		ConnectTimeout: tea.Int(10000),
	}

	dataStr := teautils.ToJSONString(data)
	var response *cloudauth.CompareFaceVerifyResponse
	for _, client := range certClients {
		response, err = client.CompareFaceVerify(request, runtime)
		if err != nil {
			logger.Errorf("CertifyCompare error: %v", err)
			return
		}
		if response == nil {
			logger.Errorf("CertifyCompare response nil")
			continue
		}
		respStr := teautils.ToJSONString(response)
		logger.Infof("CertifyCompare response: %v, %v", *dataStr, *respStr)
		if tea.BoolValue(teautils.EqualString(tea.String("500"), response.Code)) {
			continue
		}
		result = &CertifyQueryResult{
			Code:    *response.Code,
			Message: *response.Message,
		}
		if *response.Code == "200" {
			passed := false
			if *response.ResultObject.Passed == "T" {
				passed = true
			}
			result.Passed = passed
			result.CertifyId = *response.ResultObject.CertifyId
			result.VerifyScore = *response.ResultObject.VerifyScore
		} else {
			logger.Errorf("CertifyCompare response: %v", *respStr)
		}
		return
	}
	return
}
