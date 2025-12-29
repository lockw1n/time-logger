import React from "react";
import TimesheetTable from "./components/TimesheetTable";
import TimeLogModal from "./components/TimeLogModal";
import InvoiceGenerator from "./components/InvoiceGenerator";
import { useTimesheet } from "./hooks/useTimesheet";
import { useActivities } from "./hooks/useActivities";
import { useTimeLogForm } from "./hooks/useTimeLogForm";

export default function App() {
    const { days, rows, totals, rangeLabel, goToNextWeek, goToPreviousWeek, refresh } =
        useTimesheet();
    const { activities, loading: activitiesLoading, error: activitiesError } = useActivities(1);
    const { openNew, openFromCell, modalProps } = useTimeLogForm({
        activities,
        onSaved: refresh,
        defaultDate: new Date(),
    });

    return (
        <div className="min-h-screen p-6 flex flex-col items-center relative">
            <h1 className="text-3xl font-bold mb-6 text-gray-100 flex items-center gap-2">
                ⏱️ Time Logger <span className="text-gray-400">– Timesheet</span>
            </h1>

            <div className="w-full max-w-6xl flex items-center justify-between gap-3 mb-4 text-gray-200">
                <div className="flex items-center gap-3">
                    <button
                        className="px-3 py-1 rounded bg-gray-800 border border-gray-700 hover:bg-gray-700 transition"
                        onClick={goToPreviousWeek}
                    >
                        ← Previous week
                    </button>
                    <div className="px-4 py-1 rounded bg-gray-800 border border-gray-700 text-sm min-w-[220px] text-center">
                        {rangeLabel}
                    </div>
                    <button
                        className="px-3 py-1 rounded bg-gray-800 border border-gray-700 hover:bg-gray-700 transition"
                        onClick={goToNextWeek}
                    >
                        Next week →
                    </button>
                </div>
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
                    totalsByTicket={totals}
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
