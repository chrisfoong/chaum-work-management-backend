-- 0001_init: the 16 entities of the production ER (docs/02_domain_model.md).
-- NOT APPLIED. Review before running against any database.
--
-- Conventions: UUID keys (gen_random_uuid, built into PostgreSQL 13+), money as
-- NUMERIC, timestamps as timestamptz (UTC). user_id = worker_id (docs/04 §1.1).
-- Choices taken on open decisions are marked with D#; see docs/04 §2.

-- USER. "user" is reserved in PostgreSQL, so the table is "users".
-- Web users (supervisor, assistant): user_id is set to their Supabase Auth user id (D13).
-- Workers: line_id holds the LINE user id used to resolve LIFF logins.
CREATE TABLE users (
    user_id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name      text NOT NULL,
    last_name       text NOT NULL,
    phone_number    text,
    role            text NOT NULL CHECK (role IN ('supervisor', 'assistant', 'worker')),
    line_id         text UNIQUE,
    daily_wage      numeric(12, 2) CHECK (daily_wage >= 0),
    bank_name       text,
    bank_account_no text,
    is_active       boolean NOT NULL DEFAULT true,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

-- WORKER: shares USER's identity (user_id = worker_id); only worker-specific data.
CREATE TABLE worker (
    user_id      uuid PRIMARY KEY REFERENCES users (user_id),
    is_available boolean NOT NULL DEFAULT true
);

-- CONTRACT_TOR. user_id = USER MANAGES (R02). status value set is not specified.
CREATE TABLE contract_tor (
    tor_id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           uuid NOT NULL REFERENCES users (user_id),
    project_name      text NOT NULL,
    contract_no       text NOT NULL UNIQUE,
    partner_agency    text,
    contract_value    numeric(14, 2) CHECK (contract_value >= 0),
    start_date        date,
    end_date          date,
    contract_file_url text,
    status            text,
    created_at        timestamptz NOT NULL DEFAULT now(),
    updated_at        timestamptz NOT NULL DEFAULT now(),
    CHECK (end_date IS NULL OR start_date IS NULL OR end_date >= start_date)
);

-- LOCATION. 1S rejects a duplicate location name.
CREATE TABLE location (
    location_id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    location_name text NOT NULL UNIQUE,
    address       text,
    latitude      numeric(9, 6),
    longitude     numeric(9, 6)
);

-- TOR_LOCATION_ASSIGNMENT: contract + location pairing (R07 INCLUDES, R08 HOSTS).
CREATE TABLE tor_location_assignment (
    assignment_id    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tor_id           uuid NOT NULL REFERENCES contract_tor (tor_id),
    location_id      uuid NOT NULL REFERENCES location (location_id),
    required_workers integer NOT NULL CHECK (required_workers >= 0)
);

-- WORK_SCHEDULE: one worker, one date, one assignment (R10, R15).
-- 3A rule: no duplicate schedule for the same worker and date.
CREATE TABLE work_schedule (
    schedule_id      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id    uuid NOT NULL REFERENCES tor_location_assignment (assignment_id),
    worker_id        uuid NOT NULL REFERENCES worker (user_id),
    work_date        date NOT NULL,
    shift_start_time time NOT NULL,
    shift_status     text,
    UNIQUE (worker_id, work_date)
);

-- ATTENDANCE (R16 RECORDS, R17 LOGS, R18 REPLACES).
-- D11: unique/UPSERT key is (schedule_id, worker_id).
-- D11: replacement_worker_id is set only when the replacement worker accepts;
--      the exact population flow is decided in the attendance slice (P7).
-- status is NULL until the 4S daily job classifies the row.
CREATE TABLE attendance (
    attendance_id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_id           uuid NOT NULL REFERENCES work_schedule (schedule_id),
    worker_id             uuid NOT NULL REFERENCES worker (user_id),
    replacement_worker_id uuid REFERENCES worker (user_id),
    work_date             date NOT NULL,
    check_in              timestamptz,
    check_out             timestamptz,
    status                text CHECK (status IN ('on_time', 'late', 'leave', 'absent')),
    UNIQUE (schedule_id, worker_id),
    CHECK (check_out IS NULL OR check_in IS NULL OR check_out >= check_in)
);

-- LEAVE_REQUEST (R19 SUBMITS).
-- D14: keyed by user_id for now; switching to schedule_id is to be discussed.
-- The 2W duplicate-leave check is enforced in the service, not here (whether a
-- rejected leave may be re-submitted is not specified).
CREATE TABLE leave_request (
    request_id        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           uuid NOT NULL REFERENCES worker (user_id),
    leave_no          text NOT NULL UNIQUE,
    leave_date        date NOT NULL,
    reason            text NOT NULL,
    is_advance_notice boolean NOT NULL,
    status            text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    created_at        timestamptz NOT NULL DEFAULT now()
);

-- WORK_EVIDENCE (R11 COLLECTS, R20 CAPTURES).
-- D2: linked by assignment_id only (per the team's data dictionary); no schedule_id.
CREATE TABLE work_evidence (
    evidence_id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id uuid NOT NULL REFERENCES tor_location_assignment (assignment_id),
    worker_id     uuid NOT NULL REFERENCES worker (user_id),
    photo_url     text NOT NULL,
    description   text,
    submitted_at  timestamptz NOT NULL DEFAULT now()
);

-- PAYROLL (R01 RECEIVES via worker_id).
-- D1: no payer column (USER PAYS PAYROLL not stored).
-- net_wage = base_wage - total_deduction, stored as written by 5S.
-- D8 (re-run guard) is open: no unique key added.
CREATE TABLE payroll (
    payroll_id      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    worker_id       uuid NOT NULL REFERENCES worker (user_id),
    payroll_slip_no text UNIQUE,
    period_start    date NOT NULL,
    period_end      date NOT NULL,
    base_wage       numeric(12, 2) NOT NULL,
    total_deduction numeric(12, 2) NOT NULL DEFAULT 0,
    net_wage        numeric(12, 2) NOT NULL,
    is_paid         boolean NOT NULL DEFAULT false,
    created_at      timestamptz NOT NULL DEFAULT now(),
    CHECK (period_end >= period_start),
    -- Formula fixed by docs/05. TODO(decision-10): negative net wage, refunds and
    -- reversals are still open and may change this rule.
    CHECK (net_wage = base_wage - total_deduction)
);

-- DEDUCTION_TRANSACTION (R21 INCURS, R22 TRIGGERS, R23 APPLIES).
-- D3: payroll_id references the payroll primary key; NULL until the deduction is
--     applied to a payroll (5S creates deductions before the payroll row).
CREATE TABLE deduction_transaction (
    deduction_id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    worker_id      uuid NOT NULL REFERENCES worker (user_id),
    attendance_id  uuid NOT NULL REFERENCES attendance (attendance_id),
    payroll_id     uuid REFERENCES payroll (payroll_id),
    penalty_amount numeric(12, 2) NOT NULL CHECK (penalty_amount >= 0),
    reason         text NOT NULL,
    created_at     timestamptz NOT NULL DEFAULT now()
);

-- EQUIPMENT: master data; 1S finds or creates equipment by name.
CREATE TABLE equipment (
    equipment_id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    equipment_name text NOT NULL UNIQUE,
    is_active      boolean NOT NULL DEFAULT true
);

-- EQUIPMENT_REQUISITION (R03 REQUESTS via user_id, R04 REVIEWS via reviewed_by, R09 REQUIRES).
CREATE TABLE equipment_requisition (
    requisition_id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    requisition_no   text NOT NULL UNIQUE,
    assignment_id    uuid NOT NULL REFERENCES tor_location_assignment (assignment_id),
    user_id          uuid NOT NULL REFERENCES users (user_id),
    requisition_type text NOT NULL CHECK (requisition_type IN ('tor_base', 'additional')),
    status           text NOT NULL CHECK (status IN (
                         'pending_survey', 'pending_supervisor', 'pending_procurement',
                         'approved', 'completed', 'rejected')),
    reason           text,
    reject_reason    text,
    created_at       timestamptz NOT NULL DEFAULT now(),
    reviewed_by      uuid REFERENCES users (user_id),
    reviewed_at      timestamptz
);

-- REQUISITION_ITEM (R12 CONTAINS, R13 IS_LISTED_IN). actual_qty is cumulative (2A, 7A).
CREATE TABLE requisition_item (
    item_id        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    requisition_id uuid NOT NULL REFERENCES equipment_requisition (requisition_id),
    equipment_id   uuid NOT NULL REFERENCES equipment (equipment_id),
    required_qty   integer NOT NULL CHECK (required_qty >= 0),
    existing_qty   integer CHECK (existing_qty >= 0),
    to_buy_qty     integer CHECK (to_buy_qty >= 0),
    actual_qty     integer NOT NULL DEFAULT 0 CHECK (actual_qty >= 0),
    actual_price   numeric(12, 2) CHECK (actual_price >= 0),
    remark         text
);

-- EXPENSE_CLAIM (R05 CLAIMS via user_id, R14 GENERATES via requisition_id; many per requisition).
-- Only 'actual_expense' counts as material cost. transfer_ref_no is unique when present (2S).
CREATE TABLE expense_claim (
    expense_id        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    requisition_id    uuid NOT NULL REFERENCES equipment_requisition (requisition_id),
    user_id           uuid NOT NULL REFERENCES users (user_id),
    expense_no        text UNIQUE,
    expense_type      text NOT NULL CHECK (expense_type IN ('fund_transfer', 'actual_expense')),
    total_amount      numeric(14, 2) NOT NULL CHECK (total_amount >= 0),
    receipt_photo_url text,
    transfer_ref_no   text UNIQUE,
    created_at        timestamptz NOT NULL DEFAULT now(),
    -- 2S: a fund transfer requires its transfer reference.
    CHECK (expense_type <> 'fund_transfer' OR transfer_ref_no IS NOT NULL)
);

-- COMPANY_INVOICE (R06 BILLS). exat_deduction_amount keeps the ER spelling.
-- billing_month stores the first day of the month. status value set is not specified.
CREATE TABLE company_invoice (
    invoice_id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tor_id                uuid NOT NULL REFERENCES contract_tor (tor_id),
    invoice_no            text NOT NULL UNIQUE,
    billing_month         date NOT NULL CHECK (EXTRACT(DAY FROM billing_month) = 1),
    expected_amount       numeric(14, 2),
    net_received          numeric(14, 2),
    exat_deduction_amount numeric(14, 2),
    deduction_reason      text,
    status                text
);

-- Indexes on foreign keys used by SD queries (docs/03). Lookups already covered by a
-- primary key or UNIQUE constraint are not repeated:
--   work_schedule (worker_id, work_date): 1W, 2W, 3W, 3A, 5W
--   attendance (schedule_id, worker_id): 3W, 4S findCheckIn
CREATE INDEX idx_tor_location_assignment_tor_id ON tor_location_assignment (tor_id);       -- 3A findTORAssignments(torId); 6S per-project totals
CREATE INDEX idx_work_schedule_assignment_id ON work_schedule (assignment_id);             -- 4A findPendingLeaveRequests/checkAffectedSchedule(assignmentId); 8A checkScheduleRecipientAndArea
CREATE INDEX idx_work_schedule_work_date ON work_schedule (work_date);                     -- 4S findScheduledWorkers(work_date); 4A findReplacementCandidates(leave_date)
CREATE INDEX idx_attendance_worker_id_work_date ON attendance (worker_id, work_date);      -- 4W findTodayOpenAttendance(userId)
CREATE INDEX idx_leave_request_user_id_leave_date ON leave_request (user_id, leave_date);  -- 2W checkDuplicateLeave; 4S countApprovedAdvanceLeave
CREATE INDEX idx_payroll_worker_id_period_start ON payroll (worker_id, period_start);      -- 7W findPayslipsByMonth(userId, period_month)
CREATE INDEX idx_deduction_transaction_worker_id ON deduction_transaction (worker_id);     -- 5S findDeductionItems/getTotalDeduction; 7W findDeductionDetails
CREATE INDEX idx_equipment_requisition_assignment_id ON equipment_requisition (assignment_id); -- 5A/7A/8A find*Requests(assignmentId); 6S material via assignment
CREATE INDEX idx_equipment_requisition_status ON equipment_requisition (status);           -- 1A/2S/2A/3S list requisitions by status
CREATE INDEX idx_requisition_item_requisition_id ON requisition_item (requisition_id);     -- 1A, 2S, 2A, 7A (incl. lockRequestItems), 6W item lookups
CREATE INDEX idx_expense_claim_requisition_id ON expense_claim (requisition_id);           -- 6S material cost join
CREATE INDEX idx_company_invoice_tor_id ON company_invoice (tor_id);                       -- 6S getTotalRevenue(torId)

-- Row-level security ON with NO policies: the Supabase Data API (anon/publishable and
-- authenticated keys) can neither read nor write these tables. The backend connects as
-- the table owner `postgres`, which has BYPASSRLS (Supabase RLS docs, checked 2026-10-05).
-- Do not add policies without a decision; do not use FORCE ROW LEVEL SECURITY.
-- Supabase Storage is NOT covered: bucket access is set by policies on storage.objects.
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE worker ENABLE ROW LEVEL SECURITY;
ALTER TABLE contract_tor ENABLE ROW LEVEL SECURITY;
ALTER TABLE location ENABLE ROW LEVEL SECURITY;
ALTER TABLE tor_location_assignment ENABLE ROW LEVEL SECURITY;
ALTER TABLE work_schedule ENABLE ROW LEVEL SECURITY;
ALTER TABLE attendance ENABLE ROW LEVEL SECURITY;
ALTER TABLE leave_request ENABLE ROW LEVEL SECURITY;
ALTER TABLE work_evidence ENABLE ROW LEVEL SECURITY;
ALTER TABLE payroll ENABLE ROW LEVEL SECURITY;
ALTER TABLE deduction_transaction ENABLE ROW LEVEL SECURITY;
ALTER TABLE equipment ENABLE ROW LEVEL SECURITY;
ALTER TABLE equipment_requisition ENABLE ROW LEVEL SECURITY;
ALTER TABLE requisition_item ENABLE ROW LEVEL SECURITY;
ALTER TABLE expense_claim ENABLE ROW LEVEL SECURITY;
ALTER TABLE company_invoice ENABLE ROW LEVEL SECURITY;
