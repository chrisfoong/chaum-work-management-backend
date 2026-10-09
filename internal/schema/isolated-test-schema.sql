-- ISOLATED TEST DATABASE ONLY. Never run on Supabase.
-- เปิดใช้งาน Extension สำหรับ UUID (จำเป็นสำหรับ Supabase)

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

  

-- ==========================================

-- 1. CREATE ENUM TYPES

-- ==========================================

CREATE TYPE role_enum AS ENUM ('supervisor', 'assistant', 'worker');

CREATE TYPE leave_status_enum AS ENUM ('pending', 'approved', 'rejected');

CREATE TYPE contract_status_enum AS ENUM ('registered', 'active', 'complete', 'cancelled');

CREATE TYPE expense_type_enum AS ENUM ('fund_transfer', 'actual_expense');

CREATE TYPE invoice_status_enum AS ENUM ('pending', 'paid');

CREATE TYPE shift_status_enum AS ENUM ('scheduled', 'completed', 'cancelled');

CREATE TYPE attendance_status_enum AS ENUM ('on_time', 'late', 'absent', 'leave');

CREATE TYPE requisition_status_enum AS ENUM ('pending_survey', 'pending_procurement', 'pending_fund', 'pending_approval', 'completed', 'rejected');

CREATE TYPE requisition_type_enum AS ENUM ('tor_base', 'additional');

  

-- ==========================================

-- 2. CREATE TABLES (ลำดับการสร้างเพื่อไม่ให้ติดปัญหา Foreign Key)

-- ==========================================

  

-- 2.1 ตารางผู้ใช้งาน (ใช้ "USER" ในเครื่องหมายคำพูดเพราะเป็น Reserved Keyword ของ PostgreSQL)

CREATE TABLE "USER" (

user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

first_name VARCHAR(100) NOT NULL,

last_name VARCHAR(100) NOT NULL,

role role_enum NOT NULL DEFAULT 'worker',

phone_number VARCHAR(10) UNIQUE NOT NULL,

line_id VARCHAR(50) UNIQUE NOT NULL,

bank_name VARCHAR(100) NOT NULL,

bank_account_no VARCHAR(50) NOT NULL,

daily_wage DECIMAL(10,2) DEFAULT 400.00,

is_active BOOLEAN NOT NULL DEFAULT TRUE,

created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP

);

  

-- 2.1.1 ตารางพนักงาน (WORKER) - สกัดออกมาเพื่อรองรับ FK ที่ระบุว่ามาจาก WORKER

CREATE TABLE WORKER (

worker_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

user_id UUID NOT NULL REFERENCES "USER"(user_id) ON DELETE CASCADE,

is_available BOOLEAN NOT NULL DEFAULT TRUE

);

  

-- 2.2 ตารางพื้นที่ปฏิบัติงาน

CREATE TABLE LOCATION (

location_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

location_name VARCHAR(255) NOT NULL,

address TEXT,

latitude DECIMAL(10,7),

longitude DECIMAL(10,7)

);

  

-- 2.3 ตารางอุปกรณ์และเครื่องมือ

CREATE TABLE EQUIPMENT (

equipment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

equipment_name VARCHAR(255) NOT NULL,

is_active BOOLEAN NOT NULL DEFAULT TRUE

);

  

-- 2.4 ตารางบันทึกการขอลาหยุด

CREATE TABLE LEAVE_REQUEST (

request_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

leave_no VARCHAR(50) UNIQUE NOT NULL,

user_id UUID NOT NULL REFERENCES "USER"(user_id) ON DELETE CASCADE,

leave_date DATE NOT NULL,

reason TEXT NOT NULL,

is_advance_notice BOOLEAN NOT NULL DEFAULT FALSE,

status leave_status_enum NOT NULL DEFAULT 'pending',

created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP

);

  

-- 2.5 ตารางสัญญาจ้าง

CREATE TABLE CONTRACT_TOR (

tor_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

contract_no VARCHAR(100) UNIQUE NOT NULL,

project_name VARCHAR(255) NOT NULL,

user_id UUID NOT NULL REFERENCES "USER"(user_id),

partner_agency VARCHAR(255) NOT NULL,

start_date DATE NOT NULL,

end_date DATE NOT NULL,

contract_value DECIMAL(12,2) NOT NULL,

contract_file_url VARCHAR(255),

status contract_status_enum NOT NULL DEFAULT 'registered',

created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP

);

  

-- 2.6 ตารางใบวางบิล

CREATE TABLE COMPANY_INVOICE (

invoice_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

invoice_no VARCHAR(50) UNIQUE NOT NULL,

tor_id UUID NOT NULL REFERENCES CONTRACT_TOR(tor_id) ON DELETE CASCADE,

billing_month VARCHAR(7) NOT NULL,

expected_amount DECIMAL(12,2) NOT NULL,

net_received DECIMAL(12,2),

exat_deduction_amount DECIMAL(12,2) DEFAULT 0.00,

deduction_reason TEXT,

status invoice_status_enum NOT NULL DEFAULT 'pending'

);

  

-- 2.7 ตารางการมอบหมายสถานที่ในสัญญา

CREATE TABLE TOR_LOCATION_ASSIGNMENT (

assignment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

tor_id UUID NOT NULL REFERENCES CONTRACT_TOR(tor_id) ON DELETE CASCADE,

location_id UUID NOT NULL REFERENCES LOCATION(location_id) ON DELETE RESTRICT,

required_workers INT NOT NULL DEFAULT 0

);

  

