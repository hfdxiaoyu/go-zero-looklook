package logic

import (
	"context"

	"looklook/app/usercenter/cmd/rpc/internal/svc"
	"looklook/app/usercenter/cmd/rpc/usercenter"
	"looklook/app/usercenter/model"
	"looklook/common/tool"
	"looklook/common/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

var ErrGenerateTokenError = xerr.NewErrMsg("生成token失败")
var ErrUsernamePwdError = xerr.NewErrMsg("账号或密码不正确")

// Login 函数接受 LoginReq 结构体，根据身份验证类型对用户进行身份验证，使用GenerateTokenLogic 生成令牌，并返回带有访问令牌、访问过期时间和刷新时间的 LoginResp
// 结构体。
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// `Login` 函数负责处理登录逻辑。它采用“LoginReq”对象作为输入，其中包含身份验证类型、身份验证密钥（例如手机号码）和密码。
func (l *LoginLogic) Login(in *usercenter.LoginReq) (*usercenter.LoginResp, error) {

	var userId int64
	var err error
	switch in.AuthType {
	case model.UserAuthTypeSystem:
		userId, err = l.loginByMobile(in.AuthKey, in.Password)
	default:
		return nil, xerr.NewErrCode(xerr.SERVER_COMMON_ERROR)
	}
	if err != nil {
		return nil, err
	}

	//2、Generate the token, so that the service doesn't call rpc internally
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&usercenter.GenerateTokenReq{
		UserId: userId,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "GenerateToken userId : %d", userId)
	}

	return &usercenter.LoginResp{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}

// “loginByMobile”函数负责验证用户的手机号码和密码。它使用手机号码从 UserModel 检索用户信息并检查密码是否匹配。如果手机号码或密码不正确，则返回错误。如果验证成功，则返回用户的
// ID。
func (l *LoginLogic) loginByMobile(mobile, password string) (int64, error) {

	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, mobile)
	if err != nil && err != model.ErrNotFound {
		return 0, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "根据手机号查询用户信息失败，mobile:%s,err:%v", mobile, err)
	}
	if user == nil {
		return 0, errors.Wrapf(ErrUserNoExistsError, "mobile:%s", mobile)
	}

	if !(tool.Md5ByString(password) == user.Password) {
		return 0, errors.Wrap(ErrUsernamePwdError, "密码匹配出错")
	}

	return user.Id, nil
}

func (l *LoginLogic) loginBySmallWx() error {
	return nil
}
