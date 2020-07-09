package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"word/pkg/app"
	"word/pkg/database"
)

var (
	jwtConf jwtConfig
	jwtTokenSign = loadConfig()
)

// Claims 生成token的结构体
type (
	Claims struct {
		ID        interface{} // 唯一id
		CheckData string      // 验证信息
		jwt.StandardClaims
	}
	jwtConfig struct {
		Token string `toml:"jwt_token" mapstructure:"jwt_token" env:"JWT_TOKEN"`
	}
)

func loadConfig() []byte{
	app.InitConfig("application", &jwtConf)
	return []byte(jwtConf.Token)
}

// AuthInterface 参与 jwt 数据表结构体需要实现这些接口
type AuthInterface interface {
	database.Table
	GetTopic() interface{}                         // 返回唯一信息
	FindByTopic(topic interface{}) AuthInterface   // 根据唯一信息标识获取数据信息, 比如根据用户id获取用户信息,需要注意传入的数据类型
	GetCheckData() string                          // 获取验证信息, jwt加密时, 改信息会一起进行加密, 解密时会解出来然后调用 Check 验证该信息的正确性, 如果是其他数据类型直接转string，比如是个结构体或者map, 直接转为json string
	Check(ctx *gin.Context, checkData string) bool // 验证信息
	ExpiredAt() int64                              // 返回过期时间,时间戳
}

// Auth jwt认证对象
type Auth struct {
	// 在整个gin.Context 上线文中的 Get 操作的key名,可以获得 AuthEntity
	ContextKey string
	// JwtHeaderKey jwt token 在HTTP请求中的header名
	HeaderKey string
	// 默认为 false, 如果为 true , 将验证不通过后也会继续往下执行
	Continue   bool
	AuthEntity AuthInterface
}

// VerifyAuthAndContinue 验证用户有效性并继续往下执行
func VerifyAuthAndContinue(contextKey, headerKey string, AuthEntity AuthInterface) gin.HandlerFunc {
	var auth = &Auth{AuthEntity: AuthEntity, ContextKey: contextKey, HeaderKey: headerKey, Continue: true}
	return auth.verifyAuth
}

// VerifyAuth 验证用户有效性
func VerifyAuth(contextKey, headerKey string, AuthEntity AuthInterface) gin.HandlerFunc {
	var auth = &Auth{AuthEntity: AuthEntity, ContextKey: contextKey, HeaderKey: headerKey}
	return auth.verifyAuth
}

// VerifyCustomAuth 自定义
func VerifyCustomAuth(auth *Auth) gin.HandlerFunc {
	return auth.verifyAuth
}

// verifyAuth 验证用户有效性
func (auth *Auth) verifyAuth(c *gin.Context) {
	token := c.GetHeader(auth.HeaderKey)
	if token != "" {
		claims, err := ParseToken(token)
		if err == nil {
			var entity = auth.AuthEntity.FindByTopic(claims.ID)
			if entity.Check(c, claims.CheckData) {
				c.Set(auth.ContextKey, entity) // 向下设置用户信息,控制器可直接获取
				c.Header(auth.HeaderKey, token)
				c.Next()
				return
			}

			if !auth.Continue {
				app.NewResponse(app.AuthFail, "auth failed", app.AuthFailMessage).End(c, http.StatusUnauthorized)
				c.Abort()
				return
			}
			c.Next()
			return
		}
	}

	if !auth.Continue {
		app.NewResponse(app.AuthFail, "auth failed", app.AuthFailMessage).End(c, http.StatusUnauthorized)
		c.Abort()
		return
	}
	c.Next()
}

func newClaims(entity AuthInterface) Claims {
	return Claims{
		entity.GetTopic(),
		entity.GetCheckData(),
		jwt.StandardClaims{
			ExpiresAt: entity.ExpiredAt(),
		},
	}
}

// NewToken 根据传入的结构体(非空结构体)返回一个token
func NewToken(entity AuthInterface) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims(entity))
	rs, err := token.SignedString(jwtTokenSign)
	if err != nil {
		return "", err
	}
	return rs, nil
}

// ParseToken 根据传入 token 得到 Claims 信息
func ParseToken(sign string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(sign, &Claims{}, func(token *jwt.Token) (i interface{}, e error) {
		return jwtTokenSign, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("can't decode token info")
}
