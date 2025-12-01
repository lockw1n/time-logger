import React, { useEffect, useState } from "react";
import { getEntries, createEntry, updateEntry, deleteEntry } from "./api/entries";
import { getStartOfWeek, getWeekDays } from "./utils/date";
import TimesheetTable from "./components/TimesheetTable";

const formatRangeDate = (date) => {
    const d = new Date(date);
    const day = String(d.getDate()).padStart(2, "0");
    const month = String(d.getMonth() + 1).padStart(2, "0");
    const year = d.getFullYear();
    return `${day}.${month}.${year}`;
};

const parseHoursInput = (input) => {
    if (!input) return NaN;
    const val = String(input).trim().toLowerCase();
    if (!val) return NaN;

    // Match patterns like "2h", "15m", "2h 30m"
    const timeRegex = /(\d+(?:\.\d+)?)\s*(h|m)/g;
    let match;
    let totalHours = 0;
    let found = false;
    while ((match = timeRegex.exec(val)) !== null) {
        found = true;
        const num = parseFloat(match[1]);
        if (Number.isNaN(num)) continue;
        if (match[2] === "h") totalHours += num;
        if (match[2] === "m") totalHours += num / 60;
    }
    if (found) return totalHours;

    // Fallback: plain numeric hours (supports dot)
    if (/^-?\d+(\.\d+)?$/.test(val) || /^-?\d+(,\d+)?$/.test(val)) {
        const numeric = parseFloat(val.replace(",", "."));
        return Number.isNaN(numeric) ? NaN : numeric;
    }

    return NaN;
};

const LABEL_COLORS = {
    feature: "bg-blue-500",
    bug: "bg-red-500",
    meeting: "bg-green-500",
    research: "bg-yellow-400",
};

