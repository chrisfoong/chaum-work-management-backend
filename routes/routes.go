package routes

import (
	"chrisfoong/chaum-work-management-backend/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter builds the Gin engine and registers every route group.
//
// This is a route-only skeleton built straight off the ER diagram
// (SA1-ER_FINAL): every entity gets a resource group, and every
// relationship diamond that reads naturally as "parent has many
// children" gets a nested GET route. See README.md for the full
// entity/relationship writeup and the reasoning behind each nested
// route below.
//
// No class modeling yet: handlers read/write plain JSON maps, and
// nothing is wired to MongoDB yet (see config/database.go).
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/health", handlers.HealthCheck)

	api := r.Group("/api")
	{
		registerUserRoutes(api)
		registerWorkerRoutes(api)
		registerContractRoutes(api)
		registerLocationRoutes(api)
		registerAssignmentRoutes(api)
		registerScheduleRoutes(api)
		registerAttendanceRoutes(api)
		registerEvidenceRoutes(api)
		registerPayrollRoutes(api)
		registerDeductionRoutes(api)
		registerLeaveRequestRoutes(api)
		registerExpenseClaimRoutes(api)
		registerInvoiceRoutes(api)
		registerEquipmentRoutes(api)
		registerEquipmentRequisitionRoutes(api)
	}

	return r
}

// USER
func registerUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.GET("", handlers.ListUsers)
		users.POST("", handlers.CreateUser)
		users.GET("/:id", handlers.GetUser)
		users.PUT("/:id", handlers.UpdateUser)
		users.DELETE("/:id", handlers.DeleteUser)

		users.GET("/:id/payroll", handlers.GetUserPayroll)                              // RECEIVES
		users.GET("/:id/expense-claims", handlers.GetUserExpenseClaims)                 // CLAIMS
		users.GET("/:id/equipment-requisitions", handlers.GetUserEquipmentRequisitions) // REQUEST
	}
}

// WORKER (specialization of USER)
func registerWorkerRoutes(rg *gin.RouterGroup) {
	workers := rg.Group("/workers")
	{
		workers.GET("", handlers.ListWorkers)
		workers.POST("", handlers.CreateWorker)
		workers.GET("/:id", handlers.GetWorker)
		workers.PUT("/:id", handlers.UpdateWorker)
		workers.DELETE("/:id", handlers.DeleteWorker)

		workers.GET("/:id/schedules", handlers.GetWorkerSchedules)          // IS_ASSIGNED_TO
		workers.GET("/:id/attendance", handlers.GetWorkerAttendance)        // LOGS
		workers.GET("/:id/evidence", handlers.GetWorkerEvidence)            // CAPTURES
		workers.GET("/:id/deductions", handlers.GetWorkerDeductions)        // INCURS
		workers.GET("/:id/leave-requests", handlers.GetWorkerLeaveRequests) // SUBMITS
	}
}

// CONTRACT_TOR
func registerContractRoutes(rg *gin.RouterGroup) {
	contracts := rg.Group("/contracts")
	{
		contracts.GET("", handlers.ListContracts)
		contracts.POST("", handlers.CreateContract)
		contracts.GET("/:id", handlers.GetContract)
		contracts.PUT("/:id", handlers.UpdateContract)
		contracts.DELETE("/:id", handlers.DeleteContract)

		contracts.GET("/:id/assignments", handlers.GetContractAssignments)      // INCLUDES
		contracts.GET("/:id/invoices", handlers.GetContractInvoices)            // BILLS
		contracts.GET("/:id/expense-claims", handlers.GetContractExpenseClaims) // HAS_EXPENSE
	}
}

// LOCATION
func registerLocationRoutes(rg *gin.RouterGroup) {
	locations := rg.Group("/locations")
	{
		locations.GET("", handlers.ListLocations)
		locations.POST("", handlers.CreateLocation)
		locations.GET("/:id", handlers.GetLocation)
		locations.PUT("/:id", handlers.UpdateLocation)
		locations.DELETE("/:id", handlers.DeleteLocation)

		locations.GET("/:id/assignments", handlers.GetLocationAssignments) // HOSTS
	}
}

// TOR_LOCATION_ASSIGNMENT
func registerAssignmentRoutes(rg *gin.RouterGroup) {
	assignments := rg.Group("/assignments")
	{
		assignments.GET("", handlers.ListAssignments)
		assignments.POST("", handlers.CreateAssignment)
		assignments.GET("/:id", handlers.GetAssignment)
		assignments.PUT("/:id", handlers.UpdateAssignment)
		assignments.DELETE("/:id", handlers.DeleteAssignment)

		assignments.GET("/:id/schedules", handlers.GetAssignmentSchedules)                          // GENERATES
		assignments.GET("/:id/evidence", handlers.GetAssignmentEvidence)                            // COLLECTS
		assignments.GET("/:id/equipment-requisitions", handlers.GetAssignmentEquipmentRequisitions) // REQUIRES
	}
}

// WORK_SCHEDULE
func registerScheduleRoutes(rg *gin.RouterGroup) {
	schedules := rg.Group("/schedules")
	{
		schedules.GET("", handlers.ListSchedules)
		schedules.POST("", handlers.CreateSchedule)
		schedules.GET("/:id", handlers.GetSchedule)
		schedules.PUT("/:id", handlers.UpdateSchedule)
		schedules.DELETE("/:id", handlers.DeleteSchedule)

		schedules.GET("/:id/attendance", handlers.GetScheduleAttendance) // RECORDS (1:1)
	}
}

