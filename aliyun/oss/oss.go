package oss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/xhigher/hzgo/consts"
	"github.com/xhigher/hzgo/env"
	"github.com/xhigher/hzgo/logger"
	"hash"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssConfig struct {
	Endpoint            string `json:"endpoint" ini:"endpoint"`
	InternalEndpoint    string `json:"internal_endpoint" ini:"internal_endpoint"`
	Protocol            string `json:"protocol" ini:"protocol"`
	AccessKeyId         string `json:"access_key_id" ini:"accesss_key_id"`
	AccessKeySecret     string `json:"accesss_key_secret" ini:"accesss_key_secret"`
	DataAccessKeyId     string `json:"data_access_key_id" ini:"data_accesss_key_id"`
	DataAccessKeySecret string `json:"data_accesss_key_secret" ini:"data_accesss_key_secret"`
	OssDomain           string `json:"oss_domain" ini:"oss_domain"`
	OssInternalDomain   string `json:"oss_internal_domain" ini:"oss_internal_domain"`
	OssVpcDomain        string `json:"oss_vpc_domain" ini:"oss_vpc_domain"`
	CDNDomain           string `json:"cdn_domain" ini:"cdn_domain"`

	ImageCDNDomain string `json:"image_cdn_domain" ini:"image_cdn_domain"`
	VideoCDNDomain string `json:"video_cdn_domain" ini:"image_cdn_domain"`
	AudioCDNDomain string `json:"audio_cdn_domain" ini:"image_cdn_domain"`

	ImageBucket string `json:"image_bucket" ini:"image_bucket"`
	VideoBucket string `json:"video_bucket" ini:"video_bucket"`
	AudioBucket string `json:"audio_bucket" ini:"audio_bucket"`
	DataBucket  string `json:"data_bucket" ini:"data_bucket"`

	TestImageBucket string `json:"test_image_bucket" ini:"test_image_bucket"`
	TestVideoBucket string `json:"test_video_bucket" ini:"test_video_bucket"`
	TestAudioBucket string `json:"test_audio_bucket" ini:"test_audio_bucket"`
	TestDataBucket  string `json:"test_data_bucket" ini:"test_data_bucket"`

	ExpireTime int64 `json:"expire_time" ini:"expire_time"`
}

type OssPolicyConfig struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type OssPolicyData struct {
	Accessid  string `json:"accessid"`
	Host      string `json:"host"`
	Expire    int64  `json:"expire"`
	Signature string `json:"signature"`
	Policy    string `json:"policy"`
	Dir       string `json:"dir"`
	ObjectAcl string `json:"object_acl"`
	Callback  string `json:"callback"`
}

type BucketObject struct {
	Bucket string
	Key    string
	Id     string
	Type   string
}

var (
	ossConfig *OssConfig

	ossClient         *oss.Client
	ossInternalClient *oss.Client
)

func Init(conf *OssConfig) (err error){
	ossConfig = conf
	if ossClient == nil {
		ossClient, err = oss.New(GetOssEndpoint(), ossConfig.AccessKeyId, ossConfig.AccessKeySecret)
		if err != nil {
			logger.Errorf("oss init error: %v", err)
			return
		}
		if len(conf.InternalEndpoint) > 0 {
			if ossInternalClient == nil {
				if env.IsProd() {
					ossInternalClient, err = oss.New(GetOssInternalEndpoint(), ossConfig.DataAccessKeyId, ossConfig.DataAccessKeySecret)
					if err != nil {
						logger.Errorf("oss init error: %v", err)
						return
					}
				} else {
					ossInternalClient = ossClient
				}
			}
		}
	}
	return
}

func GetOssConfig() OssConfig {
	return *ossConfig
}

func GetOssImageBucketName() string {
	if env.IsProd() {
		return ossConfig.ImageBucket
	}
	return ossConfig.TestImageBucket
}

func GetOssAudioBucketName() string {
	if env.IsProd() {
		return ossConfig.AudioBucket
	}
	return ossConfig.TestAudioBucket
}

func GetOssVideoBucketName() string {
	if env.IsProd() {
		return ossConfig.VideoBucket
	}
	return ossConfig.TestVideoBucket
}

func GetOssDataBucketName() string {
	if env.IsProd() {
		return ossConfig.DataBucket
	}
	return ossConfig.TestDataBucket
}

func GetOssBucketName(typ consts.MediaType) string {
	switch typ {
	case consts.MediaTypeImage:
		return GetOssImageBucketName()
	case consts.MediaTypeVideo:
		return GetOssVideoBucketName()
	case consts.MediaTypeAudio:
		return GetOssAudioBucketName()
	case consts.MediaTypeData:
		return GetOssDataBucketName()
	default:
		return ""
	}
}

