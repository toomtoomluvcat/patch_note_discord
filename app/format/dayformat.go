package format

import (
	"fmt"
	"time"
)

var thaiDays = [...]string{
	"อาทิตย์", "จันทร์", "อังคาร", "พุธ", "พฤหัสบดี", "ศุกร์", "เสาร์",
}
var thaiMonths = [...]string{
	"มกราคม", "กุมภาพันธ์", "มีนาคม", "เมษายน", "พฤษภาคม", "มิถุนายน",
	"กรกฎาคม", "สิงหาคม", "กันยายน", "ตุลาคม", "พฤศจิกายน", "ธันวาคม",
}

func ThaiDate(title string, t time.Time, taskID uint, userID string) string {
	weeekday := thaiDays[t.Weekday()]
	month := thaiMonths[int(t.Month()-1)]
	return fmt.Sprintf(
		"✅ บันทึกงานของเธอเรียบร้อยแล้ว! <@%s>\n"+
			"━━━━━━━━━━━━━━\n"+
			"📌 หัวข้อ: **%s**\n"+
			"📆 กำหนดส่ง: **วัน%s ที่ %d %s**\n"+
			"⏰ เวลา: **%02d:%02d น.**\n"+
			"🔢 TaskID: `%d`\n"+
			"━━━━━━━━━━━━━━",
		userID, title, weeekday, t.Day(), month, t.Hour(), t.Minute(), taskID,
	)
}
func ShowThaiDate(taskID int, order int, title string, t time.Time) string {
	thaiZone, _ := time.LoadLocation("Asia/Bangkok")
	thaiDate := t.In(thaiZone)
	weekday := thaiDays[thaiDate.Weekday()]
	month := thaiMonths[int(thaiDate.Month()-1)]
	return fmt.Sprintf(
		"\n━━━━━━━━━━━━━━\n"+
			"🗂️ งานที่ %d: **%s**\n"+
			"📅 กำหนดเสร็จ: **วัน%s ที่ %d %s**\n"+
			"⏰ เวลา: **%02d:%02d**\n"+
			"🔢 TaskID: `%d`\n",
		order, title, weekday, thaiDate.Day(), month, thaiDate.Hour(), thaiDate.Minute(), taskID,
	)

}

func NotifyTaskDueSoon(title string, t time.Time, taskID uint) string {
	thaiZone, _ := time.LoadLocation("Asia/Bangkok")
	thaiDate := t.In(thaiZone)
	weekday := thaiDays[t.Weekday()]
	month := thaiMonths[int(t.Month()-1)]
	return fmt.Sprintf(
		`✅ แจ้งเตือนงานของเธอใกล้ถึงเวลาแล้ว
━━━━━━━━━━━━━━
📌 หัวข้อ: %s
📆 กำหนดส่ง: วัน%s ที่ %d %s
⏰ เวลา: %02d:%02d น.
🔢 TaskID: %d
━━━━━━━━━━━━━━
`,
		title, weekday, t.Day(), month, thaiDate.Hour(), thaiDate.Minute(), taskID,
	)
}
