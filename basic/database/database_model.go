package database

import (
	dbconn "backend/modules/dbConn"
)

const EnvkeyMember dbconn.Envkey = "MEMBER_POSTGRES_URL"
const EnvkeyBackstage dbconn.Envkey = "BACKSTAGE_POSTGRES_URL"

// postgres使用的參數設定
const MEMBER dbconn.DBName = "member"
const BACKSTAGE dbconn.DBName = "backstage"
