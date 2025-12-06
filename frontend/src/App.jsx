import React, { useState } from "react";
import { createEntry, updateEntry, deleteEntry } from "./api/entries";
import TimesheetTable from "./components/TimesheetTable";
import TimeLogModal from "./components/TimeLogModal";
import InvoiceGenerator from "./components/InvoiceGenerator";
import { useTimesheet } from "./hooks/useTimesheet";
import { toYMD } from "./utils/date";
import { parseHoursInput } from "./utils/time";

export default function App() {
    const { days, rows, totals, rangeLabel, goToNextWeek, goToPreviousWeek, refresh } = useTimesheet();

    const [adding, setAdding] = useState(false);
    const [activeEntryId, setActiveEntryId] = useState(null);
    const [newTicket, setNewTicket] = useState("");
    const [newLabel, setNewLabel] = useState("feature");
    const [newHours, setNewHours] = useState("");
    const [newDate, setNewDate] = useState(() => toYMD(new Date()));
    const [errors, setErrors] = useState({});
    const [fieldsLocked, setFieldsLocked] = useState(false);

    const openLogForm = ({ ticket = "", label = "feature", date = new Date(), entry = null, locked = false }) => {
        setActiveEntryId(entry?.id ?? null);
        setNewTicket(ticket);
        setNewLabel(label || "feature");
        setNewHours(entry ? String(entry.hours ?? "") : "");
        setNewDate(toYMD(date));
        setErrors({});
        setFieldsLocked(locked);
        setAdding(true);
    };

    const resetFormState = () => {
        setActiveEntryId(null);
        setNewTicket("");
        setNewLabel("feature");
        setNewHours("");
        setNewDate(toYMD(new Date()));
        setErrors({});
        setFieldsLocked(false);
    };

    const handleAddRow = async (ticket, label, hours, dateValue) => {
        const nextErrors = {};
        const trimmedTicket = (ticket || "").trim();
        const trimmedLabel = (label || "").trim();
        const hoursNum = parseHoursInput(hours);
        const hasDate = Boolean(dateValue);

        if (!trimmedTicket) nextErrors.ticket = true;
        if (!trimmedLabel) nextErrors.label = true;
        if (Number.isNaN(hoursNum) || hoursNum < 0 || hoursNum > 24) nextErrors.hours = true;
        if (!hasDate) nextErrors.date = true;

        if (Object.keys(nextErrors).length) {
            setErrors(nextErrors);
            return;
        }

        setErrors({});
        const dateISO = (dateValue && dateValue.trim()) || toYMD(new Date());

        if (activeEntryId) {
            await updateEntry(activeEntryId, { id: activeEntryId, ticket: trimmedTicket, label: trimmedLabel, hours: hoursNum, date: dateISO });
        } else {
            await createEntry({ ticket: trimmedTicket, label: trimmedLabel, hours: hoursNum, date: dateISO });
        }
        setAdding(false);
        resetFormState();

        refresh();
    };

    const handleDeleteActiveEntry = async () => {
        if (!activeEntryId) return;
        await deleteEntry(activeEntryId);
        setAdding(false);
        resetFormState();
        refresh();
    };

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
                    onClick={() => openLogForm({})}
                >
                    Log time
                </button>
            </div>

            <div className="relative w-full max-w-6xl">
                <TimesheetTable
                    days={days}
                    rows={rows}
                    totalsByTicket={totals}
                    onCellOpen={({ ticket, label, date, entry }) => openLogForm({ ticket, label, date, entry, locked: true })}
                />
            </div>

            <TimeLogModal
                open={adding}
                ticket={newTicket}
                label={newLabel}
                hours={newHours}
                date={newDate}
                errors={errors}
                locked={fieldsLocked}
                canDelete={Boolean(activeEntryId)}
                onChangeTicket={(val) => {
                    setNewTicket(val);
                    setErrors((prev) => ({ ...prev, ticket: false }));
                }}
                onChangeLabel={(val) => {
                    setNewLabel(val);
                    setErrors((prev) => ({ ...prev, label: false }));
                }}
                onChangeHours={(val) => {
                    setNewHours(val);
                    setErrors((prev) => ({ ...prev, hours: false }));
                }}
                onChangeDate={(val) => {
                    setNewDate(val);
                    setErrors((prev) => ({ ...prev, date: false }));
                }}
                onCancel={() => {
                    setAdding(false);
                    resetFormState();
                }}
                onSave={() => handleAddRow(newTicket, newLabel, newHours, newDate)}
                onDelete={handleDeleteActiveEntry}
            />

            <InvoiceGenerator />
        </div>
    );
}
