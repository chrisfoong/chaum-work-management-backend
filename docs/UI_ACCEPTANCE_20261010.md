# ผลตรวจ Use Case และหน้าจอ 10 ตุลาคม 2569

ยังไม่ผ่าน End to end ครบ 22 use cases และยังไม่ได้ภาพใหม่ครบทุกกรณี ห้ามใช้จำนวน automated tests แทนการยืนยันรับ LINE หรือการทดสอบมือถือจริง

Backend: 242 tests/subtests ผ่าน, fail 0, skip 0 บน PostgreSQL localhost ฐาน chaum_ui_acceptance_20261010 แยกจาก Supabase; go vet/build ผ่าน Frontend 33 tests ผ่าน ไม่มี fail/skip, lint/typecheck/format/build ผ่านก่อนแก้ข้อความ LINE เพิ่มเติม; tests หลังแก้ LINE ผ่านอีกครั้ง

## การแก้หน้าจอ

ลบคำ Backend API Snapshot รหัส use case และลิงก์ preview จาก production login; preview ยังระบุว่าไม่บันทึกข้อมูลเพื่อไม่ให้สับสน แสดงสถานะ LINE เป็นผลส่งคำขอ ไม่ยืนยันการรับ เปลี่ยนเครื่องหมายสำเร็จเป็นคำเตือนเมื่อแจ้งเตือนล้มเหลว Dashboard นับเฉพาะสัญญา active ที่อยู่ในช่วงวัน และแสดงคำขอที่ยังไม่ completed/rejected พร้อมเลขที่สัญญา

## ผลรายกรณี

| Use case | หน้าที่ | หลักฐานและขอบเขต | Automated test ที่ผ่าน |
|---|---|---|---|
| 1S | สร้าง TOR และแนบหลักฐาน | ผ่าน UI และ Storage ก่อนหน้า ภาพใหม่ 1S-saved.jpg | TestIntegrationMVP; TestContractPNGVerificationChecksBytesOwnerAndSize |
| 2S | บันทึกเงินจัดซื้อ | ผ่าน UI 100/150 บาทและ PNG ไม่มีโอนเงินจริง | TestIntegrationPurchaseRoundsDeliveryAndPaidReport |
| 3S | อนุมัติจัดซื้อเพิ่มเติม | ผ่าน UI คำขอ 002 ยังไม่ยืนยันรับ LINE เหตุการณ์นี้ | TestIntegrationUsecaseCompletion; TestMutationRoleBoundaries |
| 4S | เวลาเข้างานและเงินหัก | เช็กอินจริงมีสถานะสาย กรณีขาดและเกณฑ์หักผ่าน automated tests | TestPenaltyBoundaries; TestIntegrationAtomicAttendanceAndRollback |
| 5S | ประมวลผลค่าจ้าง | ผ่าน UI รอบ 16–30 กันยายน สลิปยังไม่จ่าย การจ่ายและรับแจ้งจริงค้าง | TestIntegrationPayrollRetainsScheduledWorkAfterRoleChange |
| 6S | รายงานต้นทุนกำไร | ผ่าน UI และเปิดตรวจ PDF/CSV ยอด 0/0/300/-300 ตรงกัน | TestProfitPDFPaginationAndFinancialFields; TestIntegrationPurchaseRoundsDeliveryAndPaidReport |
| 1A | สำรวจอุปกรณ์ตั้งต้น | ผ่าน UI คำขอ 001 | TestIntegrationMVP |
| 2A | บันทึกจัดซื้ออุปกรณ์ตั้งต้น | ผ่าน UI พร้อมใบเสร็จ 1 หน่วย 100 บาท | TestIntegrationMVP |
| 3A | จัดตารางคนงาน | ผ่าน UI และตรวจฐานวันที่ 10/11 ตุลาคม | TestIntegrationMVP; TestPeriodsAndOvernight |
| 4A | พิจารณาลาและคนแทน | ผ่าน isolated integration test UI บัญชีเดิมยังค้าง | TestIntegrationEmergencyLeaveAndAtomicReplacement |
| 5A | ตรวจเทียบ TOR | ผ่าน UI คำขอ 002 | TestIntegrationUsecaseCompletion |
| 6A | พิจารณาคำขออุปกรณ์ | ผ่าน UI ทั้งซื้อและไม่ซื้อ ยังไม่ยืนยันรับ LINE ทุกเหตุการณ์ | TestIntegrationUsecaseCompletion |
| 7A | ซื้อเพิ่มเติมหลายรอบ | ผ่าน UI 50+150 บาท ปริมาณ 2 ราคาเฉลี่ย 100 บาท | TestIntegrationPurchaseRoundsDeliveryAndPaidReport |
| 8A | บันทึกส่งมอบ | ผ่าน UI และบันทึก PNG คำขอ 002 ยังรอยืนยันรับ LINE เฉพาะเหตุการณ์ | TestIntegrationFrontendReadModelsAndDeliveryGuards |
| 9A | อ่านและส่งสรุปปฏิบัติการ | อ่านรายงานผ่าน UI ส่งแบบกดเองยังไม่ผ่านผู้รับจริง มี LINE บัญชีเดียว | TestIntegrationOperationsSummary; TestIntegrationManualOperationsNotification |
| 1W | ดูตารางงาน | มือถือใช้งานกะเดิมได้ อ่าน ongoing/ownership ผ่าน automated tests | TestIntegrationWorkerReadOngoingOvernightShift |
| 2W | แจ้งลา | ผ่าน isolated integration test ยังไม่ได้ยืนยันกดส่งจากมือถือรอบนี้ | TestIntegrationEmergencyLeaveAndAtomicReplacement |
| 3W | GPS และ QR เช็กอิน | ผู้ใช้ยืนยันมือถือจริง ตรวจ check_in ในฐานแล้ว ไม่ใช้ GPS mock เป็นหลักฐาน | TestGPSAndQRBoth; TestIntegrationAtomicAttendanceAndRollback |
| 4W | รูปผลงานและเช็กเอาท์ | ผู้ใช้ยืนยันแล้ว พบ check_out และ Storage PNG HTTP 200 | TestIntegrationUsecaseCompletion; TestStorageDoesNotAcceptOtherOwnersOrRedirect |
| 5W | แจ้งอุปกรณ์ขาด | ผ่าน isolated integration test ยังไม่ได้ยืนยันมือถือส่งหลังเช็กอินจริง | TestIntegrationUsecaseCompletion; TestWorkerEquipmentReasonLimit |
| 6W | รับข้อความ LINE | ข้อความ TEST โดยตรงได้รับแล้ว แต่ยังไม่ยืนยันรับทุกเหตุการณ์ของ use case | TestPushUsesAuthenticatedPostAndRejectsFailures; TestLINEBoundedRetriesReuseKey |
| 7W | ดูสลิปค่าจ้าง | ownership และ paid-only ผ่าน automated tests UI สลิปที่จ่ายแล้วของบัญชีเดิมค้าง | TestIntegrationUsecaseCompletion |

