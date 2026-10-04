# Workflow operation inventory (per sequence diagram)

Generated from the class/SD mapping of the historical class package. For each sequence diagram (SD) it lists the participants and the operations they receive, in the original naming.

> **Read with care.** Names here are from the *historical* snapshot. Apply the overrides in `docs/02_domain_model.md` §5 and `docs/04_decisions_and_corrections.md` §1 (payroll names, 5W status, 5S/6S/7W/2A corrections). `…` means omitted literal arguments, not a parameter. Operations list *what is received*; they are not SQL. Original call text is in `docs/references/Class_Reference_Mapping_HISTORICAL.md`.

Role legend: Boundary = UI, Control = controller, Repository = data access, Service = infrastructure/job.


## SD-1S — Create TOR contract and work scope

**Web UI** (Boundary)
- `clickCreateNewTOR()`
- `redirectToContractFormPage()`
- `redirectToScopeFormPage()`
- `submitScopeData( scopeData)`
- `redirectToContractReviewPage()`
- `clickConfirmContract()`
- `redirectToContractDetailPage(torId)`

**DB Connector** (Service)
- `executeQuery(sql)`

**LINE Notification Service** (Service)
- `sendNewContractNotification(torId, projectName)`

**ContractForm Controller** (Control)
- `openContractForm()`
- `submitContractInfo(contractData)`
- `validateContractFormat(contractData)`
- `submitScopeData(scopeData)`
- `validateScopeFormat(scopeData)`

**ConfirmContract Controller** (Control)
- `confirmContract(contractData, scopeData)`

**TOR Repository** (Repository)
- `checkDuplicateContractNo(contractNo)`
- `validateLocation(locationId)`
- `checkDuplicateLocationName( locationName)`
- `createContract( contractData)`
- `createLocation( locationName, address)`
- `createTORLocationAssignment( torId, locationId, requiredWorkers)`
- `createInitialRequisition( assignmentId, userId, requisitionNo)`
- `findEquipmentByName( equipmentName)`
- `createEquipment( equipmentName)`
- `createRequisitionItem( requisitionId, equipmentId, remark, requiredQty)`


## SD-2S — Transfer funds for procurement

**Web UI** (Boundary)
- `selectRequisition(requisitionId)`
- `showSuccessMessage(…)`
- `redirectToDashboard()`
- `clickPendingFundingMenu()`
- `redirectToTransferEvidencePage()`
- `submitTransferEvidence(total_amount, transfer_ref_no, receipt_photo_url)`
- `showConfirmationPopup()`
- `clickConfirmTransfer()`

**DB Connector** (Service)
- `executeQuery(sql)`

**LINE Notification Service** (Service)
- `sendProcurementFundingNotification(ผู้ดูแลงาน, projectName)`

**TransferForm Controller** (Control)
- `getPendingFundingRequests()`
- `loadTransferForm(requisitionId)`
- `validateTransferEvidence(transferData)`
- `validateTransferFormat(transferData)`

**ConfirmTransfer Controller** (Control)
- `confirmFundTransfer(requisitionId, transferData)`

**Funding Repository** (Repository)
- `findPendingFundingRequests()`
- `findMissingEquipment(requisitionId)`
- `checkDuplicateTransferRef(transfer_ref_no)`
- `createFundTransferExpense(requisitionId, userId, transferData)`
- `updateRequisitionStatus(…)`


## SD-3S — Review/approve equipment requisition

**Web UI** (Boundary)
- `selectRequisition(requisitionId)`
- `showSuccessMessage(…)`
- `clickRequisitionMenu()`
- `redirectToRequisitionReviewPage()`
- `selectDecision(approve | reject)`
- `showRequiredRejectionReasonField()`
- `enterRejectionReason(rejectionReason)`
- `submitReviewDecision(requisitionId, decision, rejectionReason)`
- `redirectToRequisitionListPage()`

