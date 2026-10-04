> **HISTORICAL — built from the OLD ER. Class/operation provenance and original SD call text are useful; field names for payroll, leave user_id, PAYS, 7W/2A/5W/5S/6S calls are OUTDATED. docs/02 §5 and docs/04 §1.3 override this file.**

# Reference mapping — ER-first class draft

The 64 classes contain 16 ER domain classes and 48 SD software participants. ER names are preserved. SD labels retain their original class names/spaces; `dbconn:DB Connector` is displayed as its type, DB Connector. No methods are added to the ER entities.

Each operation belongs to its receiving SD lifeline. Literal argument values are represented by … in the diagrams; this does not introduce a parameter. Original calls are recorded below. Repeated SQL calls share `executeQuery(sql)`; the positional SQL notation in SD-3.1S is recorded as source evidence, not another overload.

Full software definitions show the union of supported operations. Document workflow views show the applicable subset of the same class. Context boxes are explicitly labelled and omit members. All members remain covered across the parts.

Domain diagram relationship IDs R01–R23 identify the ER relationship names. The 1/* labels mean maximum cardinality only. U1 reproduces the visible USER–WORKER connection without claiming inheritance.

## Web UI

Role: Boundary. Sources: SD-1A, SD-1S, SD-2A, SD-2S, SD-3.1S, SD-3A, SD-3S, SD-4A, SD-4S, SD-5A, SD-5S, SD-6A, SD-6S, SD-7A, SD-8A, SD-9A.


| Diagram operation | Supporting SDs |
|---|---|
| `clickTORRequirementsMenu()` | SD-1A |
| `selectRequisition(requisitionId)` | SD-1A, SD-2A, SD-2S, SD-3.1S, SD-3S |
| `redirectToEquipmentPreparationPage()` | SD-1A |
| `enterExistingQuantities(existing_qty for each item)` | SD-1A |
| `clickConfirmSiteSurvey()` | SD-1A |
| `redirectToMainDashboard()` | SD-1A |
| `showSuccessMessage(…)` | SD-1A, SD-2A, SD-2S, SD-3.1S, SD-3A, SD-3S, SD-5S, SD-8A |
| `clickEquipmentProcurementMenu()` | SD-2A |
| `redirectToProcurementFormPage()` | SD-2A |
| `enterProcurementData(itemQuantities, unitPrices)` | SD-2A |
| `attachReceiptPhoto()` | SD-2A |
| `clickConfirmProcurement()` | SD-2A |
| `redirectToDashboard()` | SD-2A, SD-2S |
| `clickWorkScheduleMenu()` | SD-3A |
| `redirectToCreateSchedulePage()` | SD-3A |
| `enterScheduleData(assignmentId, work_date, shift_start_time, workerIds)` | SD-3A |
| `clickReviewSchedule()` | SD-3A |
| `showScheduleReviewPage()` | SD-3A |
| `clickConfirmSchedule()` | SD-3A |
| `redirectToScheduleListPage()` | SD-3A |
| `clickLeaveAndReplacementMenu()` | SD-4A |
| `selectLeaveRequest(requestId)` | SD-4A |
| `redirectToLeaveReviewPage()` | SD-4A |
| `chooseDecision(approve \| reject)` | SD-4A |
| `selectReplacement(replacementWorkerId)` | SD-4A |
| `clickReviewDecision()` | SD-4A, SD-6A |
| `showDecisionReviewPage()` | SD-4A, SD-6A |
| `showWarning(…)` | SD-4A, SD-9A |
| `clickConfirmDecision()` | SD-4A, SD-6A |
| `redirectToLeaveRequestListPage()` | SD-4A |
| `showResultMessage(decisionResult)` | SD-4A, SD-6A |
| `clickAdditionalEquipmentReviewMenu()` | SD-5A |
| `selectExistingRequest(requisitionId)` | SD-5A |
| `showReadOnlyRequestAndReferences()` | SD-5A |
| `clickConsiderActionMethod(requisitionId)` | SD-5A |
| `redirectTo6A(requisitionId)` | SD-5A |
| `showReadOnlyOriginalRequest()` | SD-6A |
| `chooseDecision(purchase \| no_purchase)` | SD-6A |
| `showReasonOrGuidanceFields(decision)` | SD-6A |
| `enterPurchaseReasonOrGuidance(text)` | SD-6A |
| `redirectAfterDecision()` | SD-6A |
| `clickApprovedAdditionalPurchasesMenu()` | SD-7A |
| `selectOriginalRequest(requisitionId)` | SD-7A |
| `redirectToAdditionalPurchaseForm()` | SD-7A |
| `enterCurrentPurchaseData(itemQuantities, unitPrices, receiptIfCostPositive)` | SD-7A |
| `clickReviewPurchase()` | SD-7A |
| `showPurchaseReviewPage()` | SD-7A |
| `clickConfirmPurchase()` | SD-7A |
| `redirectAfterPurchase(status)` | SD-7A |
| `showPurchaseResultMessage(status)` | SD-7A |
| `clickEquipmentDeliveryMenu()` | SD-8A |
| `selectDeliveryRequest(requisitionId)` | SD-8A |
| `redirectToDeliveryForm()` | SD-8A |
| `enterDeliveryData(scheduleId, deliveredQuantities, description, photos)` | SD-8A |
| `clickReviewDelivery()` | SD-8A |
| `showDeliveryReviewPage()` | SD-8A |
| `clickConfirmDelivery()` | SD-8A |
| `redirectToDeliveryListPage()` | SD-8A |
| `showReceivedSummary(project, area, summaryMessage)` | SD-9A |
| `openAndReadSummary()` | SD-9A |
| `selectArea(assignmentId)` | SD-9A |
| `showContractAndAreaDetails()` | SD-9A |
| `chooseContinueWork()` | SD-9A |
| `openAreaDashboard(assignmentId)` | SD-9A |
| `clickCreateNewTOR()` | SD-1S |
| `redirectToContractFormPage()` | SD-1S |
| `redirectToScopeFormPage()` | SD-1S |
| `submitScopeData( scopeData)` | SD-1S |
| `redirectToContractReviewPage()` | SD-1S |
| `clickConfirmContract()` | SD-1S |
| `redirectToContractDetailPage(torId)` | SD-1S |
| `clickPendingFundingMenu()` | SD-2S |
| `redirectToTransferEvidencePage()` | SD-2S |
| `submitTransferEvidence(total_amount, transfer_ref_no, receipt_photo_url)` | SD-2S |
| `showConfirmationPopup()` | SD-2S |
| `clickConfirmTransfer()` | SD-2S |
| `clickRequisitionMenu()` | SD-3.1S, SD-3S |
| `redirectToRequisitionReviewPage()` | SD-3.1S, SD-3S |
| `selectDecision(approve \| reject)` | SD-3.1S, SD-3S |
| `showRequiredRejectionReasonField()` | SD-3.1S, SD-3S |
| `enterRejectionReason(rejectionReason)` | SD-3.1S, SD-3S |
| `submitReviewDecision(requisitionId, decision, rejectionReason)` | SD-3.1S, SD-3S |
| `redirectToRequisitionListPage()` | SD-3.1S, SD-3S |
| `displayAttendanceSummary(processingSummary)` | SD-4S |
| `reviewAttendanceSummary()` | SD-4S |
| `confirmProcessingResults()` | SD-4S |
| `clickCreatePayrollPeriod()` | SD-5S |
| `showPayrollPeriodForm()` | SD-5S |
| `submitPayrollPeriod(period_start, period_end)` | SD-5S |
| `redirectToPayrollSummaryPage()` | SD-5S |
| `reviewPayrollSummary()` | SD-5S |
| `openFinancialReportMenu()` | SD-6S |
| `showProjectAndDateFilters()` | SD-6S |
| `submitProfitFilters(torId, dateStart, dateEnd)` | SD-6S |
| `redirectToNetProfitReportPage()` | SD-6S |
| `renderFinancialDashboard(totalRevenue, totalLaborCost, totalMaterialCost, netProfit)` | SD-6S |
| `reviewFinancialReport()` | SD-6S |
| `confirmClosingAndContinue()` | SD-6S |
| `requestReportPDF()` | SD-6S |

Original calls and receiver evidence:

- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-50`: `clickTORRequirementsMenu()` → Web UI.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-66`: `selectRequisition(requisitionId)` → Web UI.
- SD-2A, edge `N13osf1N4320fgIz80b4-78`: `selectRequisition(requisitionId)` → Web UI.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-53`: `selectRequisition(requisitionId)` → Web UI.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-44`: `selectRequisition(requisitionId)` → Web UI.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-54`: `selectRequisition(requisitionId)` → Web UI.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-76`: `redirectToEquipmentPreparationPage()` → Web UI.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-80`: `enterExistingQuantities(existing_qty for each item)` → Web UI.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-84`: `clickConfirmSiteSurvey()` → Web UI.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-130`: `redirectToMainDashboard()` → Web UI.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-132`: `showSuccessMessage("ยืนยันการรับทราบรายการอุปกรณ์สำเร็จ")` → Web UI.
- SD-2A, edge `N13osf1N4320fgIz80b4-154`: `showSuccessMessage("จัดหาอุปกรณ์ครบถ้วน พร้อมดำเนินการจัดตารางงาน")` → Web UI.
- SD-2A, edge `N13osf1N4320fgIz80b4-156`: `showSuccessMessage("ส่งเรื่องให้ผู้ควบคุมงานจัดหาเพิ่มเติมเรียบร้อย")` → Web UI.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-119`: `showSuccessMessage("บันทึกการโอนเงินสำเร็จ")` → Web UI.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-80`: `showSuccessMessage("บันทึกผลการพิจารณาสำเร็จ")` → Web UI.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-273`: `showSuccessMessage("บันทึกตารางงานสำเร็จ")` → Web UI.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-118`: `showSuccessMessage("บันทึกผลการพิจารณาสำเร็จ")` → Web UI.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-162`: `showSuccessMessage("ประมวลผลสำเร็จ")` → Web UI.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-140`: `showSuccessMessage("บันทึกหลักฐานการส่งมอบและแจ้งผู้รับแล้ว")` → Web UI.
- SD-2A, edge `N13osf1N4320fgIz80b4-62`: `clickEquipmentProcurementMenu()` → Web UI.
- SD-2A, edge `N13osf1N4320fgIz80b4-92`: `redirectToProcurementFormPage()` → Web UI.
- SD-2A, edge `N13osf1N4320fgIz80b4-96`: `enterProcurementData(itemQuantities, unitPrices)` → Web UI.
- SD-2A, edge `N13osf1N4320fgIz80b4-100`: `attachReceiptPhoto()` → Web UI.
- SD-2A, edge `N13osf1N4320fgIz80b4-104`: `clickConfirmProcurement()` → Web UI.
- SD-2A, edge `N13osf1N4320fgIz80b4-152`: `redirectToDashboard()` → Web UI.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-117`: `redirectToDashboard()` → Web UI.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-195`: `clickWorkScheduleMenu()` → Web UI.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-209`: `redirectToCreateSchedulePage()` → Web UI.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-213`: `enterScheduleData(assignmentId, work_date, shift_start_time, workerIds)` → Web UI.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-217`: `clickReviewSchedule()` → Web UI.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-247`: `showScheduleReviewPage()` → Web UI.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-251`: `clickConfirmSchedule()` → Web UI.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-271`: `redirectToScheduleListPage()` → Web UI.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-65`: `clickLeaveAndReplacementMenu()` → Web UI.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-85`: `selectLeaveRequest(requestId)` → Web UI.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-99`: `redirectToLeaveReviewPage()` → Web UI.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-103`: `chooseDecision(approve \| reject)` → Web UI.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-119`: `selectReplacement(replacementWorkerId)` → Web UI.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-123`: `clickReviewDecision()` → Web UI.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-70`: `clickReviewDecision()` → Web UI.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-153`: `showDecisionReviewPage()` → Web UI.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-100`: `showDecisionReviewPage()` → Web UI.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-155`: `showWarning("ยังไม่ได้จัดพนักงานแทน ต้องแจ้งผู้ควบคุมงานเพื่อติดตาม")` → Web UI.
- SD-9A, edge `ZDx1L2W36KSnwggz6YFt-71`: `showWarning("กรุณาติดต่อผู้ควบคุมงานเพื่อตรวจสอบสัญญาก่อนดำเนินงานต่อ")` → Web UI.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-159`: `clickConfirmDecision()` → Web UI.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-104`: `clickConfirmDecision()` → Web UI.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-191`: `redirectToLeaveRequestListPage()` → Web UI.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-193`: `showResultMessage(decisionResult)` → Web UI.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-136`: `showResultMessage(decisionResult)` → Web UI.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-30`: `clickAdditionalEquipmentReviewMenu()` → Web UI.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-50`: `selectExistingRequest(requisitionId)` → Web UI.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-84`: `showReadOnlyRequestAndReferences()` → Web UI.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-88`: `clickConsiderActionMethod(requisitionId)` → Web UI.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-90`: `redirectTo6A(requisitionId)` → Web UI.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-56`: `showReadOnlyOriginalRequest()` → Web UI.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-60`: `chooseDecision(purchase \| no_purchase)` → Web UI.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-64`: `showReasonOrGuidanceFields(decision)` → Web UI.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-66`: `enterPurchaseReasonOrGuidance(text)` → Web UI.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-134`: `redirectAfterDecision()` → Web UI.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-58`: `clickApprovedAdditionalPurchasesMenu()` → Web UI.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-78`: `selectOriginalRequest(requisitionId)` → Web UI.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-92`: `redirectToAdditionalPurchaseForm()` → Web UI.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-96`: `enterCurrentPurchaseData(itemQuantities, unitPrices, receiptIfCostPositive)` → Web UI.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-100`: `clickReviewPurchase()` → Web UI.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-126`: `showPurchaseReviewPage()` → Web UI.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-130`: `clickConfirmPurchase()` → Web UI.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-186`: `redirectAfterPurchase(status)` → Web UI.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-188`: `showPurchaseResultMessage(status)` → Web UI.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-42`: `clickEquipmentDeliveryMenu()` → Web UI.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-62`: `selectDeliveryRequest(requisitionId)` → Web UI.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-76`: `redirectToDeliveryForm()` → Web UI.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-80`: `enterDeliveryData(scheduleId, deliveredQuantities, description, photos)` → Web UI.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-84`: `clickReviewDelivery()` → Web UI.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-106`: `showDeliveryReviewPage()` → Web UI.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-110`: `clickConfirmDelivery()` → Web UI.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-138`: `redirectToDeliveryListPage()` → Web UI.
- SD-9A, edge `ZDx1L2W36KSnwggz6YFt-27`: `showReceivedSummary(project, area, summaryMessage)` → Web UI.
- SD-9A, edge `ZDx1L2W36KSnwggz6YFt-29`: `openAndReadSummary()` → Web UI.
- SD-9A, edge `ZDx1L2W36KSnwggz6YFt-33`: `selectArea(assignmentId)` → Web UI.
- SD-9A, edge `ZDx1L2W36KSnwggz6YFt-55`: `showContractAndAreaDetails()` → Web UI.
- SD-9A, edge `ZDx1L2W36KSnwggz6YFt-59`: `chooseContinueWork()` → Web UI.
- SD-9A, edge `ZDx1L2W36KSnwggz6YFt-69`: `openAreaDashboard(assignmentId)` → Web UI.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-252`: `clickCreateNewTOR()` → Web UI.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-258`: `redirectToContractFormPage()` → Web UI.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-278`: `redirectToScopeFormPage()` → Web UI.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-282`: `submitScopeData( scopeData)` → Web UI.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-306`: `redirectToContractReviewPage()` → Web UI.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-310`: `clickConfirmContract()` → Web UI.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-374`: `redirectToContractDetailPage(torId)` → Web UI.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-37`: `clickPendingFundingMenu()` → Web UI.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-67`: `redirectToTransferEvidencePage()` → Web UI.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-71`: `submitTransferEvidence(total_amount, transfer_ref_no, receipt_photo_url)` → Web UI.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-89`: `showConfirmationPopup()` → Web UI.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-93`: `clickConfirmTransfer()` → Web UI.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-38`: `clickRequisitionMenu()` → Web UI.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-42`: `clickRequisitionMenu()` → Web UI.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-49`: `redirectToRequisitionReviewPage()` → Web UI.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-64`: `redirectToRequisitionReviewPage()` → Web UI.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-52`: `selectDecision(approve \| reject)` → Web UI.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-68`: `selectDecision(approve \| reject)` → Web UI.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-53`: `showRequiredRejectionReasonField()` → Web UI.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-72`: `showRequiredRejectionReasonField()` → Web UI.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-55`: `enterRejectionReason(rejectionReason)` → Web UI.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-74`: `enterRejectionReason(rejectionReason)` → Web UI.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-57`: `submitReviewDecision(requisitionId, decision, rejectionReason)` → Web UI.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-78`: `submitReviewDecision(requisitionId, decision, rejectionReason)` → Web UI.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-78`: `redirectToRequisitionListPage()` → Web UI.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-116`: `redirectToRequisitionListPage()` → Web UI.
- SD-4S, edge `Na5WVMECRJByJS5ewn7f-123`: `displayAttendanceSummary(processingSummary)` → Web UI.
- SD-4S, edge `Na5WVMECRJByJS5ewn7f-127`: `reviewAttendanceSummary()` → Web UI.
- SD-4S, edge `Na5WVMECRJByJS5ewn7f-131`: `confirmProcessingResults()` → Web UI.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-96`: `clickCreatePayrollPeriod()` → Web UI.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-98`: `showPayrollPeriodForm()` → Web UI.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-102`: `submitPayrollPeriod(period_start, period_end)` → Web UI.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-160`: `redirectToPayrollSummaryPage()` → Web UI.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-166`: `reviewPayrollSummary()` → Web UI.
- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-153`: `openFinancialReportMenu()` → Web UI.
- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-155`: `showProjectAndDateFilters()` → Web UI.
- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-159`: `submitProfitFilters(torId, dateStart, dateEnd)` → Web UI.
- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-193`: `redirectToNetProfitReportPage()` → Web UI.
- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-195`: `renderFinancialDashboard(totalRevenue, totalLaborCost, totalMaterialCost, netProfit)` → Web UI.
- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-199`: `reviewFinancialReport()` → Web UI.
- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-203`: `confirmClosingAndContinue()` → Web UI.
- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-207`: `requestReportPDF()` → Web UI.
## EquipmentSurvey Controller

Role: Control. Sources: SD-1A.


| Diagram operation | Supporting SDs |
|---|---|
| `getPendingSurveyRequisitions()` | SD-1A |
| `loadEquipmentRequirements(requisitionId)` | SD-1A |
| `processSiteSurvey(requisitionId, itemQuantities)` | SD-1A |
| `validateExistingQuantities(itemQuantities)` | SD-1A |
| `calculateToBuyQty(required_qty, existing_qty)` | SD-1A |
| `selectRequisitionStatus(…)` | SD-1A |

Original calls and receiver evidence:

- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-52`: `getPendingSurveyRequisitions()` → EquipmentSurvey Controller.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-68`: `loadEquipmentRequirements(requisitionId)` → EquipmentSurvey Controller.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-86`: `processSiteSurvey(requisitionId, itemQuantities)` → EquipmentSurvey Controller.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-88`: `validateExistingQuantities(itemQuantities)` → EquipmentSurvey Controller.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-100`: `calculateToBuyQty(required_qty, existing_qty)` → EquipmentSurvey Controller.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-112`: `selectRequisitionStatus('completed')` → EquipmentSurvey Controller.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-116`: `selectRequisitionStatus('pending_procurement')` → EquipmentSurvey Controller.
## Requisition Repository

Role: Repository. Sources: SD-1A, SD-2A, SD-3.1S, SD-3S.


| Diagram operation | Supporting SDs |
|---|---|
| `findPendingSurveyRequisitions()` | SD-1A |
| `findEquipmentRequirements(requisitionId)` | SD-1A |
| `checkRequisitionItemExists(itemId)` | SD-1A, SD-2A |
| `updateSurveyItem(itemId, existing_qty)` | SD-1A |
| `updateRequisitionStatus(requisitionId, status)` | SD-1A |
| `findProcurementRequisitions()` | SD-2A |
| `findRemainingEquipment(requisitionId)` | SD-2A |
| `createActualExpense(expense_no, requisitionId, userId, total_amount, receipt_photo_url)` | SD-2A |
| `updateProcuredItem(itemId, quantity, unitPrice)` | SD-2A |
| `updateRequisitionStatus(requisitionId)` | SD-2A |
| `findPendingRequisitions()` | SD-3.1S, SD-3S |
| `findRequisitionDetails(requisitionId)` | SD-3.1S, SD-3S |
| `checkPendingRequisition(requisitionId)` | SD-3.1S, SD-3S |
| `approveRequisition(requisitionId)` | SD-3.1S, SD-3S |
| `rejectRequisition(requisitionId)` | SD-3.1S, SD-3S |

Original calls and receiver evidence:

- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-54`: `findPendingSurveyRequisitions()` → Requisition Repository.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-70`: `findEquipmentRequirements(requisitionId)` → Requisition Repository.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-92`: `checkRequisitionItemExists(itemId)` → Requisition Repository.
- SD-2A, edge `N13osf1N4320fgIz80b4-112`: `checkRequisitionItemExists(itemId)` → Requisition Repository.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-104`: `updateSurveyItem(itemId, existing_qty)` → Requisition Repository.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-120`: `updateRequisitionStatus(requisitionId, status)` → Requisition Repository.
- SD-2A, edge `N13osf1N4320fgIz80b4-66`: `findProcurementRequisitions()` → Requisition Repository.
- SD-2A, edge `N13osf1N4320fgIz80b4-82`: `findRemainingEquipment(requisitionId)` → Requisition Repository.
- SD-2A, edge `N13osf1N4320fgIz80b4-124`: `createActualExpense(expense_no, requisitionId, userId, total_amount, receipt_photo_url)` → Requisition Repository.
- SD-2A, edge `N13osf1N4320fgIz80b4-132`: `updateProcuredItem(itemId, quantity, unitPrice)` → Requisition Repository.
- SD-2A, edge `N13osf1N4320fgIz80b4-140`: `updateRequisitionStatus(requisitionId)` → Requisition Repository.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-40`: `findPendingRequisitions()` → Requisition Repository.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-46`: `findPendingRequisitions()` → Requisition Repository.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-46`: `findRequisitionDetails(requisitionId)` → Requisition Repository.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-58`: `findRequisitionDetails(requisitionId)` → Requisition Repository.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-61`: `checkPendingRequisition(requisitionId)` → Requisition Repository.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-86`: `checkPendingRequisition(requisitionId)` → Requisition Repository.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-65`: `approveRequisition(requisitionId)` → Requisition Repository.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-94`: `approveRequisition(requisitionId)` → Requisition Repository.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-71`: `rejectRequisition(requisitionId)` → Requisition Repository.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-104`: `rejectRequisition(requisitionId)` → Requisition Repository.
## DB Connector

Role: Service. Sources: SD-1A, SD-1S, SD-1W, SD-2A, SD-2S, SD-2W, SD-3.1S, SD-3A, SD-3S, SD-3W, SD-4A, SD-4S, SD-4W, SD-5A, SD-5S, SD-5W, SD-6A, SD-6S, SD-6W, SD-7A, SD-7W, SD-8A, SD-9A.


| Diagram operation | Supporting SDs |
|---|---|
| `executeQuery(sql)` | SD-1A, SD-1S, SD-1W, SD-2A, SD-2S, SD-2W, SD-3.1S, SD-3A, SD-3S, SD-3W, SD-4A, SD-4S, SD-4W, SD-5A, SD-5S, SD-5W, SD-6A, SD-6S, SD-6W, SD-7A, SD-7W, SD-8A, SD-9A |

Original calls and receiver evidence:

- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-56`: `executeQuery(sql = "SELECT ... FROM EQUIPMENT_REQUISITION ... WHERE er.status = 'pending_survey';")` → DB Connector.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-94`: `executeQuery(sql = "SELECT COUNT(item_id) FROM REQUISITION_ITEM WHERE item_id = $1;")` → DB Connector.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-122`: `executeQuery(sql = "UPDATE EQUIPMENT_REQUISITION SET status = $1 WHERE requisition_id = $2;")` → DB Connector.
- SD-1A, edge `GPv0WRyF4ufz6NQO2nEs-106`: `executeQuery(sql = "UPDATE REQUISITION_ITEM SET existing_qty = $1, to_buy_qty = GREATEST(required_qty - $1, 0) WHERE item_id = $3;")` → DB Connector.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-316`: `executeQuery(sql = " INSERT INTO CONTRACT_TOR ... ")` → DB Connector.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-356`: `executeQuery(sql = " INSERT INTO EQUIPMENT ... ")` → DB Connector.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-340`: `executeQuery(sql = " INSERT INTO EQUIPMENT_REQUISITION ... ")` → DB Connector.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-324`: `executeQuery(sql = " INSERT INTO LOCATION ... ")` → DB Connector.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-364`: `executeQuery(sql = " INSERT INTO REQUISITION_ITEM ... ")` → DB Connector.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-332`: `executeQuery(sql = " INSERT INTO TOR_LOCATION_ASSIGNMENT ... ")` → DB Connector.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-290`: `executeQuery(sql = " SELECT ... ")` → DB Connector.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-298`: `executeQuery(sql = " SELECT ... ")` → DB Connector.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-270`: `executeQuery(sql = " SELECT COUNT ...")` → DB Connector.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-348`: `executeQuery(sql = " SELECT equipment_id ..." )` → DB Connector.
- SD-1W, edge `CF87uobWHUu3A-gUhq7B-35`: `executeQuery(sql = "SELECT ws.schedule_id, ws.work_date, ws.shift_start_time, ws.shift_status, l.location_name, l.address, ct.project_name FROM WORK_SCHEDULE ws JOIN WORKER w ... JOIN TOR_LOCATION_ASSIGNMENT tla ... JOIN LOCATION l ... JOIN CONTRACT_TOR ct ... WHERE w.user_id = $1 AND ws.work_date >= CURRENT_DATE ORDER BY ws.work_date ASC, ws.shift_start_time ASC;")` → DB Connector.
- SD-2A, edge `N13osf1N4320fgIz80b4-126`: `executeQuery(sql = "INSERT INTO EXPENSE_CLAIM (..., expense_type) VALUES (..., 'actual_expense') RETURNING expense_id;")` → DB Connector.
- SD-2A, edge `N13osf1N4320fgIz80b4-68`: `executeQuery(sql = "SELECT ... FROM EQUIPMENT_REQUISITION ... WHERE er.status IN ('pending_procurement', 'approved');")` → DB Connector.
- SD-2A, edge `N13osf1N4320fgIz80b4-114`: `executeQuery(sql = "SELECT COUNT(item_id) FROM REQUISITION_ITEM WHERE item_id = $1;")` → DB Connector.
- SD-2A, edge `N13osf1N4320fgIz80b4-84`: `executeQuery(sql = "SELECT item_id, equipment_name, (to_buy_qty - actual_qty) AS remaining_qty FROM REQUISITION_ITEM WHERE requisition_id = $1 AND (to_buy_qty - actual_qty) > 0;")` → DB Connector.
- SD-2A, edge `N13osf1N4320fgIz80b4-142`: `executeQuery(sql = "UPDATE EQUIPMENT_REQUISITION ... CASE WHEN EXISTS (... ri.actual_qty < ri.to_buy_qty) THEN 'pending_supervisor' ELSE 'completed' END WHERE er.requisition_id = $1;")` → DB Connector.
- SD-2A, edge `N13osf1N4320fgIz80b4-134`: `executeQuery(sql = "UPDATE REQUISITION_ITEM SET actual_qty = actual_qty + $1, actual_price = $2 WHERE item_id = $3;")` → DB Connector.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-99`: `executeQuery(sql = "INSERT INTO EXPENSE_CLAIM ... 'fund_transfer' ... RETURNING expense_id;")` → DB Connector.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-59`: `executeQuery(sql = "SELECT ... FROM REQUISITION_ITEM ... WHERE (ri.to_buy_qty - ri.actual_qty) > 0;")` → DB Connector.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-43`: `executeQuery(sql = "SELECT ... WHERE er.status = 'pending_supervisor';")` → DB Connector.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-81`: `executeQuery(sql = "SELECT COUNT(expense_id) FROM EXPENSE_CLAIM WHERE transfer_ref_no = $1;")` → DB Connector.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-107`: `executeQuery(sql = "UPDATE EQUIPMENT_REQUISITION SET status = 'pending_procurement' WHERE requisition_id = $1;")` → DB Connector.
- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-73`: `executeQuery(sql = "INSERT INTO LEAVE_REQUEST (user_id, leave_date, reason, status, is_advance_notice) VALUES ($1, $2, $3, 'pending', $4) RETURNING request_id;")` → DB Connector.
- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-61`: `executeQuery(sql = "SELECT COUNT(request_id) FROM LEAVE_REQUEST WHERE user_id = $1 AND leave_date = $2;")` → DB Connector.
- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-39`: `executeQuery(sql = "SELECT ws.work_date, l.location_name FROM WORK_SCHEDULE ws JOIN WORKER w ... JOIN TOR_LOCATION_ASSIGNMENT tla ... JOIN LOCATION l ... WHERE w.user_id = $1 AND ws.work_date >= CURRENT_DATE;")` → DB Connector.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-62`: `executeQuery(Q3S.1: SELECT COUNT(requisition_id) FROM equipment_requisition WHERE requisition_id = $1 AND status = 'pending_supervisor';)` → DB Connector.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-66`: `executeQuery(Q3S.2: UPDATE equipment_requisition SET status = 'approved' WHERE requisition_id = $1;)` → DB Connector.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-72`: `executeQuery(Q3S.3: UPDATE equipment_requisition SET status = 'rejected' WHERE requisition_id = $1;)` → DB Connector.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-261`: `executeQuery(sql = "INSERT INTO WORK_SCHEDULE (...) VALUES ($1, $2, $3, $4, 'scheduled') RETURNING schedule_id;")` → DB Connector.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-201`: `executeQuery(sql = "SELECT ... FROM TOR_LOCATION_ASSIGNMENT ... WHERE tla.tor_id = $1;")` → DB Connector.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-239`: `executeQuery(sql = "SELECT COUNT(schedule_id) FROM WORK_SCHEDULE WHERE worker_id = $1 AND work_date = $2;")` → DB Connector.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-231`: `executeQuery(sql = "SELECT COUNT(w.worker_id) FROM WORKER w JOIN "USER" u ... WHERE w.worker_id = $1 AND w.is_available = true AND u.is_active = true;")` → DB Connector.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-88`: `executeQuery(sql = "SELECT COUNT(requisition_id) FROM equipment_requisition WHERE requisition_id = $1 AND status = 'pending_supervisor';")` → DB Connector.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-96`: `executeQuery(sql = "UPDATE equipment_requisition SET status = 'approved' WHERE requisition_id = $1;")` → DB Connector.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-106`: `executeQuery(sql = "UPDATE equipment_requisition SET status = 'rejected' WHERE requisition_id = $1;")` → DB Connector.
- SD-3W, edge `E7oHi6kEsZcUCtm5OUxS-73`: `executeQuery(sql = "INSERT INTO ATTENDANCE (schedule_id, worker_id, work_date, check_in) VALUES ($1, $2, CURRENT_DATE, CURRENT_TIMESTAMP) ON CONFLICT (schedule_id, worker_id) DO UPDATE SET check_in = EXCLUDED.check_in RETURNING attendance_id;")` → DB Connector.
- SD-3W, edge `E7oHi6kEsZcUCtm5OUxS-65`: `executeQuery(sql = "SELECT COUNT(attendance_id) FROM ATTENDANCE a JOIN WORKER w ... WHERE w.user_id = $1 AND a.schedule_id = $2 AND a.check_in IS NOT NULL;")` → DB Connector.
- SD-3W, edge `E7oHi6kEsZcUCtm5OUxS-39`: `executeQuery(sql = "SELECT ws.schedule_id, ws.shift_start_time, l.location_name, ct.project_name FROM WORK_SCHEDULE ws JOIN WORKER w ... JOIN TOR_LOCATION_ASSIGNMENT tla ... JOIN LOCATION l ... JOIN CONTRACT_TOR ct ... WHERE w.user_id = $1 AND ws.work_date = CURRENT_DATE;")` → DB Connector.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-177`: `executeQuery(sql = "INSERT INTO WORK_SCHEDULE (...) SELECT assignment_id, $1, work_date, shift_start_time, 'scheduled' FROM WORK_SCHEDULE WHERE schedule_id = $2 RETURNING schedule_id;")` → DB Connector.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-91`: `executeQuery(sql = "SELECT ... FROM LEAVE_REQUEST ... WHERE lr.request_id = $1 AND ws.assignment_id = $2 AND lr.status = 'pending' AND ws.shift_status = 'scheduled';")` → DB Connector.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-75`: `executeQuery(sql = "SELECT ... FROM LEAVE_REQUEST ... WHERE lr.status = 'pending' AND ws.assignment_id = $1 AND ws.shift_status = 'scheduled' ORDER BY ...;")` → DB Connector.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-111`: `executeQuery(sql = "SELECT ... FROM WORKER w JOIN "USER" u ... WHERE w.is_available = true AND u.is_active = true AND w.worker_id <> $1 AND NOT EXISTS (... WORK_SCHEDULE ...) AND NOT EXISTS (... approved leave ...) ORDER BY ...;")` → DB Connector.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-133`: `executeQuery(sql = "SELECT COUNT(request_id) FROM LEAVE_REQUEST WHERE request_id = $1 AND status = 'pending';")` → DB Connector.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-141`: `executeQuery(sql = "SELECT COUNT(ws.schedule_id) FROM WORK_SCHEDULE ... WHERE lr.request_id = $1 AND ws.schedule_id = $2 AND ws.assignment_id = $3 AND ws.shift_status = 'scheduled';")` → DB Connector.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-169`: `executeQuery(sql = "UPDATE LEAVE_REQUEST SET status = $1 WHERE request_id = $2 AND status = 'pending' RETURNING request_id;")` → DB Connector.
- SD-4S, edge `Na5WVMECRJByJS5ewn7f-75`: `executeQuery(sql = "SELECT COUNT(request_id) FROM leave_request WHERE user_id = $2 AND leave_date = $3 AND is_advance_notice = true AND status = 'approved';")` → DB Connector.
- SD-4S, edge `Na5WVMECRJByJS5ewn7f-67`: `executeQuery(sql = "SELECT check_in FROM attendance WHERE schedule_id = $1 AND worker_id = $2;")` → DB Connector.
- SD-4S, edge `Na5WVMECRJByJS5ewn7f-83`: `executeQuery(sql = "UPDATE attendance SET status = 'absent' WHERE attendance_id = $4;")` → DB Connector.
- SD-4S, edge `Na5WVMECRJByJS5ewn7f-111`: `executeQuery(sql = "UPDATE attendance SET status = 'late' WHERE attendance_id = $4;")` → DB Connector.
- SD-4S, edge `Na5WVMECRJByJS5ewn7f-91`: `executeQuery(sql = "UPDATE attendance SET status = 'leave' WHERE attendance_id = $4;")` → DB Connector.
- SD-4S, edge `Na5WVMECRJByJS5ewn7f-103`: `executeQuery(sql = "UPDATE attendance SET status = 'on_time' WHERE attendance_id = $4;")` → DB Connector.
- SD-4W, edge `4W7N71NcVjHcfgNTzLRr-61`: `executeQuery(sql = "INSERT INTO WORK_EVIDENCE (worker_id, schedule_id, description, photo_url, submitted_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP) RETURNING evidence_id;")` → DB Connector.
- SD-4W, edge `4W7N71NcVjHcfgNTzLRr-35`: `executeQuery(sql = "SELECT a.attendance_id, a.schedule_id, w.worker_id FROM ATTENDANCE a JOIN WORKER w ON a.worker_id = w.worker_id WHERE w.user_id = $1 AND a.work_date = CURRENT_DATE AND a.check_in IS NOT NULL AND a.check_out IS NULL;")` → DB Connector.
- SD-4W, edge `4W7N71NcVjHcfgNTzLRr-69`: `executeQuery(sql = "UPDATE ATTENDANCE SET check_out = CURRENT_TIMESTAMP WHERE attendance_id = $1;")` → DB Connector.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-40`: `executeQuery(sql = "SELECT ... FROM EQUIPMENT_REQUISITION ... WHERE er.assignment_id = $1 AND er.requisition_type = 'additional' AND er.status = 'pending_survey' ORDER BY er.created_at;")` → DB Connector.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-60`: `executeQuery(sql = "SELECT ... FROM EQUIPMENT_REQUISITION ... WHERE er.requisition_id = $1 AND er.assignment_id = $2 AND er.requisition_type = 'additional' AND er.status = 'pending_survey' ORDER BY ri.item_id;")` → DB Connector.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-68`: `executeQuery(sql = "SELECT ... GREATEST(bi.to_buy_qty - COALESCE(bi.actual_qty,0),0) AS remaining_qty FROM EQUIPMENT_REQUISITION ... AND base.requisition_type = 'tor_base' ... WHERE current_req.requisition_id = $1;")` → DB Connector.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-76`: `executeQuery(sql = "SELECT DISTINCT ... FROM EQUIPMENT_REQUISITION ... other.requisition_id <> current_req.requisition_id ... WHERE current_req.requisition_id = $1 AND other.status IN ('pending_survey', 'pending_procurement','pending','approved', 'pending_supervisor');")` → DB Connector.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-132`: `executeQuery(sql = "INSERT INTO deduction_transaction (worker_id, penalty_amount, reason) VALUES ($1, $2, $3);")` → DB Connector.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-152`: `executeQuery(sql = "INSERT INTO payroll (..., is_paid) VALUES (..., false) RETURNING payroll_id;")` → DB Connector.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-124`: `executeQuery(sql = "SELECT ... CASE ... END AS penalty_amount, ... AS penalty_reason FROM attendance ...;")` → DB Connector.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-140`: `executeQuery(sql = "SELECT COALESCE(SUM(penalty_amount), 0) AS total_deduction FROM deduction_transaction WHERE worker_id = $1 AND created_at BETWEEN $2 AND $3;")` → DB Connector.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-112`: `executeQuery(sql = "SELECT COUNT(a.attendance_id) AS work_days, u.daily_wage ... WHERE a.status IN ('on_time', 'late') ...;")` → DB Connector.
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-83`: `executeQuery(sql = "INSERT INTO EQUIPMENT_REQUISITION (requisition_no, user_id, assignment_id, reason, status, requisition_type) VALUES ($1, $2, $3, $4, 'pending', 'additional') RETURNING requisition_id;")` → DB Connector.
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-91`: `executeQuery(sql = "INSERT INTO REQUISITION_ITEM (requisition_id, equipment_id, required_qty, to_buy_qty) VALUES ($1, $2, $3, $3);")` → DB Connector.
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-75`: `executeQuery(sql = "SELECT COUNT(equipment_id) FROM EQUIPMENT WHERE equipment_id = $1 AND is_active = true;")` → DB Connector.
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-49`: `executeQuery(sql = "SELECT equipment_id, equipment_name FROM EQUIPMENT WHERE is_active = true ORDER BY equipment_name ASC;")` → DB Connector.
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-41`: `executeQuery(sql = "SELECT ws.assignment_id, l.location_name, u.user_id FROM ATTENDANCE a JOIN WORK_SCHEDULE ws ... JOIN WORKER w ... JOIN "USER" u ... JOIN TOR_LOCATION_ASSIGNMENT tla ... JOIN LOCATION l ... WHERE u.user_id = $1 AND a.work_date = CURRENT_DATE AND a.check_in IS NOT NULL AND a.check_out IS NULL;")` → DB Connector.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-44`: `executeQuery(sql = "SELECT ... u.line_id FROM EQUIPMENT_REQUISITION er JOIN "USER" u ... WHERE er.requisition_id = $1 AND er.assignment_id = $2 AND er.requisition_type = 'additional' AND er.status = 'pending_survey';")` → DB Connector.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-80`: `executeQuery(sql = "SELECT er.requisition_id, er.user_id, er.assignment_id FROM EQUIPMENT_REQUISITION ... WHERE er.requisition_id = $1 AND er.assignment_id = $2 AND er.requisition_type = 'additional' AND er.status = 'pending_survey';")` → DB Connector.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-88`: `executeQuery(sql = "SELECT ri.item_id, ri.equipment_id, ri.required_qty, ri.to_buy_qty, e.is_active FROM REQUISITION_ITEM ri JOIN EQUIPMENT e ... WHERE ri.requisition_id = $1;")` → DB Connector.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-114`: `executeQuery(sql = "UPDATE EQUIPMENT_REQUISITION SET status = 'pending_supervisor', reason = $1 WHERE requisition_id = $2 AND assignment_id = $3 AND requisition_type = 'additional' AND status = 'pending_survey' RETURNING requisition_id, requisition_no, user_id;")` → DB Connector.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-124`: `executeQuery(sql = "UPDATE EQUIPMENT_REQUISITION SET status = 'rejected', reviewed_by = $1, reviewed_at = CURRENT_TIMESTAMP WHERE requisition_id = $2 AND assignment_id = $3 AND requisition_type = 'additional' AND status = 'pending_survey' RETURNING requisition_id, requisition_no, user_id;")` → DB Connector.
- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-181`: `executeQuery(sql = "SELECT COALESCE(SUM(ec.total_amount), 0) AS total_material_cost FROM expense_claim ec JOIN equipment_requisition er ON ec.expense_id = er.expense_id ... WHERE tla.tor_id = $1 AND er.status = 'approved' AND ec.created_at BETWEEN $2 AND $3;")` → DB Connector.
- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-173`: `executeQuery(sql = "SELECT COALESCE(SUM(net_pay), 0) AS total_labor_cost FROM payroll WHERE period_start >= $1 AND period_end <= $2 AND is_paid = true;")` → DB Connector.
- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-165`: `executeQuery(sql = "SELECT COALESCE(SUM(net_received), 0) AS total_revenue FROM company_invoice WHERE tor_id = $1 AND status = 'paid' AND billing_month = $2;")` → DB Connector.
- SD-6W, edge `OQYhqyozvix0y_K6DQpe-21`: `executeQuery(sql = "SELECT u.line_id, er.status, e.equipment_name, ri.actual_qty FROM EQUIPMENT_REQUISITION er JOIN REQUISITION_ITEM ri ... JOIN EQUIPMENT e ... JOIN "USER" u ... WHERE er.requisition_id = $1;")` → DB Connector.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-164`: `executeQuery(sql = "INSERT INTO EXPENSE_CLAIM (requisition_id, user_id, expense_no, total_amount, receipt_photo_url, expense_type) VALUES ($1, $2, $3, $4, $5, 'actual_expense') RETURNING expense_id;")` → DB Connector.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-68`: `executeQuery(sql = "SELECT ... FROM EQUIPMENT_REQUISITION ... WHERE er.assignment_id = $1 AND er.requisition_type = 'additional' AND (er.status = 'approved' OR (... pending_supervisor, prior approval and fund_transfer ...)) AND EXISTS (... actual_qty < to_buy_qty) ORDER BY er.created_at;")` → DB Connector.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-144`: `executeQuery(sql = "SELECT item_id, to_buy_qty, COALESCE(actual_qty,0) AS actual_qty, COALESCE(actual_price,0) AS actual_price FROM REQUISITION_ITEM WHERE requisition_id = $1 ORDER BY item_id FOR UPDATE;")` → DB Connector.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-114`: `executeQuery(sql = "SELECT item_id, to_buy_qty, COALESCE(actual_qty,0) AS actual_qty, COALESCE(actual_price,0) AS actual_price, GREATEST(to_buy_qty - COALESCE(actual_qty,0),0) AS remaining_qty FROM REQUISITION_ITEM WHERE item_id = $1 AND requisition_id = $2;")` → DB Connector.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-136`: `executeQuery(sql = "SELECT requisition_id, assignment_id, requisition_type, status, reviewed_by, reviewed_at FROM EQUIPMENT_REQUISITION WHERE requisition_id = $1 FOR UPDATE;")` → DB Connector.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-84`: `executeQuery(sql = "SELECT ri.item_id, e.equipment_name, ri.to_buy_qty, COALESCE(ri.actual_qty,0) AS actual_qty, GREATEST(ri.to_buy_qty - COALESCE(ri.actual_qty,0),0) AS remaining_qty FROM REQUISITION_ITEM ... WHERE ri.requisition_id = $1 AND ri.to_buy_qty > COALESCE(ri.actual_qty,0) ORDER BY ri.item_id;")` → DB Connector.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-172`: `executeQuery(sql = "UPDATE EQUIPMENT_REQUISITION er SET status = CASE WHEN EXISTS (... COALESCE(ri.actual_qty,0) < ri.to_buy_qty) THEN 'pending_supervisor' ELSE 'completed' END WHERE er.requisition_id = $1 RETURNING status;")` → DB Connector.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-156`: `executeQuery(sql = "UPDATE REQUISITION_ITEM SET actual_price = (... weighted average ...), actual_qty = COALESCE(actual_qty,0) + $1 WHERE item_id = $3 AND requisition_id = $4 AND $1 > 0 AND COALESCE(actual_qty,0) + $1 <= to_buy_qty RETURNING item_id;")` → DB Connector.
- SD-7W, edge `uUxBqSW2nSQnUcCSPHx_-58`: `executeQuery(sql = "SELECT dt.amount, dt.reason, a.work_date FROM DEDUCTION_TRANSACTION dt JOIN ATTENDANCE a ON dt.attendance_id = a.attendance_id JOIN WORKER w ON a.worker_id = w.worker_id WHERE w.user_id = $1 AND a.work_date BETWEEN $2 AND $3;")` → DB Connector.
- SD-7W, edge `uUxBqSW2nSQnUcCSPHx_-50`: `executeQuery(sql = "SELECT p.payroll_id, p.period_start, p.period_end, p.total_wage, p.total_deduction, p.net_wage, p.status FROM PAYROLL p JOIN WORKER w ON p.worker_id = w.worker_id WHERE w.user_id = $1 AND TO_CHAR(p.period_start, 'YYYY-MM') = $2 ORDER BY p.period_start DESC;")` → DB Connector.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-120`: `executeQuery(sql = "INSERT INTO WORK_EVIDENCE (schedule_id, worker_id, description, photo_url, submitted_at) SELECT ws.schedule_id, ws.worker_id, $3, $4, CURRENT_TIMESTAMP FROM WORK_SCHEDULE ... WHERE ws.schedule_id = $1 AND er.requisition_id = $2 AND er.requisition_type = 'additional' AND er.status = 'completed' RETURNING evidence_id;")` → DB Connector.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-52`: `executeQuery(sql = "SELECT ... FROM EQUIPMENT_REQUISITION ... WHERE er.assignment_id = $1 AND er.requisition_type = 'additional' AND er.status = 'completed' ORDER BY er.created_at;")` → DB Connector.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-68`: `executeQuery(sql = "SELECT er.requisition_no, er.user_id, ... u.line_id, ri.item_id, e.equipment_name, COALESCE(ri.actual_qty,0) AS actual_qty FROM EQUIPMENT_REQUISITION ... WHERE er.requisition_id = $1 AND er.status = 'completed' ORDER BY ri.item_id;")` → DB Connector.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-128`: `executeQuery(sql = "SELECT u.line_id, er.requisition_no, e.equipment_name, ri.actual_qty FROM EQUIPMENT_REQUISITION ... WHERE er.requisition_id = $1 AND er.requisition_type = 'additional' AND er.status = 'completed';")` → DB Connector.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-98`: `executeQuery(sql = "SELECT ws.schedule_id, ws.worker_id FROM WORK_SCHEDULE ... WHERE ws.schedule_id = $1 AND er.requisition_id = $2 AND er.requisition_type = 'additional' AND er.status = 'completed' AND ws.work_date <= CURRENT_DATE;")` → DB Connector.
- SD-9A, edge `ZDx1L2W36KSnwggz6YFt-47`: `executeQuery(sql = "SELECT ct.tor_id, ct.contract_no, ct.project_name, ct.start_date, ct.end_date, ct.status, tla.assignment_id, tla.required_workers, l.location_id, l.location_name FROM TOR_LOCATION_ASSIGNMENT tla JOIN CONTRACT_TOR ct ... JOIN LOCATION l ... WHERE tla.assignment_id = $1;")` → DB Connector.
## LINE Notification Service

Role: Service. Sources: SD-1A, SD-1S, SD-2A, SD-2S, SD-2W, SD-3.1S, SD-3S, SD-3W, SD-5W, SD-6A, SD-6W, SD-8A.


| Diagram operation | Supporting SDs |
|---|---|
| `sendShortageNotification(ผู้ควบคุมงาน)` | SD-2A |
| `sendGuidanceToOriginalRequester(line_id, originalItems, guidanceText)` | SD-6A |
| `sendActualDeliverySummary(line_id, requisition_no, itemsAndQuantities)` | SD-8A |
| `sendNewContractNotification(torId, projectName)` | SD-1S |
| `sendProcurementFundingNotification(ผู้ดูแลงาน, projectName)` | SD-2S |
| `sendApprovalNotification(ผู้สร้างคำขอ)` | SD-3.1S, SD-3S |
| `sendRejectionNotification(ผู้สร้างคำขอ, rejectionReason)` | SD-3.1S, SD-3S |
| `notifyAssistantOfLeave(requestId)` | SD-2W |
| `notifyAssistantOfAdditionalRequest(location_name)` | SD-5W |
| `sendEquipmentReceivedPush(line_id, equipment_name, actual_qty)` | SD-6W |
| `sendNoPurchaseGuidancePush(line_id, equipment_name)` | SD-6W |

Original calls and receiver evidence:

- SD-2A, edge `N13osf1N4320fgIz80b4-148`: `sendShortageNotification(ผู้ควบคุมงาน)` → LINE Notification Service.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-130`: `sendGuidanceToOriginalRequester(line_id, originalItems, guidanceText)` → LINE Notification Service.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-134`: `sendActualDeliverySummary(line_id, requisition_no, itemsAndQuantities)` → LINE Notification Service.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-370`: `sendNewContractNotification(torId, projectName)` → LINE Notification Service.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-113`: `sendProcurementFundingNotification(ผู้ดูแลงาน, projectName)` → LINE Notification Service.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-70`: `sendApprovalNotification(ผู้สร้างคำขอ)` → LINE Notification Service.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-102`: `sendApprovalNotification(ผู้สร้างคำขอ)` → LINE Notification Service.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-76`: `sendRejectionNotification(ผู้สร้างคำขอ, rejectionReason)` → LINE Notification Service.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-112`: `sendRejectionNotification(ผู้สร้างคำขอ, rejectionReason)` → LINE Notification Service.
- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-81`: `notifyAssistantOfLeave(requestId)` → LINE Notification Service.
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-97`: `notifyAssistantOfAdditionalRequest(location_name)` → LINE Notification Service.
- SD-6W, edge `OQYhqyozvix0y_K6DQpe-27`: `sendEquipmentReceivedPush(line_id, equipment_name, actual_qty)` → LINE Notification Service.
- SD-6W, edge `OQYhqyozvix0y_K6DQpe-31`: `sendNoPurchaseGuidancePush(line_id, equipment_name)` → LINE Notification Service.
## Procurement Controller

Role: Control. Sources: SD-2A.


| Diagram operation | Supporting SDs |
|---|---|
| `getProcurementRequisitions()` | SD-2A |
| `loadProcurementForm(requisitionId)` | SD-2A |
| `processProcurement(requisitionId, procurementData)` | SD-2A |
| `validateProcurementData(procurementData)` | SD-2A |
| `determineProcurementStatus(itemQuantities)` | SD-2A |

Original calls and receiver evidence:

- SD-2A, edge `N13osf1N4320fgIz80b4-64`: `getProcurementRequisitions()` → Procurement Controller.
- SD-2A, edge `N13osf1N4320fgIz80b4-80`: `loadProcurementForm(requisitionId)` → Procurement Controller.
- SD-2A, edge `N13osf1N4320fgIz80b4-106`: `processProcurement(requisitionId, procurementData)` → Procurement Controller.
- SD-2A, edge `N13osf1N4320fgIz80b4-108`: `validateProcurementData(procurementData)` → Procurement Controller.
- SD-2A, edge `N13osf1N4320fgIz80b4-120`: `determineProcurementStatus(itemQuantities)` → Procurement Controller.
## WorkSchedule Controller

Role: Control. Sources: SD-1W, SD-3A.


| Diagram operation | Supporting SDs |
|---|---|
| `loadWorkScheduleForm(torId)` | SD-3A |
| `reviewSchedule(scheduleData)` | SD-3A |
| `validateScheduleFormat(scheduleData)` | SD-3A |
| `validateAreaExistenceAndPermission(assignmentId)` | SD-3A |
| `confirmSchedule(scheduleData)` | SD-3A |
| `recheckLatestScheduleData(scheduleData)` | SD-3A |
| `notifyAssignedWorkers(workerIds)` | SD-3A |
| `loadMySchedule(LINEUserId)` | SD-1W |
| `authenticateAndResolveLinkedUser(tokenOrLINEUserId)` | SD-1W |

Original calls and receiver evidence:

- SD-3A, edge `brZw5Uc51kGho3DvyTKA-197`: `loadWorkScheduleForm(torId)` → WorkSchedule Controller.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-219`: `reviewSchedule(scheduleData)` → WorkSchedule Controller.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-221`: `validateScheduleFormat(scheduleData)` → WorkSchedule Controller.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-225`: `validateAreaExistenceAndPermission(assignmentId)` → WorkSchedule Controller.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-253`: `confirmSchedule(scheduleData)` → WorkSchedule Controller.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-255`: `recheckLatestScheduleData(scheduleData)` → WorkSchedule Controller.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-267`: `notifyAssignedWorkers(workerIds)` → WorkSchedule Controller.
- SD-1W, edge `CF87uobWHUu3A-gUhq7B-27`: `loadMySchedule(LINEUserId)` → WorkSchedule Controller.
- SD-1W, edge `CF87uobWHUu3A-gUhq7B-29`: `authenticateAndResolveLinkedUser(tokenOrLINEUserId)` → WorkSchedule Controller.
## WorkSchedule Repository

Role: Repository. Sources: SD-1W, SD-3A.


| Diagram operation | Supporting SDs |
|---|---|
| `findTORAssignments(torId)` | SD-3A |
| `checkWorkerAvailable(workerId)` | SD-3A |
| `checkDuplicateSchedule(workerId, work_date)` | SD-3A |
| `createWorkSchedule(assignmentId, workerId, work_date, shift_start_time)` | SD-3A |
| `findUpcomingSchedules(userId)` | SD-1W |

Original calls and receiver evidence:

- SD-3A, edge `brZw5Uc51kGho3DvyTKA-199`: `findTORAssignments(torId)` → WorkSchedule Repository.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-229`: `checkWorkerAvailable(workerId)` → WorkSchedule Repository.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-237`: `checkDuplicateSchedule(workerId, work_date)` → WorkSchedule Repository.
- SD-3A, edge `brZw5Uc51kGho3DvyTKA-259`: `createWorkSchedule(assignmentId, workerId, work_date, shift_start_time)` → WorkSchedule Repository.
- SD-1W, edge `CF87uobWHUu3A-gUhq7B-33`: `findUpcomingSchedules(userId)` → WorkSchedule Repository.
## LeaveReview Controller

Role: Control. Sources: SD-4A.


| Diagram operation | Supporting SDs |
|---|---|
| `loadPendingLeaveRequests(assignmentId)` | SD-4A |
| `checkAreaPermission(assignmentId)` | SD-4A |
| `loadLeaveReviewForm(requestId, assignmentId)` | SD-4A |
| `loadReplacementCandidates(leavingWorkerId, leave_date)` | SD-4A |
| `reviewLeaveDecision(reviewData)` | SD-4A |
| `validateLeaveDecision(reviewData)` | SD-4A |
| `recheckSelectedReplacement(reviewData)` | SD-4A |
| `confirmLeaveDecision(reviewData)` | SD-4A |
| `recheckLatestData(reviewData)` | SD-4A |
| `notifySupervisorReplacementMissing()` | SD-4A |
| `notifyLeavingWorker(decision)` | SD-4A |
| `notifyReplacementWorker(replacementWorkerId)` | SD-4A |

Original calls and receiver evidence:

- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-67`: `loadPendingLeaveRequests(assignmentId)` → LeaveReview Controller.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-69`: `checkAreaPermission(assignmentId)` → LeaveReview Controller.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-87`: `loadLeaveReviewForm(requestId, assignmentId)` → LeaveReview Controller.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-107`: `loadReplacementCandidates(leavingWorkerId, leave_date)` → LeaveReview Controller.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-125`: `reviewLeaveDecision(reviewData)` → LeaveReview Controller.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-127`: `validateLeaveDecision(reviewData)` → LeaveReview Controller.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-147`: `recheckSelectedReplacement(reviewData)` → LeaveReview Controller.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-161`: `confirmLeaveDecision(reviewData)` → LeaveReview Controller.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-163`: `recheckLatestData(reviewData)` → LeaveReview Controller.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-183`: `notifySupervisorReplacementMissing()` → LeaveReview Controller.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-185`: `notifyLeavingWorker(decision)` → LeaveReview Controller.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-187`: `notifyReplacementWorker(replacementWorkerId)` → LeaveReview Controller.
## LeaveReview Repository

Role: Repository. Sources: SD-4A.


| Diagram operation | Supporting SDs |
|---|---|
| `findPendingLeaveRequests(assignmentId)` | SD-4A |
| `findLeaveDetails(requestId, assignmentId)` | SD-4A |
| `findReplacementCandidates(leavingWorkerId, leave_date)` | SD-4A |
| `checkPendingLeave(requestId)` | SD-4A |
| `checkAffectedSchedule(requestId, scheduleId, assignmentId)` | SD-4A |
| `updateLeaveStatus(requestId, approved \| rejected)` | SD-4A |
| `createReplacementSchedule(replacementWorkerId, scheduleId)` | SD-4A |

Original calls and receiver evidence:

- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-73`: `findPendingLeaveRequests(assignmentId)` → LeaveReview Repository.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-89`: `findLeaveDetails(requestId, assignmentId)` → LeaveReview Repository.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-109`: `findReplacementCandidates(leavingWorkerId, leave_date)` → LeaveReview Repository.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-131`: `checkPendingLeave(requestId)` → LeaveReview Repository.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-139`: `checkAffectedSchedule(requestId, scheduleId, assignmentId)` → LeaveReview Repository.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-167`: `updateLeaveStatus(requestId, approved \| rejected)` → LeaveReview Repository.
- SD-4A, edge `g_5Bvjasd6hpiueaxGk9-175`: `createReplacementSchedule(replacementWorkerId, scheduleId)` → LeaveReview Repository.
## EquipmentReview Controller

Role: Control. Sources: SD-5A.


| Diagram operation | Supporting SDs |
|---|---|
| `loadAdditionalRequests(assignmentId)` | SD-5A |
| `checkAreaPermission(assignmentId)` | SD-5A |
| `loadEquipmentReview(requisitionId, assignmentId)` | SD-5A |
| `validateRequestId(requisitionId)` | SD-5A |
| `reloadOriginalRequestFor6A(requisitionId)` | SD-5A |
| `reloadOriginalRequestFromDatabase(requisitionId)` | SD-5A |

Original calls and receiver evidence:

- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-32`: `loadAdditionalRequests(assignmentId)` → EquipmentReview Controller.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-34`: `checkAreaPermission(assignmentId)` → EquipmentReview Controller.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-52`: `loadEquipmentReview(requisitionId, assignmentId)` → EquipmentReview Controller.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-54`: `validateRequestId(requisitionId)` → EquipmentReview Controller.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-92`: `reloadOriginalRequestFor6A(requisitionId)` → EquipmentReview Controller.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-94`: `reloadOriginalRequestFromDatabase(requisitionId)` → EquipmentReview Controller.
## EquipmentReview Repository

Role: Repository. Sources: SD-5A.


| Diagram operation | Supporting SDs |
|---|---|
| `findAdditionalRequests(assignmentId)` | SD-5A |
| `findOriginalRequest(requisitionId, assignmentId)` | SD-5A |
| `findTORReference(requisitionId)` | SD-5A |
| `findOtherActiveRequests(requisitionId)` | SD-5A |

Original calls and receiver evidence:

- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-38`: `findAdditionalRequests(assignmentId)` → EquipmentReview Repository.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-58`: `findOriginalRequest(requisitionId, assignmentId)` → EquipmentReview Repository.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-66`: `findTORReference(requisitionId)` → EquipmentReview Repository.
- SD-5A, edge `aTN_IOUyUpbm8zimzSvL-74`: `findOtherActiveRequests(requisitionId)` → EquipmentReview Repository.
## EquipmentDecision Controller

Role: Control. Sources: SD-6A.


| Diagram operation | Supporting SDs |
|---|---|
| `loadOriginalRequestFrom5A(requisitionId, assignmentId)` | SD-6A |
| `loadOriginalProjectAndEquipmentData(requisitionId)` | SD-6A |
| `reviewEquipmentDecision(decisionData)` | SD-6A |
| `validateDecisionData(decisionData)` | SD-6A |
| `checkOtherRequestsCoveringItems(requisitionId)` | SD-6A |
| `confirmEquipmentDecision(decisionData)` | SD-6A |
| `recheckLatestRequestData(requisitionId)` | SD-6A |
| `notifySupervisor(…)` | SD-6A |

Original calls and receiver evidence:

- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-40`: `loadOriginalRequestFrom5A(requisitionId, assignmentId)` → EquipmentDecision Controller.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-50`: `loadOriginalProjectAndEquipmentData(requisitionId)` → EquipmentDecision Controller.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-72`: `reviewEquipmentDecision(decisionData)` → EquipmentDecision Controller.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-74`: `validateDecisionData(decisionData)` → EquipmentDecision Controller.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-94`: `checkOtherRequestsCoveringItems(requisitionId)` → EquipmentDecision Controller.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-106`: `confirmEquipmentDecision(decisionData)` → EquipmentDecision Controller.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-108`: `recheckLatestRequestData(requisitionId)` → EquipmentDecision Controller.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-120`: `notifySupervisor("มีคำขอจัดซื้อเพิ่มเติม [requisition_no เดิม] รอพิจารณา")` → EquipmentDecision Controller.
## EquipmentDecision Repository

Role: Repository. Sources: SD-6A.


| Diagram operation | Supporting SDs |
|---|---|
| `findOriginalRequestAndRecipient(requisitionId, assignmentId)` | SD-6A |
| `checkOriginalRequestAndArea(requisitionId, assignmentId)` | SD-6A |
| `checkOriginalItems(requisitionId)` | SD-6A |
| `forwardOriginalRequest(requisitionId, assignmentId, purchaseReason)` | SD-6A |
| `rejectOriginalRequest(requisitionId, assignmentId, assistantUserId)` | SD-6A |

Original calls and receiver evidence:

- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-42`: `findOriginalRequestAndRecipient(requisitionId, assignmentId)` → EquipmentDecision Repository.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-78`: `checkOriginalRequestAndArea(requisitionId, assignmentId)` → EquipmentDecision Repository.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-86`: `checkOriginalItems(requisitionId)` → EquipmentDecision Repository.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-112`: `forwardOriginalRequest(requisitionId, assignmentId, purchaseReason)` → EquipmentDecision Repository.
- SD-6A, edge `vDi9Eg9ZTxjohCE4Oyga-122`: `rejectOriginalRequest(requisitionId, assignmentId, assistantUserId)` → EquipmentDecision Repository.
## AdditionalPurchase Controller

Role: Control. Sources: SD-7A.


| Diagram operation | Supporting SDs |
|---|---|
| `loadEligibleAdditionalRequests(assignmentId)` | SD-7A |
| `checkAreaPermission(assignmentId)` | SD-7A |
| `loadAdditionalPurchaseForm(requisitionId)` | SD-7A |
| `reviewAdditionalPurchase(requisitionId, purchaseData)` | SD-7A |
| `validatePurchaseData(purchaseData)` | SD-7A |
| `checkRequestEligibilityAndAreaPermission(requisitionId)` | SD-7A |
| `calculateReviewTotals(purchaseData)` | SD-7A |
| `confirmAdditionalPurchase(requisitionId, purchaseData)` | SD-7A |
| `recheckLatestDataAndPreventDuplicateSubmission(purchaseData)` | SD-7A |
| `commitTransaction()` | SD-7A |
| `notifySupervisorOfShortage(requisition_no, nextStep = 2S)` | SD-7A |

Original calls and receiver evidence:

- SD-7A, edge `BdMNqnIujAcWbsTTLP88-60`: `loadEligibleAdditionalRequests(assignmentId)` → AdditionalPurchase Controller.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-62`: `checkAreaPermission(assignmentId)` → AdditionalPurchase Controller.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-80`: `loadAdditionalPurchaseForm(requisitionId)` → AdditionalPurchase Controller.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-102`: `reviewAdditionalPurchase(requisitionId, purchaseData)` → AdditionalPurchase Controller.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-104`: `validatePurchaseData(purchaseData)` → AdditionalPurchase Controller.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-108`: `checkRequestEligibilityAndAreaPermission(requisitionId)` → AdditionalPurchase Controller.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-120`: `calculateReviewTotals(purchaseData)` → AdditionalPurchase Controller.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-132`: `confirmAdditionalPurchase(requisitionId, purchaseData)` → AdditionalPurchase Controller.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-150`: `recheckLatestDataAndPreventDuplicateSubmission(purchaseData)` → AdditionalPurchase Controller.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-178`: `commitTransaction()` → AdditionalPurchase Controller.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-182`: `notifySupervisorOfShortage(requisition_no, nextStep = 2S)` → AdditionalPurchase Controller.
## AdditionalPurchase Repository

Role: Repository. Sources: SD-7A.


| Diagram operation | Supporting SDs |
|---|---|
| `findEligibleAdditionalRequests(assignmentId)` | SD-7A |
| `findRemainingItems(requisitionId)` | SD-7A |
| `checkItemAndLatestRemaining(itemId, requisitionId)` | SD-7A |
| `lockOriginalRequest(requisitionId)` | SD-7A |
| `lockRequestItems(requisitionId)` | SD-7A |
| `updateCumulativePurchase(itemId, requisitionId, quantity, unitPrice)` | SD-7A |
| `createActualExpense(requisitionId, assistantUserId, currentTotal, receipt_photo_url)` | SD-7A |
| `updateRequestStatusFromSavedQuantities(requisitionId)` | SD-7A |

Original calls and receiver evidence:

- SD-7A, edge `BdMNqnIujAcWbsTTLP88-66`: `findEligibleAdditionalRequests(assignmentId)` → AdditionalPurchase Repository.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-82`: `findRemainingItems(requisitionId)` → AdditionalPurchase Repository.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-112`: `checkItemAndLatestRemaining(itemId, requisitionId)` → AdditionalPurchase Repository.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-134`: `lockOriginalRequest(requisitionId)` → AdditionalPurchase Repository.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-142`: `lockRequestItems(requisitionId)` → AdditionalPurchase Repository.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-154`: `updateCumulativePurchase(itemId, requisitionId, quantity, unitPrice)` → AdditionalPurchase Repository.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-162`: `createActualExpense(requisitionId, assistantUserId, currentTotal, receipt_photo_url)` → AdditionalPurchase Repository.
- SD-7A, edge `BdMNqnIujAcWbsTTLP88-170`: `updateRequestStatusFromSavedQuantities(requisitionId)` → AdditionalPurchase Repository.
## EquipmentDelivery Controller

Role: Control. Sources: SD-8A.


| Diagram operation | Supporting SDs |
|---|---|
| `loadCompletedAdditionalRequests(assignmentId)` | SD-8A |
| `checkAreaPermission(assignmentId)` | SD-8A |
| `loadDeliveryForm(requisitionId)` | SD-8A |
| `reviewDelivery(requisitionId, deliveryData)` | SD-8A |
| `validateDeliveryFormat(deliveryData)` | SD-8A |
| `checkOriginalRecipientAreaAndFullQuantities(deliveryData)` | SD-8A |
| `confirmDelivery(requisitionId, deliveryData)` | SD-8A |
| `recheckLatestDeliveryAndPreventDuplicateEvidence(deliveryData)` | SD-8A |

Original calls and receiver evidence:

- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-44`: `loadCompletedAdditionalRequests(assignmentId)` → EquipmentDelivery Controller.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-46`: `checkAreaPermission(assignmentId)` → EquipmentDelivery Controller.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-64`: `loadDeliveryForm(requisitionId)` → EquipmentDelivery Controller.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-86`: `reviewDelivery(requisitionId, deliveryData)` → EquipmentDelivery Controller.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-88`: `validateDeliveryFormat(deliveryData)` → EquipmentDelivery Controller.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-92`: `checkOriginalRecipientAreaAndFullQuantities(deliveryData)` → EquipmentDelivery Controller.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-112`: `confirmDelivery(requisitionId, deliveryData)` → EquipmentDelivery Controller.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-114`: `recheckLatestDeliveryAndPreventDuplicateEvidence(deliveryData)` → EquipmentDelivery Controller.
## EquipmentDelivery Repository

Role: Repository. Sources: SD-8A.


| Diagram operation | Supporting SDs |
|---|---|
| `findCompletedAdditionalRequests(assignmentId)` | SD-8A |
| `findItemsAndOriginalRecipient(requisitionId)` | SD-8A |
| `checkScheduleRecipientAndArea(scheduleId, requisitionId)` | SD-8A |
| `saveDeliveryEvidence(scheduleId, requisitionId, description, photo_url)` | SD-8A |
| `findDeliveryNotificationData(requisitionId)` | SD-8A |

Original calls and receiver evidence:

- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-50`: `findCompletedAdditionalRequests(assignmentId)` → EquipmentDelivery Repository.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-66`: `findItemsAndOriginalRecipient(requisitionId)` → EquipmentDelivery Repository.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-96`: `checkScheduleRecipientAndArea(scheduleId, requisitionId)` → EquipmentDelivery Repository.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-118`: `saveDeliveryEvidence(scheduleId, requisitionId, description, photo_url)` → EquipmentDelivery Repository.
- SD-8A, edge `nSBGD8L_2IM09fm7yx4Z-126`: `findDeliveryNotificationData(requisitionId)` → EquipmentDelivery Repository.
## WorkContinuation Controller

Role: Control. Sources: SD-9A.


| Diagram operation | Supporting SDs |
|---|---|
| `loadContractAndArea(assignmentId)` | SD-9A |
| `validateAssignmentId(assignmentId)` | SD-9A |
| `checkAreaAccessPermission(assignmentId)` | SD-9A |
| `checkContractForContinuation(contractData)` | SD-9A |
| `evaluateContractStatusAndEndDate(contractData)` | SD-9A |

Original calls and receiver evidence:

- SD-9A, edge `ZDx1L2W36KSnwggz6YFt-35`: `loadContractAndArea(assignmentId)` → WorkContinuation Controller.
- SD-9A, edge `ZDx1L2W36KSnwggz6YFt-37`: `validateAssignmentId(assignmentId)` → WorkContinuation Controller.
- SD-9A, edge `ZDx1L2W36KSnwggz6YFt-41`: `checkAreaAccessPermission(assignmentId)` → WorkContinuation Controller.
- SD-9A, edge `ZDx1L2W36KSnwggz6YFt-61`: `checkContractForContinuation(contractData)` → WorkContinuation Controller.
- SD-9A, edge `ZDx1L2W36KSnwggz6YFt-63`: `evaluateContractStatusAndEndDate(contractData)` → WorkContinuation Controller.
## Contract Repository

Role: Repository. Sources: SD-9A.


| Diagram operation | Supporting SDs |
|---|---|
| `findContractAndArea(assignmentId)` | SD-9A |

Original calls and receiver evidence:

- SD-9A, edge `ZDx1L2W36KSnwggz6YFt-45`: `findContractAndArea(assignmentId)` → Contract Repository.
## ContractForm Controller

Role: Control. Sources: SD-1S.


| Diagram operation | Supporting SDs |
|---|---|
| `openContractForm()` | SD-1S |
| `submitContractInfo(contractData)` | SD-1S |
| `validateContractFormat(contractData)` | SD-1S |
| `submitScopeData(scopeData)` | SD-1S |
| `validateScopeFormat(scopeData)` | SD-1S |

Original calls and receiver evidence:

- SD-1S, edge `H62lu0E5nTIolGZWRA4C-254`: `openContractForm()` → ContractForm Controller.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-264`: `submitContractInfo(contractData)` → ContractForm Controller.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-266`: `validateContractFormat(contractData)` → ContractForm Controller.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-284`: `submitScopeData(scopeData)` → ContractForm Controller.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-286`: `validateScopeFormat(scopeData)` → ContractForm Controller.
## ConfirmContract Controller

Role: Control. Sources: SD-1S.


| Diagram operation | Supporting SDs |
|---|---|
| `confirmContract(contractData, scopeData)` | SD-1S |

Original calls and receiver evidence:

- SD-1S, edge `H62lu0E5nTIolGZWRA4C-312`: `confirmContract(contractData, scopeData)` → ConfirmContract Controller.
## TOR Repository

Role: Repository. Sources: SD-1S.


| Diagram operation | Supporting SDs |
|---|---|
| `checkDuplicateContractNo(contractNo)` | SD-1S |
| `validateLocation(locationId)` | SD-1S |
| `checkDuplicateLocationName( locationName)` | SD-1S |
| `createContract( contractData)` | SD-1S |
| `createLocation( locationName, address)` | SD-1S |
| `createTORLocationAssignment( torId, locationId, requiredWorkers)` | SD-1S |
| `createInitialRequisition( assignmentId, userId, requisitionNo)` | SD-1S |
| `findEquipmentByName( equipmentName)` | SD-1S |
| `createEquipment( equipmentName)` | SD-1S |
| `createRequisitionItem( requisitionId, equipmentId, remark, requiredQty)` | SD-1S |

Original calls and receiver evidence:

- SD-1S, edge `H62lu0E5nTIolGZWRA4C-268`: `checkDuplicateContractNo(contractNo)` → TOR Repository.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-288`: `validateLocation(locationId)` → TOR Repository.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-296`: `checkDuplicateLocationName( locationName)` → TOR Repository.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-314`: `createContract( contractData)` → TOR Repository.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-322`: `createLocation( locationName, address)` → TOR Repository.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-330`: `createTORLocationAssignment( torId, locationId, requiredWorkers)` → TOR Repository.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-338`: `createInitialRequisition( assignmentId, userId, requisitionNo)` → TOR Repository.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-346`: `findEquipmentByName( equipmentName)` → TOR Repository.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-354`: `createEquipment( equipmentName)` → TOR Repository.
- SD-1S, edge `H62lu0E5nTIolGZWRA4C-362`: `createRequisitionItem( requisitionId, equipmentId, remark, requiredQty)` → TOR Repository.
## TransferForm Controller

Role: Control. Sources: SD-2S.


| Diagram operation | Supporting SDs |
|---|---|
| `getPendingFundingRequests()` | SD-2S |
| `loadTransferForm(requisitionId)` | SD-2S |
| `validateTransferEvidence(transferData)` | SD-2S |
| `validateTransferFormat(transferData)` | SD-2S |

Original calls and receiver evidence:

- SD-2S, edge `ts-TrVpdrIT73yhabX2p-39`: `getPendingFundingRequests()` → TransferForm Controller.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-55`: `loadTransferForm(requisitionId)` → TransferForm Controller.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-73`: `validateTransferEvidence(transferData)` → TransferForm Controller.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-75`: `validateTransferFormat(transferData)` → TransferForm Controller.
## ConfirmTransfer Controller

Role: Control. Sources: SD-2S.


| Diagram operation | Supporting SDs |
|---|---|
| `confirmFundTransfer(requisitionId, transferData)` | SD-2S |

Original calls and receiver evidence:

- SD-2S, edge `ts-TrVpdrIT73yhabX2p-95`: `confirmFundTransfer(requisitionId, transferData)` → ConfirmTransfer Controller.
## Funding Repository

Role: Repository. Sources: SD-2S.


| Diagram operation | Supporting SDs |
|---|---|
| `findPendingFundingRequests()` | SD-2S |
| `findMissingEquipment(requisitionId)` | SD-2S |
| `checkDuplicateTransferRef(transfer_ref_no)` | SD-2S |
| `createFundTransferExpense(requisitionId, userId, transferData)` | SD-2S |
| `updateRequisitionStatus(…)` | SD-2S |

Original calls and receiver evidence:

- SD-2S, edge `ts-TrVpdrIT73yhabX2p-41`: `findPendingFundingRequests()` → Funding Repository.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-57`: `findMissingEquipment(requisitionId)` → Funding Repository.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-79`: `checkDuplicateTransferRef(transfer_ref_no)` → Funding Repository.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-97`: `createFundTransferExpense(requisitionId, userId, transferData)` → Funding Repository.
- SD-2S, edge `ts-TrVpdrIT73yhabX2p-105`: `updateRequisitionStatus(requisitionId, 'pending_procurement')` → Funding Repository.
## ReviewRequisition Controller

Role: Control. Sources: SD-3.1S, SD-3S.


| Diagram operation | Supporting SDs |
|---|---|
| `getPendingRequisitions()` | SD-3.1S, SD-3S |
| `loadRequisitionReview(requisitionId)` | SD-3.1S, SD-3S |
| `processReviewDecision(requisitionId, decision, rejectionReason)` | SD-3.1S, SD-3S |
| `validateReviewFormat(decision, rejectionReason)` | SD-3.1S, SD-3S |

Original calls and receiver evidence:

- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-39`: `getPendingRequisitions()` → ReviewRequisition Controller.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-44`: `getPendingRequisitions()` → ReviewRequisition Controller.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-45`: `loadRequisitionReview(requisitionId)` → ReviewRequisition Controller.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-56`: `loadRequisitionReview(requisitionId)` → ReviewRequisition Controller.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-58`: `processReviewDecision(requisitionId, decision, rejectionReason)` → ReviewRequisition Controller.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-80`: `processReviewDecision(requisitionId, decision, rejectionReason)` → ReviewRequisition Controller.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-59`: `validateReviewFormat(decision, rejectionReason)` → ReviewRequisition Controller.
- SD-3S, edge `imd1lA4TALJ79PT9acw6-82`: `validateReviewFormat(decision, rejectionReason)` → ReviewRequisition Controller.
## Background Job

Role: Service. Sources: SD-2W, SD-3.1S, SD-3S.


| Diagram operation | Supporting SDs |
|---|---|
| `enqueueApprovalNotification(ผู้สร้างคำขอ)` | SD-3.1S |
| `enqueueRejectionNotification(ผู้สร้างคำขอ, rejectionReason)` | SD-3.1S |
| `sendLeaveNotification(requestId)` | SD-2W |

Original calls and receiver evidence:

- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-69`: `enqueueApprovalNotification(ผู้สร้างคำขอ)` → Background Job.
- SD-3.1S, edge `zjOQ5TT1hnlVahdh1K2d-75`: `enqueueRejectionNotification(ผู้สร้างคำขอ, rejectionReason)` → Background Job.
- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-79`: `sendLeaveNotification(requestId)` → Background Job.
## Automatic Job

Role: Service. Sources: SD-4S.


No received operations are invented for this class.
## Attendance Controller

Role: Control. Sources: SD-4S.


| Diagram operation | Supporting SDs |
|---|---|
| `processDailyAttendance(work_date)` | SD-4S |
| `compareCheckInWithShiftStart(checkIn, shiftStartTime)` | SD-4S |

Original calls and receiver evidence:

- SD-4S, edge `Na5WVMECRJByJS5ewn7f-59`: `processDailyAttendance(work_date)` → Attendance Controller.
- SD-4S, edge `Na5WVMECRJByJS5ewn7f-97`: `compareCheckInWithShiftStart(checkIn, shiftStartTime)` → Attendance Controller.
## Attendance Repository

Role: Repository. Sources: SD-3W, SD-4S.


| Diagram operation | Supporting SDs |
|---|---|
| `findScheduledWorkers(work_date)` | SD-4S |
| `findCheckIn(scheduleId, workerId)` | SD-4S |
| `countApprovedAdvanceLeave(userId, work_date)` | SD-4S |
| `updateAttendanceStatus(…)` | SD-4S |
| `findTodayAssignedShift(userId)` | SD-3W |
| `checkDuplicateCheckIn(userId, scheduleId)` | SD-3W |
| `upsertSubstituteAttendance(scheduleId, workerId)` | SD-3W |

Original calls and receiver evidence:

- SD-4S, edge `Na5WVMECRJByJS5ewn7f-61`: `findScheduledWorkers(work_date)` → Attendance Repository.
- SD-4S, edge `Na5WVMECRJByJS5ewn7f-65`: `findCheckIn(scheduleId, workerId)` → Attendance Repository.
- SD-4S, edge `Na5WVMECRJByJS5ewn7f-73`: `countApprovedAdvanceLeave(userId, work_date)` → Attendance Repository.
- SD-4S, edge `Na5WVMECRJByJS5ewn7f-81`: `updateAttendanceStatus(attendanceId, 'absent')` → Attendance Repository.
- SD-4S, edge `Na5WVMECRJByJS5ewn7f-109`: `updateAttendanceStatus(attendanceId, 'late')` → Attendance Repository.
- SD-4S, edge `Na5WVMECRJByJS5ewn7f-89`: `updateAttendanceStatus(attendanceId, 'leave')` → Attendance Repository.
- SD-4S, edge `Na5WVMECRJByJS5ewn7f-101`: `updateAttendanceStatus(attendanceId, 'on_time')` → Attendance Repository.
- SD-3W, edge `E7oHi6kEsZcUCtm5OUxS-37`: `findTodayAssignedShift(userId)` → Attendance Repository.
- SD-3W, edge `E7oHi6kEsZcUCtm5OUxS-63`: `checkDuplicateCheckIn(userId, scheduleId)` → Attendance Repository.
- SD-3W, edge `E7oHi6kEsZcUCtm5OUxS-71`: `upsertSubstituteAttendance(scheduleId, workerId)` → Attendance Repository.
## PayrollSummary Controller (5S)

Role: Control. Sources: SD-4S.


| Diagram operation | Supporting SDs |
|---|---|
| `startPeriodSummary5S()` | SD-4S |

Original calls and receiver evidence:

- SD-4S, edge `Na5WVMECRJByJS5ewn7f-119`: `startPeriodSummary5S()` → PayrollSummary Controller (5S).
## Payroll Controller

Role: Control. Sources: SD-5S.


| Diagram operation | Supporting SDs |
|---|---|
| `processPayroll(period_start, period_end)` | SD-5S |
| `calculateBaseWage(workDays, dailyWage)` | SD-5S |
| `calculateNetPay(baseWage, totalDeduction)` | SD-5S |

Original calls and receiver evidence:

- SD-5S, edge `g_NC1AvDctcExgY8cTaH-104`: `processPayroll(period_start, period_end)` → Payroll Controller.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-118`: `calculateBaseWage(workDays, dailyWage)` → Payroll Controller.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-146`: `calculateNetPay(baseWage, totalDeduction)` → Payroll Controller.
## Payroll Repository

Role: Repository. Sources: SD-5S.


| Diagram operation | Supporting SDs |
|---|---|
| `findWorkersInPeriod(period_start, period_end)` | SD-5S |
| `getBaseWageData(workerId, period_start, period_end)` | SD-5S |
| `findDeductionItems(workerId, period_start, period_end)` | SD-5S |
| `createDeduction(workerId, penaltyAmount, penaltyReason)` | SD-5S |
| `getTotalDeduction(workerId, period_start, period_end)` | SD-5S |
| `createPayroll(workerId, period_start, period_end, baseWage, totalDeduction, netPay)` | SD-5S |

Original calls and receiver evidence:

- SD-5S, edge `g_NC1AvDctcExgY8cTaH-106`: `findWorkersInPeriod(period_start, period_end)` → Payroll Repository.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-110`: `getBaseWageData(workerId, period_start, period_end)` → Payroll Repository.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-122`: `findDeductionItems(workerId, period_start, period_end)` → Payroll Repository.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-130`: `createDeduction(workerId, penaltyAmount, penaltyReason)` → Payroll Repository.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-138`: `getTotalDeduction(workerId, period_start, period_end)` → Payroll Repository.
- SD-5S, edge `g_NC1AvDctcExgY8cTaH-150`: `createPayroll(workerId, period_start, period_end, baseWage, totalDeduction, netPay)` → Payroll Repository.
## FinancialReport Controller

Role: Control. Sources: SD-6S.


| Diagram operation | Supporting SDs |
|---|---|
| `calculateFinancialReport(torId, dateStart, dateEnd)` | SD-6S |
| `calculateNetProfit(totalRevenue, totalLaborCost, totalMaterialCost)` | SD-6S |
| `exportFinancialReportPDF(reportData)` | SD-6S |
| `generateReportPDF(reportData)` | SD-6S |

Original calls and receiver evidence:

- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-161`: `calculateFinancialReport(torId, dateStart, dateEnd)` → FinancialReport Controller.
- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-187`: `calculateNetProfit(totalRevenue, totalLaborCost, totalMaterialCost)` → FinancialReport Controller.
- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-209`: `exportFinancialReportPDF(reportData)` → FinancialReport Controller.
- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-211`: `generateReportPDF(reportData)` → FinancialReport Controller.
## Financial Repository

Role: Repository. Sources: SD-6S.


| Diagram operation | Supporting SDs |
|---|---|
| `getTotalRevenue(torId )` | SD-6S |
| `getTotalLaborCost(dateStart, dateEnd)` | SD-6S |
| `getTotalMaterialCost(torId, dateStart, dateEnd)` | SD-6S |

Original calls and receiver evidence:

- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-163`: `getTotalRevenue(torId )` → Financial Repository.
- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-171`: `getTotalLaborCost(dateStart, dateEnd)` → Financial Repository.
- SD-6S, edge `bIPx0u1nT4T6iZqJpZoo-179`: `getTotalMaterialCost(torId, dateStart, dateEnd)` → Financial Repository.
## LINE Mini App (LIFF)

Role: Boundary. Sources: SD-1W, SD-2W, SD-3W, SD-4W, SD-5W, SD-7W.


| Diagram operation | Supporting SDs |
|---|---|
| `tapMyScheduleRichMenu()` | SD-1W |
| `openLIFFWindow()` | SD-1W |
| `showScrollableScheduleCardsOrList()` | SD-1W |
| `optionallyTapScheduleCard(scheduleId, nextAction)` | SD-1W |
| `openLeaveScreen2W(scheduleId)` | SD-1W |
| `openPreWorkConfirmationScreen3W(scheduleId)` | SD-1W |
| `tapLeaveButton()` | SD-2W |
| `showLeaveDateAndReasonForm()` | SD-2W |
| `confirmLeaveRequest(leave_date, reason)` | SD-2W |
| `redirectToLIFFHome()` | SD-2W, SD-3W, SD-4W, SD-5W |
| `showToast(…)` | SD-2W, SD-3W, SD-4W, SD-5W |
| `tapConfirmSubstituteWork()` | SD-3W |
| `showAssignedLocationAndStartTime()` | SD-3W |
| `tapConfirmAndCheckIn()` | SD-3W |
| `tapWorkReportAndCheckOut()` | SD-4W |
| `showDescriptionAndPhotoForm()` | SD-4W |
| `confirmWorkReport(description, photo)` | SD-4W |
| `tapInsufficientEquipment()` | SD-5W |
| `showEquipmentQuantityAndReasonForm()` | SD-5W |
| `sendEquipmentRequest(equipmentId, quantity, reason)` | SD-5W |
| `tapMyPayslip()` | SD-7W |
| `showPeriodMonthSearchForm()` | SD-7W |
| `searchPayslip(period_month)` | SD-7W |
| `showDigitalPayslipCard(period_start, period_end, total_wage, total_deduction, deductionDetails, net_wage)` | SD-7W |
| `closePayslipWindow()` | SD-7W |
| `returnToLIFFHome()` | SD-7W |

Original calls and receiver evidence:

- SD-1W, edge `CF87uobWHUu3A-gUhq7B-23`: `tapMyScheduleRichMenu()` → LINE Mini App (LIFF).
- SD-1W, edge `CF87uobWHUu3A-gUhq7B-25`: `openLIFFWindow()` → LINE Mini App (LIFF).
- SD-1W, edge `CF87uobWHUu3A-gUhq7B-43`: `showScrollableScheduleCardsOrList()` → LINE Mini App (LIFF).
- SD-1W, edge `CF87uobWHUu3A-gUhq7B-47`: `optionallyTapScheduleCard(scheduleId, nextAction)` → LINE Mini App (LIFF).
- SD-1W, edge `CF87uobWHUu3A-gUhq7B-49`: `openLeaveScreen2W(scheduleId)` → LINE Mini App (LIFF).
- SD-1W, edge `CF87uobWHUu3A-gUhq7B-51`: `openPreWorkConfirmationScreen3W(scheduleId)` → LINE Mini App (LIFF).
- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-29`: `tapLeaveButton()` → LINE Mini App (LIFF).
- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-47`: `showLeaveDateAndReasonForm()` → LINE Mini App (LIFF).
- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-51`: `confirmLeaveRequest(leave_date, reason)` → LINE Mini App (LIFF).
- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-85`: `redirectToLIFFHome()` → LINE Mini App (LIFF).
- SD-3W, edge `E7oHi6kEsZcUCtm5OUxS-81`: `redirectToLIFFHome()` → LINE Mini App (LIFF).
- SD-4W, edge `4W7N71NcVjHcfgNTzLRr-77`: `redirectToLIFFHome()` → LINE Mini App (LIFF).
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-101`: `redirectToLIFFHome()` → LINE Mini App (LIFF).
- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-87`: `showToast("ส่งคำขอลาเรียบร้อยแล้ว กรุณารอการอนุมัติ")` → LINE Mini App (LIFF).
- SD-3W, edge `E7oHi6kEsZcUCtm5OUxS-83`: `showToast("ยืนยันรับหน้าที่แทนและเช็คอินสำเร็จ เวลา [HH:mm]")` → LINE Mini App (LIFF).
- SD-4W, edge `4W7N71NcVjHcfgNTzLRr-79`: `showToast("ส่งผลการปฏิบัติงานและเช็คเอาท์สำเร็จ เวลา [HH:mm]")` → LINE Mini App (LIFF).
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-103`: `showToast("ส่งคำขออุปกรณ์เพิ่มเติมสำเร็จ กรุณารอผู้ดูแลงานตรวจสอบ")` → LINE Mini App (LIFF).
- SD-3W, edge `E7oHi6kEsZcUCtm5OUxS-29`: `tapConfirmSubstituteWork()` → LINE Mini App (LIFF).
- SD-3W, edge `E7oHi6kEsZcUCtm5OUxS-47`: `showAssignedLocationAndStartTime()` → LINE Mini App (LIFF).
- SD-3W, edge `E7oHi6kEsZcUCtm5OUxS-51`: `tapConfirmAndCheckIn()` → LINE Mini App (LIFF).
- SD-4W, edge `4W7N71NcVjHcfgNTzLRr-25`: `tapWorkReportAndCheckOut()` → LINE Mini App (LIFF).
- SD-4W, edge `4W7N71NcVjHcfgNTzLRr-43`: `showDescriptionAndPhotoForm()` → LINE Mini App (LIFF).
- SD-4W, edge `4W7N71NcVjHcfgNTzLRr-47`: `confirmWorkReport(description, photo)` → LINE Mini App (LIFF).
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-31`: `tapInsufficientEquipment()` → LINE Mini App (LIFF).
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-57`: `showEquipmentQuantityAndReasonForm()` → LINE Mini App (LIFF).
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-61`: `sendEquipmentRequest(equipmentId, quantity, reason)` → LINE Mini App (LIFF).
- SD-7W, edge `uUxBqSW2nSQnUcCSPHx_-26`: `tapMyPayslip()` → LINE Mini App (LIFF).
- SD-7W, edge `uUxBqSW2nSQnUcCSPHx_-36`: `showPeriodMonthSearchForm()` → LINE Mini App (LIFF).
- SD-7W, edge `uUxBqSW2nSQnUcCSPHx_-40`: `searchPayslip(period_month)` → LINE Mini App (LIFF).
- SD-7W, edge `uUxBqSW2nSQnUcCSPHx_-66`: `showDigitalPayslipCard(period_start, period_end, total_wage, total_deduction, deductionDetails, net_wage)` → LINE Mini App (LIFF).
- SD-7W, edge `uUxBqSW2nSQnUcCSPHx_-70`: `closePayslipWindow()` → LINE Mini App (LIFF).
- SD-7W, edge `uUxBqSW2nSQnUcCSPHx_-72`: `returnToLIFFHome()` → LINE Mini App (LIFF).
## LeaveRequest Controller

Role: Control. Sources: SD-2W.


| Diagram operation | Supporting SDs |
|---|---|
| `loadLeaveDateOptions(lineIdentity)` | SD-2W |
| `verifyLinkedLINEIdentity(lineIdentity)` | SD-2W |
| `submitLeaveRequest(userId, leave_date, reason)` | SD-2W |
| `validateLeaveDateAndReason(leave_date, reason)` | SD-2W |
| `calculateAdvanceNotice(currentTimestamp, leave_date)` | SD-2W |

Original calls and receiver evidence:

- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-31`: `loadLeaveDateOptions(lineIdentity)` → LeaveRequest Controller.
- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-33`: `verifyLinkedLINEIdentity(lineIdentity)` → LeaveRequest Controller.
- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-53`: `submitLeaveRequest(userId, leave_date, reason)` → LeaveRequest Controller.
- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-55`: `validateLeaveDateAndReason(leave_date, reason)` → LeaveRequest Controller.
- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-67`: `calculateAdvanceNotice(currentTimestamp, leave_date)` → LeaveRequest Controller.
## LeaveRequest Repository

Role: Repository. Sources: SD-2W.


| Diagram operation | Supporting SDs |
|---|---|
| `findScheduledDates(userId)` | SD-2W |
| `checkDuplicateLeave(userId, leave_date)` | SD-2W |
| `createPendingLeave(userId, leave_date, reason, is_advance_notice)` | SD-2W |

Original calls and receiver evidence:

- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-37`: `findScheduledDates(userId)` → LeaveRequest Repository.
- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-59`: `checkDuplicateLeave(userId, leave_date)` → LeaveRequest Repository.
- SD-2W, edge `j2ubfQiAtmlUmQNyKNjS-71`: `createPendingLeave(userId, leave_date, reason, is_advance_notice)` → LeaveRequest Repository.
## SubstituteCheckIn Controller

Role: Control. Sources: SD-3W.


| Diagram operation | Supporting SDs |
|---|---|
| `loadTodayAssignedShift(lineIdentity)` | SD-3W |
| `verifyLinkedLINEIdentity(lineIdentity)` | SD-3W |
| `confirmSubstituteCheckIn(lineUserId, scheduleId)` | SD-3W |
| `validateLINEUserIdAndScheduleUUID(lineUserId, scheduleId)` | SD-3W |
| `checkTodayShiftMatchesWorker(todayAssignedShift)` | SD-3W |

Original calls and receiver evidence:

- SD-3W, edge `E7oHi6kEsZcUCtm5OUxS-31`: `loadTodayAssignedShift(lineIdentity)` → SubstituteCheckIn Controller.
- SD-3W, edge `E7oHi6kEsZcUCtm5OUxS-33`: `verifyLinkedLINEIdentity(lineIdentity)` → SubstituteCheckIn Controller.
- SD-3W, edge `E7oHi6kEsZcUCtm5OUxS-53`: `confirmSubstituteCheckIn(lineUserId, scheduleId)` → SubstituteCheckIn Controller.
- SD-3W, edge `E7oHi6kEsZcUCtm5OUxS-55`: `validateLINEUserIdAndScheduleUUID(lineUserId, scheduleId)` → SubstituteCheckIn Controller.
- SD-3W, edge `E7oHi6kEsZcUCtm5OUxS-59`: `checkTodayShiftMatchesWorker(todayAssignedShift)` → SubstituteCheckIn Controller.
## WorkReport Controller

Role: Control. Sources: SD-4W.


| Diagram operation | Supporting SDs |
|---|---|
| `loadOpenCheckIn(lineIdentity)` | SD-4W |
| `verifyLinkedLINEIdentity(lineIdentity)` | SD-4W |
| `submitWorkReport(description, photo_url)` | SD-4W |
| `validateDescriptionAndImage(description, photo_url)` | SD-4W |
| `checkExactlyOneOpenAttendance(openAttendance)` | SD-4W |

Original calls and receiver evidence:

- SD-4W, edge `4W7N71NcVjHcfgNTzLRr-27`: `loadOpenCheckIn(lineIdentity)` → WorkReport Controller.
- SD-4W, edge `4W7N71NcVjHcfgNTzLRr-29`: `verifyLinkedLINEIdentity(lineIdentity)` → WorkReport Controller.
- SD-4W, edge `4W7N71NcVjHcfgNTzLRr-49`: `submitWorkReport(description, photo_url)` → WorkReport Controller.
- SD-4W, edge `4W7N71NcVjHcfgNTzLRr-51`: `validateDescriptionAndImage(description, photo_url)` → WorkReport Controller.
- SD-4W, edge `4W7N71NcVjHcfgNTzLRr-55`: `checkExactlyOneOpenAttendance(openAttendance)` → WorkReport Controller.
## WorkReport Repository

Role: Repository. Sources: SD-4W.


| Diagram operation | Supporting SDs |
|---|---|
| `findTodayOpenAttendance(userId)` | SD-4W |
| `saveWorkEvidence(workerId, scheduleId, description, photo_url)` | SD-4W |
| `updateCheckOut(attendanceId)` | SD-4W |

Original calls and receiver evidence:

- SD-4W, edge `4W7N71NcVjHcfgNTzLRr-33`: `findTodayOpenAttendance(userId)` → WorkReport Repository.
- SD-4W, edge `4W7N71NcVjHcfgNTzLRr-59`: `saveWorkEvidence(workerId, scheduleId, description, photo_url)` → WorkReport Repository.
- SD-4W, edge `4W7N71NcVjHcfgNTzLRr-67`: `updateCheckOut(attendanceId)` → WorkReport Repository.
## EquipmentRequest Controller

Role: Control. Sources: SD-5W.


| Diagram operation | Supporting SDs |
|---|---|
| `loadEquipmentRequestForm(lineIdentity)` | SD-5W |
| `verifyLinkedLINEIdentity(lineIdentity)` | SD-5W |
| `submitAdditionalRequest(equipmentId, quantity, reason)` | SD-5W |
| `validateEquipmentUUIDQuantityAndReason(equipmentId, quantity, reason)` | SD-5W |
| `checkCurrentShiftHasData(currentAreaAndRequester)` | SD-5W |

Original calls and receiver evidence:

- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-33`: `loadEquipmentRequestForm(lineIdentity)` → EquipmentRequest Controller.
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-35`: `verifyLinkedLINEIdentity(lineIdentity)` → EquipmentRequest Controller.
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-63`: `submitAdditionalRequest(equipmentId, quantity, reason)` → EquipmentRequest Controller.
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-65`: `validateEquipmentUUIDQuantityAndReason(equipmentId, quantity, reason)` → EquipmentRequest Controller.
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-69`: `checkCurrentShiftHasData(currentAreaAndRequester)` → EquipmentRequest Controller.
## EquipmentRequest Repository

Role: Repository. Sources: SD-5W.


| Diagram operation | Supporting SDs |
|---|---|
| `findCurrentWorkingArea(userId)` | SD-5W |
| `findActiveEquipment()` | SD-5W |
| `checkActiveEquipmentExists(equipmentId)` | SD-5W |
| `createAdditionalRequest(autoRequisitionNo, userId, assignmentId, reason)` | SD-5W |
| `createRequestItem(requisitionId, equipmentId, quantity)` | SD-5W |

Original calls and receiver evidence:

- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-39`: `findCurrentWorkingArea(userId)` → EquipmentRequest Repository.
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-47`: `findActiveEquipment()` → EquipmentRequest Repository.
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-73`: `checkActiveEquipmentExists(equipmentId)` → EquipmentRequest Repository.
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-81`: `createAdditionalRequest(autoRequisitionNo, userId, assignmentId, reason)` → EquipmentRequest Repository.
- SD-5W, edge `7qI0YqkbS6rHrinXjDqz-89`: `createRequestItem(requisitionId, equipmentId, quantity)` → EquipmentRequest Repository.
## EquipmentResult Controller

Role: Control. Sources: SD-6W.


No received operations are invented for this class.
## EquipmentResult Repository

Role: Repository. Sources: SD-6W.


| Diagram operation | Supporting SDs |
|---|---|
| `findEquipmentRequestResult(requisitionId)` | SD-6W |

Original calls and receiver evidence:

- SD-6W, edge `OQYhqyozvix0y_K6DQpe-19`: `findEquipmentRequestResult(requisitionId)` → EquipmentResult Repository.
## Payslip Controller

Role: Control. Sources: SD-7W.


| Diagram operation | Supporting SDs |
|---|---|
| `openPayslipSearch(lineIdentity)` | SD-7W |
| `verifyLinkedLINEIdentity(lineIdentity)` | SD-7W |
| `getPayslip(userId, period_month)` | SD-7W |
| `validatePeriodMonth(period_month)` | SD-7W |

Original calls and receiver evidence:

- SD-7W, edge `uUxBqSW2nSQnUcCSPHx_-28`: `openPayslipSearch(lineIdentity)` → Payslip Controller.
- SD-7W, edge `uUxBqSW2nSQnUcCSPHx_-30`: `verifyLinkedLINEIdentity(lineIdentity)` → Payslip Controller.
- SD-7W, edge `uUxBqSW2nSQnUcCSPHx_-42`: `getPayslip(userId, period_month)` → Payslip Controller.
- SD-7W, edge `uUxBqSW2nSQnUcCSPHx_-44`: `validatePeriodMonth(period_month)` → Payslip Controller.
## Payslip Repository

Role: Repository. Sources: SD-7W.


| Diagram operation | Supporting SDs |
|---|---|
| `findPayslipsByMonth(userId, period_month)` | SD-7W |
| `findDeductionDetails(userId, period_start, period_end)` | SD-7W |

Original calls and receiver evidence:

- SD-7W, edge `uUxBqSW2nSQnUcCSPHx_-48`: `findPayslipsByMonth(userId, period_month)` → Payslip Repository.
- SD-7W, edge `uUxBqSW2nSQnUcCSPHx_-56`: `findDeductionDetails(userId, period_start, period_end)` → Payslip Repository.
## USER

Role: Domain. Sources: ER latest(ปรับทะแยง).

| ER attribute | Identifier / derived | Source cell |
|---|---|---|
| `first_name` | ordinary | `8tqF58g8p9KXpsFFBl3G-7` |
| `user_id` | key | `8tqF58g8p9KXpsFFBl3G-5` |
| `phone_number` | ordinary | `8tqF58g8p9KXpsFFBl3G-8` |
| `last_name` | ordinary | `8tqF58g8p9KXpsFFBl3G-6` |
| `created_at` | ordinary | `8tqF58g8p9KXpsFFBl3G-21` |
| `updated_at` | ordinary | `8tqF58g8p9KXpsFFBl3G-23` |
| `is_active` | ordinary | `8tqF58g8p9KXpsFFBl3G-22` |
| `bank_account_no` | ordinary | `8tqF58g8p9KXpsFFBl3G-61` |
| `daily_wage` | ordinary | `8tqF58g8p9KXpsFFBl3G-66` |
| `bank_name` | ordinary | `8tqF58g8p9KXpsFFBl3G-68` |
| `line_id` | ordinary | `8tqF58g8p9KXpsFFBl3G-207` |
| `role` | ordinary | `8tqF58g8p9KXpsFFBl3G-341` |

No received operations are invented for this class.
## WORKER

Role: Domain. Sources: ER latest(ปรับทะแยง).

| ER attribute | Identifier / derived | Source cell |
|---|---|---|
| `is_available` | ordinary | `8tqF58g8p9KXpsFFBl3G-16` |

No received operations are invented for this class.
## PAYROLL

Role: Domain. Sources: ER latest(ปรับทะแยง).

| ER attribute | Identifier / derived | Source cell |
|---|---|---|
| `period_start` | ordinary | `pVqMR6C_7T5EkJBiJlV9-24` |
| `period_end` | ordinary | `pVqMR6C_7T5EkJBiJlV9-25` |
| `payroll_id` | key | `8tqF58g8p9KXpsFFBl3G-20` |
| `base_wage` | ordinary | `8tqF58g8p9KXpsFFBl3G-50` |
| `deduction` | ordinary | `8tqF58g8p9KXpsFFBl3G-52` |
| `is_paid` | ordinary | `8tqF58g8p9KXpsFFBl3G-55` |
| `created_at` | ordinary | `8tqF58g8p9KXpsFFBl3G-57` |
| `net_pay` | derived | `8tqF58g8p9KXpsFFBl3G-53` |
| `payroll_slip_no` | ordinary | `8tqF58g8p9KXpsFFBl3G-334` |

No received operations are invented for this class.
## CONTRACT_TOR

Role: Domain. Sources: ER latest(ปรับทะแยง).

| ER attribute | Identifier / derived | Source cell |
|---|---|---|
| `start_date` | ordinary | `8tqF58g8p9KXpsFFBl3G-35` |
| `tor_id` | key | `8tqF58g8p9KXpsFFBl3G-31` |
| `project_name` | ordinary | `8tqF58g8p9KXpsFFBl3G-33` |
| `end_date` | ordinary | `8tqF58g8p9KXpsFFBl3G-37` |
| `contract_value` | ordinary | `8tqF58g8p9KXpsFFBl3G-39` |
| `created_at` | ordinary | `8tqF58g8p9KXpsFFBl3G-41` |
| `updated_at` | ordinary | `8tqF58g8p9KXpsFFBl3G-43` |
| `partner_agency` | ordinary | `8tqF58g8p9KXpsFFBl3G-308` |
| `contract_no` | ordinary | `8tqF58g8p9KXpsFFBl3G-310` |
| `contract_file_url` | ordinary | `8tqF58g8p9KXpsFFBl3G-312` |
| `status` | ordinary | `rlVyrjYtN52kc_-kk-_f-1` |

No received operations are invented for this class.
## ATTENDANCE

Role: Domain. Sources: ER latest(ปรับทะแยง).

| ER attribute | Identifier / derived | Source cell |
|---|---|---|
| `attendance_id` | key | `8tqF58g8p9KXpsFFBl3G-72` |
| `work_date` | ordinary | `8tqF58g8p9KXpsFFBl3G-74` |
| `check_in` | ordinary | `8tqF58g8p9KXpsFFBl3G-76` |
| `check_out` | ordinary | `8tqF58g8p9KXpsFFBl3G-78` |
| `status` | ordinary | `8tqF58g8p9KXpsFFBl3G-80` |

No received operations are invented for this class.
## WORK_SCHEDULE

Role: Domain. Sources: ER latest(ปรับทะแยง).

| ER attribute | Identifier / derived | Source cell |
|---|---|---|
| `schedule_id` | key | `8tqF58g8p9KXpsFFBl3G-94` |
| `work_date` | ordinary | `8tqF58g8p9KXpsFFBl3G-96` |
| `shift_status` | ordinary | `8tqF58g8p9KXpsFFBl3G-98` |
| `shift_start_time` | ordinary | `8tqF58g8p9KXpsFFBl3G-271` |

No received operations are invented for this class.
## LEAVE_REQUEST

Role: Domain. Sources: ER latest(ปรับทะแยง).

| ER attribute | Identifier / derived | Source cell |
|---|---|---|
| `request_id` | key | `8tqF58g8p9KXpsFFBl3G-110` |
| `user_id` | ordinary | `8tqF58g8p9KXpsFFBl3G-112` |
| `leave_date` | ordinary | `8tqF58g8p9KXpsFFBl3G-114` |
| `is_advance_notice` | ordinary | `8tqF58g8p9KXpsFFBl3G-116` |
| `reason` | ordinary | `8tqF58g8p9KXpsFFBl3G-118` |
| `created_at` | ordinary | `8tqF58g8p9KXpsFFBl3G-120` |
| `status` | ordinary | `8tqF58g8p9KXpsFFBl3G-122` |
| `leave_no` | ordinary | `8tqF58g8p9KXpsFFBl3G-336` |

No received operations are invented for this class.
## TOR_LOCATION_ASSIGNMENT

Role: Domain. Sources: ER latest(ปรับทะแยง).

| ER attribute | Identifier / derived | Source cell |
|---|---|---|
| `required_workers` | ordinary | `8tqF58g8p9KXpsFFBl3G-132` |
| `assignment_id` | key | `8tqF58g8p9KXpsFFBl3G-130` |

No received operations are invented for this class.
## LOCATION

Role: Domain. Sources: ER latest(ปรับทะแยง).

| ER attribute | Identifier / derived | Source cell |
|---|---|---|
| `location_id` | key | `8tqF58g8p9KXpsFFBl3G-138` |
| `address` | ordinary | `8tqF58g8p9KXpsFFBl3G-140` |
| `latitude` | ordinary | `8tqF58g8p9KXpsFFBl3G-142` |
| `location_name` | ordinary | `8tqF58g8p9KXpsFFBl3G-144` |
| `longitude` | ordinary | `8tqF58g8p9KXpsFFBl3G-146` |

No received operations are invented for this class.
## WORK_EVIDENCE

Role: Domain. Sources: ER latest(ปรับทะแยง).

| ER attribute | Identifier / derived | Source cell |
|---|---|---|
| `evidence_id` | key | `8tqF58g8p9KXpsFFBl3G-159` |
| `photo_url` | ordinary | `8tqF58g8p9KXpsFFBl3G-161` |
| `submitted_at` | ordinary | `8tqF58g8p9KXpsFFBl3G-163` |
| `description` | ordinary | `8tqF58g8p9KXpsFFBl3G-165` |

No received operations are invented for this class.
## DEDUCTION_TRANSACTION

Role: Domain. Sources: ER latest(ปรับทะแยง).

| ER attribute | Identifier / derived | Source cell |
|---|---|---|
| `deduction_id` | key | `8tqF58g8p9KXpsFFBl3G-185` |
| `penalty_amount` | ordinary | `8tqF58g8p9KXpsFFBl3G-190` |
| `reason` | ordinary | `8tqF58g8p9KXpsFFBl3G-192` |
| `created_at` | ordinary | `8tqF58g8p9KXpsFFBl3G-194` |

No received operations are invented for this class.
## EQUIPMENT_REQUISITION

Role: Domain. Sources: ER latest(ปรับทะแยง).

| ER attribute | Identifier / derived | Source cell |
|---|---|---|
| `requisition_type` | ordinary | `qRc-ApRebziGKMpHV3Vu-1` |
| `requisition_id` | key | `8tqF58g8p9KXpsFFBl3G-211` |
| `created_at` | ordinary | `8tqF58g8p9KXpsFFBl3G-318` |
| `reason` | ordinary | `8tqF58g8p9KXpsFFBl3G-326` |
| `requisition_no` | ordinary | `8tqF58g8p9KXpsFFBl3G-330` |
| `status` | ordinary | `8tqF58g8p9KXpsFFBl3G-343` |
| `reviewed_by` | ordinary | `8tqF58g8p9KXpsFFBl3G-349` |
| `reviewed_at` | ordinary | `8tqF58g8p9KXpsFFBl3G-351` |

No received operations are invented for this class.
## REQUISITION_ITEM

Role: Domain. Sources: ER latest(ปรับทะแยง).

| ER attribute | Identifier / derived | Source cell |
|---|---|---|
| `required_qty` | ordinary | `8tqF58g8p9KXpsFFBl3G-218` |
| `remark` | ordinary | `8tqF58g8p9KXpsFFBl3G-220` |
| `actual_qty` | ordinary | `8tqF58g8p9KXpsFFBl3G-320` |
| `actual_price` | ordinary | `8tqF58g8p9KXpsFFBl3G-322` |
| `item_id` | key | `8tqF58g8p9KXpsFFBl3G-324` |
| `existing_qty` | ordinary | `8tqF58g8p9KXpsFFBl3G-338` |
| `to_buy_qty` | ordinary | `8tqF58g8p9KXpsFFBl3G-340` |

No received operations are invented for this class.
## EQUIPMENT

Role: Domain. Sources: ER latest(ปรับทะแยง).

| ER attribute | Identifier / derived | Source cell |
|---|---|---|
| `equipment_name` | ordinary | `8tqF58g8p9KXpsFFBl3G-200` |
| `equipment_id` | key | `8tqF58g8p9KXpsFFBl3G-202` |
| `is_active` | ordinary | `8tqF58g8p9KXpsFFBl3G-229` |

No received operations are invented for this class.
## EXPENSE_CLAIM

Role: Domain. Sources: ER latest(ปรับทะแยง).

| ER attribute | Identifier / derived | Source cell |
|---|---|---|
| `receipt_photo_url` | ordinary | `8tqF58g8p9KXpsFFBl3G-241` |
| `expense_id` | key | `8tqF58g8p9KXpsFFBl3G-240` |
| `total_amount` | ordinary | `8tqF58g8p9KXpsFFBl3G-243` |
| `transfer_ref_no` | ordinary | `8tqF58g8p9KXpsFFBl3G-245` |
| `created_at` | ordinary | `pVqMR6C_7T5EkJBiJlV9-14` |
| `expense_no` | ordinary | `8tqF58g8p9KXpsFFBl3G-332` |
| `expense_type` | ordinary | `qRc-ApRebziGKMpHV3Vu-16` |

No received operations are invented for this class.
## COMPANY_INVOICE

Role: Domain. Sources: ER latest(ปรับทะแยง).

| ER attribute | Identifier / derived | Source cell |
|---|---|---|
| `billing_month` | ordinary | `8tqF58g8p9KXpsFFBl3G-260` |
| `net_received` | ordinary | `8tqF58g8p9KXpsFFBl3G-265` |
| `exat_deduction_amount` | ordinary | `8tqF58g8p9KXpsFFBl3G-263` |
| `deduction_reason` | ordinary | `8tqF58g8p9KXpsFFBl3G-264` |
| `invoice_id` | key | `8tqF58g8p9KXpsFFBl3G-258` |
| `expected_amount` | ordinary | `8tqF58g8p9KXpsFFBl3G-262` |
| `status` | ordinary | `8tqF58g8p9KXpsFFBl3G-267` |
| `invoice_no` | ordinary | `8tqF58g8p9KXpsFFBl3G-328` |

No received operations are invented for this class.

## Domain association key

| ID | ER association | Max sides |
|---|---|---|
| R01 | USER — RECEIVES — PAYROLL | 1 : M |
| R02 | USER — MANAGES — CONTRACT_TOR | 1 : M |
| R03 | USER — REQUESTS — EQUIPMENT_REQUISITION | 1 : M |
| R04 | USER — REVIEWS — EQUIPMENT_REQUISITION | 1 : M |
| R05 | USER — CLAIMS — EXPENSE_CLAIM | 1 : M |
| R06 | CONTRACT_TOR — BILLS — COMPANY_INVOICE | 1 : M |
| R07 | CONTRACT_TOR — INCLUDES — TOR_LOCATION_ASSIGNMENT | 1 : M |
| R08 | LOCATION — HOSTS — TOR_LOCATION_ASSIGNMENT | 1 : M |
| R09 | TOR_LOCATION_ASSIGNMENT — REQUIRES — EQUIPMENT_REQUISITION | 1 : M |
| R10 | TOR_LOCATION_ASSIGNMENT — GENERATES — WORK_SCHEDULE | 1 : M |
| R11 | TOR_LOCATION_ASSIGNMENT — COLLECTS — WORK_EVIDENCE | 1 : M |
| R12 | EQUIPMENT_REQUISITION — CONTAINS — REQUISITION_ITEM | 1 : M |
| R13 | EQUIPMENT — IS_LISTED_IN — REQUISITION_ITEM | 1 : M |
| R14 | EQUIPMENT_REQUISITION — GENERATES — EXPENSE_CLAIM | 1 : M |
| R15 | WORKER — IS_ASSIGNED_TO — WORK_SCHEDULE | 1 : M |
| R16 | WORK_SCHEDULE — RECORDS — ATTENDANCE | 1 : 1 |
| R17 | WORKER — LOGS — ATTENDANCE | 1 : M |
| R18 | WORKER — REPLACES — ATTENDANCE | 1 : M |
| R19 | WORKER — SUBMITS — LEAVE_REQUEST | 1 : M |
| R20 | WORKER — CAPTURES — WORK_EVIDENCE | 1 : M |
| R21 | WORKER — INCURS — DEDUCTION_TRANSACTION | 1 : M |
| R22 | ATTENDANCE — TRIGGERS — DEDUCTION_TRANSACTION | 1 : M |
| R23 | PAYROLL — APPLIES — DEDUCTION_TRANSACTION | 1 : M |

U1: USER — WORKER, `(p, e)` unresolved. Four malformed ER connector attachments were interpreted using the visible intended endpoints. Original ER source was not edited.
