package service

import (
	"dockertest/connect"
	"dockertest/format"
	"dockertest/schema"
	"fmt"
	"time"
)

func CalendarCreate(title string, assignment time.Time, userID string) string {
	calendar := schema.Calendar{
		Title:       title,
		DueDateTime: assignment,
		AccountID:   userID,
	}
	if result := connect.DB.Model(&schema.Calendar{}).
		Create(&calendar); result.Error != nil {
		return "⚠️ เกิดข้อผิดพลาดขึ้น ลองตรวจสอบข้อมูลให้ดีแล้วลองใหม่อีกครั้งนะ"

	}

	return format.ThaiDate(title, assignment, calendar.ID, userID)
}

func GetTaskByPeriod(userID string, period int64) string {

	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc)

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24*time.Hour*time.Duration(period) - time.Nanosecond)

	var task []schema.Calendar
	if err := connect.DB.Order("due_date_time ASC").Where("account_id = ? AND due_date_time >= ? AND due_date_time <= ?", userID, startOfDay, endOfDay).Find(&task).Error; err != nil {
		fmt.Println(err)
	}

	if len(task) == 0 {
		return "🔍 ดูเหมือนว่าจะยังไม่มีงานในช่วงเวลาที่เธอเลือกนะ "
	}
	var result string
	for i, t := range task {
		result += format.ShowThaiDate(int(t.ID), i+1, t.Title, t.DueDateTime)
	}

	var periodType string
	if period == 1 {
		periodType = "วันนี้"
	} else if period == 7 {
		periodType = "สัปดาห์นี้"
	} else {
		periodType = "เดือนนี้"
	}

	return fmt.Sprintf("ได้เลย! นี่คืองานของ <@%s>  ภายใน%s\n.%s", userID, periodType, result)
}

func DeleteTaskByID(userID string, taskID uint) string {
	var task schema.Calendar

	if err := connect.DB.Where("account_id = ? AND id = ?", userID, taskID).First(&task).Error; err != nil {
		return "🚫 ไม่พบงานที่เธอต้องการลบ ลองเช็ค TaskID ให้แน่ใจอีกครั้งนะ"

	}

	if result := connect.DB.Where("account_id = ? AND ID = ?", userID, taskID).Delete(&schema.Calendar{}); result.Error != nil {
		return "⚠️ เกิดข้อผิดพลาดขึ้น ลองตรวจสอบข้อมูลให้ดีแล้วลองใหม่อีกครั้งนะ"

	}
	return fmt.Sprintf(
		"✅ ลบงานเรียบร้อยแล้ว <@%s>!\n"+
			"━━━━━━━━━━━━━━\n"+
			"🔢 TaskID: `%d`\n"+
			"📝 หัวข้อ: **%s**\n"+
			"🗑️ สถานะ: ถูกลบเรียบร้อยแล้ว",
		userID, taskID, task.Title,
	)

}

func GetTaskByTimeRange(startTime time.Time, endTime time.Time, userID string) string {
	var task []schema.Calendar

	if err := connect.DB.Order("due_date_time ASC").Where("account_id = ? AND due_date_time >= ? AND due_date_time <= ?", userID, startTime, endTime).Find(&task).Error; err != nil {
		fmt.Println(err)
	}

	if len(task) == 0 {
		return "🔍 ดูเหมือนว่าจะยังไม่มีงานในช่วงเวลาที่เธอเลือกนะ "

	}
	var result string
	for i, t := range task {
		result += format.ShowThaiDate(int(t.ID), i+1, t.Title, t.DueDateTime)
	}
	return fmt.Sprintf("ได้เลย! นี่คืองานของ <@%s> ตามที่เธอขอมา! \n.%s", userID, result)

}
