package auth

import (
	"backend/basic/access"
	"backend/basic/apiprotocol"
	"backend/basic/cookies"
	"backend/basic/database"
	"backend/basic/locals"
	basicauth "backend/modules/auth/basicAuth"
	jwthandler "backend/modules/jwtHandler"
	systemparam "backend/modules/systemParam"
	twostep "backend/modules/twoStep"
	validatehandler "backend/modules/validateHandler"
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/lithammer/shortuuid/v4"
)

func SetRouter(router fiber.Router) {
	auth := router.Group("/auth")
	{
		auth.Post("/profile", jwthandler.New(), profileHandler)

		auth.Post("/updateProfileName", jwthandler.New(), validatehandler.New(&UpdateProfileNameReq{}), updateProfileHandler(ProfileTypeName))

		auth.Post("/updateProfileAddress", jwthandler.New(), validatehandler.New(&UpdateProfileAddressReq{}), updateProfileHandler(ProfileTypeAddress))

		auth.Post("/updateProfilePhone", jwthandler.New(), validatehandler.New(&UpdateProfilePhoneReq{}), updateProfileHandler(ProfileTypePhone))

		auth.Post("/checkToken", jwthandler.New(), checkTokenHandler)

		auth.Post("/logout", jwthandler.New(), func(c *fiber.Ctx) error {
			cookies.ClearJWT(c)
			return c.Status(fiber.StatusOK).JSON(apiprotocol.Success10000.ToRes())
		})

		auth.Post("/accountRegister", validatehandler.New(&AccountRegisterReq{}), accountRegisterHandler)

		auth.Post("/accountRegisterCheck", validatehandler.New(&AccountRegisterCheckReq{}), accountRegisterCheckHandler)

		auth.Post("/accountRegisterVerify", validatehandler.New(&AccountRegisterVerifyReq{}), accountRegisterVerifyHandler)

		auth.Post("/accountLogin", validatehandler.New(&LoginReq{}), accountLoginHandler)

		auth.Post("/googleLogin", validatehandler.New(&GoogleLoginReq{}), googleLoginHandler)

		auth.Post("/fbLogin", validatehandler.New(&FBLoginReq{}), fbLoginHandler)

		auth.Post("/changePassword", jwthandler.New(), validatehandler.New(&ChangePasswordReq{}), changePasswordHandler)

		auth.Post("/forgetPassword", validatehandler.New(&ForgetPasswordReq{}), forgetPasswordHandler)

		auth.Post("/forgetPasswordVerify", validatehandler.New(&ForgetPasswordVerifyReq{}), forgetPasswordVerifyHandler)

		auth.Post("/forgetChangePassword", validatehandler.New(&ForgetChangePasswordReq{}), forgetChangePasswordHandler)

		auth.Post("/delete", jwthandler.New(), deleteHandler)
	}
}

func profileHandler(c *fiber.Ctx) error {
	res := profileRes{}
	memberUID := locals.GetMemberUID(c)
	if memberUID == "" {
		return apiprotocol.Auth200107.ToRes()
	}

	// 取得連線
	db, err := database.MEMBER.DB()
	if err != nil {
		return apiprotocol.Database900501.ToRes().Err(err.Error()).ToErr()
	}

	if resData, err := getProfile(db, memberUID); err != nil {
		return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
	} else {
		res = resData
		if cfg.IsUseLogAll {
			res.IsUseLog = true
		}
	}

	return c.Status(fiber.StatusOK).JSON(apiprotocol.APIResponse{
		Data:      res,
		RetStatus: *apiprotocol.Success10000.ToRet(),
	})
}

func updateProfileHandler(profileType string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		memberUID := locals.GetMemberUID(c)
		if memberUID == "" {
			return apiprotocol.Auth200107.ToRes()
		}

		// 取得連線
		db, err := database.MEMBER.DB()
		if err != nil {
			return apiprotocol.Database900501.ToRes().Err(err.Error()).ToErr()
		}

		switch profileType {
		case ProfileTypeName:
			req := UpdateProfileNameReq{}
			if err := c.BodyParser(&req); err != nil {
				return apiprotocol.Request900101.ToRes().Err(err.Error()).ToErr()
			}
			if err := updateProfileName(db, memberUID, req.Name); err != nil {
				return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
			}
		case ProfileTypeAddress:
			req := UpdateProfileAddressReq{}
			if err := c.BodyParser(&req); err != nil {
				return apiprotocol.Request900101.ToRes().Err(err.Error()).ToErr()
			}
			if err := updateProfileAddress(db, memberUID, req.Address, req.Country); err != nil {
				return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
			}
		case ProfileTypePhone:
			req := UpdateProfilePhoneReq{}
			if err := c.BodyParser(&req); err != nil {
				return apiprotocol.Request900101.ToRes().Err(err.Error()).ToErr()
			}
			if err := updateProfilePhone(db, memberUID, req.Areacode, req.Phone); err != nil {
				return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
			}
		default:
			return apiprotocol.Request900101.ToRes().Err(err.Error()).ToErr()
		}

		return c.Status(fiber.StatusOK).JSON(apiprotocol.StatusResponse{
			RetStatus: *apiprotocol.Success10000.ToRet(),
		})
	}
}

