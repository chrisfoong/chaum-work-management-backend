-- Reverts 0001_init. NOT APPLIED. Destroys all data in these tables.
-- Indexes, CHECK constraints and row-level security are dropped with their tables.
DROP TABLE IF EXISTS company_invoice;
DROP TABLE IF EXISTS expense_claim;
DROP TABLE IF EXISTS requisition_item;
DROP TABLE IF EXISTS equipment_requisition;
DROP TABLE IF EXISTS equipment;
DROP TABLE IF EXISTS deduction_transaction;
DROP TABLE IF EXISTS payroll;
DROP TABLE IF EXISTS work_evidence;
DROP TABLE IF EXISTS leave_request;
DROP TABLE IF EXISTS attendance;
DROP TABLE IF EXISTS work_schedule;
DROP TABLE IF EXISTS tor_location_assignment;
DROP TABLE IF EXISTS location;
DROP TABLE IF EXISTS contract_tor;
DROP TABLE IF EXISTS worker;
DROP TABLE IF EXISTS users;