**Requisition Repository** (Repository)
- `findPendingRequisitions()`
- `findRequisitionDetails(requisitionId)`
- `checkPendingRequisition(requisitionId)`
- `approveRequisition(requisitionId)`
- `rejectRequisition(requisitionId)`

**DB Connector** (Service)
- `executeQuery(sql)`

**LINE Notification Service** (Service)
- `sendApprovalNotification(ผู้สร้างคำขอ)`
- `sendRejectionNotification(ผู้สร้างคำขอ, rejectionReason)`

**ReviewRequisition Controller** (Control)
- `getPendingRequisitions()`
- `loadRequisitionReview(requisitionId)`
- `processReviewDecision(requisitionId, decision, rejectionReason)`
- `validateReviewFormat(decision, rejectionReason)`


## SD-3.1S — Review requisition (detail/variant of 3S)

**Web UI** (Boundary)
- `selectRequisition(requisitionId)`
- `showSuccessMessage(…)`
- `clickRequisitionMenu()`
- `redirectToRequisitionReviewPage()`
- `selectDecision(approve | reject)`
- `showRequiredRejectionReasonField()`
- `enterRejectionReason(rejectionReason)`
- `submitReviewDecision(requisitionId, decision, rejectionReason)`
- `redirectToRequisitionListPage()`

**Requisition Repository** (Repository)
- `findPendingRequisitions()`
- `findRequisitionDetails(requisitionId)`
- `checkPendingRequisition(requisitionId)`
- `approveRequisition(requisitionId)`
- `rejectRequisition(requisitionId)`

**DB Connector** (Service)
- `executeQuery(sql)`

**LINE Notification Service** (Service)
- `sendApprovalNotification(ผู้สร้างคำขอ)`
- `sendRejectionNotification(ผู้สร้างคำขอ, rejectionReason)`

**ReviewRequisition Controller** (Control)
- `getPendingRequisitions()`
- `loadRequisitionReview(requisitionId)`
- `processReviewDecision(requisitionId, decision, rejectionReason)`
- `validateReviewFormat(decision, rejectionReason)`

**Background Job** (Service)
- `enqueueApprovalNotification(ผู้สร้างคำขอ)`
- `enqueueRejectionNotification(ผู้สร้างคำขอ, rejectionReason)`


## SD-4S — Daily attendance processing (job) and trigger of payroll summary

**Web UI** (Boundary)
- `displayAttendanceSummary(processingSummary)`
- `reviewAttendanceSummary()`
- `confirmProcessingResults()`

**DB Connector** (Service)
- `executeQuery(sql)`

**Attendance Controller** (Control)
- `processDailyAttendance(work_date)`
- `compareCheckInWithShiftStart(checkIn, shiftStartTime)`

**Attendance Repository** (Repository)
- `findScheduledWorkers(work_date)`
- `findCheckIn(scheduleId, workerId)`
- `countApprovedAdvanceLeave(userId, work_date)`
- `updateAttendanceStatus(…)`

**PayrollSummary Controller (5S)** (Control)
- `startPeriodSummary5S()`


## SD-5S — Payroll and deduction processing

**Web UI** (Boundary)
- `showSuccessMessage(…)`
- `clickCreatePayrollPeriod()`
- `showPayrollPeriodForm()`
- `submitPayrollPeriod(period_start, period_end)`
- `redirectToPayrollSummaryPage()`
- `reviewPayrollSummary()`

**DB Connector** (Service)
- `executeQuery(sql)`

**Payroll Controller** (Control)
- `processPayroll(period_start, period_end)`
- `calculateBaseWage(workDays, dailyWage)`
- `calculateNetPay(baseWage, totalDeduction)`

**Payroll Repository** (Repository)
- `findWorkersInPeriod(period_start, period_end)`
- `getBaseWageData(workerId, period_start, period_end)`
- `findDeductionItems(workerId, period_start, period_end)`
- `createDeduction(workerId, penaltyAmount, penaltyReason)`
- `getTotalDeduction(workerId, period_start, period_end)`
- `createPayroll(workerId, period_start, period_end, baseWage, totalDeduction, netPay)`