func checkTokenHandler(c *fiber.Ctx) error {
	return c.JSON(apiprotocol.Success10000.ToRes())
}

func changePasswordHandler(c *fiber.Ctx) error {
	req := ChangePasswordReq{}
	if err := c.BodyParser(&req); err != nil {
		return apiprotocol.Request900101.ToRes().Err(err.Error()).ToErr()
	}

	memberUID := locals.GetMemberUID(c)
	if memberUID == "" {
		return apiprotocol.Auth200107.ToRes()
	}

	// 取得連線
	db, err := database.MEMBER.DB()
	if err != nil {
		return apiprotocol.Database900501.ToRes().Err(err.Error()).ToErr()
	}

	// 修改密碼
	statusCode, err := changePassword(db, memberUID, req.OldPassword, req.NewPassword)
	if err != nil {
		return statusCode.ToRes().Err(err.Error()).ToErr()
	}

	return c.JSON(statusCode.ToRes())
}

func forgetPasswordHandler(c *fiber.Ctx) error {
	req, res := ForgetPasswordReq{}, ForgetPasswordRes{}
	if err := c.BodyParser(&req); err != nil {
		return apiprotocol.Request900101.ToRes().Err(err.Error()).ToErr()
	}

	// 取得連線
	db, err := database.MEMBER.DB()
	if err != nil {
		return apiprotocol.Database900501.ToRes().Err(err.Error()).ToErr()
	}

	loginData, err := getLoginDataByAccount(db, req.Mail)
	if err != nil && err != sql.ErrNoRows {
		return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
	}

	if loginData.UID != "" {
		twoStepData := twostep.NewMail(TwoStepForgetPassword, req.Mail)
		res.UID = twoStepData.UID
		if err := addTwostepMail(db, twoStepData.UID, twoStepData.Code, twoStepData.Action, twoStepData.Mail); err != nil {
			return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
		}

		mailContent := fmt.Sprintf(`
		<p>Reset password verification code as follows</p>
		<p><b>%s</b></p>
	`, twoStepData.Code)
		if err := twoStepData.SendMail("8mb verification code", mailContent); err != nil {
			return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
		}
	} else {
		res.UID = shortuuid.New()
	}

	return c.Status(fiber.StatusOK).JSON(apiprotocol.APIResponse{
		Data:      res,
		RetStatus: *apiprotocol.Success10000.ToRet(),
	})
}

func forgetPasswordVerifyHandler(c *fiber.Ctx) error {
	req := ForgetPasswordVerifyReq{}
	if err := c.BodyParser(&req); err != nil {
		return apiprotocol.Request900101.ToRes().Err(err.Error()).ToErr()
	}

	// 取得連線
	db, err := database.MEMBER.DB()
	if err != nil {
		return apiprotocol.Database900501.ToRes().Err(err.Error()).ToErr()
	}

	if data, err := getTwostep(db, req.UID); err != nil {
		if err == sql.ErrNoRows {
			return apiprotocol.Auth200110.ToRes().Err(err.Error()).ToErr()
		} else {
			return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
		}
	} else {
		if data.Action != TwoStepForgetPassword {
			return apiprotocol.Auth200116.ToRes()
		}
		if data.Count >= MaxTwoStepCount {
			return apiprotocol.Auth200113.ToRes()
		}
		// 增加輸入次數
		if err := addTwostepCount(db, req.UID, data.Count+1); err != nil {
			return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
		}
		if req.Code != data.Code {
			return apiprotocol.Auth200111.ToRes()
		}
		twoStepData := twostep.TwoStepData{
			Action: TwoStepForgetPassword,
			UID:    data.UID,
			Code:   data.Code,
			Mail:   data.Mail,
		}

		if err := twoStepData.VerifySuccess(c); err != nil {
			return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
		}
	}

	return c.Status(fiber.StatusOK).JSON(apiprotocol.Success10000.ToRes())
}

