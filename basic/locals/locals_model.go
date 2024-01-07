package locals

// Key : 使用於locals的key
type LocalsKey = string

// KeyMemberUID : memberUID使用的KeyNames
const KeyMemberUID LocalsKey = "memberUID"

// KeyMemberInfo : memberInfo使用的KeyNames
const KeyMemberInfo LocalsKey = "memberInfo"

// KeyJWTToken : 於存放jwt的token的key name
const KeyJWTToken LocalsKey = "jwtToken"

type MemberInfo struct {
	UID      string
	GroupUID string
}
type MemberStatus struct {
	UID      string
	GroupUID string
	Flag     string
}
