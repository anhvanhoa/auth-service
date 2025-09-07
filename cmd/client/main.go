package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serverAddress = "localhost:50050" // Default gRPC server address
)

type AuthGRPCClient struct {
	authClient proto_auth.AuthServiceClient
	conn       *grpc.ClientConn
}

func NewAuthGRPCClient(address string) (*AuthGRPCClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}

	return &AuthGRPCClient{
		authClient: proto_auth.NewAuthServiceClient(conn),
		conn:       conn,
	}, nil
}

func (c *AuthGRPCClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// Login Test
func (c *AuthGRPCClient) TestLogin() {
	fmt.Println("\n=== Test Login ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter email or phone: ")
	emailOrPhone, _ := reader.ReadString('\n')
	emailOrPhone = strings.TrimSpace(emailOrPhone)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("Enter OS (web/mobile): ")
	os, _ := reader.ReadString('\n')
	os = strings.TrimSpace(os)
	if os == "" {
		os = "web"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.Login(ctx, &proto_auth.LoginRequest{
		EmailOrPhone: emailOrPhone,
		Password:     password,
		Os:           os,
	})
	if err != nil {
		fmt.Printf("Error calling Login: %v\n", err)
		return
	}

	fmt.Printf("Login successful!\n")
	fmt.Printf("Message: %s\n", resp.Message)
	fmt.Printf("Access Token: %s\n", resp.AccessToken)
	fmt.Printf("Refresh Token: %s\n", resp.RefreshToken)
	fmt.Printf("User Info:\n")
	fmt.Printf("  ID: %s\n", resp.User.Id)
	fmt.Printf("  Email: %s\n", resp.User.Email)
	fmt.Printf("  Phone: %s\n", resp.User.Phone)
	fmt.Printf("  Full Name: %s\n", resp.User.FullName)
	fmt.Printf("  Avatar: %s\n", resp.User.Avatar)
	fmt.Printf("  Bio: %s\n", resp.User.Bio)
	fmt.Printf("  Address: %s\n", resp.User.Address)
	if resp.User.Birthday != nil {
		fmt.Printf("  Birthday: %s\n", resp.User.Birthday.AsTime().Format(time.RFC3339))
	}
}

// Register Test
func (c *AuthGRPCClient) TestRegister() {
	fmt.Println("\n=== Test Register ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Enter full name: ")
	fullName, _ := reader.ReadString('\n')
	fullName = strings.TrimSpace(fullName)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("Enter confirm password: ")
	confirmPassword, _ := reader.ReadString('\n')
	confirmPassword = strings.TrimSpace(confirmPassword)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.authClient.Register(ctx, &proto_auth.RegisterRequest{
		Email:           email,
		FullName:        fullName,
		Password:        password,
		ConfirmPassword: confirmPassword,
	})
	if err != nil {
		fmt.Printf("Error calling Register: %v\n", err)
		return
	}

	fmt.Printf("Register successful!\n")
	fmt.Printf("Message: %s\n", resp.Message)
	fmt.Printf("Verification Token: %s\n", resp.Token)
	fmt.Printf("User Info:\n")
	fmt.Printf("  ID: %s\n", resp.User.Id)
	fmt.Printf("  Email: %s\n", resp.User.Email)
	fmt.Printf("  Phone: %s\n", resp.User.Phone)
	fmt.Printf("  Full Name: %s\n", resp.User.FullName)
	fmt.Printf("  Avatar: %s\n", resp.User.Avatar)
	fmt.Printf("  Bio: %s\n", resp.User.Bio)
	fmt.Printf("  Address: %s\n", resp.User.Address)
	if resp.User.Birthday != nil {
		fmt.Printf("  Birthday: %s\n", resp.User.Birthday.AsTime().Format(time.RFC3339))
	}
}

// Check Token Test
func (c *AuthGRPCClient) TestCheckToken() {
	fmt.Println("\n=== Test Check Token ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter token to check: ")
	token, _ := reader.ReadString('\n')
	token = strings.TrimSpace(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.CheckToken(ctx, &proto_auth.CheckTokenRequest{
		Token: token,
	})
	if err != nil {
		fmt.Printf("Error calling CheckToken: %v\n", err)
		return
	}

	fmt.Printf("Token check result:\n")
	fmt.Printf("Message: %s\n", resp.Message)
	fmt.Printf("Valid: %t\n", resp.Data)
}

// Refresh Token Test
func (c *AuthGRPCClient) TestRefreshToken() {
	fmt.Println("\n=== Test Refresh Token ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter refresh token: ")
	refreshToken, _ := reader.ReadString('\n')
	refreshToken = strings.TrimSpace(refreshToken)

	fmt.Print("Enter OS (web/mobile): ")
	os, _ := reader.ReadString('\n')
	os = strings.TrimSpace(os)
	if os == "" {
		os = "web"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.RefreshToken(ctx, &proto_auth.RefreshTokenRequest{
		RefreshToken: refreshToken,
		Os:           os,
	})
	if err != nil {
		fmt.Printf("Error calling RefreshToken: %v\n", err)
		return
	}

	fmt.Printf("Token refresh successful!\n")
	fmt.Printf("Message: %s\n", resp.Message)
	fmt.Printf("New Access Token: %s\n", resp.AccessToken)
	fmt.Printf("New Refresh Token: %s\n", resp.RefreshToken)
}

// Logout Test
func (c *AuthGRPCClient) TestLogout() {
	fmt.Println("\n=== Test Logout ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter access token to logout: ")
	token, _ := reader.ReadString('\n')
	token = strings.TrimSpace(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.Logout(ctx, &proto_auth.LogoutRequest{
		Token: token,
	})
	if err != nil {
		fmt.Printf("Error calling Logout: %v\n", err)
		return
	}

	fmt.Printf("Logout successful!\n")
	fmt.Printf("Message: %s\n", resp.Message)
}

// Verify Account Test
func (c *AuthGRPCClient) TestVerifyAccount() {
	fmt.Println("\n=== Test Verify Account ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter verification token: ")
	token, _ := reader.ReadString('\n')
	token = strings.TrimSpace(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.VerifyAccount(ctx, &proto_auth.VerifyAccountRequest{
		Token: token,
	})
	if err != nil {
		fmt.Printf("Error calling VerifyAccount: %v\n", err)
		return
	}

	fmt.Printf("Account verification successful!\n")
	fmt.Printf("Message: %s\n", resp.Message)
}

// Forgot Password Test
func (c *AuthGRPCClient) TestForgotPassword() {
	fmt.Println("\n=== Test Forgot Password ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Enter OS (web/mobile): ")
	os, _ := reader.ReadString('\n')
	os = strings.TrimSpace(os)
	if os == "" {
		os = "web"
	}

	fmt.Print("Enter method (0=Code, 1=Token): ")
	methodStr, _ := reader.ReadString('\n')
	methodStr = strings.TrimSpace(methodStr)

	var method proto_auth.ForgotPasswordType
	if methodStr == "1" {
		method = proto_auth.ForgotPasswordType_FORGOT_PASSWORD_TYPE_TOKEN
	} else {
		method = proto_auth.ForgotPasswordType_FORGOT_PASSWORD_TYPE_UNSPECIFIED
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.ForgotPassword(ctx, &proto_auth.ForgotPasswordRequest{
		Email:  email,
		Os:     os,
		Method: method,
	})
	if err != nil {
		fmt.Printf("Error calling ForgotPassword: %v\n", err)
		return
	}

	fmt.Printf("Forgot password request successful!\n")
	fmt.Printf("Message: %s\n", resp.Message)
	fmt.Printf("Token: %s\n", resp.Token)
	fmt.Printf("Code: %s\n", resp.Code)
	fmt.Printf("User Info:\n")
	fmt.Printf("  ID: %s\n", resp.User.Id)
	fmt.Printf("  Email: %s\n", resp.User.Email)
	fmt.Printf("  Phone: %s\n", resp.User.Phone)
	fmt.Printf("  Full Name: %s\n", resp.User.FullName)
}

// Check Code Test
func (c *AuthGRPCClient) TestCheckCode() {
	fmt.Println("\n=== Test Check Code ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter verification code: ")
	code, _ := reader.ReadString('\n')
	code = strings.TrimSpace(code)

	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.CheckCode(ctx, &proto_auth.CheckCodeRequest{
		Code:  code,
		Email: email,
	})
	if err != nil {
		fmt.Printf("Error calling CheckCode: %v\n", err)
		return
	}

	fmt.Printf("Code check result:\n")
	fmt.Printf("Message: %s\n", resp.Message)
	fmt.Printf("Valid: %t\n", resp.Valid)
}

// Reset Password By Code Test
func (c *AuthGRPCClient) TestResetPasswordByCode() {
	fmt.Println("\n=== Test Reset Password By Code ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter verification code: ")
	code, _ := reader.ReadString('\n')
	code = strings.TrimSpace(code)

	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Enter new password: ")
	newPassword, _ := reader.ReadString('\n')
	newPassword = strings.TrimSpace(newPassword)

	fmt.Print("Enter confirm password: ")
	confirmPassword, _ := reader.ReadString('\n')
	confirmPassword = strings.TrimSpace(confirmPassword)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.ResetPasswordByCode(ctx, &proto_auth.ResetPasswordByCodeRequest{
		Code:            code,
		Email:           email,
		NewPassword:     newPassword,
		ConfirmPassword: confirmPassword,
	})
	if err != nil {
		fmt.Printf("Error calling ResetPasswordByCode: %v\n", err)
		return
	}

	fmt.Printf("Password reset successful!\n")
	fmt.Printf("Message: %s\n", resp.Message)
}

// Reset Password By Token Test
func (c *AuthGRPCClient) TestResetPasswordByToken() {
	fmt.Println("\n=== Test Reset Password By Token ===")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter reset token: ")
	token, _ := reader.ReadString('\n')
	token = strings.TrimSpace(token)

	fmt.Print("Enter new password: ")
	newPassword, _ := reader.ReadString('\n')
	newPassword = strings.TrimSpace(newPassword)

	fmt.Print("Enter confirm password: ")
	confirmPassword, _ := reader.ReadString('\n')
	confirmPassword = strings.TrimSpace(confirmPassword)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.ResetPasswordByToken(ctx, &proto_auth.ResetPasswordByTokenRequest{
		Token:           token,
		NewPassword:     newPassword,
		ConfirmPassword: confirmPassword,
	})
	if err != nil {
		fmt.Printf("Error calling ResetPasswordByToken: %v\n", err)
		return
	}

	fmt.Printf("Password reset successful!\n")
	fmt.Printf("Message: %s\n", resp.Message)
}

func printMenu() {
	fmt.Println("\n=== gRPC Auth Service Test Client ===")
	fmt.Println("1. Authentication Tests")
	fmt.Println("  1.1 Login")
	fmt.Println("  1.2 Register")
	fmt.Println("  1.3 Check Token")
	fmt.Println("  1.4 Refresh Token")
	fmt.Println("  1.5 Logout")
	fmt.Println("2. Account Verification Tests")
	fmt.Println("  2.1 Verify Account")
	fmt.Println("  2.2 Check Code")
	fmt.Println("3. Password Management Tests")
	fmt.Println("  3.1 Forgot Password")
	fmt.Println("  3.2 Reset Password By Code")
	fmt.Println("  3.3 Reset Password By Token")
	fmt.Println("0. Exit")
	fmt.Print("Enter your choice: ")
}

func main() {
	// Get server address from command line or use default
	address := serverAddress
	if len(os.Args) > 1 {
		address = os.Args[1]
	}

	fmt.Printf("Connecting to gRPC Auth server at %s...\n", address)
	client, err := NewAuthGRPCClient(address)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer client.Close()

	fmt.Println("Connected successfully!")

	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1.1":
			client.TestLogin()
		case "1.2":
			client.TestRegister()
		case "1.3":
			client.TestCheckToken()
		case "1.4":
			client.TestRefreshToken()
		case "1.5":
			client.TestLogout()
		case "2.1":
			client.TestVerifyAccount()
		case "2.2":
			client.TestCheckCode()
		case "3.1":
			client.TestForgotPassword()
		case "3.2":
			client.TestResetPasswordByCode()
		case "3.3":
			client.TestResetPasswordByToken()
		case "0":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
