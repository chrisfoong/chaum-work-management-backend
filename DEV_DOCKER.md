# เปิด Backend

รันแบบพัฒนาใน repo นี้:

```bash
go run ./cmd/server
```

หรือใช้ Docker Desktop (Linux containers) และ `.env` ที่ตั้งค่าจาก `.env.example` แล้ว:

```bash
docker compose build
docker compose up --build -d
docker compose ps
docker compose logs -f backend
docker compose down
```

ตรวจ health ที่ http://127.0.0.1:8080/health ไม่รัน native และ Docker บนพอร์ต 8080 พร้อมกัน สามารถกำหนด BACKEND_PORT เป็นพอร์ตอื่นสำหรับ host ได้ แต่ container ยังคงใช้ 8080

Dockerfile ใช้ Go 1.26 build เป็น Linux binary และรันโดย user ที่ไม่ใช่ root พร้อม CA certificates/timezone สำหรับ Supabase/LINE/Timestamps `.env` ไม่ถูก COPY เข้า image Compose โหลด config ตอนรัน เปลี่ยน config แล้วใช้ `docker compose up -d --force-recreate`

ไม่มี database container, migration, seed หรือ reset ใน stack นี้ ใช้ Supabase ที่ระบุใน DATABASE_URL ตามปกติ งาน background ของ Backend อาจประมวลผลข้อมูลเมื่อเปิด server

ถ้าต้องการเปิด Frontend ด้วย ให้วางสอง repo ข้างกัน แล้วรัน Compose จาก repo Frontend ดู [คู่มือ Frontend](../chaum-work-management-frontend/DEV_DOCKER.md) Frontend ใช้ service hostname `backend` ภายใน Docker network และแสดงเว็บที่พอร์ต 5500