## SD-6S — Financial report (revenue, labor, material, net profit)

**Web UI** (Boundary)
- `openFinancialReportMenu()`
- `showProjectAndDateFilters()`
- `submitProfitFilters(torId, dateStart, dateEnd)`
- `redirectToNetProfitReportPage()`
- `renderFinancialDashboard(totalRevenue, totalLaborCost, totalMaterialCost, netProfit)`
- `reviewFinancialReport()`
- `confirmClosingAndContinue()`
- `requestReportPDF()`

**DB Connector** (Service)
- `executeQuery(sql)`

**FinancialReport Controller** (Control)
- `calculateFinancialReport(torId, dateStart, dateEnd)`
- `calculateNetProfit(totalRevenue, totalLaborCost, totalMaterialCost)`
- `exportFinancialReportPDF(reportData)`
- `generateReportPDF(reportData)`

**Financial Repository** (Repository)
- `getTotalRevenue(torId )`
- `getTotalLaborCost(dateStart, dateEnd)`
- `getTotalMaterialCost(torId, dateStart, dateEnd)`


## SD-1A — Site survey of equipment

**Web UI** (Boundary)
- `clickTORRequirementsMenu()`
- `selectRequisition(requisitionId)`
- `redirectToEquipmentPreparationPage()`
- `enterExistingQuantities(existing_qty for each item)`
- `clickConfirmSiteSurvey()`
- `redirectToMainDashboard()`
- `showSuccessMessage(…)`

**EquipmentSurvey Controller** (Control)
- `getPendingSurveyRequisitions()`
- `loadEquipmentRequirements(requisitionId)`
- `processSiteSurvey(requisitionId, itemQuantities)`
- `validateExistingQuantities(itemQuantities)`
- `calculateToBuyQty(required_qty, existing_qty)`
- `selectRequisitionStatus(…)`

**Requisition Repository** (Repository)
- `findPendingSurveyRequisitions()`
- `findEquipmentRequirements(requisitionId)`
- `checkRequisitionItemExists(itemId)`
- `updateSurveyItem(itemId, existing_qty)`
- `updateRequisitionStatus(requisitionId, status)`

**DB Connector** (Service)
- `executeQuery(sql)`


## SD-2A — Equipment procurement

**Web UI** (Boundary)
- `selectRequisition(requisitionId)`
- `showSuccessMessage(…)`
- `clickEquipmentProcurementMenu()`
- `redirectToProcurementFormPage()`
- `enterProcurementData(itemQuantities, unitPrices)`
- `attachReceiptPhoto()`
- `clickConfirmProcurement()`
- `redirectToDashboard()`

**Requisition Repository** (Repository)
- `checkRequisitionItemExists(itemId)`
- `findProcurementRequisitions()`
- `findRemainingEquipment(requisitionId)`
- `createActualExpense(expense_no, requisitionId, userId, total_amount, receipt_photo_url)`
- `updateProcuredItem(itemId, quantity, unitPrice)`
- `updateRequisitionStatus(requisitionId)`

**DB Connector** (Service)
- `executeQuery(sql)`

**LINE Notification Service** (Service)
- `sendShortageNotification(ผู้ควบคุมงาน)`

**Procurement Controller** (Control)
- `getProcurementRequisitions()`
- `loadProcurementForm(requisitionId)`
- `processProcurement(requisitionId, procurementData)`
- `validateProcurementData(procurementData)`
- `determineProcurementStatus(itemQuantities)`


## SD-3A — Create work schedule

**Web UI** (Boundary)
- `showSuccessMessage(…)`
- `clickWorkScheduleMenu()`
- `redirectToCreateSchedulePage()`
- `enterScheduleData(assignmentId, work_date, shift_start_time, workerIds)`
- `clickReviewSchedule()`
- `showScheduleReviewPage()`
- `clickConfirmSchedule()`
- `redirectToScheduleListPage()`

