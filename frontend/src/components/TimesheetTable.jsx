import React from "react";
import TimesheetCell from "./TimesheetCell";

const toUtcKey = (date) => {
    const utc = new Date(date);
    return utc.toISOString().split("T")[0]; // YYYY-MM-DD (UTC)
};

const isWeekend = (date) => {
    const day = date.getDay();
    return day === 6 || day === 0;
};

export default function TimesheetTable({ days, rows, onCellOpen, onDeleteRow }) {
    return (
        <div className="overflow-x-auto shadow-lg rounded-xl bg-gray-800 border border-gray-700">
            <table className="w-full border-collapse text-sm">
                <thead className="bg-gray-700 text-gray-200 sticky top-0">
                <tr>
                    <th className="px-4 py-2 text-left min-w-[140px]">Ticket</th>
                    {days.map((d, idx) => {
                        const weekend = isWeekend(d);
                        const startOfSecondWeek = idx === 7;
                        const separatorClass = startOfSecondWeek ? "border-l-8 border-l-slate-400" : "";
                        return (
                            <th
                                key={toUtcKey(d)}
                                className={`px-3 py-2 text-center w-20 ${weekend ? "bg-slate-700 text-slate-100" : ""} ${separatorClass}`}
                            >
                                <div className="flex flex-col items-center leading-tight">
                                    <span className="font-semibold">
                                        {d.toLocaleDateString(undefined, { weekday: "short" })}
                                    </span>
                                    <span className="text-xs text-gray-300">
                                        {d.toLocaleDateString(undefined, { month: "short" })} {d.getDate()}
                                    </span>
                                </div>
                            </th>
                        );
                    })}
                    <th className="px-2 py-2 text-center">Delete</th>
                </tr>
                </thead>

                <tbody>
                {rows.map((row) => (
                    <tr
                        key={row.ticket}
                        className="border-t border-gray-700 hover:bg-gray-750 transition-colors"
                    >
                        <td className="px-4 py-2 font-medium text-blue-300 min-w-[140px]">{row.ticket}</td>

                        {days.map((d, idx) => {
                            const key = toUtcKey(d);
                            const entry = row.cells[key];
                            const weekend = isWeekend(d);
                            const startOfSecondWeek = idx === 7;
                            const separatorClass = startOfSecondWeek ? "border-l-8 border-l-slate-400" : "";

                            return (
                <TimesheetCell
                    key={key}
                    entry={entry}
                    label={row.label}
                    weekend={weekend}
                    extraClass={separatorClass}
                    onOpen={() => onCellOpen({ ticket: row.ticket, label: row.label, date: d, entry })}
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
