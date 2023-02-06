package middlewares

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xhigher/hzgo/config"
	"github.com/xhigher/hzgo/defines"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/resp"
	"github.com/xhigher/hzgo/utils"
)

const (
	BaseParamsHeaderKey   = "X-BaseParams"
)

var (
	ErrEmptyBaseParamsHeader = errors.New("base params header is empty")
)

type SecSign struct {

	signature *utils.Signature
}


func NewSecSign(conf *config.SecConfig) *SecSign {
	if len(conf.SignSecret) == 0 {
		panic(ErrMissingSecretKey)
	}


	signature := utils.NewSignature(conf.SignSecret, conf.SignKeyName)
	signature.UseReflect(true)
	mw := &SecSign{
		signature: signature,
	}
	return mw
}

func (mw *SecSign) Verify() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		baseParams, err := mw.getBaseParamsFromHeader(ctx, c)
		if err != nil {
			logger.Errorf("getBaseParamsFromHeader error: %v", err)
			resp.ReplyErrorIllegal(c)
			return
		}
		logger.Infof("getBaseParamsFromHeader params: %v", baseParams)
		if len(baseParams.Sign) != 32 {
			resp.ReplyErrorIllegal(c)
			return
		}

		ok, err := mw.signature.Verify(baseParams)
		if err != nil || !ok {
			logger.Errorf("getBaseParamsFromHeader error: %v", err)
			resp.ReplyErrorIllegal(c)
			return
		}

		mw.setBaseParams(c, baseParams)

		c.Next(ctx)
	}
}

func (mw *SecSign) getBaseParamsFromHeader(ctx context.Context, c *app.RequestContext) (params *defines.BaseParams, err error) {
	baseParamsHeader := c.Request.Header.Get(BaseParamsHeaderKey)

	if len(baseParamsHeader)== 0{
		err =ErrEmptyBaseParamsHeader
		return
	}

	bytes, err := base64.StdEncoding.DecodeString(baseParamsHeader)
	if err != nil{
		return
	}

	err = json.Unmarshal(bytes, &params)
	if err != nil{
		return
	}
	return
}

func (mw *SecSign) setBaseParams(c *app.RequestContext, params *defines.BaseParams) {
	c.Set("BASE_PARAMS", params)
}

func (mw *SecSign) GetBaseParams(c *app.RequestContext) *defines.BaseParams {
	if params, ok := c.Get("BASE_PARAMS"); ok {
		return params.(*defines.BaseParams)
	}
	return nil
}