**DB Connector** (Service)
- `executeQuery(sql)`

**WorkSchedule Controller** (Control)
- `loadWorkScheduleForm(torId)`
- `reviewSchedule(scheduleData)`
- `validateScheduleFormat(scheduleData)`
- `validateAreaExistenceAndPermission(assignmentId)`
- `confirmSchedule(scheduleData)`
- `recheckLatestScheduleData(scheduleData)`
- `notifyAssignedWorkers(workerIds)`

**WorkSchedule Repository** (Repository)
- `findTORAssignments(torId)`
- `checkWorkerAvailable(workerId)`
- `checkDuplicateSchedule(workerId, work_date)`
- `createWorkSchedule(assignmentId, workerId, work_date, shift_start_time)`


## SD-4A — Leave review and replacement

**Web UI** (Boundary)
- `clickLeaveAndReplacementMenu()`
- `selectLeaveRequest(requestId)`
- `redirectToLeaveReviewPage()`
- `chooseDecision(approve | reject)`
- `selectReplacement(replacementWorkerId)`
- `clickReviewDecision()`
- `showDecisionReviewPage()`
- `showWarning(…)`
- `clickConfirmDecision()`
- `redirectToLeaveRequestListPage()`
- `showResultMessage(decisionResult)`

**DB Connector** (Service)
- `executeQuery(sql)`

**LeaveReview Controller** (Control)
- `loadPendingLeaveRequests(assignmentId)`
- `checkAreaPermission(assignmentId)`
- `loadLeaveReviewForm(requestId, assignmentId)`
- `loadReplacementCandidates(leavingWorkerId, leave_date)`
- `reviewLeaveDecision(reviewData)`
- `validateLeaveDecision(reviewData)`
- `recheckSelectedReplacement(reviewData)`
- `confirmLeaveDecision(reviewData)`
- `recheckLatestData(reviewData)`
- `notifySupervisorReplacementMissing()`
- `notifyLeavingWorker(decision)`
- `notifyReplacementWorker(replacementWorkerId)`

**LeaveReview Repository** (Repository)
- `findPendingLeaveRequests(assignmentId)`
- `findLeaveDetails(requestId, assignmentId)`
- `findReplacementCandidates(leavingWorkerId, leave_date)`
- `checkPendingLeave(requestId)`
- `checkAffectedSchedule(requestId, scheduleId, assignmentId)`
- `updateLeaveStatus(requestId, approved | rejected)`
- `createReplacementSchedule(replacementWorkerId, scheduleId)`


## SD-5A — Review additional equipment request

**Web UI** (Boundary)
- `clickAdditionalEquipmentReviewMenu()`
- `selectExistingRequest(requisitionId)`
- `showReadOnlyRequestAndReferences()`
- `clickConsiderActionMethod(requisitionId)`
- `redirectTo6A(requisitionId)`

**DB Connector** (Service)
- `executeQuery(sql)`

**EquipmentReview Controller** (Control)
- `loadAdditionalRequests(assignmentId)`
- `checkAreaPermission(assignmentId)`
- `loadEquipmentReview(requisitionId, assignmentId)`
- `validateRequestId(requisitionId)`
- `reloadOriginalRequestFor6A(requisitionId)`
- `reloadOriginalRequestFromDatabase(requisitionId)`

**EquipmentReview Repository** (Repository)
- `findAdditionalRequests(assignmentId)`
- `findOriginalRequest(requisitionId, assignmentId)`
- `findTORReference(requisitionId)`
- `findOtherActiveRequests(requisitionId)`


## SD-6A — Equipment decision (forward/reject)

**Web UI** (Boundary)
- `clickReviewDecision()`
- `showDecisionReviewPage()`
- `clickConfirmDecision()`
- `showResultMessage(decisionResult)`
- `showReadOnlyOriginalRequest()`
- `chooseDecision(purchase | no_purchase)`
- `showReasonOrGuidanceFields(decision)`
- `enterPurchaseReasonOrGuidance(text)`
- `redirectAfterDecision()`