func forgetChangePasswordHandler(c *fiber.Ctx) error {
	req := ForgetChangePasswordReq{}
	if err := c.BodyParser(&req); err != nil {
		return apiprotocol.Request900101.ToRes().Err(err.Error()).ToErr()
	}

	// 取得連線
	db, err := database.MEMBER.DB()
	if err != nil {
		return apiprotocol.Database900501.ToRes().Err(err.Error()).ToErr()
	}

	if twoStepData, err := twostep.GetTokenData(c, TwoStepForgetPassword); err != nil {
		return apiprotocol.Auth200110.ToRes().Err(err.Error()).ToErr()
	} else {
		if twoStepData.Action != TwoStepForgetPassword {
			return apiprotocol.Auth200116.ToRes().Err("兩段式驗證與執行的動作不符").ToErr()
		}

		accountData, err := getLoginDataByAccount(db, twoStepData.Mail)
		if err != nil && err != sql.ErrNoRows {
			return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
		}

		hashNewPassword, err := basicauth.HashPassword(req.NewPassword)
		if err != nil {
			return apiprotocol.Auth200108.ToRes().Err(err.Error()).ToErr()
		}

		if err := editAccountPassword(db, accountData.UID, hashNewPassword); err != nil {
			return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
		}
	}

	return c.Status(fiber.StatusOK).JSON(apiprotocol.Success10000.ToRes())
}

func accountLoginHandler(c *fiber.Ctx) error {
	req, res := LoginReq{}, LoginRes{}
	if err := c.BodyParser(&req); err != nil {
		return apiprotocol.Request900101.ToRes().Err(err.Error()).ToErr()
	}

	// 登入
	uid, isUseLog, statusCode, err := loginAccount(req.Account, req.Password)
	if err != nil {
		return statusCode.ToRes().Err(err.Error()).ToErr()
	}

	// 登入成功後執行jwt登入成功流程
	if err := access.LoginSuccess(c, uid); err != nil {
		return apiprotocol.JWT900601.ToRes().Err(err.Error()).ToErr()
	}
	res.IsUseLog = isUseLog

	return c.Status(fiber.StatusOK).JSON(apiprotocol.APIResponse{
		Data:      res,
		RetStatus: *statusCode.ToRet(),
	})
}

func googleLoginHandler(c *fiber.Ctx) error {
	req := GoogleLoginReq{}
	if err := c.BodyParser(&req); err != nil {
		return apiprotocol.Request900101.ToRes().Err(err.Error()).ToErr()
	}

	// 登入
	uid, statusCode, err := loginGoogle(req.Token)
	if err != nil {
		return statusCode.ToRes().Err(err.Error()).ToErr()
	}

	// 登入成功後執行jwt登入成功流程
	if err := access.LoginSuccess(c, uid); err != nil {
		return apiprotocol.JWT900601.ToRes().Err(err.Error()).ToErr()
	}

	return c.Status(fiber.StatusOK).JSON(apiprotocol.Success10000.ToRes())
}

func fbLoginHandler(c *fiber.Ctx) error {
	req := FBLoginReq{}
	if err := c.BodyParser(&req); err != nil {
		return apiprotocol.Request900101.ToRes().Err(err.Error()).ToErr()
	}

	// 登入
	uid, statusCode, err := loginFB(req.Token)
	if err != nil {
		return statusCode.ToRes().Err(err.Error()).ToErr()
	}

	// 登入成功後執行jwt登入成功流程
	if err := access.LoginSuccess(c, uid); err != nil {
		return apiprotocol.JWT900601.ToRes().Err(err.Error()).ToErr()
	}

	return c.Status(fiber.StatusOK).JSON(apiprotocol.Success10000.ToRes())
}

