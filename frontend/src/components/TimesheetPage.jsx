import React from "react";
import TimesheetTable from "./TimesheetTable";
import TimeLogModal from "./TimeLogModal";
import InvoiceGenerator from "./InvoiceGenerator";
import WeekNavigator from "./WeekNavigator";
import CompanySelectionGate from "./CompanySelectionGate";
import { useTimesheet } from "../hooks/useTimesheet";
import { useActivities } from "../hooks/useActivities";
import { useTimeLogForm } from "../hooks/useTimeLogForm";
import { useCompany } from "../context/CompanyContext";

function TimesheetContent({ selectedCompanyId }) {
    const {
        days,
        rows,
        totalsPerDayMinutes,
        overallMinutes,
        rangeLabel,
        goToNextWeek,
        goToPreviousWeek,
        refresh,
    } = useTimesheet(selectedCompanyId);
    const {
        activities,
        loading: activitiesLoading,
        error: activitiesError,
    } = useActivities(selectedCompanyId);
    const { openNew, openFromCell, modalProps } = useTimeLogForm({
        activities,
        onSaved: refresh,
        defaultDate: new Date(),
    });

    return (
        <div className="min-h-screen p-6 flex flex-col items-center relative">
            <h1 className="text-2xl font-semibold mb-6 text-gray-100">Ledvix</h1>

            <div className="w-full max-w-6xl flex items-center justify-between gap-3 mb-4 text-gray-200">
                <WeekNavigator
                    label={rangeLabel}
                    onPrev={goToPreviousWeek}
                    onNext={goToNextWeek}
                />
                <button
                    className="px-4 py-2 rounded bg-blue-600 hover:bg-blue-700 text-white font-medium shadow-sm transition"
                    onClick={openNew}
                >
                    Log time
                </button>
            </div>

            <div className="relative w-full max-w-6xl">
                <TimesheetTable
                    days={days}
                    rows={rows}
                    totalsPerDayMinutes={totalsPerDayMinutes}
                    overallMinutes={overallMinutes}
                    onCellOpen={openFromCell}
                />
            </div>

            <TimeLogModal
                {...modalProps}
                loadingActivities={activitiesLoading}
                activityError={activitiesError}
            />

            <InvoiceGenerator />
        </div>
    );
}

export default function TimesheetPage() {
    const { selectedCompanyId } = useCompany();

    if (selectedCompanyId === null) {
        return (
            <div className="min-h-screen p-6 flex items-center justify-center bg-black/60">
                <CompanySelectionGate />
            </div>
        );
    }

    return <TimesheetContent selectedCompanyId={selectedCompanyId} />;
}