-- 2.8 ตารางคำขอจัดซื้ออุปกรณ์

CREATE TABLE EQUIPMENT_REQUISITION (

requisition_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

requisition_no VARCHAR(50) UNIQUE NOT NULL,

requested_by UUID NOT NULL REFERENCES "USER"(user_id),

assignment_id UUID NOT NULL REFERENCES TOR_LOCATION_ASSIGNMENT(assignment_id) ON DELETE CASCADE,

reason TEXT NOT NULL,

status requisition_status_enum NOT NULL DEFAULT 'pending_survey',

requisition_type requisition_type_enum NOT NULL,

reviewed_by UUID REFERENCES "USER"(user_id),

reviewed_at TIMESTAMPTZ,

created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP

);

  

-- 2.9 ตารางรายการอุปกรณ์ในใบจัดซื้อ

CREATE TABLE REQUISITION_ITEM (

item_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

requisition_id UUID NOT NULL REFERENCES EQUIPMENT_REQUISITION(requisition_id) ON DELETE CASCADE,

equipment_id UUID NOT NULL REFERENCES EQUIPMENT(equipment_id) ON DELETE RESTRICT,

required_qty INT NOT NULL DEFAULT 1,

existing_qty INT,

-- Database จะคำนวณ to_buy_qty ให้อัตโนมัติ ป้องกันความผิดพลาดจากการ Insert ฝั่ง Backend

to_buy_qty INT GENERATED ALWAYS AS (GREATEST(required_qty - COALESCE(existing_qty, 0), 0)) STORED,

actual_qty INT,

actual_price DECIMAL(10,2),

remark TEXT

);

  

-- 2.10 ตารางรายการสำรองจ่าย

CREATE TABLE EXPENSE_CLAIM (

expense_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

expense_no VARCHAR(50) UNIQUE NOT NULL,

expense_type expense_type_enum NOT NULL,

user_id UUID REFERENCES "USER"(user_id),

requisition_id UUID REFERENCES EQUIPMENT_REQUISITION(requisition_id) ON DELETE CASCADE,

total_amount DECIMAL(10,2) NOT NULL CHECK (total_amount > 0),

transfer_ref_no VARCHAR(100) UNIQUE, -- ถอด NOT NULL ออกเพื่อรองรับใบเสร็จเงินสด

receipt_photo_url VARCHAR(255) NOT NULL,

created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP

);

  

-- 2.11 ตารางการจัดกะตารางงาน

CREATE TABLE WORK_SCHEDULE (

schedule_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

assignment_id UUID NOT NULL REFERENCES TOR_LOCATION_ASSIGNMENT(assignment_id) ON DELETE CASCADE,

worker_id UUID NOT NULL REFERENCES WORKER(worker_id) ON DELETE CASCADE,

work_date DATE NOT NULL,

shift_start_time TIME NOT NULL,

shift_status shift_status_enum NOT NULL DEFAULT 'scheduled',

-- ล็อคป้องกันพนักงาน 1 คน มีตารางงานเกิน 1 กะในวันเดียวกัน

UNIQUE (worker_id, work_date)

);

  

-- 2.12 ตารางบันทึกการเข้าออกงาน

CREATE TABLE ATTENDANCE (

attendance_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

schedule_id UUID NOT NULL UNIQUE REFERENCES WORK_SCHEDULE(schedule_id) ON DELETE CASCADE, -- 1-to-1

worker_id UUID NOT NULL REFERENCES WORKER(worker_id) ON DELETE CASCADE,

work_date DATE NOT NULL,

check_in TIMESTAMPTZ,

check_out TIMESTAMPTZ,

status attendance_status_enum

);

  

-- 2.13 ตารางบันทึกหลักฐานการทำงาน

CREATE TABLE WORK_EVIDENCE (

evidence_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

schedule_id UUID NOT NULL REFERENCES WORK_SCHEDULE(schedule_id) ON DELETE CASCADE,

worker_id UUID NOT NULL REFERENCES WORKER(worker_id) ON DELETE CASCADE,

description TEXT,

photo_url VARCHAR(255) NOT NULL,

submitted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP

);

  

-- 2.14 ตารางรอบการจ่ายเงิน

CREATE TABLE PAYROLL (

payroll_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

payroll_slip_no VARCHAR(50) UNIQUE NOT NULL,

user_id UUID NOT NULL REFERENCES "USER"(user_id) ON DELETE CASCADE,

managed_by_id UUID NOT NULL REFERENCES "USER"(user_id),

period_start DATE NOT NULL,

period_end DATE NOT NULL,

base_wage DECIMAL(10,2) NOT NULL,

total_deduction DECIMAL(10,2) NOT NULL DEFAULT 0.00,

net_wage DECIMAL(10,2) NOT NULL,

is_paid BOOLEAN NOT NULL DEFAULT FALSE,

created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP

);

  

-- 2.15 ตารางประวัติการหักเงิน

CREATE TABLE DEDUCTION_TRANSACTION (

deduction_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

worker_id UUID NOT NULL REFERENCES WORKER(worker_id) ON DELETE CASCADE,

attendance_id UUID REFERENCES ATTENDANCE(attendance_id) ON DELETE SET NULL,

payroll_id UUID REFERENCES PAYROLL(payroll_id) ON DELETE SET NULL, -- Nullable

penalty_amount DECIMAL(10,2) NOT NULL,

reason TEXT NOT NULL,

created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP

);