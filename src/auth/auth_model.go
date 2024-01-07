package auth

var cfg Config

const TwoStepRegister = "register"
const TwoStepForgetPassword = "forgetPassword"
const MaxTwoStepCount = 5
const MaxErrorCount = 5
const MemberStatusNormal = "N"
const MemberStatusSuspend = "S"
const MemberStatusDelete = "D"

const ProfileTypeName = "name"
const ProfileTypeAddress = "address"
const ProfileTypePhone = "phone"

type Config struct {
	IsTwoStep     bool
	IsUseLogAll   bool
	PasswordHash  string
	TempJWTSecret string
}

type (
	profileRes  = ProfileData
	ProfileData struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		Address   string `json:"address"`
		Areacode  string `json:"areacode"`
		Phone     string `json:"phone"`
		Country   string `json:"country"`
		HasGoogle bool   `json:"hasGoogle"`
		HasFB     bool   `json:"hasFB"`
		IsUseLog  bool   `json:"isUseLog"`
	}
)

type (
	UpdateProfileNameReq struct {
		Name string `json:"name" validate:"required,excludesall={}!?=\\"`
	}
	UpdateProfileAddressReq struct {
		Address string `json:"address" validate:"required,excludesall={}!?=\\"`
		Country string `json:"country" validate:"required,excludesall={}!?=\\"`
	}
	UpdateProfilePhoneReq struct {
		Areacode string `json:"areacode" validate:"required,startswith=+,containsany=number+,endsnotwith=+,max=4"`
		Phone    string `json:"phone" validate:"required,number,max=15"`
	}
)

type (
	LoginReq struct {
		Account  string `json:"account" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	LoginRes struct {
		IsUseLog bool `json:"isUseLog"`
	}
	accountLoginData struct {
		UID        string
		Account    string
		Password   string
		ErrorCount int
		Status     string
	}
)

type (
	GoogleLoginReq struct {
		Token string `json:"token"`
	}
	googleLoginData struct {
		UID string
	}
)

type (
	FBLoginReq struct {
		Token string `json:"token" validate:"required"`
	}
	fbLoginData struct {
		UID string
	}
)

type (
	AccountRegisterReq struct {
		Account  string `json:"account" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	AccountRegisterCheckReq struct {
		Account string `json:"account" validate:"required,email"`
	}
	AccountRegisterCheckRes struct {
		UID string `json:"uid"`
	}
	AccountRegisterVerifyReq struct {
		UID  string `json:"uid" validate:"required"`
		Code string `json:"code" validate:"required"`
	}
	twoStepMailData struct {
		UID    string
		Code   string
		Count  int
		Action string
		Mail   string
	}
)

type (
	ChangePasswordReq struct {
		OldPassword string `json:"oldPassword" validate:"required"`
		NewPassword string `json:"newPassword" validate:"required"`
	}
)

type (
	ForgetPasswordReq struct {
		Mail string `json:"mail" validate:"required,email"`
	}
	ForgetPasswordRes struct {
		UID string `json:"uid" validate:"required,email"`
	}
)

type (
	ForgetPasswordVerifyReq struct {
		UID  string `json:"uid" validate:"required,min=20"`
		Code string `json:"code" validate:"required,min=6"`
	}
)

type (
	ForgetChangePasswordReq struct {
		NewPassword string `json:"newPassword" validate:"required"`
	}
)

type (
	BindMailReq struct {
		Mail string `json:"mail"`
	}
	BindMailRes struct {
		UID string `json:"uid"`
	}
)

type (
	BindMailVerifyReq struct {
		UID  string `json:"uid"`
		Code string `json:"code"`
	}
	TwostepMail struct {
		Code string
	}
)