func GetCDNDomain(typ consts.MediaType) string {
	if env.IsProd() {
		switch typ {
		case consts.MediaTypeImage:
			return ossConfig.ImageCDNDomain
		case consts.MediaTypeVideo:
			return ossConfig.VideoCDNDomain
		case consts.MediaTypeAudio:
			return ossConfig.AudioCDNDomain
		default:
			return fmt.Sprintf("%s.%s", GetOssBucketName(typ), ossConfig.CDNDomain)
		}
	}
	return fmt.Sprintf("%s.%s", GetOssBucketName(typ), ossConfig.CDNDomain)
}

func GetOssBucketMediaUrl(typ consts.MediaType, biz consts.MediaBiz, id string) string {
	if biz.Exp > 0 {
		bucketName := GetOssBucketName(typ)
		bucket, err := GetOssBucket(bucketName)
		if err == nil {
			signedURL, err := bucket.SignURL(GetOssObjectKey(biz, id), oss.HTTPGet, time.Now().Add(biz.Exp).Unix())
			if err == nil {
				return signedURL
			}
		}
	}
	return fmt.Sprintf("%s://%s/%s/%s", ossConfig.Protocol, GetCDNDomain(typ), biz.Dir, id)
}

func GetOssBucketMediaPrivateUrl(typ consts.MediaType, biz consts.MediaBiz, id string) string {
	bucketName := GetOssBucketName(typ)
	bucket, err := GetOssBucket(bucketName)
	if err == nil {
		signedURL, err := bucket.SignURL(GetOssObjectKey(biz, id), oss.HTTPGet, time.Now().Add(time.Hour*24*3).Unix())
		if err == nil {
			return signedURL
		}
	}
	return fmt.Sprintf("%s://%s/%s/%s", ossConfig.Protocol, GetCDNDomain(typ), biz.Dir, id)
}

func GetOssObjectKey(biz consts.MediaBiz, id string) string {
	return fmt.Sprintf("%s/%s", biz.Dir, id)
}

func GetOssEndpoint() string {
	return fmt.Sprintf("%s://%s", ossConfig.Protocol, ossConfig.Endpoint)
}

func GetOssHost(bucket string) string {
	return fmt.Sprintf("%s://%s.%s", ossConfig.Protocol, bucket, ossConfig.Endpoint)
}

func MakeSignature(content string) string {
	mac := hmac.New(sha1.New, []byte(ossConfig.AccessKeySecret))
	mac.Write([]byte(content))
	return fmt.Sprintf("OSS %s:%s", ossConfig.AccessKeyId, base64.StdEncoding.EncodeToString(mac.Sum(nil)))
}

func GetOssClient() *oss.Client {
	return ossClient
}

func GetOssBucket(name string) (*oss.Bucket, error) {
	return ossClient.Bucket(name)
}

func GetOssInternalEndpoint() string {
	return fmt.Sprintf("%s://%s", ossConfig.Protocol, ossConfig.InternalEndpoint)
}

//https://help.aliyun.com/document_detail/31988.htm?spm=a2c4g.11186623.0.0.16072d50v2WB7B#section-d5z-1ww-wdb
func MakePolicyData(typ consts.MediaType, biz string) (data *OssPolicyData, err error) {
	bucket := GetOssBucketName(typ)
	expireTime := time.Now().Unix() + ossConfig.ExpireTime
	dir := fmt.Sprintf("%s/", biz)

	//create post policy json
	var config OssPolicyConfig
	config.Expiration = time.Unix(expireTime, 0).UTC().Format("2006-01-02T15:04:05Z")
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, dir)
	config.Conditions = append(config.Conditions, condition)

	configBytes, _ := json.Marshal(config)
	policy := base64.StdEncoding.EncodeToString(configBytes)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(ossConfig.AccessKeySecret))
	io.WriteString(h, policy)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	data = &OssPolicyData{
		Accessid:  ossConfig.AccessKeyId,
		Host:      GetOssHost(bucket),
		Expire:    expireTime,
		Signature: signature,
		Dir:       dir,
		Policy:    policy,
		ObjectAcl: "public-read",
	}
	return
}

func SetObjectPublicRead(typ consts.MediaType, biz consts.MediaBiz, id string) {
	bucket, err := GetOssBucket(GetOssBucketName(typ))
	if err != nil {
		logger.Errorf("GetOssBucket err:%v, %v", typ, err)
		return
	}
	err = bucket.SetObjectACL(GetOssObjectKey(biz, id), oss.ACLPublicRead)
	if err != nil {
		logger.Errorf("error %v", err)
		return
	}
	return
}

