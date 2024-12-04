package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Register registers custom validation tags
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

		// 注册自定义错误消息
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

// validatePassword checks if the password meets security requirements
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// 至少包含一个大写字母
	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	// 至少包含一个小写字母
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	// 至少包含一个数字
	hasNumber := strings.ContainsAny(password, "0123456789")
	// 至少包含一个特殊字符
	hasSpecial := strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:,.<>?")

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateUsername performs additional username validation
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	// 不允许以特殊字符开头
	if matched, _ := regexp.MatchString(`^[^a-zA-Z]`, username); matched {
		return false
	}

	// 不允许连续的下划线
	if strings.Contains(username, "__") {
		return false
	}

	// 不允许特殊字符结尾
	if matched, _ := regexp.MatchString(`[^a-zA-Z0-9]$`, username); matched {
		return false
	}

	return true
}

// validateEmailDomain performs email domain validation
func validateEmailDomain(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	domain := strings.ToLower(parts[1])

	// 检查域名格式
	if matched, _ := regexp.MatchString(`^[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,}$`, domain); !matched {
		return false
	}

	// 域名黑名单
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

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of ValidationError
type ValidationErrors []ValidationError

// FormatError formats validator.ValidationErrors into ValidationErrors
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

// getErrorMsg returns a user-friendly error message for validation errors
func getErrorMsg(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Minimum length is %s", e.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s", e.Param())
	case "alphanum":
		return "Must contain only letters and numbers"
	case "password":
		return "Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character"
	case "username_valid":
		return "Username must start with a letter, cannot contain consecutive special characters, and must end with a letter or number"
	case "email_domain":
		return "Email domain is not allowed or invalid"
	default:
		return fmt.Sprintf("Validation failed on condition: %s", e.Tag())
	}
}
