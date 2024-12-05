package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Register 注册自定义验证标签和错误消息
// 它会向 Gin 的验证器中添加自定义的验证规则
func Register() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义验证函数
		if err := v.RegisterValidation("password", validatePassword); err != nil {
			return err
		}
		if err := v.RegisterValidation("username_valid", validateUsername); err != nil {
			return err
		}
		if err := v.RegisterValidation("email_domain", validateEmailDomain); err != nil {
			return err
		}

		// 注册自定义错误消息处理函数
		// 使用 JSON 标签作为字段名，以保持一致性
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return fld.Name
			}
			return name
		})

		return nil
	}
	return fmt.Errorf("failed to register validator")
}

// validatePassword 检查密码是否满足安全要求
// 密码必须包含：
// - 至少一个大写字母
// - 至少一个小写字母
// - 至少一个数字
// - 至少一个特殊字符
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// 检查是否包含大写字母
	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	// 检查是否包含小写字母
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	// 检查是否包含数字
	hasNumber := strings.ContainsAny(password, "0123456789")
	// 检查是否包含特殊字符
	hasSpecial := strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:,.<>?")

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateUsername 执行用户名的附加验证
// 用户名必须：
// - 以字母开头
// - 不能包含连续的下划线
// - 不能以特殊字符结尾
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	// 检查是否以字母开头
	if matched, _ := regexp.MatchString(`^[^a-zA-Z]`, username); matched {
		return false
	}

	// 检查是否包含连续的下划线
	if strings.Contains(username, "__") {
		return false
	}

	// 检查是否以特殊字符结尾
	if matched, _ := regexp.MatchString(`[^a-zA-Z0-9]$`, username); matched {
		return false
	}

	return true
}

// validateEmailDomain 执行邮箱域名验证
// 它会：
// - 检查域名格式是否正确
// - 检查域名是否在黑名单中
func validateEmailDomain(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	domain := strings.ToLower(parts[1])

	// 验证域名格式
	if matched, _ := regexp.MatchString(`^[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,}$`, domain); !matched {
		return false
	}

	// 检查域名是否在黑名单中
	blacklist := []string{
		"example.com",
		"test.com",
		"invalid.com",
	}

	for _, blocked := range blacklist {
		if domain == blocked {
			return false
		}
	}

	return true
}

// ValidationError 表示验证错误
// 它包含字段名和错误消息
type ValidationError struct {
	Field   string `json:"field"`   // 验证失败的字段名
	Message string `json:"message"` // 错误描述信息
}

// ValidationErrors 是验证错误的集合
type ValidationErrors []ValidationError

// FormatError 将验证器的错误转换为用户友好的格式
// 它会将每个验证错误转换为包含字段名和错误消息的结构
func FormatError(err error) ValidationErrors {
	var errors ValidationErrors

	validatorErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors
	}

	for _, e := range validatorErrs {
		message := getErrorMsg(e)
		errors = append(errors, ValidationError{
			Field:   e.Field(),
			Message: message,
		})
	}

	return errors
}

// getErrorMsg 返回验证错误的用户友好消息
// 它根据验证标签类型返回相应的错误描述
func getErrorMsg(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "此字段为必填项"
	case "email":
		return "邮箱格式不正确"
	case "min":
		return fmt.Sprintf("最小长度为 %s", e.Param())
	case "max":
		return fmt.Sprintf("最大长度为 %s", e.Param())
	case "alphanum":
		return "只能包含字母和数字"
	case "password":
		return "密码必须包含至少一个大写字母、一个小写字母、一个数字和一个特殊字符"
	case "username_valid":
		return "用户名必须以字母开头，不能包含连续的特殊字符，且必须以字母或数字结尾"
	case "email_domain":
		return "邮箱域名不被允许或格式不正确"
	default:
		return fmt.Sprintf("验证失败：%s", e.Tag())
	}
}