func accountRegisterHandler(c *fiber.Ctx) error {
	req := AccountRegisterReq{}
	if err := c.BodyParser(&req); err != nil {
		return apiprotocol.Request900101.ToRes().Err(err.Error()).ToErr()
	}

	// 取得連線
	db, err := database.MEMBER.DB()
	if err != nil {
		return apiprotocol.Database900501.ToRes().Err(err.Error()).ToErr()
	}
	tx, err := database.MEMBER.TX()
	if err != nil {
		return apiprotocol.Database900501.ToRes().Err(err.Error()).ToErr()
	}
	defer func() {
		if err := tx.Commit(); err != nil {
			fmt.Println(err)
		}
	}()

	// 檢查是否已兩段式驗證
	if systemparam.NEED_TWOSTEP.Get() == "true" {
		if twoStepData, err := twostep.GetTokenData(c, TwoStepRegister); err != nil {
			return apiprotocol.Auth200110.ToRes().Err(err.Error()).ToErr()
		} else {
			if twoStepData.Mail != req.Account {
				return apiprotocol.Auth200111.ToRes().Err("兩步驟token與欲註冊的信箱不符").ToErr()
			} else if twoStepData.Action != TwoStepRegister {
				return apiprotocol.Auth200116.ToRes().Err("兩段式驗證與執行的動作不符").ToErr()
			}
		}
	}

	// 註冊
	uid, statusCode, err := registerAccount(db, tx, req.Account, req.Password)
	if err != nil {
		return statusCode.ToRes().Err(err.Error()).ToErr()
	}

	// 登入成功後執行jwt登入成功流程
	if err := access.LoginSuccess(c, uid); err != nil {
		return apiprotocol.JWT900601.ToRes().Err(err.Error()).ToErr()
	}

	return c.Status(fiber.StatusOK).JSON(statusCode.ToRes())
}

func accountRegisterCheckHandler(c *fiber.Ctx) error {
	req, res := AccountRegisterCheckReq{}, AccountRegisterCheckRes{}
	if err := c.BodyParser(&req); err != nil {
		return apiprotocol.Request900101.ToRes().Err(err.Error()).ToErr()
	}

	twoStepData := twostep.NewMail(TwoStepRegister, req.Account)
	res.UID = twoStepData.UID

	// 取得連線
	db, err := database.MEMBER.DB()
	if err != nil {
		return apiprotocol.Database900501.ToRes().Err(err.Error()).ToErr()
	} else {
		if err := addTwostepMail(db, twoStepData.UID, twoStepData.Code, twoStepData.Action, twoStepData.Mail); err != nil {
			return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
		}
	}

	mailContent := fmt.Sprintf(`
		<p>The registration verification code</p>
		<p><b>%s</b></p>
	`, twoStepData.Code)

	if err := twoStepData.SendMail("8mb verification code", mailContent); err != nil {
		return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
	}

	return c.Status(fiber.StatusOK).JSON(apiprotocol.APIResponse{
		Data:      res,
		RetStatus: *apiprotocol.Success10000.ToRet(),
	})
}

func accountRegisterVerifyHandler(c *fiber.Ctx) error {
	req := AccountRegisterVerifyReq{}
	if err := c.BodyParser(&req); err != nil {
		return apiprotocol.Request900101.ToRes().Err(err.Error()).ToErr()
	}

	// 取得連線
	db, err := database.MEMBER.DB()
	if err != nil {
		return apiprotocol.Database900501.ToRes().Err(err.Error()).ToErr()
	}

	if data, err := getTwostep(db, req.UID); err != nil {
		if err == sql.ErrNoRows {
			return apiprotocol.Auth200110.ToRes().Err(err.Error()).ToErr()
		} else {
			return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
		}
	} else {
		if data.Action != TwoStepRegister {
			return apiprotocol.Auth200116.ToRes()
		}
		if data.Count >= MaxTwoStepCount {
			return apiprotocol.Auth200113.ToRes()
		}
		// 增加輸入次數
		if err := addTwostepCount(db, req.UID, data.Count+1); err != nil {
			return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
		}
		if req.Code != data.Code {
			return apiprotocol.Auth200111.ToRes()
		}
		twoStepData := twostep.TwoStepData{
			Action: TwoStepRegister,
			UID:    data.UID,
			Code:   data.Code,
			Mail:   data.Mail,
		}

		if err := twoStepData.VerifySuccess(c); err != nil {
			return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
		}
	}

	return c.Status(fiber.StatusOK).JSON(apiprotocol.Success10000.ToRes())
}

func deleteHandler(c *fiber.Ctx) error {
	memberUID := locals.GetMemberUID(c)
	if memberUID == "" {
		return apiprotocol.Auth200107.ToRes()
	}

	// 取得連線
	if db, err := database.MEMBER.DB(); err != nil {
		return apiprotocol.Database900501.ToRes().Err(err.Error()).ToErr()
	} else {
		if err := updateStatus(db, memberUID, MemberStatusDelete); err != nil {
			return apiprotocol.Database900502.ToRes().Err(err.Error()).ToErr()
		}
	}

	return c.Status(fiber.StatusOK).JSON(apiprotocol.Success10000.ToRes())
}