**DB Connector** (Service)
- `executeQuery(sql)`

**LINE Notification Service** (Service)
- `sendGuidanceToOriginalRequester(line_id, originalItems, guidanceText)`

**EquipmentDecision Controller** (Control)
- `loadOriginalRequestFrom5A(requisitionId, assignmentId)`
- `loadOriginalProjectAndEquipmentData(requisitionId)`
- `reviewEquipmentDecision(decisionData)`
- `validateDecisionData(decisionData)`
- `checkOtherRequestsCoveringItems(requisitionId)`
- `confirmEquipmentDecision(decisionData)`
- `recheckLatestRequestData(requisitionId)`
- `notifySupervisor(…)`

**EquipmentDecision Repository** (Repository)
- `findOriginalRequestAndRecipient(requisitionId, assignmentId)`
- `checkOriginalRequestAndArea(requisitionId, assignmentId)`
- `checkOriginalItems(requisitionId)`
- `forwardOriginalRequest(requisitionId, assignmentId, purchaseReason)`
- `rejectOriginalRequest(requisitionId, assignmentId, assistantUserId)`


## SD-7A — Additional equipment purchase

**Web UI** (Boundary)
- `clickApprovedAdditionalPurchasesMenu()`
- `selectOriginalRequest(requisitionId)`
- `redirectToAdditionalPurchaseForm()`
- `enterCurrentPurchaseData(itemQuantities, unitPrices, receiptIfCostPositive)`
- `clickReviewPurchase()`
- `showPurchaseReviewPage()`
- `clickConfirmPurchase()`
- `redirectAfterPurchase(status)`
- `showPurchaseResultMessage(status)`

**DB Connector** (Service)
- `executeQuery(sql)`

**AdditionalPurchase Controller** (Control)
- `loadEligibleAdditionalRequests(assignmentId)`
- `checkAreaPermission(assignmentId)`
- `loadAdditionalPurchaseForm(requisitionId)`
- `reviewAdditionalPurchase(requisitionId, purchaseData)`
- `validatePurchaseData(purchaseData)`
- `checkRequestEligibilityAndAreaPermission(requisitionId)`
- `calculateReviewTotals(purchaseData)`
- `confirmAdditionalPurchase(requisitionId, purchaseData)`
- `recheckLatestDataAndPreventDuplicateSubmission(purchaseData)`
- `commitTransaction()`
- `notifySupervisorOfShortage(requisition_no, nextStep = 2S)`

**AdditionalPurchase Repository** (Repository)
- `findEligibleAdditionalRequests(assignmentId)`
- `findRemainingItems(requisitionId)`
- `checkItemAndLatestRemaining(itemId, requisitionId)`
- `lockOriginalRequest(requisitionId)`
- `lockRequestItems(requisitionId)`
- `updateCumulativePurchase(itemId, requisitionId, quantity, unitPrice)`
- `createActualExpense(requisitionId, assistantUserId, currentTotal, receipt_photo_url)`
- `updateRequestStatusFromSavedQuantities(requisitionId)`


## SD-8A — Equipment delivery evidence

**Web UI** (Boundary)
- `showSuccessMessage(…)`
- `clickEquipmentDeliveryMenu()`
- `selectDeliveryRequest(requisitionId)`
- `redirectToDeliveryForm()`
- `enterDeliveryData(scheduleId, deliveredQuantities, description, photos)`
- `clickReviewDelivery()`
- `showDeliveryReviewPage()`
- `clickConfirmDelivery()`
- `redirectToDeliveryListPage()`

**DB Connector** (Service)
- `executeQuery(sql)`

**LINE Notification Service** (Service)
- `sendActualDeliverySummary(line_id, requisition_no, itemsAndQuantities)`

