package auth

import (
	"backend/basic/database"
	fbauth "backend/modules/auth/fbAuth"
	googleauth "backend/modules/auth/googleAuth"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

// getProfile : 以帳號取得使用者登入用的資料
func getProfile(db *sql.DB, uid string) (ProfileData, error) {
	pssq := database.PSSQ()
	data := ProfileData{}

	where := sq.Eq{
		"mm.uid": uid,
	}
	query := pssq.Select(`
		mm.name,
		mm.address,
		mm.email,
		mm.areacode,
		mm.phone,
		mm.country,
		ag.gid,
		af.fid,
		um.islogging
	`).From(`member as mm
	`).LeftJoin(`auth_account as aa ON mm.uid = aa.uid
	`).LeftJoin(`auth_google as ag ON mm.uid = ag.uid
	`).LeftJoin(`auth_fb as af ON mm.uid = af.uid
	`).LeftJoin(`uselog_member as um ON mm.uid = um.member_uid
	`).Where(where)
	sqll, args, sqErr := query.ToSql()
	if sqErr != nil {
		return data, sqErr
	}
	rawGID := sql.NullString{}
	rawFID := sql.NullString{}
	rawAddress := sql.NullString{}
	rawAreacode := sql.NullString{}
	rawPhone := sql.NullString{}
	rawCountry := sql.NullString{}
	rawIslogging := sql.NullBool{}
	switch err := db.QueryRow(sqll, args...).Scan(
		&data.Name,
		&rawAddress,
		&data.Email,
		&rawAreacode,
		&rawPhone,
		&rawCountry,
		&rawGID,
		&rawFID,
		&rawIslogging,
	); err {
	case sql.ErrNoRows:
		return data, err
	case nil:
		data.HasGoogle = rawGID.Valid
		data.HasFB = rawFID.Valid
		if rawAddress.Valid {
			data.Address = rawAddress.String
		}
		if rawAreacode.Valid {
			data.Areacode = rawAreacode.String
		}
		if rawPhone.Valid {
			data.Phone = rawPhone.String
		}
		if rawCountry.Valid {
			data.Country = rawCountry.String
		}
		if rawIslogging.Valid {
			data.IsUseLog = rawIslogging.Bool
		} else {
			data.IsUseLog = false
		}
		return data, nil
	default:
		return data, err
	}
}

// updateProfileName : 修改名稱
func updateProfileName(db *sql.DB, uid string, name string) error {
	pssq := database.PSSQ()
	query := pssq.Update("member").Set("name", name)

	sqll, args, sqErr := query.Where(sq.Eq{"uid": uid}).ToSql()
	if sqErr != nil {
		return sqErr
	}

	_, execErr := db.Exec(sqll, args...)

	return execErr
}

// updateProfileAddress : 修改住址與地區
func updateProfileAddress(db *sql.DB, uid, address, country string) error {
	pssq := database.PSSQ()
	query := pssq.Update("member").Set("address", address).Set("country", country)

	sqll, args, sqErr := query.Where(sq.Eq{"uid": uid}).ToSql()
	if sqErr != nil {
		return sqErr
	}

	_, execErr := db.Exec(sqll, args...)

	return execErr
}

// updateProfilePhone : 修改資料
func updateProfilePhone(db *sql.DB, uid, areacode, phone string) error {
	pssq := database.PSSQ()
	query := pssq.Update("member").Set("areacode", areacode).Set("phone", phone)

	sqll, args, sqErr := query.Where(sq.Eq{"uid": uid}).ToSql()
	if sqErr != nil {
		return sqErr
	}

	_, execErr := db.Exec(sqll, args...)

	return execErr
}

// getLoginDataByAccount : 以帳號取得使用者登入用的資料
func getLoginDataByAccount(db *sql.DB, account string) (accountLoginData, error) {
	pssq := database.PSSQ()
	data := accountLoginData{}

	where := sq.Eq{
		"aa.account": account,
	}
	query := pssq.Select(`
		m.uid,
		aa.account,
		aa.password,
		aa.errcount,
		m.status
	`).From(`member as m
	`).LeftJoin(`auth_account as aa ON m.uid = aa.uid
	`).Where(where)
	sqll, args, sqErr := query.ToSql()
	if sqErr != nil {
		return data, sqErr
	}

	switch err := db.QueryRow(sqll, args...).Scan(
		&data.UID,
		&data.Account,
		&data.Password,
		&data.ErrorCount,
		&data.Status,
	); err {
	case sql.ErrNoRows:
		return data, err
	case nil:
		return data, nil
	default:
		return data, err
	}
}

// getLoginDataByGoogleID : 以Google取得使用者資料
func getLoginDataByGoogleID(db *sql.DB, gid string) (googleLoginData, error) {
	pssq := database.PSSQ()
	data := googleLoginData{}

	where := sq.Eq{
		"ga.gid": gid,
	}
	query := pssq.Select(`
		mm.uid
	`).From(`member as mm
	`).LeftJoin(`google_auth as ga ON mm.uid = ga.uid
	`).Where(where)
	sqll, args, sqErr := query.ToSql()
	if sqErr != nil {
		return data, sqErr
	}

	switch err := db.QueryRow(sqll, args...).Scan(
		&data.UID,
	); err {
	case sql.ErrNoRows:
		return data, err
	case nil:
		return data, nil
	default:
		return data, err
	}
}

// getLoginDataByFBID : 以FB取得使用者資料
func getLoginDataByFBID(db *sql.DB, fid string) (fbLoginData, error) {
	pssq := database.PSSQ()
	data := fbLoginData{}

	where := sq.Eq{
		"fa.fid": fid,
	}
	query := pssq.Select(`
		mm.uid
	`).From(`member as mm
	`).LeftJoin(`auth_fb as fa ON mm.uid = fa.uid
	`).Where(where)
	sqll, args, sqErr := query.ToSql()
	if sqErr != nil {
		return data, sqErr
	}

	switch err := db.QueryRow(sqll, args...).Scan(
		&data.UID,
	); err {
	case sql.ErrNoRows:
		return data, err
	case nil:
		return data, nil
	default:
		return data, err
	}
}

// addMember : 新增會員資料
func addMember(tx *sql.Tx, newUid string, name string, email string) error {
	pssq := database.PSSQ()
	sqll, args, sqErr := pssq.Insert("member").Columns("uid", "name", "email").Values(newUid, name, email).ToSql()
	if sqErr != nil {
		return sqErr
	}

	_, execErr := tx.Exec(sqll, args...)
	if execErr != nil {
		return execErr
	}

	return execErr
}

// getAccountDataByUID : 以uid取得使用者登入用的帳號資料
func getAccountDataByUID(db *sql.DB, uid string) (accountLoginData, error) {
	pssq := database.PSSQ()
	data := accountLoginData{}

	where := sq.Eq{
		"ma.uid": uid,
	}
	query := pssq.Select(`
		mm.uid,
		ma.account,
		ma.password
	`).From(`member as mm
	`).LeftJoin(`auth_account as ma ON mm.uid = ma.uid
	`).Where(where)
	sqll, args, sqErr := query.ToSql()
	if sqErr != nil {
		return data, sqErr
	}

	switch err := db.QueryRow(sqll, args...).Scan(
		&data.UID,
		&data.Account,
		&data.Password,
	); err {
	case sql.ErrNoRows:
		return data, err
	case nil:
		return data, nil
	default:
		return data, err
	}
}

// editAccountPassword : 修改帳號密碼
func editAccountPassword(db *sql.DB, uid string, password string) error {
	pssq := database.PSSQ()
	sqll, args, sqErr := pssq.Update("auth_account").SetMap(map[string]interface{}{
		"password": password,
	}).Where(sq.Eq{"uid": uid}).ToSql()
	if sqErr != nil {
		return sqErr
	}

	_, execErr := db.Exec(sqll, args...)
	if execErr != nil {
		return execErr
	}

	return execErr
}

// updateErrorCount : 修改錯誤次數
func updateErrorCount(db *sql.DB, uid string, count int) error {
	pssq := database.PSSQ()
	sqll, args, sqErr := pssq.Update("auth_account").SetMap(map[string]interface{}{
		"errcount": count,
	}).Where(sq.Eq{"uid": uid}).ToSql()
	if sqErr != nil {
		return sqErr
	}

	_, execErr := db.Exec(sqll, args...)
	if execErr != nil {
		return execErr
	}

	return execErr
}

// updateStatus : 修改狀態
func updateStatus(db *sql.DB, uid string, status string) error {
	pssq := database.PSSQ()
	query := pssq.Update("member").
		Set("status", status)

	sqll, args, sqErr := query.Where(sq.Eq{"uid": uid}).ToSql()
	if sqErr != nil {
		return sqErr
	}

	_, execErr := db.Exec(sqll, args...)

	return execErr
}

// addGoogleAuth : 新增會員Google資料
func addGoogleAuth(tx *sql.Tx, uid string, scopeData googleauth.ScopeEmailData) error {
	pssq := database.PSSQ()
	sqll, args, sqErr := pssq.Insert("auth_google").Columns("uid", "gid", "gmail", "picture").Values(uid, scopeData.ID, scopeData.Email, scopeData.Picture).ToSql()
	if sqErr != nil {
		return sqErr
	}

	_, execErr := tx.Exec(sqll, args...)

	return execErr
}

// addFBAuth : 新增會員FB資料
func addFBAuth(tx *sql.Tx, uid string, profileData fbauth.ProfileData) error {
	pssq := database.PSSQ()
	sqll, args, sqErr := pssq.Insert("auth_fb").Columns("uid", "fid", "email", "picture").Values(uid, profileData.ID, profileData.Email, profileData.Picture).ToSql()
	if sqErr != nil {
		return sqErr
	}

	_, execErr := tx.Exec(sqll, args...)

	return execErr
}

// addAccountAuth : 新增會員Account資料
func addAccountAuth(tx *sql.Tx, uid string, account string, password string) error {
	pssq := database.PSSQ()
	sqll, args, sqErr := pssq.Insert("auth_account").Columns("uid", "account", "password").Values(uid, account, password).ToSql()
	if sqErr != nil {
		return sqErr
	}

	_, execErr := tx.Exec(sqll, args...)

	return execErr
}

// addTwostepMail : 新增兩步驟驗證信件資料
func addTwostepMail(db *sql.DB, uid string, code string, action string, mail string) error {
	pssq := database.PSSQ()
	sqll, args, sqErr := pssq.Insert("twostep").Columns("uid", "code", "action", "mail").Values(uid, code, action, mail).ToSql()
	if sqErr != nil {
		return sqErr
	}

	_, execErr := db.Exec(sqll, args...)

	return execErr
}

// getTwostep : 取得兩步驟驗證信件資料
func getTwostep(db *sql.DB, uid string) (twoStepMailData, error) {
	pssq := database.PSSQ()
	data := twoStepMailData{
		UID: uid,
	}

	query := pssq.Select(`
		code,
		action,
		count,
		mail
	`).From(`twostep
	`).Where(sq.Eq{
		"uid": uid,
	})
	sqll, args, sqErr := query.ToSql()
	if sqErr != nil {
		return data, sqErr
	}

	switch err := db.QueryRow(sqll, args...).Scan(
		&data.Code,
		&data.Action,
		&data.Count,
		&data.Mail,
	); err {
	case sql.ErrNoRows:
		return data, err
	case nil:
		return data, nil
	default:
		return data, err
	}
}

// addTwostepCount : 兩步驟驗證信件增加嘗試次數
func addTwostepCount(db *sql.DB, uid string, count int) error {
	pssq := database.PSSQ()
	sqll, args, sqErr := pssq.Update("twostep").SetMap(map[string]interface{}{
		"count": count,
	}).Where(sq.Eq{"uid": uid}).ToSql()
	if sqErr != nil {
		return sqErr
	}

	_, execErr := db.Exec(sqll, args...)
	if execErr != nil {
		return execErr
	}

	return execErr
}