## สิ่งที่ติด

LINE callback ใน Codex browser แจ้ง Failed to fetch หลัง token หมดอายุและ login ใหม่ ตรวจใน LINE Developers พบ Endpoint HTTPS โดเมนเดิมและ scopes openid/profile ถูกต้อง ทดลอง production Docker แล้วยังไม่สำเร็จ ยังไม่ทราบสาเหตุแน่ชัด ไม่ bypass การยืนยันตัวตน ไม่สร้างบัญชี LINE ซ้ำเพื่ออ้างว่า 9A ส่งถึง Assistant จริง

## รูปใหม่ที่ตรวจแล้ว

- dashboard.jpg: Supervisor Dashboard ข้อมูลจริงหลังปรับข้อความ สัญญาที่ดำเนินงาน 1 คำขอค้าง 2
- 1S-form.jpg: ฟอร์มสร้าง TOR เป็นภาพฟอร์ม ไม่ถือว่ากดบันทึกผ่านจากภาพนี้
- 1S-saved.jpg: รายละเอียดสัญญาที่บันทึกแล้ว 9910100001
- 4S-attendance.jpg: session หมดอายุ เป็นหลักฐานข้อผิดพลาด ไม่ถือเป็น Normal Case ที่ผ่าน

เอกสาร Google Docs แท็บ User interface ยังคงภาพประวัติเดิม และเพิ่มผลล่าสุดกับภาพใหม่เฉพาะที่ตรวจแล้ว

## ทะเบียนชุดข้อมูลและบริการ

ชื่อชุดทดลองเปลี่ยนเฉพาะ project_name/partner_agency/location_name ให้เป็นภาษางาน ไม่เปลี่ยนรหัส พิกัด เวลาเข้างานหรือหลักฐาน สัญญา f3370bbf-ca14-4506-89bc-cda9808635c2 เลข 9910100001 ยังคงเป็นข้อมูลทดลอง ไม่มีการโอนเงินจริง

เตรียมคำขอ disposable ใหม่ 004 ID 44385bd3-8d30-4c4d-83e4-d89f0972ef8f และ 005 ID f22d8e13-f3a8-4e77-a649-a5d7639cb67a แต่ยังไม่ดำเนินขั้นตอน UI เพราะ LINE callback ติด อย่านับเป็นการสำรวจ/ซื้อผ่าน

Docker frontend 5501 backend 8081; native dev 5500 backend 8080 ยังเปิดอยู่ โดเมน ngrok เดิมเปลี่ยนให้เข้า production Docker 5501 ระหว่างวินิจฉัย callback


คำขอ disposable ใหม่ 004/005 ถูกล้างแบบจำกัดรหัสและตรวจว่าไม่มีธุรกรรมเงินแล้ว ไม่ใช่การสำรวจหรือจัดซื้อผ่าน จึงไม่ทิ้งคำขอค้างที่ขวางพื้นที่เดิม บัญชีเดิมคืนเป็น Supervisor หลังตรวจยืนยัน