**EquipmentDelivery Controller** (Control)
- `loadCompletedAdditionalRequests(assignmentId)`
- `checkAreaPermission(assignmentId)`
- `loadDeliveryForm(requisitionId)`
- `reviewDelivery(requisitionId, deliveryData)`
- `validateDeliveryFormat(deliveryData)`
- `checkOriginalRecipientAreaAndFullQuantities(deliveryData)`
- `confirmDelivery(requisitionId, deliveryData)`
- `recheckLatestDeliveryAndPreventDuplicateEvidence(deliveryData)`

**EquipmentDelivery Repository** (Repository)
- `findCompletedAdditionalRequests(assignmentId)`
- `findItemsAndOriginalRecipient(requisitionId)`
- `checkScheduleRecipientAndArea(scheduleId, requisitionId)`
- `saveDeliveryEvidence(scheduleId, requisitionId, description, photo_url)`
- `findDeliveryNotificationData(requisitionId)`


## SD-9A — Work continuation (contract check)

**Web UI** (Boundary)
- `showWarning(…)`
- `showReceivedSummary(project, area, summaryMessage)`
- `openAndReadSummary()`
- `selectArea(assignmentId)`
- `showContractAndAreaDetails()`
- `chooseContinueWork()`
- `openAreaDashboard(assignmentId)`

**DB Connector** (Service)
- `executeQuery(sql)`

**WorkContinuation Controller** (Control)
- `loadContractAndArea(assignmentId)`
- `validateAssignmentId(assignmentId)`
- `checkAreaAccessPermission(assignmentId)`
- `checkContractForContinuation(contractData)`
- `evaluateContractStatusAndEndDate(contractData)`

**Contract Repository** (Repository)
- `findContractAndArea(assignmentId)`


## SD-1W — View my schedule

**DB Connector** (Service)
- `executeQuery(sql)`

**WorkSchedule Controller** (Control)
- `loadMySchedule(LINEUserId)`
- `authenticateAndResolveLinkedUser(tokenOrLINEUserId)`

**WorkSchedule Repository** (Repository)
- `findUpcomingSchedules(userId)`

**LINE Mini App (LIFF)** (Boundary)
- `tapMyScheduleRichMenu()`
- `openLIFFWindow()`
- `showScrollableScheduleCardsOrList()`
- `optionallyTapScheduleCard(scheduleId, nextAction)`
- `openLeaveScreen2W(scheduleId)`
- `openPreWorkConfirmationScreen3W(scheduleId)`


## SD-2W — Submit leave request

**DB Connector** (Service)
- `executeQuery(sql)`

**LINE Notification Service** (Service)
- `notifyAssistantOfLeave(requestId)`

**Background Job** (Service)
- `sendLeaveNotification(requestId)`

**LINE Mini App (LIFF)** (Boundary)
- `tapLeaveButton()`
- `showLeaveDateAndReasonForm()`
- `confirmLeaveRequest(leave_date, reason)`
- `redirectToLIFFHome()`
- `showToast(…)`

**LeaveRequest Controller** (Control)
- `loadLeaveDateOptions(lineIdentity)`
- `verifyLinkedLINEIdentity(lineIdentity)`
- `submitLeaveRequest(userId, leave_date, reason)`
- `validateLeaveDateAndReason(leave_date, reason)`
- `calculateAdvanceNotice(currentTimestamp, leave_date)`

**LeaveRequest Repository** (Repository)
- `findScheduledDates(userId)`
- `checkDuplicateLeave(userId, leave_date)`
- `createPendingLeave(userId, leave_date, reason, is_advance_notice)`


## SD-3W — Substitute check-in

**DB Connector** (Service)
- `executeQuery(sql)`

**Attendance Repository** (Repository)
- `findTodayAssignedShift(userId)`
- `checkDuplicateCheckIn(userId, scheduleId)`
- `upsertSubstituteAttendance(scheduleId, workerId)`

