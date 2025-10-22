package grpcservice

import "github.com/anhvanhoa/service-core/domain/oops"

var (
	ErrMailServiceNotAvailable     = oops.New("Định vụ mail không khả dụng")
	ErrRegisterUsecaseNotAvailable = oops.New("Định vụ đăng ký không khả dụng")
	ErrUuidGeneratorNotAvailable   = oops.New("Định vụ UUID generator không khả dụng")
	ErrPasswordMatchNotMatch       = oops.New("Mật khẩu và xác nhận mật khẩu không khớp")
)
