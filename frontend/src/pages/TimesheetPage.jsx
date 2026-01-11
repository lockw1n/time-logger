import React, { useEffect } from "react";
import TimesheetTable from "../components/TimesheetTable";
import TimeLogModal from "../components/TimeLogModal";
import InvoiceGenerator from "../components/InvoiceGenerator";
import WeekNavigator from "../components/WeekNavigator";
import CompanySelectionGate from "../components/CompanySelectionGate";
import TimesheetSkeleton from "../components/TimesheetSkeleton";
import { useTimesheet } from "../hooks/useTimesheet";
import { useActivities } from "../hooks/useActivities";
import { useTimeLogForm } from "../hooks/useTimeLogForm";
import { useCompany } from "../context/CompanyContext";

function TimesheetContent({
    selectedCompanyId,
    isSwitchingCompany,
    isLoading,
    days,
    rows,
    totalsPerDayMinutes,
    overallMinutes,
    rangeLabel,
    goToNextWeek,
    goToPreviousWeek,
    refresh,
}) {
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

    const hasNoActivities = !activitiesLoading && activities.length === 0;
    const hasNoEntries =
        !isLoading && activities.length > 0 && rows.length === 0 && !isSwitchingCompany;
    const showSkeleton =
        isLoading && !isSwitchingCompany && activities.length > 0 && !hasNoEntries;
    const handleCellOpen = isSwitchingCompany ? () => {} : openFromCell;

    return (
        <div className="min-h-screen p-6 flex flex-col items-center relative">
            <h1 className="text-2xl font-semibold mb-6 text-gray-100">Ledvix</h1>

            <div className="w-full max-w-6xl flex items-center justify-between gap-3 mb-4 text-gray-200">
                <WeekNavigator
                    label={rangeLabel}
                    onPrev={goToPreviousWeek}
                    onNext={goToNextWeek}
                    disabled={isSwitchingCompany || isLoading}
                />
                <button
                    className="px-4 py-2 rounded bg-blue-600 hover:bg-blue-700 text-white font-medium shadow-sm transition"
                    onClick={openNew}
                    disabled={isSwitchingCompany || isLoading}
                >
                    Log time
                </button>
            </div>

            <div className="relative w-full max-w-6xl">
                {hasNoActivities ? (
                    <div className="shadow-lg rounded-xl bg-gray-800 border border-gray-700 min-h-[260px] flex items-center justify-center text-center px-6">
                        <div className="max-w-md">
                            <h2 className="text-lg font-semibold text-gray-100">
                                No activities available
                            </h2>
                            <p className="text-sm text-gray-300 mt-2">
                                You can’t log time until at least one activity exists for this company.
                            </p>
                        </div>
                    </div>
                ) : showSkeleton ? (
                    <TimesheetSkeleton />
                ) : hasNoEntries ? (
                    <div className="shadow-lg rounded-xl bg-gray-800 border border-gray-700 min-h-[260px] flex items-center justify-center text-center px-6">
                        <div className="max-w-md">
                            <h2 className="text-lg font-semibold text-gray-100">
                                No time entries yet
                            </h2>
                            <p className="text-sm text-gray-300 mt-2">
                                There are no logged hours for this period.
                            </p>
                            <button
                                className="mt-4 px-4 py-2 rounded bg-blue-600 hover:bg-blue-700 text-white font-medium shadow-sm transition"
                                onClick={openNew}
                            >
                                Log time
                            </button>
                        </div>
                    </div>
                ) : (
                    <TimesheetTable
                        days={days}
                        rows={rows}
                        totalsPerDayMinutes={totalsPerDayMinutes}
                        overallMinutes={overallMinutes}
                        onCellOpen={handleCellOpen}
                    />
                )}
            </div>

            <TimeLogModal
                {...modalProps}
                loadingActivities={activitiesLoading}
                activityError={activitiesError}
            />

            <InvoiceGenerator overallMinutes={overallMinutes} />

            {isSwitchingCompany ? (
                <div className="absolute inset-0 bg-black/20 flex items-start justify-center pt-24 text-sm text-gray-200 pointer-events-auto">
                    Switching company…
                </div>
            ) : null}
        </div>
    );
}

export default function TimesheetPage() {
    const { selectedCompanyId, isSwitchingCompany, finishSwitchCompany } = useCompany();
    const {
        days,
        rows,
        totalsPerDayMinutes,
        overallMinutes,
        rangeLabel,
        goToNextWeek,
        goToPreviousWeek,
        refresh,
        isLoading,
    } = useTimesheet(selectedCompanyId);

    useEffect(() => {
        if (isSwitchingCompany && !isLoading) {
            finishSwitchCompany();
        }
    }, [finishSwitchCompany, isLoading, isSwitchingCompany]);

    if (selectedCompanyId === null) {
        return (
            <div className="min-h-screen p-6 flex items-center justify-center bg-black/60">
                <CompanySelectionGate />
            </div>
        );
    }

    return (
        <TimesheetContent
            selectedCompanyId={selectedCompanyId}
            isSwitchingCompany={isSwitchingCompany}
            isLoading={isLoading}
            days={days}
            rows={rows}
            totalsPerDayMinutes={totalsPerDayMinutes}
            overallMinutes={overallMinutes}
            rangeLabel={rangeLabel}
            goToNextWeek={goToNextWeek}
            goToPreviousWeek={goToPreviousWeek}
            refresh={refresh}
        />
    );
}