**LINE Mini App (LIFF)** (Boundary)
- `redirectToLIFFHome()`
- `showToast(…)`
- `tapConfirmSubstituteWork()`
- `showAssignedLocationAndStartTime()`
- `tapConfirmAndCheckIn()`

**SubstituteCheckIn Controller** (Control)
- `loadTodayAssignedShift(lineIdentity)`
- `verifyLinkedLINEIdentity(lineIdentity)`
- `confirmSubstituteCheckIn(lineUserId, scheduleId)`
- `validateLINEUserIdAndScheduleUUID(lineUserId, scheduleId)`
- `checkTodayShiftMatchesWorker(todayAssignedShift)`


## SD-4W — Work report and check-out

**DB Connector** (Service)
- `executeQuery(sql)`

**LINE Mini App (LIFF)** (Boundary)
- `redirectToLIFFHome()`
- `showToast(…)`
- `tapWorkReportAndCheckOut()`
- `showDescriptionAndPhotoForm()`
- `confirmWorkReport(description, photo)`

**WorkReport Controller** (Control)
- `loadOpenCheckIn(lineIdentity)`
- `verifyLinkedLINEIdentity(lineIdentity)`
- `submitWorkReport(description, photo_url)`
- `validateDescriptionAndImage(description, photo_url)`
- `checkExactlyOneOpenAttendance(openAttendance)`

**WorkReport Repository** (Repository)
- `findTodayOpenAttendance(userId)`
- `saveWorkEvidence(workerId, scheduleId, description, photo_url)`
- `updateCheckOut(attendanceId)`


## SD-5W — Report missing equipment (additional request)

**DB Connector** (Service)
- `executeQuery(sql)`

**LINE Notification Service** (Service)
- `notifyAssistantOfAdditionalRequest(location_name)`

**LINE Mini App (LIFF)** (Boundary)
- `redirectToLIFFHome()`
- `showToast(…)`
- `tapInsufficientEquipment()`
- `showEquipmentQuantityAndReasonForm()`
- `sendEquipmentRequest(equipmentId, quantity, reason)`

**EquipmentRequest Controller** (Control)
- `loadEquipmentRequestForm(lineIdentity)`
- `verifyLinkedLINEIdentity(lineIdentity)`
- `submitAdditionalRequest(equipmentId, quantity, reason)`
- `validateEquipmentUUIDQuantityAndReason(equipmentId, quantity, reason)`
- `checkCurrentShiftHasData(currentAreaAndRequester)`

**EquipmentRequest Repository** (Repository)
- `findCurrentWorkingArea(userId)`
- `findActiveEquipment()`
- `checkActiveEquipmentExists(equipmentId)`
- `createAdditionalRequest(autoRequisitionNo, userId, assignmentId, reason)`
- `createRequestItem(requisitionId, equipmentId, quantity)`


## SD-6W — View equipment request result

**DB Connector** (Service)
- `executeQuery(sql)`

**LINE Notification Service** (Service)
- `sendEquipmentReceivedPush(line_id, equipment_name, actual_qty)`
- `sendNoPurchaseGuidancePush(line_id, equipment_name)`

**EquipmentResult Repository** (Repository)
- `findEquipmentRequestResult(requisitionId)`


## SD-7W — View payslip

**DB Connector** (Service)
- `executeQuery(sql)`

**LINE Mini App (LIFF)** (Boundary)
- `tapMyPayslip()`
- `showPeriodMonthSearchForm()`
- `searchPayslip(period_month)`
- `showDigitalPayslipCard(period_start, period_end, total_wage, total_deduction, deductionDetails, net_wage)`
- `closePayslipWindow()`
- `returnToLIFFHome()`

**Payslip Controller** (Control)
- `openPayslipSearch(lineIdentity)`
- `verifyLinkedLINEIdentity(lineIdentity)`
- `getPayslip(userId, period_month)`
- `validatePeriodMonth(period_month)`

**Payslip Repository** (Repository)
- `findPayslipsByMonth(userId, period_month)`
- `findDeductionDetails(userId, period_start, period_end)`
