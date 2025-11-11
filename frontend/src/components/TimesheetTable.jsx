import React from "react";
import TimesheetCell from "./TimesheetCell";

const toUtcKey = (date) => {
    const utc = new Date(date);
    return utc.toISOString().split("T")[0]; // YYYY-MM-DD (UTC)
};

export default function TimesheetTable({ days, rows, onCellChange, onDeleteRow }) {
    return (
        <div className="overflow-x-auto shadow-lg rounded-xl bg-gray-800 border border-gray-700">
            <table className="w-full border-collapse text-sm">
                <thead className="bg-gray-700 text-gray-200 sticky top-0">
                <tr>
                    <th className="px-4 py-2 text-left">Ticket</th>
                    {days.map((d) => (
                        <th key={toUtcKey(d)} className="px-3 py-2 text-center">
                            {d.toLocaleDateString(undefined, {
                                weekday: "short",
                                day: "numeric",
                            })}
                        </th>
                    ))}
                    <th className="px-2 py-2 text-center">Delete</th>
                </tr>
                </thead>

                <tbody>
                {rows.map((row) => (
                    <tr
                        key={row.ticket}
                        className="border-t border-gray-700 hover:bg-gray-750 transition-colors"
                    >
                        <td className="px-4 py-2 font-medium text-blue-300">{row.ticket}</td>

                        {days.map((d) => {
                            const key = toUtcKey(d);
                            const entry = row.cells[key];

                            return (
                                <TimesheetCell
                                    key={key}
                                    entry={entry}
                                    label={row.label}
                                    onChange={(hours) => {
                                        const utcDate = new Date(d).toISOString();
                                        onCellChange(row.ticket, row.label, utcDate, entry, hours);
                                    }}
                                />
                            );
                        })}

                        <td className="px-2 py-2 text-center">
                            <button
                                onClick={() => onDeleteRow(row.ticket)}
                                className="text-gray-400 hover:text-red-500 transition"
                            >
                                🗑️
                            </button>
                        </td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
}
