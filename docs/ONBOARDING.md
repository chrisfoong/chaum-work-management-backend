# เพิ่มผู้ใช้ด้วยลิงก์เชิญ LINE โดยไม่แก้ ER

## วิธีใช้

1. Supervisor เปิด **ข้อมูลหลัก → บุคลากร → เพิ่มบุคลากร**
2. กรอกข้อมูลบุคลากร เลือก Supervisor / Assistant / Worker และค่าจ้าง
3. กด **สร้างลิงก์เชิญ** แล้วคัดลอกส่งให้เจ้าของบัญชีโดยตรง
4. ผู้รับเปิดลิงก์ใน LINE บนมือถือ ล็อกอิน และตรวจชื่อกับบทบาท
5. ผู้รับกด **ยืนยันว่าเป็นบัญชีของฉันและเชื่อม LINE** แล้วเข้าใช้งาน

Worker เปิดผ่าน Worker LIFF; Supervisor/Assistant เปิดผ่าน Web LIFF ตามลิงก์ที่ระบบสร้าง
Worker ต้องได้รับการจัดกะก่อนเช็กอิน คอมต้องเปิด Docker และ HTTPS tunnel ค้างไว้

## ลิงก์และข้อมูลเดิม

- ไม่มีตาราง/column/enum/migration ใหม่ และไม่ต้องแก้ ER
- เตรียม USER เดิมเป็น is_active=false; line_id เก็บ `INV:` ตามด้วย SHA-256 digest
  แบบ base64url (47 ตัวอักษร อยู่ใน varchar(50) เดิม) ไม่เก็บ raw invitation
- เมื่อผูกสำเร็จ เปลี่ยน line_id เป็น LINE subject ที่ Backend ตรวจจาก ID token จริง
  เปิด is_active=true และสร้าง WORKER หนึ่งรายการเฉพาะ role Worker
- บทบาทและค่าจ้างเลือกโดย Supervisor ผู้รับไม่ส่ง role/LINE subject เอง
- ลิงก์สุ่ม 256 bits หมดอายุ 24 ชั่วโมง ใช้ได้ครั้งเดียว และไม่มี LINE token ในลิงก์
- คัดลอกลิงก์ได้ตอนสร้างเท่านั้น หากลืมหรือหมดอายุ ใช้ **ออกลิงก์เชิญใหม่** ในตาราง
  ลิงก์เก่าจะถูกยกเลิกทันที บัญชีที่ผูก LINE แล้วหรือถูกระงับออกลิงก์เชิญใหม่ไม่ได้
- ก่อนผูก LINE ตารางแสดง **รอผูก LINE** แยกจาก **ระงับ** การแก้ข้อมูลบุคลากร
  ไม่สามารถเปิดใช้งานข้ามขั้นตอนผูก LINE ได้
- ข้อมูลธนาคารยังเก็บใน USER ตาม schema เดิม; ไม่แสดง token/digest/ธนาคารในรายการบุคลากร
- LINE หนึ่งบัญชีผูกกับ USER หนึ่งรายการ ผู้ที่ผูกอยู่แล้วไม่สามารถเปลี่ยนตัวตน/role
  ด้วยการกดลิงก์ของอีกคนได้

## API

- `POST /api/web/invitations`: Supervisor เท่านั้น สร้าง inactive USER และคืนลิงก์ข้อมูล
- `POST /api/web/users/:id/invitation`: Supervisor เท่านั้น เปลี่ยนรหัสสำหรับบัญชีที่ยังไม่ผูก
- `GET /api/{web|liff}/onboarding`: ตรวจ LINE token แล้วคืนสถานะของผู้เรียกเอง
- `POST /api/{web|liff}/invitation/preview`: token ใน JSON body; ต้องตรวจ LINE ก่อน
- `POST /api/{web|liff}/invitation/claim`: ผูก verified LINE subject ภายใน transaction

หน้าแอปวางรหัสใน URL fragment และเก็บเฉพาะรหัสเชิญชั่วคราวใน sessionStorage ระหว่าง
LINE redirect; LINE ID token ไม่ถูกเก็บใน storage URL parameters/fragment ถูกกู้หลัง
liff.init ตาม [LINE documentation](https://developers.line.biz/en/tips/2026/07/16/liff-url-additional-info/)
ฝั่ง API รับรหัสใน body และ logger เดิมบันทึกเพียง route ไม่บันทึก token/body

## Verification

Tests ตรวจ token หมดอายุ/รูปแบบ/digest, LINE token ที่ผิด, ปลอม role/subject, สิทธิ์
สร้างและออกลิงก์ใหม่, ลิงก์เก่าหลัง renew, ผูกซ้ำ/ผูกพร้อมกัน, LINE ที่มีบัญชีแล้ว,
ป้องกันเปิดใช้งานเอง, หนึ่ง WORKER และทั้งสาม role โดยใช้ isolated loopback PostgreSQL
schema เดิมเท่านั้น Frontend tests ตรวจ LIFF state/fragment/channel และ proxy allowlist

การผ่าน tests แบบ mock และ isolated DB ไม่ใช่หลักฐานว่าคนใหม่ล็อกอิน LINE จริงผ่านแล้ว
การตรวจบนมือถือจริงยังต้องมี LINE อีกบัญชีของผู้รับเชิญ ไม่ใช้การสลับ role ของบัญชีเดิม
แทนหลักฐาน onboarding

## ยกเลิกลิงก์เชิญ

Supervisor กด **ยกเลิกลิงก์เชิญ** ในตารางบุคลากร แล้วยืนยัน ลิงก์ปัจจุบันจะใช้ไม่ได้
ทันที โดยยังเก็บบุคลากรไว้เป็นบัญชีที่ยังไม่เปิดใช้งาน ตารางแสดง **ยกเลิกลิงก์แล้ว**
และยังออกลิงก์ใหม่ได้ภายหลัง การออกลิงก์ใหม่เองก็ยกเลิกลิงก์เก่าอัตโนมัติ

ใช้ `POST /api/web/users/:id/invitation/cancel` สำหรับ Supervisor เท่านั้น บันทึก digest
ใน USER.line_id ด้วย prefix CNL: แทน INV: โดยไม่เปลี่ยน schema บัญชีที่ผูก LINE แล้ว
ยกเลิกผ่าน endpoint นี้ไม่ได้ และ cancel/claim ล็อก USER เดียวกันเพื่อป้องกัน race
การทดสอบ isolated DB ครอบ repeat cancel, cancelled-link replay, reissue, unauthorized
roles, manual activation prevention และการ cancel/claim พร้อมกัน
