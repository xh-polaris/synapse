package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/xh-polaris/synapse/biz/domain/thirdparty/dal/model"
	"github.com/xh-polaris/synapse/biz/domain/thirdparty/entity"
	"github.com/xh-polaris/synapse/biz/domain/thirdparty/repo"
	"github.com/xh-polaris/synapse/biz/infra/contract/id"
	"github.com/xh-polaris/synapse/biz/pkg/errorx"
	"github.com/xh-polaris/synapse/biz/pkg/logs"
	"github.com/xh-polaris/synapse/biz/types/errno"
)

const (
	BUPTTicketCheckURL = "https://auth.bupt.edu.cn/authserver/p3/serviceValidate?ticket=%s&service=https%%3A%%2F%%2Finnospark.xhpolaris.com%%2Fauth%%2Fthirdparty-login%%3Fthirdparty%%3Dbupt"
)

type Component struct {
	ThirdPartyUserRepo repo.ThirdPartyUserRepo
	IdGen              id.IDGenerator
}

func NewThirdPartyDomain(ctx context.Context, c *Component) ThirdParty {
	return &thirdPartyImpl{Component: c}
}

type thirdPartyImpl struct {
	*Component
}

func (t *thirdPartyImpl) BUPTLogin(ctx context.Context, ticket string) (*entity.ThirdPartyUser, error) {
	// 获取编号
	code, err := buptExtractCode(ticket)
	if err != nil {
		logs.Error(err.Error())
		return nil, errorx.New(errno.ErrThirdPartyLogin)
	}
	// 登录或注册
	u, err := t.ThirdPartyUserRepo.FindOrCreate(ctx, t.IdGen.GenID(ctx), code, "bupt")
	if err != nil {
		return nil, errorx.New(errno.ErrThirdPartyLogin)
	}
	return thirdPartyUserPO2VO(u), nil
}

func buptExtractCode(ticket string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(BUPTTicketCheckURL, ticket))
	if err != nil {
		return "", fmt.Errorf("请求CAS服务失败: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %w", err)
	}

	var casResponse BUPTCASServiceResponse
	if err := xml.Unmarshal(body, &casResponse); err != nil {
		return "", fmt.Errorf("解析XML响应失败: %w", err)
	}

	// 处理认证失败情况
	if casResponse.AuthenticationFailure != nil {
		return "", fmt.Errorf("CAS认证失败: %s (code: %s)",
			casResponse.AuthenticationFailure.Message,
			casResponse.AuthenticationFailure.Code)
	}

	// 处理认证成功情况
	if casResponse.AuthenticationSuccess != nil {
		return casResponse.AuthenticationSuccess.User, nil
	}

	return "", fmt.Errorf("无效的CAS响应")
}

func thirdPartyUserPO2VO(user *model.ThirdPartyUser) *entity.ThirdPartyUser {
	return &entity.ThirdPartyUser{
		ID:        user.ID.Hex(),
		Code:      user.Code,
		App:       user.App,
		CreatedAt: time.UnixMilli(user.CreatedAt),
		UpdatedAt: time.UnixMilli(user.UpdatedAt),
	}
}

// BUPTCASServiceResponse CAS 认证响应结构体
type BUPTCASServiceResponse struct {
	XMLName               xml.Name                   `xml:"serviceResponse"`
	AuthenticationSuccess *BUPTAuthenticationSuccess `xml:"authenticationSuccess"`
	AuthenticationFailure *BUPTAuthenticationFailure `xml:"authenticationFailure"`
}

type BUPTAuthenticationSuccess struct {
	User       string          `xml:"user"`
	Attributes *BUPTAttributes `xml:"attributes"`
}

type BUPTAttributes struct {
	Name           string `xml:"name"`
	EmployeeNumber string `xml:"employeeNumber"`
}

type BUPTAuthenticationFailure struct {
	Code    string `xml:"code,attr"`
	Message string `xml:",chardata"`
}
