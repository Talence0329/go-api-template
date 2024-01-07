package locals

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// 統一處理API中使用的locals

// GetMemberInfo : 透過Locals取得使用者資訊
func GetMemberInfo(ctx *fiber.Ctx) (memberInfo MemberInfo) {
	if ctx.Locals(KeyMemberInfo) != nil {
		switch ctx.Locals(KeyMemberInfo).(type) {
		case MemberInfo:
			memberInfo = ctx.Locals(KeyMemberInfo).(MemberInfo)
		}
		return memberInfo
	}
	return MemberInfo{}
}

// SetUserInfo : 透過Locals設定使用者資訊
func SetUserInfo(ctx *fiber.Ctx, memberInfo MemberInfo) {
	ctx.Locals(KeyMemberInfo, memberInfo)
}

// GetJwt : 透過Locals取得jwt
func GetJwt(ctx *fiber.Ctx) *jwt.Token {
	if ctx.Locals(KeyJWTToken) != nil {
		switch ctx.Locals(KeyJWTToken).(type) {
		case *jwt.Token:
			return ctx.Locals(KeyJWTToken).(*jwt.Token)
		}
	}
	return nil
}

// SetMemberUID : 透過Locals設定MemberUID
func SetMemberUID(ctx *fiber.Ctx, uid string) {
	ctx.Locals(KeyMemberUID, uid)
}

// GetMemberUID : 透過Locals取得使用者ID
func GetMemberUID(ctx *fiber.Ctx) (memberUID string) {
	if ctx.Locals(KeyMemberUID) != nil {
		switch ctx.Locals(KeyMemberUID).(type) {
		case string:
			memberUID = ctx.Locals(KeyMemberUID).(string)
		}
		return memberUID
	}
	return ""
}
