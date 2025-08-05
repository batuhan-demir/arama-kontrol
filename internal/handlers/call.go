package handlers

import (
	"arama-kontrol/internal/dal"
	"arama-kontrol/pkg/database"
	"arama-kontrol/pkg/file"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetCalls(c *fiber.Ctx) error {
	callerFilter := c.Query("caller")
	statusFilter := c.Query("status")
	mineFilter := c.Query("mine") // "true" or "false"
	order := c.Query("order")     // "asc" or "desc" for timestamp

	var calls []dal.Call

	query := database.DB.Model(&dal.Call{})
	if callerFilter != "" {
		query = query.Where("caller_num LIKE ?", "%"+callerFilter+"%")
	}
	if statusFilter != "" {
		query = query.Where("call_status = ?", statusFilter)
	}

	// Order by timestamp
	if order == "asc" {
		query = query.Order("started_at ASC")
	} else {
		query = query.Order("started_at DESC") // default to newest first
	}

	if mineFilter == "true" {
		// Get the email of the authenticated user
		userEmail := c.Locals("user").(jwt.MapClaims)["email"].(string)
		query = query.Where("personel = ?", userEmail)
	}

	res := query.Find(&calls)

	// Tüm numaraları tek seferde yükle (performans için)
	var numbers []dal.Number
	database.DB.Find(&numbers)

	// Map oluştur: number -> name
	numberMap := make(map[string]string)
	for _, num := range numbers {
		numberMap[num.Number] = num.Name
	}

	// populate caller names and redirects
	for i, call := range calls {
		// Arayan kişinin ismini kontrol et ve ekle
		if call.CallerNum != "" {
			if name, exists := numberMap[call.CallerNum]; exists {
				calls[i].CallerNum = call.CallerNum + "  - " + name
			}
		}

		// Redirect edilen numaraların isimlerini ekle
		if call.Redirects != nil {
			for j, redirect := range call.Redirects {
				if name, exists := numberMap[redirect]; exists {
					call.Redirects[j] += " - " + name
				}
			}
			calls[i].Redirects = call.Redirects
		}
	}

	if res.Error != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "An error occurred while fetching calls",
			"error":   res.Error.Error(),
		})
	}

	return c.JSON(&fiber.Map{
		"success": "true",
		"data":    calls,
	})

}

func CallCallback(c *fiber.Ctx) error {
	body := new(dal.CreateCall)

	c.BodyParser(body)

	if body.CustomerNum == "" {
		// try to parse from dal.CreateCallCDR
		body2 := new(dal.CreateCallCDR)

		c.BodyParser(body2)
		fmt.Println("CDR Callback received:", body2)

		// download the call record and save it if it exists
		if body2.CallRecord != "" {
			err := file.Download(body2.CallRecord, body2.CallId+".wav")
			if err != nil {
				fmt.Printf("Failed to download call record: %v\n", err)
				// Don't set CallRecord URL if download failed
			} else {
				body.CallRecord = fmt.Sprintf("%s/files/%s.wav", os.Getenv("ORIGIN"), body2.CallId)
			}
		} else {
			fmt.Println("No call record URL provided")
		}

		body.Scenario = body2.Scenario
		body.CallId = body2.CallId
		body.CustomerNum = body2.CustomerNum
		body.IncomingNumber = body2.IncomingNumber
		body.Timestamp = body2.Timestamp
	}

	fmt.Println("Call Callback received:", body)

	//check if this call already exists
	var existingCall dal.Call
	database.DB.First(&existingCall, "call_id = ?", body.CallId)

	if existingCall.CallId != "" {
		// Call already exists, update its redirects and events

		if existingCall.Redirects == nil {
			existingCall.Redirects = []string{}
		}
		if body.InternalNum != "" {
			// if array doesnt include the number
			includes := false
			for _, redirect := range existingCall.Redirects {
				if redirect == body.InternalNum {
					includes = true
					break
				}
			}
			if !includes {
				existingCall.Redirects = append(existingCall.Redirects, body.InternalNum)
			}
		}
		if existingCall.Events == nil {
			existingCall.Events = dal.JSONBArray{}
		}
		newEvent := dal.JSONB{
			"Scenario":  body.Scenario,
			"Timestamp": body.Timestamp,
		}
		existingCall.Events = append(existingCall.Events, newEvent)

		if body.Scenario == "Answer" {
			existingCall.CallStatus = "answered"
		}
		if body.Scenario == "Hangup" {
			existingCall.EndedAt = body.Timestamp
		}
		if body.CallRecord != "" {
			existingCall.CallRecord = body.CallRecord
		}

		res := database.DB.Save(&existingCall)

		if res.Error != nil {
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"message": "An error occurred while updating the call",
				"error":   res.Error.Error(),
			})
		}

		return c.JSON(&fiber.Map{
			"success": true,
			"message": "Call updated successfully",
			"data":    existingCall,
		})
	}

	var addedBy string
	var personel string

	if body.Scenario == "NewManualCall" {
		addedBy = c.Locals("user").(jwt.MapClaims)["email"].(string)
		personel = addedBy // manuel arama ekleyenin kendisi sorumlu
	} else {
		addedBy = "system"

		personel = getNextUser()
	}

	newCall := dal.Call{
		CallId:    body.CallId,
		CallerNum: body.CustomerNum,
		Events: []dal.JSONB{
			{
				"Scenario":  body.Scenario,
				"Timestamp": body.Timestamp,
			},
		},
		StartedAt:  body.Timestamp,
		AddedBy:    addedBy,
		CallRecord: body.CallRecord,
		Personel:   personel,
	}

	res := database.DB.Create(&newCall)

	if res.Error != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "An error occured in server. Please try again later",
			"error":   res.Error.Error(),
		})
	}

	return c.Status(201).JSON(&fiber.Map{
		"success": true,
		"message": "Call Created Successfully",
		"data":    newCall,
	})
}

func UpdateCallStatus(c *fiber.Ctx) error {

	id := c.Params("id")

	newStatus := c.Params("newStatus")

	if newStatus != "answered" && newStatus != "not_answered" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "You should select one of those: 'answered', 'not_answered'",
		})
	}

	var call dal.Call

	database.DB.First(&call, "call_id = ?", id)

	call.CallStatus = newStatus

	newEvent := dal.JSONB{
		"Scenario":  fmt.Sprintf("CallStatus_%s_%s", newStatus, c.Locals("user").(jwt.MapClaims)["email"]),
		"Timestamp": time.Now(),
	}
	call.Events = append(call.Events, newEvent)

	res := database.DB.Save(&call)

	if res.Error != nil || res.RowsAffected == 0 {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "An error occured in server. Please try again later",
			"error":   res.Error.Error(),
		})
	}
	return c.JSON(&fiber.Map{
		"success": true,
		"message": "Call status updated successfully",
		"data":    call,
	})

}

func getNextUser() string {
	var users []dal.User
	database.DB.Where("is_active = ?", true).Order("Id ASC").Find(&users)

	if len(users) == 0 {
		return "system"
	}

	var lastCall dal.Call
	database.DB.Where("personel != ? AND personel != ?", "", "system").
		Order("started_at DESC").
		First(&lastCall)

	if lastCall.Personel == "" {
		return users[0].Email
	}

	var currentIndex int = -1
	for i, user := range users {
		if user.Email == lastCall.Personel {
			currentIndex = i
			break
		}
	}

	nextIndex := (currentIndex + 1) % len(users)
	return users[nextIndex].Email
}
