package auth

import (
	"backend/basic/apiprotocol"
	"backend/basic/database"
	basicauth "backend/modules/auth/basicAuth"
	fbauth "backend/modules/auth/fbAuth"
	googleauth "backend/modules/auth/googleAuth"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lithammer/shortuuid/v4"
)

// Init : 初始化，設定domain和圖片路徑
func Init(initConfig Config) error {
	cfg = initConfig
	if len(cfg.PasswordHash) == 0 {
		return errors.New("passwordHash 長度不可小於等於0")
	}

	return nil
}

// loginAccount : 登入流程
func loginAccount(account string, password string) (string, bool, apiprotocol.Code, error) {
	loginUID, isUseLog := "", false
	// 取得連線
	db, err := database.MEMBER.DB()
	if err != nil {
		return loginUID, isUseLog, apiprotocol.Database900501, err
	}

	// 抓取使用者資料並確認是否有該使用者
	if loginData, err := getLoginDataByAccount(db, account); err != nil {
		return loginUID, isUseLog, apiprotocol.Auth200101, err
	} else if loginData.UID != "" {
		switch loginData.Status {
		case MemberStatusSuspend:
			return loginUID, isUseLog, apiprotocol.Auth200103, errors.New("該會員已被停用")
		case MemberStatusDelete:
			return loginUID, isUseLog, apiprotocol.Auth200115, errors.New("該會員已刪除")
		}
		// loginCheck
		if statusCode, err := loginAccountCheck(password, loginData.Password); err != nil {
			// 密碼比對錯誤
			// 超過最大次數會直接禁用
			if loginData.ErrorCount >= MaxErrorCount {
				if err := updateStatus(db, loginData.UID, MemberStatusSuspend); err != nil {
					return loginUID, isUseLog, apiprotocol.Database900502, err
				}
				return loginUID, isUseLog, statusCode, errors.New("已超過錯誤次數，帳號被鎖定")
			}

			// 每次失敗會+1
			if err := updateErrorCount(db, loginData.UID, loginData.ErrorCount+1); err != nil {
				return loginUID, isUseLog, apiprotocol.Database900502, err
			}

			return loginUID, isUseLog, statusCode, err
		} else {
			// 密碼比對成功
			loginUID = loginData.UID

			// 失敗數清空
			if err := updateErrorCount(db, loginData.UID, 0); err != nil {
				return loginUID, isUseLog, apiprotocol.Database900502, err
			}
		}
	} else {
		// 查無使用此帳號的會員
		return loginUID, isUseLog, apiprotocol.Auth200101, errors.New("查無使用此帳號的會員")
	}

	if cfg.IsUseLogAll {
		isUseLog = true
	} else {
		if resData, err := getProfile(db, loginUID); err != nil {
			return loginUID, isUseLog, apiprotocol.Database900502, err
		} else {
			isUseLog = resData.IsUseLog
		}
	}

	return loginUID, isUseLog, apiprotocol.Success10000, nil
}

// loginAccountCheck
func loginAccountCheck(inputPassword string, password string) (apiprotocol.Code, error) {
	// checkPassword
	if err := basicauth.Check(inputPassword, password); err != nil {
		return apiprotocol.Auth200104, err
	}

	return apiprotocol.Success10000, nil
}

// loginGoogle
func loginGoogle(token string) (string, apiprotocol.Code, error) {
	var scopeData googleauth.ScopeEmailData
	// 取得連線
	db, err := database.MEMBER.DB()
	if err != nil {
		return "", apiprotocol.Database900501, err
	}

	// 向Google取得會員Google資訊
	if _scopeData, err := googleauth.GetGoogleUserInfo(token); err != nil {
		return "", apiprotocol.Auth200101, err
	} else {
		scopeData = _scopeData
	}

	// 抓取使用者資料並確認是否有該使用者
	loginData, err := getLoginDataByGoogleID(db, scopeData.ID)
	if err != nil && err == sql.ErrNoRows {
		// 無使用者的錯誤，可以自動創帳號
		// 取得連線
		tx, err := database.MEMBER.TX()
		if err != nil {
			return "", apiprotocol.Database900501, err
		}
		defer func() {
			if err := tx.Commit(); err != nil {
				fmt.Println(err)
			}
		}()
		if newUUID, err := registerGoogle(tx, scopeData); err != nil {
			return "", apiprotocol.Auth200105, err
		} else {
			return newUUID, apiprotocol.Success10000, nil
		}
	} else if err != nil {
		return "", apiprotocol.Auth200101, err
	}

	return loginData.UID, apiprotocol.Success10000, nil
}

