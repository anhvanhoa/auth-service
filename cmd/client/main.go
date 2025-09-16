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

type GRPCClient struct {
	authClient proto_auth.AuthServiceClient
	conn       *grpc.ClientConn
}

func NewGRPCClient(address string) (*GRPCClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}

	return &GRPCClient{
		authClient: proto_auth.NewAuthServiceClient(conn),
		conn:       conn,
	}, nil
}

func (c *GRPCClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// --- Helper để làm sạch input ---
func cleanInput(s string) string {
	return strings.ToValidUTF8(strings.TrimSpace(s), "")
}

// ================== Auth Service Tests ==================

func (c *GRPCClient) TestRegister() {
	fmt.Println("\n=== Test Register ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	email = cleanInput(email)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = cleanInput(password)

	fmt.Print("Enter confirm password: ")
	confirmPassword, _ := reader.ReadString('\n')
	confirmPassword = cleanInput(confirmPassword)

	fmt.Print("Enter full name: ")
	fullName, _ := reader.ReadString('\n')
	fullName = cleanInput(fullName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.Register(ctx, &proto_auth.RegisterRequest{
		Email:           email,
		Password:        password,
		ConfirmPassword: confirmPassword,
		FullName:        fullName,
	})
	if err != nil {
		fmt.Printf("Error calling Register: %v\n", err)
		return
	}

	fmt.Printf("Register result:\n")
	fmt.Printf("Message: %s\n", resp.Message)
	fmt.Printf("User ID: %s\n", resp.User.Id)
	fmt.Printf("User Email: %s\n", resp.User.Email)
	fmt.Printf("User Phone: %s\n", resp.User.Phone)
	fmt.Printf("User Full Name: %s\n", resp.User.FullName)
	fmt.Printf("Token: %s\n", resp.Token)
}

func (c *GRPCClient) TestLogin() {
	fmt.Println("\n=== Test Login ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter email or phone: ")
	emailOrPhone, _ := reader.ReadString('\n')
	emailOrPhone = cleanInput(emailOrPhone)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = cleanInput(password)

	fmt.Print("Enter OS: ")
	osName, _ := reader.ReadString('\n')
	osName = cleanInput(osName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.Login(ctx, &proto_auth.LoginRequest{
		EmailOrPhone: emailOrPhone,
		Password:     password,
		Os:           osName,
	})
	if err != nil {
		fmt.Printf("Error calling Login: %v\n", err)
		return
	}

	fmt.Printf("Login result:\n")
	fmt.Printf("User ID: %s\n", resp.User.Id)
	fmt.Printf("User Email: %s\n", resp.User.Email)
	fmt.Printf("User Phone: %s\n", resp.User.Phone)
	fmt.Printf("User Full Name: %s\n", resp.User.FullName)
	fmt.Printf("Access Token: %s\n", resp.AccessToken)
	fmt.Printf("Refresh Token: %s\n", resp.RefreshToken)
}

func (c *GRPCClient) TestVerifyAccount() {
	fmt.Println("\n=== Test Verify Account ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter verification token: ")
	token, _ := reader.ReadString('\n')
	token = cleanInput(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.VerifyAccount(ctx, &proto_auth.VerifyAccountRequest{
		Token: token,
	})
	if err != nil {
		fmt.Printf("Error calling VerifyAccount: %v\n", err)
		return
	}

	fmt.Printf("Verify Account result:\n")
	fmt.Printf("Message: %s\n", resp.Message)
}

func (c *GRPCClient) TestRefreshToken() {
	fmt.Println("\n=== Test Refresh Token ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter refresh token: ")
	refreshToken, _ := reader.ReadString('\n')
	refreshToken = cleanInput(refreshToken)

	fmt.Print("Enter OS: ")
	osName, _ := reader.ReadString('\n')
	osName = cleanInput(osName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.RefreshToken(ctx, &proto_auth.RefreshTokenRequest{
		RefreshToken: refreshToken,
		Os:           osName,
	})
	if err != nil {
		fmt.Printf("Error calling RefreshToken: %v\n", err)
		return
	}

	fmt.Printf("Refresh Token result:\n")
	fmt.Printf("Access Token: %s\n", resp.AccessToken)
	fmt.Printf("Refresh Token: %s\n", resp.RefreshToken)
}

func (c *GRPCClient) TestLogout() {
	fmt.Println("\n=== Test Logout ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter access token: ")
	token, _ := reader.ReadString('\n')
	token = cleanInput(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.Logout(ctx, &proto_auth.LogoutRequest{
		Token: token,
	})
	if err != nil {
		fmt.Printf("Error calling Logout: %v\n", err)
		return
	}

	fmt.Printf("Logout result:\n")
	fmt.Printf("Message: %s\n", resp.Message)
}

func (c *GRPCClient) TestForgotPassword() {
	fmt.Println("\n=== Test Forgot Password ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	email = cleanInput(email)

	fmt.Print("Enter method (0=Code, 1=Token): ")
	methodStr, _ := reader.ReadString('\n')
	methodStr = cleanInput(methodStr)

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
		Method: method,
		Os:     "os-test",
	})
	if err != nil {
		fmt.Printf("Error calling ForgotPassword: %v\n", err)
		return
	}

	fmt.Printf("Forgot Password result:\n")
	fmt.Printf("Message: %s\n", resp.Message)
	fmt.Printf("User ID: %s\n", resp.User.Id)
	fmt.Printf("User Email: %s\n", resp.User.Email)
	fmt.Printf("Token: %s\n", resp.Token)
	fmt.Printf("Code: %s\n", resp.Code)
}

func (c *GRPCClient) TestResetPasswordByToken() {
	fmt.Println("\n=== Test Reset Password By Token ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter token: ")
	token, _ := reader.ReadString('\n')
	token = cleanInput(token)

	fmt.Print("Enter new password: ")
	newPassword, _ := reader.ReadString('\n')
	newPassword = cleanInput(newPassword)

	fmt.Print("Enter confirm password: ")
	confirmPassword, _ := reader.ReadString('\n')
	confirmPassword = cleanInput(confirmPassword)

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

	fmt.Printf("Reset Password By Token result:\n")
	fmt.Printf("Message: %s\n", resp.Message)
}

func (c *GRPCClient) TestResetPasswordByCode() {
	fmt.Println("\n=== Test Reset Password By Code ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter code: ")
	code, _ := reader.ReadString('\n')
	code = cleanInput(code)

	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	email = cleanInput(email)

	fmt.Print("Enter new password: ")
	newPassword, _ := reader.ReadString('\n')
	newPassword = cleanInput(newPassword)

	fmt.Print("Enter confirm password: ")
	confirmPassword, _ := reader.ReadString('\n')
	confirmPassword = cleanInput(confirmPassword)

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

	fmt.Printf("Reset Password By Code result:\n")
	fmt.Printf("Message: %s\n", resp.Message)
}

func (c *GRPCClient) TestCheckToken() {
	fmt.Println("\n=== Test Check Token ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter token: ")
	token, _ := reader.ReadString('\n')
	token = cleanInput(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.authClient.CheckToken(ctx, &proto_auth.CheckTokenRequest{
		Token: token,
	})
	if err != nil {
		fmt.Printf("Error calling CheckToken: %v\n", err)
		return
	}

	fmt.Printf("Check Token result:\n")
	fmt.Printf("Valid: %t\n", resp.Data)
	fmt.Printf("Message: %s\n", resp.Message)
}

func (c *GRPCClient) TestCheckCode() {
	fmt.Println("\n=== Test Check Code ===")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter code: ")
	code, _ := reader.ReadString('\n')
	code = cleanInput(code)

	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	email = cleanInput(email)

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

	fmt.Printf("Check Code result:\n")
	fmt.Printf("Valid: %t\n", resp.Valid)
	fmt.Printf("Message: %s\n", resp.Message)
}

func printMenu() {
	fmt.Println("\n=== gRPC Auth Service Test Client ===")
	fmt.Println("1. Register")
	fmt.Println("2. Login")
	fmt.Println("3. Verify Account")
	fmt.Println("4. Refresh Token")
	fmt.Println("5. Logout")
	fmt.Println("6. Forgot Password")
	fmt.Println("7. Reset Password By Token")
	fmt.Println("8. Reset Password By Code")
	fmt.Println("9. Check Token")
	fmt.Println("10. Check Code")
	fmt.Println("0. Exit")
	fmt.Print("Enter your choice: ")
}

func main() {
	address := serverAddress
	if len(os.Args) > 1 {
		address = os.Args[1]
	}

	fmt.Printf("Connecting to gRPC server at %s...\n", address)
	client, err := NewGRPCClient(address)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer client.Close()

	fmt.Println("Connected successfully!")

	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()
		choice, _ := reader.ReadString('\n')
		choice = cleanInput(choice)

		switch choice {
		case "1":
			client.TestRegister()
		case "2":
			client.TestLogin()
		case "3":
			client.TestVerifyAccount()
		case "4":
			client.TestRefreshToken()
		case "5":
			client.TestLogout()
		case "6":
			client.TestForgotPassword()
		case "7":
			client.TestResetPasswordByToken()
		case "8":
			client.TestResetPasswordByCode()
		case "9":
			client.TestCheckToken()
		case "10":
			client.TestCheckCode()
		case "0":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
