package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"server/internal/email"
	"server/internal/files"
	"server/internal/types"
	"server/internal/utils"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
)

func (s *FiberServer) RegisterFiberRoutes() {
	origins := os.Getenv("FRONTEND_URL")
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Refresh-Token",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowCredentials: true,
	}))
	api := s.App.Group("/api")
	api.Post("/refresh-token", RefreshTokenHandler)
	api.Get("/", s.HelloWorldHandler)
	api.Post("/login", s.loginUserHandler)
	// api.Use(middleware.AuthenticationMiddleware())
	api.Get("/send-email", s.sendEmail)
	api.Get("/file-analyze", s.fileAnalyze)
	api.Get("/ope", s.ope)
	api.Post("/send-email-through-json-file", s.sendMailThroughJsonFile)

}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) sendEmail(c *fiber.Ctx) error {
	var emailToSend types.Email
	if err := c.BodyParser(&emailToSend); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := email.SendEmail(emailToSend)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Email sent"})
}

func (s *FiberServer) loginUserHandler(c *fiber.Ctx) error {
	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"status": "failed", "message": err.Error()})
	}
	if strings.ToLower(user.Email) != os.Getenv("EMAIL") || user.Password != os.Getenv("PASSWORD") {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"status": "failed", "message": "Invalid credentials"})
	}

	accessToken, refreshToken, err := utils.GenerateToken(user.Email, user.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"status": "failed", "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{"status": "success", "message": "User logged in successfully", "access_token": accessToken, "refresh_token": refreshToken})
}
func (s *FiberServer) fileAnalyze(c *fiber.Ctx) error {
	files.Analyze()
	return c.SendString("File Analyzed")
}

func (s *FiberServer) ope(c *fiber.Ctx) error {
	// Define a struct to represent the data
	type Person struct {
		Name string `json:"name"`
		Mail string `json:"email"`
	}
	// Open the output.json file
	file, err := os.Open(os.Getenv("JSON_FILE_PATH"))
	if err != nil {
		log.Printf("Error opening JSON file: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error reading JSON file")
	}
	defer file.Close()

	// Parse the JSON data
	var people []Person
	if err := json.NewDecoder(file).Decode(&people); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing JSON")
	}

	// // Log the parsed data
	// for _, person := range people {
	// 	log.Printf("First Name: %s, Last Name: %s, Mail: %s\n", person.FirstName, person.LastName, person.Mail)
	// }

	return c.Status(fiber.StatusOK).JSON(people) // Optionally return the data as a JSON response
}

func RefreshTokenHandler(c *fiber.Ctx) error {
	// Get the refresh token from the request body (or headers)
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	var request RefreshRequest
	if err := c.BodyParser(&request); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if request.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing refresh token",
		})
	}

	// Refresh the tokens
	newAccessToken, newRefreshToken, err := utils.RefreshToken(request.RefreshToken)
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired refresh token",
		})
	}

	// Return the new tokens
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken":  newAccessToken,
		"refreshToken": newRefreshToken,
	})
}

func (s *FiberServer) sendMailThroughJsonFile(c *fiber.Ctx) error {
	var emailToSend struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	type Person struct {
		Name string `json:"name"`
		Mail string `json:"email"`
	}

	if err := c.BodyParser(&emailToSend); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	file, err := os.Open(os.Getenv("JSON_FILE_PATH"))
	if err != nil {
		log.Printf("Error opening JSON file: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error reading JSON file")
	}
	defer file.Close()

	var people []Person
	if err := json.NewDecoder(file).Decode(&people); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing JSON")
	}

	concurrencyLimit := 100
	sem := make(chan struct{}, concurrencyLimit)
	var wg sync.WaitGroup
	errChan := make(chan error, len(people))

	for _, u := range people {
		sem <- struct{}{}
		wg.Add(1)

		go func(user Person) {
			defer func() { <-sem }()
			defer wg.Done()

			err := email.SendEmail(types.Email{
				EmailTo: user.Mail,
				Subject: emailToSend.Subject,
				Body:    emailToSend.Body,
			})

			errChan <- err
		}(u)
	}

	wg.Wait()
	close(errChan)

	// Collect errors
	var failedEmails []string
	for err := range errChan {
		if err != nil {
			log.Printf("Error sending email: %v", err)
			failedEmails = append(failedEmails, err.Error())
		}
	}

	if len(failedEmails) > 0 {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Some emails failed", "errors": failedEmails})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Emails sent to all users"})
}
