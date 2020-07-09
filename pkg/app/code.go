package app

// 返回自定义状态码
const (
	Success            = 0
	PermissionDenied   = 403
	NotFound           = 404
	ManyRequest        = 429
	Fail               = 500
	AuthFail           = 401
	PerNotAllowed      = 605
	AdminOwnedRole     = 606
	RoleNotExisit      = 607
	UserNotExisit      = 608
	EParaIllegal       = 799
	AccessToken        = 40001
	GameSyncErr        = 40010
	IDCardVerification = 40011 // 实名认证
)

// 自定义的一些错误消息的返回
const (
	AuthFailMessage         = "认证已失效，请重新认证"
	PermissionDeniedMessage = "PermissionDenied"
	ParamError              = "参数错误"
	CaptchaError            = "验证码错误"
	ErrUserNotFound         = "未找到用户"
	ErrPasswordError        = "密码错误"
	ErrUserLogin            = "登陆失败"
	ErrOperational          = "操作失败"
	ErrDataNotFind          = "数据不存在"
	ErrGiftDateNotFinish    = "礼包时限参数不完整"
	ErrGiftDateFail         = "礼包时限配置时间错误"
	ErrGiftDateExcelFail    = "上架礼包无法删除对应文件"
	GameSyncErrMsg          = "游戏区服同步失败"
	ErrGiftExcelNotDeleted  = "请先删除上传的礼包文件excel再更新"
	ErrGifStatus            = "请先下架礼包，才能进行编辑更新"
	ErrGameStatus           = "请先下架游戏，才能进行编辑更新"
)

// CodeMsg 错误编码映射
var CodeMsg = map[int]string{

}