// loginFB
func loginFB(token string) (string, apiprotocol.Code, error) {
	var profileData fbauth.ProfileData
	// 取得連線
	db, err := database.MEMBER.DB()
	if err != nil {
		return "", apiprotocol.Database900501, err
	}

	// 向FB取得會員FB資訊
	if _profileData, err := fbauth.GetFBUserInfo(token); err != nil {
		return "", apiprotocol.Auth200101, err
	} else {
		profileData = _profileData
	}

	// 抓取使用者資料並確認是否有該使用者
	loginData, err := getLoginDataByFBID(db, profileData.ID)
	if err != nil && err == sql.ErrNoRows {
		// 無使用者的錯誤，可以自動創帳號
		// 取得連線
		tx, err := database.MEMBER.TX()
		if err != nil {
			return "", apiprotocol.Database900501, err
		}
		defer func() {
			if err := tx.Commit(); err != nil {
				fmt.Println(err)
			}
		}()
		if newUUID, err := registerFB(tx, profileData); err != nil {
			return "", apiprotocol.Auth200105, err
		} else {
			return newUUID, apiprotocol.Success10000, nil
		}
	} else if err != nil {
		return "", apiprotocol.Auth200101, err
	}

	return loginData.UID, apiprotocol.Success10000, nil
}

// registerGoogle : google帳號註冊，先新增member再新增googleauth
func registerGoogle(tx *sql.Tx, scopeData googleauth.ScopeEmailData) (string, error) {
	newUUID := shortuuid.New()
	newMemberName := strings.Split(scopeData.Email, "@")[0]
	if err := addMember(tx, newUUID, newMemberName, scopeData.Email); err != nil {
		tx.Rollback()
		return "", err
	}
	if err := addGoogleAuth(tx, newUUID, scopeData); err != nil {
		tx.Rollback()
		return "", err
	}

	return newUUID, nil
}

// registerFB : fb帳號註冊，先新增member再新增fbauth
func registerFB(tx *sql.Tx, profileData fbauth.ProfileData) (string, error) {
	newUUID := shortuuid.New()
	newMemberName := strings.Split(profileData.Email, "@")[0]
	if err := addMember(tx, newUUID, newMemberName, profileData.Email); err != nil {
		tx.Rollback()
		return "", err
	}
	if err := addFBAuth(tx, newUUID, profileData); err != nil {
		tx.Rollback()
		return "", err
	}

	return newUUID, nil
}

// registerAccount : 帳號註冊，先檢查是否有重複account、先新增member再新增account
func registerAccount(db *sql.DB, tx *sql.Tx, account string, passsword string) (string, apiprotocol.Code, error) {
	loginData, err := getLoginDataByAccount(db, account)
	if err != nil && err != sql.ErrNoRows {
		return "", apiprotocol.Database900502, err
	}
	if loginData.UID != "" {
		return "", apiprotocol.Auth200103, errors.New("此帳號已被註冊")
	}

	newUUID := shortuuid.New()

	if err := addMember(tx, newUUID, account, account); err != nil {
		tx.Rollback()
		return "", apiprotocol.Database900502, err
	}

	hashPassword, err := basicauth.HashPassword(passsword)
	if err != nil {
		tx.Rollback()
		return "", apiprotocol.Database900502, err
	}
	if err := addAccountAuth(tx, newUUID, account, hashPassword); err != nil {
		tx.Rollback()
		return "", apiprotocol.Auth200105, err
	}

	return newUUID, apiprotocol.Success10000, nil
}

// changePassword : 修改密碼
func changePassword(db *sql.DB, uid string, oldPasssword string, newPasssword string) (apiprotocol.Code, error) {
	loginData, err := getAccountDataByUID(db, uid)
	if err != nil && err != sql.ErrNoRows {
		return apiprotocol.Database900502, err
	}
	if loginData.UID == "" {
		return apiprotocol.Auth200107, errors.New("無法取得使用者資訊")
	}

	if err := basicauth.Check(oldPasssword, loginData.Password); err != nil {
		return apiprotocol.Auth200104, err
	}

	hashNewPassword, err := basicauth.HashPassword(newPasssword)
	if err != nil {
		return apiprotocol.Auth200108, err
	}

	if err := editAccountPassword(db, uid, hashNewPassword); err != nil {
		return apiprotocol.Database900502, err
	}

	return apiprotocol.Success10000, nil
}
