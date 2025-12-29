import React from "react";
import TimesheetCell from "./TimesheetCell";

const toLocalKey = (date) => {
    const d = new Date(date);
    const pad = (n) => String(n).padStart(2, "0");
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
};

const isWeekend = (date) => {
    const day = date.getDay();
    return day === 6 || day === 0;
};

const formatHours = (val = 0) => {
    if (val === null || val === undefined) return "";
    const rounded = Math.round(val * 100) / 100;
    if (Number.isInteger(rounded)) return String(rounded);
    return rounded.toFixed(2).replace(/\.?0+$/, "");
};

export default function TimesheetTable({days, rows, totalsByTicket = {}, onCellOpen}) {
    return (
        <div className="overflow-x-auto shadow-lg rounded-xl bg-gray-800 border border-gray-700">
            <table className="w-full border-collapse text-sm table-fixed">
                <colgroup>
                    <col className="w-[140px]" />
                    <col className="w-[40px]" />
                    {days.map((d) => (
                        <col key={toLocalKey(d)} className="w-16" />
                    ))}
                    <col className="w-16" />
                </colgroup>
                <thead className="bg-gray-700 text-gray-200 sticky top-0">
                <tr>
                    <th className="px-4 py-2 text-left min-w-[140px]">Ticket</th>
                    <th className="px-2 py-2 text-center min-w-[40px]" aria-label="Activity"></th>
                    {days.map((d, idx) => {
                        const weekend = isWeekend(d);
                        const startOfSecondWeek = idx === 7;
                        const separatorClass = startOfSecondWeek ? "border-l-8 border-l-slate-400" : "";
                        return (
                            <th
                                key={toLocalKey(d)}
                                className={`px-3 py-2 text-center w-20 ${weekend ? "bg-slate-700 text-slate-100" : ""} ${separatorClass}`}
                            >
                                <div className="flex flex-col items-center leading-tight">
                                    <span className="font-semibold">
                                        {d.toLocaleDateString(undefined, {weekday: "short"})}
                                    </span>
                                    <span className="text-xs text-gray-300">
                                        {d.toLocaleDateString(undefined, {month: "short"})} {d.getDate()}
                                    </span>
                                </div>
                            </th>
                        );
                    })}
                    <th className="px-3 py-2 text-center w-20">Total</th>
                </tr>
                </thead>

                <tbody>
                {rows.map((row, rowIdx) => (
                    <tr
                        key={`${row.ticket?.id || row.ticket?.code || rowIdx}-${row.activity?.id || "activity"}`}
                        className="border-t border-gray-700 hover:bg-gray-750 transition-colors"
                    >
                        <td className="px-4 py-2 font-medium text-blue-300 min-w-[140px] whitespace-nowrap truncate">
                            {row.ticket?.code || "—"}
                        </td>
                        <td
                            className="px-2 py-2 text-gray-200 min-w-[40px]"
                            title={row.activity?.name || "No activity"}
                        >
                            <div className="flex items-center justify-center">
                                <span
                                    className="h-5 w-5 rounded-full border border-gray-600"
                                    style={{ backgroundColor: row.activity?.color || "#6b7280" }}
                                    aria-hidden="true"
                                />
                            </div>
                        </td>

                        {days.map((d, idx) => {
                            const key = toLocalKey(d);
                            const entry = row.cells[key];
                            const weekend = isWeekend(d);
                            const startOfSecondWeek = idx === 7;
                            const separatorClass = startOfSecondWeek ? "border-l-8 border-l-slate-400" : "";

                            return (
                                <TimesheetCell
                                    key={key}
                                    entry={entry}
                                    color={row.color}
                                    weekend={weekend}
                                    extraClass={separatorClass}
                                    onOpen={() => onCellOpen({ticket: row.ticket, activity: row.activity, date: d, entry})}
                                />
                            );
                        })}

                        <td className="px-3 py-2 text-center font-semibold text-gray-100">
                            {formatHours(totalsByTicket[row.ticket?.code] || 0)}
                        </td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
}
