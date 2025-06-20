# 🤖 PactchNote – Discord AI Note Assistant (with Gemini + PostgreSQL)

**🇹🇭 PactchNote คือผู้ช่วย AI สำหรับ Discord ที่จะช่วยจดบันทึกงานของคุณแบบอัตโนมัติ**  
เชื่อมต่อกับฐานข้อมูล PostgreSQL และใช้ Google Gemini API ในการสร้างโน้ตอย่างชาญฉลาด 🧠✨

**🇬🇧 PactchNote is a smart AI-powered note-taking assistant for Discord.**  
It connects to a PostgreSQL database and uses Google Gemini API to generate intelligent notes.

---

## 🚀 ความสามารถ / Features

- 📝 สร้างโน้ตจากข้อความธรรมดา / Generate notes from natural language
- 🔎 ค้นหางานตามช่วงเวลา / Search for tasks by time range
- 🗑️ ลบงานด้วยคำสั่ง `/deltask` / Delete tasks with `/deltask`
- ⏰ แจ้งเตือนก่อนเริ่มงาน 30 นาที / 30-minute pre-task reminders
- 💾 ใช้ฐานข้อมูล PostgreSQL / Uses PostgreSQL for storage
- 🧠 ใช้ Google Gemini API ในการประมวลผล / Uses Google Gemini for AI processing

---

## ⚙️ วิธีติดตั้ง / Setup

### 1. ติดตั้งสิ่งจำเป็น / Requirements

- [Go](https://go.dev/) (หากพัฒนาในเครื่อง / for local dev)
- [Docker](https://www.docker.com/) (แนะนำ / recommended)
- Discord Bot Token
- [Gemini API Key](https://makersuite.google.com/)
- PostgreSQL Database

---

### 2. สร้างไฟล์ `.env` / Create `.env` file

สร้างไฟล์ชื่อ `.env` และกำหนดค่าต่อไปนี้:  
Create a file named `.env` and set the following values:

```env
# 🔐 Discord Bot Token
Discord_Token=your_discord_bot_token

# 🗄️ PostgreSQL Configuration
Host=your_postgres_host         # โฮสต์ของฐานข้อมูล เช่น localhost หรือชื่อ container
Database=                       # ชื่อฐานข้อมูล
Username=                       # ชื่อผู้ใช้
Password=your_postgres_password # รหัสผ่าน
Port=5432                       # พอร์ตของ PostgreSQL
88