// ATTENDANCE
func registerAttendanceRoutes(rg *gin.RouterGroup) {
	attendance := rg.Group("/attendance")
	{
		attendance.GET("", handlers.ListAttendanceRecords)
		attendance.POST("", handlers.CreateAttendance)
		attendance.GET("/:id", handlers.GetAttendance)
		attendance.PUT("/:id", handlers.UpdateAttendance)
		attendance.DELETE("/:id", handlers.DeleteAttendance)

		attendance.GET("/:id/deductions", handlers.GetAttendanceDeductions) // TRIGGERS
	}
}

// WORK_EVIDENCE
func registerEvidenceRoutes(rg *gin.RouterGroup) {
	evidence := rg.Group("/evidence")
	{
		evidence.GET("", handlers.ListEvidenceRecords)
		evidence.POST("", handlers.CreateEvidence)
		evidence.GET("/:id", handlers.GetEvidence)
		evidence.PUT("/:id", handlers.UpdateEvidence)
		evidence.DELETE("/:id", handlers.DeleteEvidence)
	}
}

// PAYROLL
func registerPayrollRoutes(rg *gin.RouterGroup) {
	payroll := rg.Group("/payroll")
	{
		payroll.GET("", handlers.ListPayrollRecords)
		payroll.POST("", handlers.CreatePayroll)
		payroll.GET("/:id", handlers.GetPayroll)
		payroll.PUT("/:id", handlers.UpdatePayroll)
		payroll.DELETE("/:id", handlers.DeletePayroll)

		payroll.GET("/:id/deductions", handlers.GetPayrollDeductions) // CLEARS
	}
}

// DEDUCTION_TRANSACTION
func registerDeductionRoutes(rg *gin.RouterGroup) {
	deductions := rg.Group("/deductions")
	{
		deductions.GET("", handlers.ListDeductions)
		deductions.POST("", handlers.CreateDeduction)
		deductions.GET("/:id", handlers.GetDeduction)
		deductions.PUT("/:id", handlers.UpdateDeduction)
		deductions.DELETE("/:id", handlers.DeleteDeduction)
	}
}

// LEAVE_REQUEST
func registerLeaveRequestRoutes(rg *gin.RouterGroup) {
	leaveRequests := rg.Group("/leave-requests")
	{
		leaveRequests.GET("", handlers.ListLeaveRequests)
		leaveRequests.POST("", handlers.CreateLeaveRequest)
		leaveRequests.GET("/:id", handlers.GetLeaveRequest)
		leaveRequests.PUT("/:id", handlers.UpdateLeaveRequest)
		leaveRequests.DELETE("/:id", handlers.DeleteLeaveRequest)
	}
}

// EXPENSE_CLAIM
func registerExpenseClaimRoutes(rg *gin.RouterGroup) {
	expenseClaims := rg.Group("/expense-claims")
	{
		expenseClaims.GET("", handlers.ListExpenseClaims)
		expenseClaims.POST("", handlers.CreateExpenseClaim)
		expenseClaims.GET("/:id", handlers.GetExpenseClaim)
		expenseClaims.PUT("/:id", handlers.UpdateExpenseClaim)
		expenseClaims.DELETE("/:id", handlers.DeleteExpenseClaim)
	}
}

// COMPANY_INVOICE
func registerInvoiceRoutes(rg *gin.RouterGroup) {
	invoices := rg.Group("/invoices")
	{
		invoices.GET("", handlers.ListInvoices)
		invoices.POST("", handlers.CreateInvoice)
		invoices.GET("/:id", handlers.GetInvoice)
		invoices.PUT("/:id", handlers.UpdateInvoice)
		invoices.DELETE("/:id", handlers.DeleteInvoice)
	}
}

// EQUIPMENT
func registerEquipmentRoutes(rg *gin.RouterGroup) {
	equipment := rg.Group("/equipment")
	{
		equipment.GET("", handlers.ListEquipmentList)
		equipment.POST("", handlers.CreateEquipmentItem)
		equipment.GET("/:id", handlers.GetEquipmentItem)
		equipment.PUT("/:id", handlers.UpdateEquipmentItem)
		equipment.DELETE("/:id", handlers.DeleteEquipmentItem)

		equipment.GET("/:id/requisition-items", handlers.GetEquipmentRequisitionHistory) // IS_LISTED_IN
	}
}

// EQUIPMENT_REQUISITION (+ nested REQUISITION_ITEM)
func registerEquipmentRequisitionRoutes(rg *gin.RouterGroup) {
	requisitions := rg.Group("/equipment-requisitions")
	{
		requisitions.GET("", handlers.ListEquipmentRequisitions)
		requisitions.POST("", handlers.CreateEquipmentRequisition)
		requisitions.GET("/:id", handlers.GetEquipmentRequisition)
		requisitions.PUT("/:id", handlers.UpdateEquipmentRequisition)
		requisitions.DELETE("/:id", handlers.DeleteEquipmentRequisition)

		// CONTAINS: REQUISITION_ITEM lines that belong to this requisition
		requisitions.GET("/:id/items", handlers.ListRequisitionItems)
		requisitions.POST("/:id/items", handlers.AddRequisitionItem)
		requisitions.DELETE("/:id/items/:itemId", handlers.RemoveRequisitionItem)
	}
}