export default function App() {
    const [rows, setRows] = useState([]);
    const [days, setDays] = useState([]);
    const [adding, setAdding] = useState(false);
    const [anchorDate, setAnchorDate] = useState(() => getStartOfWeek(new Date())); // Monday of current view
    const [newTicket, setNewTicket] = useState("");
    const [newLabel, setNewLabel] = useState("feature");
    const [newHours, setNewHours] = useState("");
    const [newDate, setNewDate] = useState(() => {
        const now = new Date();
        const pad = (n) => String(n).padStart(2, "0");
        return `${now.getFullYear()}-${pad(now.getMonth() + 1)}-${pad(now.getDate())}`;
    });
    const [errors, setErrors] = useState({});
    const [fieldsLocked, setFieldsLocked] = useState(false);

    const buildRangeParams = (anchor) => {
        const start = new Date(anchor);
        start.setHours(0, 0, 0, 0);
        const end = new Date(anchor);
        end.setDate(end.getDate() + 13); // inclusive of two weeks
        end.setHours(0, 0, 0, 0);

        return {
            startStr: start.toISOString().split("T")[0],
            endStr: end.toISOString().split("T")[0],
        };
    };

    // --- Generate day list for current two-week window ---
    useEffect(() => {
        setDays(getWeekDays(anchorDate, 14));
    }, [anchorDate]);

    // --- Load entries from backend for visible range ---
    useEffect(() => {
        const fetchEntries = async () => {
            const { startStr, endStr } = buildRangeParams(anchorDate);
            const data = await getEntries({ start: startStr, end: endStr });

            // Group entries by ticket -> dateKey (UTC YYYY-MM-DD)
            const grouped = {};
            data.forEach((e) => {
                const utcKey = new Date(e.date || e.created_at).toISOString().split("T")[0];
                if (!grouped[e.ticket]) {
                    grouped[e.ticket] = {
                        ticket: e.ticket,
                        label: e.label || "feature",
                        cells: {},
                    };
                }
                grouped[e.ticket].cells[utcKey] = e;
            });
            setRows(Object.values(grouped));
        };

        fetchEntries();
    }, [anchorDate]);

    const goToPreviousWeek = () => {
        setAnchorDate((prev) => {
            const d = new Date(prev);
            d.setDate(d.getDate() - 7);
            return getStartOfWeek(d);
        });
    };

    const goToNextWeek = () => {
        setAnchorDate((prev) => {
            const d = new Date(prev);
            d.setDate(d.getDate() + 7);
            return getStartOfWeek(d);
        });
    };

    // --- Delete full ticket row ---
    const handleDeleteRow = async (ticket) => {
        const row = rows.find((r) => r.ticket === ticket);
        if (!row) return;
        const ids = Object.values(row.cells).map((c) => c.id);
        await Promise.all(ids.map((id) => deleteEntry(id)));
        setRows((prev) => prev.filter((r) => r.ticket !== ticket));
    };

    // --- Add new empty row ---
    const [activeEntryId, setActiveEntryId] = useState(null);

    const openLogForm = ({ ticket = "", label = "feature", date = new Date(), entry = null, locked = false }) => {
        const isoDate = new Date(date).toISOString().split("T")[0];
        setActiveEntryId(entry?.id ?? null);
        setNewTicket(ticket);
        setNewLabel(label || "feature");
        setNewHours(entry ? String(entry.hours ?? "") : "");
        setNewDate(isoDate);
        setErrors({});
        setFieldsLocked(locked);
        setAdding(true);
    };

    const resetFormState = () => {
        setActiveEntryId(null);
        setNewTicket("");
        setNewLabel("feature");
        setNewHours("");
        setNewDate(() => {
            const now = new Date();
            const pad = (n) => String(n).padStart(2, "0");
            return `${now.getFullYear()}-${pad(now.getMonth() + 1)}-${pad(now.getDate())}`;
        });
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

        const dateISO = dateValue ? new Date(dateValue).toISOString() : new Date().toISOString();

        if (activeEntryId) {
            await updateEntry(activeEntryId, { id: activeEntryId, ticket: trimmedTicket, label: trimmedLabel, hours: hoursNum, date: dateISO });
        } else {
            await createEntry({ ticket: trimmedTicket, label: trimmedLabel, hours: hoursNum, date: dateISO });
        }
        setAdding(false);
        resetFormState();

        const { startStr, endStr } = buildRangeParams(anchorDate);
        const data = await getEntries({ start: startStr, end: endStr });
        const grouped = {};
        data.forEach((e) => {
            const utcKey = new Date(e.date || e.created_at).toISOString().split("T")[0];
            if (!grouped[e.ticket]) {
                grouped[e.ticket] = { ticket: e.ticket, label: e.label || "feature", cells: {} };
            }
            grouped[e.ticket].cells[utcKey] = e;
        });
        setRows(Object.values(grouped));
    };

    // --- Render ---
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
                    <div className="px-4 py-1 rounded bg-gray-800 border border-gray-700 text-sm">
                        {days.length > 0
                            ? `${formatRangeDate(days[0])} – ${formatRangeDate(days[days.length - 1])}`
                            : ""}
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

            <div
                className="relative w-full max-w-6xl"
            >
                <TimesheetTable
                    days={days}
                    rows={rows}
                    onCellOpen={({ ticket, label, date, entry }) => openLogForm({ ticket, label, date, entry, locked: true })}
                    onDeleteRow={handleDeleteRow}
                />

            </div>

            {adding && (
                <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50 px-4">
                    <div className="bg-gray-800 p-6 rounded-lg shadow-2xl border border-gray-700 w-full max-w-lg">
                        <h2 className="text-xl font-semibold text-gray-100 mb-4">Log time</h2>
                        <div className="flex flex-col gap-3">
                            <input
                                type="text"
                                placeholder="Ticket key (e.g. ABC-123)"
                                value={newTicket}
                                onChange={(e) => {
                                    setNewTicket(e.target.value);
                                    setErrors((prev) => ({ ...prev, ticket: false }));
                                }}
                                readOnly={fieldsLocked}
                                className={`px-3 py-2 rounded border text-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-400 w-full ${
                                    errors.ticket ? "border-red-500 ring-red-400/60" : "border-gray-600"
                                } ${fieldsLocked ? "bg-gray-800 cursor-not-allowed" : "bg-gray-900"}`}
                            />
                            <div className="relative w-full">
                                <select
                                    value={newLabel}
                                    onChange={(e) => {
                                        setNewLabel(e.target.value);
                                        setErrors((prev) => ({ ...prev, label: false }));
                                    }}
                                    disabled={fieldsLocked}
                                    className={`label-select px-3 py-2 pr-14 rounded border text-gray-100 focus:outline-none w-full ${
                                        errors.label ? "border-red-500 ring-1 ring-red-400/60" : "border-gray-600"
                                    } ${fieldsLocked ? "bg-gray-800 cursor-not-allowed" : "bg-gray-900"}`}
                                >
                                    <option value="feature">Feature</option>
                                    <option value="bug">Bug</option>
                                    <option value="meeting">Meeting</option>
                                    <option value="research">Research</option>
                                </select>
                                <span
                                    className={`absolute right-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 rounded-sm border border-gray-600 ${LABEL_COLORS[newLabel] || "bg-gray-500"}`}
                                    aria-hidden="true"
                                />
                            </div>
                            <div className="grid grid-cols-2 gap-3">
                                <input
                                    type="text"
                                    inputMode="decimal"
                                    placeholder="Hours (e.g. 2.75 or 2h 15m)"
                                    value={newHours}
                                    onChange={(e) => {
                                        setNewHours(e.target.value);
                                        setErrors((prev) => ({ ...prev, hours: false }));
                                    }}
                                    className={`px-3 py-2 rounded bg-gray-900 border text-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-400 ${
                                        errors.hours ? "border-red-500 ring-red-400/60" : "border-gray-600"
                                    }`}
                                />
                                <input
                                    type="date"
                                    value={newDate}
                                    onChange={(e) => {
                                        setNewDate(e.target.value);
                                        setErrors((prev) => ({ ...prev, date: false }));
                                    }}
                                    className={`px-3 py-2 rounded bg-gray-900 border text-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-400 date-input ${
                                        errors.date ? "border-red-500 ring-red-400/60" : "border-gray-600"
                                    }`}
                                />
                            </div>
                        </div>
                        <div className="flex justify-end gap-3 mt-5">
                            <button
                                className="text-gray-300 hover:text-white"
                                onClick={() => {
                                    setAdding(false);
                                    resetFormState();
                                }}
                            >
                                Cancel
                            </button>
                            <button
                                className="bg-blue-600 px-4 py-2 rounded text-white font-medium hover:bg-blue-700"
                                onClick={() => handleAddRow(newTicket, newLabel, newHours, newDate)}
                            >
                                Save
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
