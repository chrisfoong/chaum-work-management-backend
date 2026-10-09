# เปิด Backend และสถานะการทดสอบ

อัปเดต 10 ตุลาคม 2026 (Asia/Bangkok) — repo chaum-work-management-backend ใช้ branch รวมงาน develop

## เปิด server

Git Bash (Terminal 1):

```bash
cd /c/SA/chaum-work-management-backend
go run ./cmd/server
```

PowerShell:

```powershell
cd C:\SA\chaum-work-management-backend
go run ./cmd/server
```

คง Terminal ไว้ระหว่างใช้งาน โปรแกรมโหลด config ตามปกติจากเครื่องผู้ใช้ คู่มือนี้ไม่มี secrets และไม่ต้องเปิดหรือแก้ .env เพื่อรันคำสั่ง หากผู้ใช้เปลี่ยน config ต้อง Ctrl+C แล้วเปิด server ใหม่

พอร์ตรอบทดสอบคือ 8080 ตรวจในอีก Terminal:

```powershell
curl.exe -i http://127.0.0.1:8080/health/ready
```

ต้องได้ HTTP 200 และ database up; LISTENING อย่างเดียวไม่ยืนยันฐานพร้อม API ที่ต้องยืนยันตัวตนเมื่อไม่มี token ต้องได้ 401

ถ้าพอร์ตถูกใช้แล้ว:

```powershell
netstat -ano | findstr :8080
Get-Process -Id <PID>
```

ตรวจว่าเป็น server เดิมก่อนหยุดด้วย Ctrl+C หรือ Stop-Process -Id <PID> ใน PowerShell; ใน Git Bash ใช้ powershell.exe -NoProfile -Command "Stop-Process -Id <PID>" อย่าหยุด process อื่นโดยเดา PID

## Config และเปิดคู่ Frontend

ชื่อเหล่านี้เป็นชื่อตัวแปร ไม่ใช่ค่าจริง:

- DATABASE_URL เป็น PostgreSQL connection string; ใช้แทน HTTPS SUPABASE_URL ไม่ได้
- SUPABASE_URL, SUPABASE_SERVICE_ROLE_KEY, STORAGE_BUCKET สำหรับ private Storage; bucket ต้องมีจริง และ service key อยู่ Backend เท่านั้น
- QR_SIGNING_SECRET ค่าลงนามอย่างน้อย 32 bytes; ไม่ใส่ใน QR หรือ Frontend
- LINE_MESSAGING_TOKEN ของ Messaging API; LINE Login channel secret ไม่ใช่ token นี้
- LINE Web/Worker channel/provider allowlist และ ASSISTANT_DASHBOARD_URL ต้องตรง deployment จริง

ผู้ใช้จัดการค่าจริงเองโดยอิง .env.example placeholders ห้าม commit secrets หรือรัน migrations/seeds/schema initialization บน Supabase เดิม

เปิด Frontend อีก Terminal ที่ 5500 ตั้ง BACKEND_URL ใน Frontend เป็น origin เช่น http://127.0.0.1:8080 ไม่เติม /api LINE จริงใช้ HTTPS tunnel ของ Frontend คง Backend, Frontend, tunnel พร้อมกัน ดูคู่มือใน frontend repo docs/SERVER_AND_TEST_STATUS.md

## หลักฐานที่ผ่าน

- 9 ต.ค.: Backend 235 tests/subtests ผ่าน fail 0 test skip 0 บน PostgreSQL disposable ที่ loopback ไม่ใช่ Supabase; vet/build ผ่าน
- Login LINE จริง Supervisor/Assistant/Worker ผ่าน โดยสลับ role user เดิมชั่วคราว; DB role และปฏิเสธแพลตฟอร์มผิดผ่าน
- อ่าน Dashboard/คำขอ TOR/9A จริงผ่าน Worker schedules/payslips อ่านได้แต่รายการว่าง
- ออก QR จาก Backend จริงและแสดงหมดอายุผ่าน; ไม่ใช่หลักฐาน GPS+สแกน QR เช็กอิน
- คืน user 11111111-1111-1111-1111-111111111111 เป็น Supervisor และลบ Worker mapping ทดสอบแล้ว ไม่มีกะทดสอบเหลือให้ใช้

## ยังไม่ได้ตรวจครบผ่านระบบจริง

| Flow              | ต้องตรวจต่อ                                                                                                                  |
| ----------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| TOR/Storage       | PNG upload → ยืนยัน TOR → อ่าน/ดาวน์โหลด private Storage; ชนิด/ขนาด/สิทธิ์                                                   |
| 1A/2S/2A          | สำรวจ → ยอดขาด → บันทึก funding/slip → partial/repeated purchase → completed; ไม่ต้องโอนเงินจริงเพื่อทดสอบการบันทึก          |
| 5W/5A/6A/3S/7A/8A | Worker ขอเพิ่ม → Assistant ตรวจ → approve/reject/no_purchase → funding/purchase → ส่งมอบ                                     |
| 3A/4A/2W          | สัญญาใช้งานได้/อุปกรณ์พร้อม/คนขั้นต่ำ → กะ → ลา → อนุมัติพร้อมคนแทนหรือปฏิเสธ                                                |
| 3W/4W             | Worker mapping+กะ → กล้อง/ตำแหน่งจริงในเกณฑ์ → QR ไม่หมดอายุ → check-in → ภาพผลงาน → check-out; duplicate/ผิดพื้นที่/หมดอายุ |
| 4S/5S/7W          | populated attendance → payroll/รายการหัก → paid flag → Worker paid slip; ตรวจทำซ้ำไม่ซ้ำยอด                                  |
| 6S                | ใบวางบิล/รับเงิน/ต้นทุน → รายงานตามช่วง → PDF/CSV ยอดตรงกัน; PDF ภาษาไทยมีข้อจำกัด                                           |
| 6W/9A             | ส่งบัญชีทดสอบที่เพิ่มเพื่อน OA → ตรวจแชตผู้รับจริง → เปิดลิงก์ Assistant; API accepted ไม่เท่ากับ delivered/read             |

เทสอัตโนมัติจำลอง LINE/Storage ไม่ยืนยันกล้อง/GPS/รับข้อความจริง รายการว่างและจำนวน tests ไม่ใช่หลักฐาน end-to-end ครบ

## ตรวจอัตโนมัติ

```powershell
go test ./...
go vet ./...
go build ./...
```

Integration ต้องตั้ง TEST_DATABASE_URL และ TEST_DATABASE_ISOLATED=yes ให้ชี้ disposable PostgreSQL ที่ loopback เท่านั้น ดู README และ docker-compose.test.yml ห้ามใช้ Supabase เดิม หากไม่มีฐานแยก integration จะ skip และยังยืนยันไม่ได้

หลักฐานรอบก่อนอยู่ workspace ผู้ทดสอบแยกจาก Git: backend-role-recheck-tests.jsonl, SUPERVISOR_LIVE_TEST_2026-10-09.md และภาพ role/QR