func SetObjectPrivate(typ consts.MediaType, biz consts.MediaBiz, id string) {
	bucket, err := GetOssBucket(GetOssBucketName(typ))
	if err != nil {
		logger.Errorf("GetOssBucket err:%v, %v", typ, err)
		return
	}
	bucket.SetObjectACL(GetOssObjectKey(biz, id), oss.ACLPrivate)
}

func CopyBucketObject(src, dest BucketObject) (err error) {
	bucket, err := GetOssBucket(dest.Bucket)
	if err != nil {
		logger.Errorf("CopyBucketObject %v, err: %v", dest, err)
		return
	}

	_, err = bucket.CopyObjectFrom(src.Bucket, src.Key, dest.Key)
	if err != nil {
		logger.Errorf("CopyBucketObject %v, %v, err: %v", src, dest, err)
		return
	}
	return
}

func CopyUrlObjectToBucket(url string, dest BucketObject) (err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()
	destBucket, err := GetOssBucket(dest.Bucket)
	if err != nil {
		logger.Errorf("CopyUrlObjectToBucket %v, %v, err: %v", url, dest, err)
		return
	}

	err = destBucket.PutObject(dest.Key, io.Reader(res.Body), oss.ContentType(res.Header.Get(oss.HTTPHeaderContentType)))
	if err != nil {
		logger.Errorf("CopyBucketStreamObject err: %v", err)
		return
	}
	destBucket.SetObjectACL(dest.Key, oss.ACLPublicRead)
	return
}

func ArchiveMedia(typ consts.MediaType, biz consts.MediaBiz, id string) (err error) {
	bucket, err := GetOssBucket(GetOssBucketName(typ))
	if err != nil {
		logger.Errorf("RestoreMedia err:%v, %v", typ, err)
		return
	}
	objectKey := GetOssObjectKey(biz, id)
	meta, err := bucket.GetObjectDetailedMeta(objectKey)
	if err != nil {
		logger.Errorf("error %v", err)
		return
	}
	switch meta.Get("X-Oss-Storage-Class") {
	case string(oss.StorageArchive), string(oss.StorageColdArchive):
		logger.Warnf("already archive")
		return nil
	}

	_, err = bucket.CopyObject(objectKey, objectKey, oss.ObjectStorageClass(oss.StorageArchive))
	if err != nil {
		fmt.Println(err)
		logger.Errorf("RestoreMedia err: %v", err)
		return
	}

	return
}

//需要完全控制权限
func DeleteMedia(typ consts.MediaType, biz consts.MediaBiz, id string) (err error) {
	bucket, err := GetOssBucket(GetOssBucketName(typ))
	if err != nil {
		logger.Errorf("DeleteMedia err:%v, %v", typ, err)
		return
	}
	objectKey := GetOssObjectKey(biz, id)
	delRes, err := bucket.DeleteObjects([]string{objectKey})
	if err != nil {
		fmt.Println(err)
		logger.Errorf("DeleteMedia err: %v", err)
		return
	}

	logger.Warnf("DeleteMedia: %v", delRes)
	return
}

func CheckMedia(typ consts.MediaType, objKey string) (url string, err error) {
	bucket, err := GetOssBucket(GetOssBucketName(typ))
	if err != nil {
		logger.Errorf("DeleteMedia err:%v, %v", typ, err)
		return
	}

	props, err := bucket.GetObjectDetailedMeta(objKey)
	if err != nil {
		if strings.Count(err.Error(), "StatusCode=404") > 0 {
			err = nil
			return
		}
		logger.Errorf("error GetObjectDetailedMeta: %v", err)
		return
	}
	//logger.Infof("props GetObjectDetailedMeta: %v", props)
	if props.Get("Content-Length") == "0" {
		return
	}
	if props.Get("ETag") == "" {
		return
	}

	url, err = bucket.SignURL(objKey, oss.HTTPGet, 3600)
	if err != nil {
		url = ""
		return
	}
	return
}

func UploadMedia(typ consts.MediaType, objKey string, reader io.Reader) (url string, err error) {
	bucket, err := GetOssBucket(GetOssBucketName(typ))
	if err != nil {
		return
	}

	err = bucket.PutObject(objKey, reader)
	if err != nil {
		logger.Errorf("error PutObject: %v", err)
		return
	}
	//logger.Infof("PutObject: %v", objKey)

	url, err = bucket.SignURL(objKey, oss.HTTPGet, 3600)
	if err != nil {
		url = ""
		return
	}
	return
}
