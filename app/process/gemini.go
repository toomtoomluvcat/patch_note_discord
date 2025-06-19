package process

import (
	"context"
	"dockertest/service"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/genai"
)

type Payload struct {
	Today    string `json:"today"`
	DateType string `json:"date_type"`
	Text     string `json:"text"`
}
type PayloadFind struct {
	Today    string `json:"today"`
	TimeNow  string `json:"time_now"`
	DateType string `json:"date_type"`
	Text     string `json:"text"`
}

func GenAI(req string, i *discordgo.InteractionCreate) string {
	err := service.CheckRegister(service.GetUserID(i))
	if err != nil {
		return fmt.Sprintf("🚫 ดูเหมือนว่าบัญชี @%s \nลองใช้คำสั่ง /register พร้อมใส่ passkey ที่ได้รับสิทธิ์ดูนะ 🗝️", service.GetUsername(i))
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  "YOUR--GEMINI--APIKEY",
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	now := time.Now().In(time.FixedZone("Asia/Bangkok", 7*60*60))
	today := now.Format("2006-01-02")
	timeNow := now.Format("15:04")
	dateType := strings.ToLower(now.Format("Mon"))

	payload := PayloadFind{
		TimeNow:  timeNow,
		Today:    today,
		DateType: dateType,
		Text:     req,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(string(jsonBytes)),
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{
					{Text: "You are an assistant that extracts task and datetime information from the provided text and converts it into a JSON object."},
					{Text: "Respond only with the JSON object inside braces {}. Do not include any explanations, code blocks, or additional text."},
					{Text: "The JSON should have the following structure:"},
					{Text: `{
	"title": "task title" ,
	"due_datetime": "YYYY-MM-DD HH:MM" (24-hour format)
}`},
					{Text: "Input may contain task description, specific date, and time."},
					{Text: "Use the current date provided as \"today\" in YYYY-MM-DD format to determine the due date."},
					{Text: "If a specific date is given, use that date directly."},
					{Text: "If no date is specified, assume the next occurrence of the target date after today."},
					{Text: "due_datetime should be complete (both date and time)."},
					{Text: "If date or time info is incomplete, respond with: { \"error\": \"Missing date or time information for task 'task title'\" }."},
					{Text: "Output only the JSON inside braces {}."},
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to generate content: %v", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(result.Text()), &data)
	if err != nil {
		return "เกิดข้อผิดพลาดฝั่ง server แจ้งแอดมินให้หน่อยนะ"
	}
	if _, ok := data["error"]; ok {
		return "งานมันกำกวมง่ะ บอกวันเวลาที่ชัดกว่านี้ได้ไหม"
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	dateTimeStr := data["due_datetime"].(string)
	dueDateTime, err := time.ParseInLocation("2006-01-02 15:04", dateTimeStr, loc)
	if err != nil {
		return "เกิดข้อผิดพลาดฝั่ง server แจ้งแอดมินให้หน่อยนะ"
	}

	message := service.CalendarCreate(data["title"].(string), dueDateTime, service.GetUserID(i))
	return message
}

func GenAISearchTimeRange(req string, i *discordgo.InteractionCreate) string {
	err := service.CheckRegister(service.GetUserID(i))
	if err != nil {
		return fmt.Sprintf("🚫 ดูเหมือนว่าบัญชี @%s \nลองใช้คำสั่ง /register พร้อมใส่ passkey ที่ได้รับสิทธิ์ดูนะ 🗝️", service.GetUsername(i))
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  "YOUR--GEMINI--APIKEY",
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	now := time.Now()
	today := now.Format("2006-01-02")
	timeNow := now.Format("15:04")
	dateType := strings.ToLower(now.Format("Mon"))

	payload := PayloadFind{
		Today:    today,
		TimeNow:  timeNow,
		DateType: dateType,
		Text:     req,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	// 🔧 System Prompt สำหรับค้นหาช่วงเวลา
	systemPrompt := &genai.Content{
		Parts: []*genai.Part{
			{Text: "You are an assistant that extracts a time range from the user's input and responds with a valid JSON object only."},
			{Text: "You will receive the current date as 'today' in YYYY-MM-DD format, and 'time_now' in HH:MM (24-hour format)."},
			{Text: "Do NOT include any explanations, markdown, or code blocks in your response."},
			{Text: "Respond only with a raw JSON object. Do not wrap it in triple backticks or any other symbols."},
			{Text: `Example request:
			{
			"today": "2025-06-16",
			"time_now": "14:45",
			"date_type": "mon",
			"text": "หางานช่วงเช้าของเมื่อวาน"
			}`},
			{Text: `Correct response format:
			{
			"start_datetime": "2025-06-15 06:00",
			"end_datetime": "2025-06-15 12:00"
			}`},
			{Text: "If the input is ambiguous or missing time information, respond with: { \"error\": \"Missing or unclear time information in user request.\" }"},
		},
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(string(jsonBytes)),
		&genai.GenerateContentConfig{
			SystemInstruction: systemPrompt,
		},
	)
	if err != nil {
		log.Fatalf("Failed to generate content: %v", err)
	}
	var data map[string]interface{}

	err = json.Unmarshal([]byte(result.Text()), &data)
	if err != nil {
		fmt.Println(err)
		return "เกิดปัญหาที่ฝั่ง server ลองติดต่อแอดมินดูนะ"
	}
	if _, ok := data["error"]; ok {
		return "ข้อความดูกำกวมระหว่างช่วงเวลา บอกเวลาที่ชัดกว่านี้ได้มั้ย"
	}
	startStr := data["start_datetime"].(string)
	endStr := data["end_datetime"].(string)

	loc, _ := time.LoadLocation("Asia/Bangkok")
	startTime, err1 := time.ParseInLocation("2006-01-02 15:04", startStr, loc)
	endTime, err2 := time.ParseInLocation("2006-01-02 15:04", endStr, loc)

	if err1 != nil || err2 != nil {
		return "เกิดปัญหาในการตีความเวลา ลองใหม่อีกครั้งนะ"
	}
	return service.GetTaskByTimeRange(startTime, endTime, service.GetUserID(i))

}
