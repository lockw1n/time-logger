import React, { useEffect, useState } from "react";
import { getEntries, createEntry, updateEntry, deleteEntry } from "./api/entries";
import TimesheetTable from "./components/TimesheetTable";

export default function App() {
    const [rows, setRows] = useState([]);
    const [days, setDays] = useState([]);
    const [showAddButton, setShowAddButton] = useState(false);
    const [adding, setAdding] = useState(false);

    // --- Generate current week days ---
    useEffect(() => {
        const today = new Date();
        const weekDays = Array.from({ length: 7 }, (_, i) => {
            const d = new Date(Date.UTC(today.getFullYear(), today.getMonth(), today.getDate() - today.getDay() + i));
            return d;
        });
        setDays(weekDays);
    }, []);

    // --- Load entries from backend ---
    useEffect(() => {
        (async () => {
            const data = await getEntries();

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
        })();
    }, []);

    // --- Handle cell edit ---
    const handleCellChange = async (ticket, label, utcDate, entry, hours) => {
        if (entry) {
            const updated = await updateEntry(entry.id, { ...entry, hours, date: utcDate });
            setRows((prev) =>
                prev.map((r) =>
                    r.ticket === ticket ? { ...r, cells: { ...r.cells, [utcDate.slice(0, 10)]: updated } } : r
                )
            );
        } else {
            const created = await createEntry({ ticket, label, hours, date: utcDate });
            setRows((prev) =>
                prev.map((r) =>
                    r.ticket === ticket ? { ...r, cells: { ...r.cells, [utcDate.slice(0, 10)]: created } } : r
                )
            );
        }
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
    const handleAddRow = async (ticket, label) => {
        if (!ticket || !label) return;
        await createEntry({ ticket, label, hours: 0, date: new Date().toISOString() });
        setAdding(false);
        setShowAddButton(false);

        const data = await getEntries();
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

            <div
                className="relative"
                onMouseEnter={() => setShowAddButton(true)}
                onMouseLeave={() => setShowAddButton(false)}
            >
                <TimesheetTable
                    days={days}
                    rows={rows}
                    onCellChange={handleCellChange}
                    onDeleteRow={handleDeleteRow}
                />

                {showAddButton && !adding && (
                    <button
                        className="absolute -bottom-8 left-1/2 transform -translate-x-1/2 bg-blue-600 text-white px-3 py-1 rounded-full shadow-md hover:bg-blue-700 transition"
                        onClick={() => setAdding(true)}
                    >
                        ➕ Add Row
                    </button>
                )}
            </div>

            {adding && (
                <div className="mt-8 flex gap-2 items-center bg-gray-800 p-4 rounded-lg shadow-md border border-gray-700">
                    <input
                        type="text"
                        placeholder="Ticket key (e.g. ABC-123)"
                        id="ticket"
                        className="px-3 py-2 rounded bg-gray-900 border border-gray-600 text-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-400"
                    />
                    <select
                        id="label"
                        className="px-3 py-2 rounded bg-gray-900 border border-gray-600 text-gray-100 focus:outline-none"
                    >
                        <option value="feature">Feature</option>
                        <option value="bug">Bug</option>
                        <option value="meeting">Meeting</option>
                        <option value="research">Research</option>
                    </select>
                    <button
                        className="bg-blue-600 px-4 py-2 rounded text-white font-medium hover:bg-blue-700"
                        onClick={() =>
                            handleAddRow(
                                document.getElementById("ticket").value,
                                document.getElementById("label").value
                            )
                        }
                    >
                        Save
                    </button>
                    <button
                        className="text-gray-400 hover:text-gray-200"
                        onClick={() => setAdding(false)}
                    >
                        ✖ Cancel
                    </button>
                </div>
            )}
        </div>
    );
}
