import React from "react";
import TimesheetCell from "./TimesheetCell";
import { getStartOfWeek } from "../utils/date";

const COLORS = {
    red: "bg-red-500/60 border-red-400/40",
    blue: "bg-blue-500/60 border-blue-400/40",
    green: "bg-green-500/60 border-green-400/40",
    yellow: "bg-yellow-500/60 border-yellow-400/40",
};

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

const getWeekKey = (date) => toLocalKey(getStartOfWeek(date));

const sumMinutesForDays = (perDayMinutes, days) =>
    days.reduce((total, day) => total + (perDayMinutes?.[toLocalKey(day)] ?? 0), 0);

export default function TimesheetTable({
    days,
    rows,
    totalsPerDayMinutes = {},
    overallMinutes = 0,
    onCellOpen,
}) {
    const firstWeekKey = days.length ? getWeekKey(days[0]) : null;
    const secondWeekStartIndex = days.findIndex((day) => getWeekKey(day) !== firstWeekKey);
    const splitIndex = secondWeekStartIndex === -1 ? days.length : secondWeekStartIndex;
    const firstWeekDays = days.slice(0, splitIndex);
    const secondWeekDays = days.slice(splitIndex);
    const separatorClass = "border-r-8 border-r-slate-400";
    return (
        <div className="shadow-lg rounded-xl bg-gray-800 border border-gray-700">
            <table className="w-full border-collapse text-sm table-fixed">
                <colgroup>
                    <col className="w-[140px]" />
                    <col className="w-[40px]" />
                    {firstWeekDays.map((d) => (
                        <col key={toLocalKey(d)} className="w-12" />
                    ))}
                    <col className="w-16" />
                    {secondWeekDays.map((d) => (
                        <col key={toLocalKey(d)} className="w-12" />
                    ))}
                    <col className="w-16" />
                </colgroup>
                <thead className="bg-gray-700 text-gray-200 sticky top-0">
                <tr>
                    <th className="px-4 py-2 text-left min-w-[140px]">Ticket</th>
                    <th className="px-2 py-2 text-center min-w-[40px]" aria-label="Activity"></th>
                    {firstWeekDays.map((d) => {
                        const weekend = isWeekend(d);
                        return (
                            <th
                                key={toLocalKey(d)}
                                className={`px-1 py-2 text-center w-12 ${weekend ? "bg-slate-700 text-slate-100" : ""}`}
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
                    <th className={`px-2 py-2 text-center w-16 ${separatorClass}`}>Total</th>
                    {secondWeekDays.map((d) => {
                        const weekend = isWeekend(d);
                        return (
                            <th
                                key={toLocalKey(d)}
                                className={`px-1 py-2 text-center w-12 ${weekend ? "bg-slate-700 text-slate-100" : ""}`}
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
                    <th className="px-2 py-2 text-center w-16">Total</th>
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
                                {(() => {
                                    const paletteKey = (row.color || "").toLowerCase();
                                    const swatchClass = COLORS[paletteKey] || "bg-gray-700/60 border-gray-600";
                                    const hasCustomColor = row.color && !COLORS[paletteKey];
                                    const swatchStyle = hasCustomColor
                                        ? { backgroundColor: row.color, borderColor: row.color }
                                        : undefined;

                                    return (
                                        <span
                                            className={`h-4 w-4 border ${swatchClass}`}
                                            style={swatchStyle}
                                            aria-hidden="true"
                                        />
                                    );
                                })()}
                            </div>
                        </td>

                        {firstWeekDays.map((d) => {
                            const key = toLocalKey(d);
                            const entry = row.cells[key];
                            const weekend = isWeekend(d);

                            return (
                                <TimesheetCell
                                    key={key}
                                    entry={entry}
                                    color={row.color}
                                    weekend={weekend}
                                    onOpen={() => onCellOpen({ticket: row.ticket, activity: row.activity, date: d, entry})}
                                />
                            );
                        })}

                        <td className={`px-2 py-2 text-center font-semibold text-gray-100 ${separatorClass}`}>
                            {formatHours(sumMinutesForDays(row.perDayMinutes, firstWeekDays) / 60)}
                        </td>

                        {secondWeekDays.map((d) => {
                            const key = toLocalKey(d);
                            const entry = row.cells[key];
                            const weekend = isWeekend(d);

                            return (
                                <TimesheetCell
                                    key={key}
                                    entry={entry}
                                    color={row.color}
                                    weekend={weekend}
                                    onOpen={() => onCellOpen({ticket: row.ticket, activity: row.activity, date: d, entry})}
                                />
                            );
                        })}

                        <td className="px-2 py-2 text-center font-semibold text-gray-100">
                            {formatHours(sumMinutesForDays(row.perDayMinutes, secondWeekDays) / 60)}
                        </td>
                    </tr>
                ))}
                <tr className="border-t border-gray-600 bg-gray-700/60">
                    <td className="px-4 py-2 font-semibold text-gray-100">Total</td>
                    <td className="px-2 py-2"></td>
                    {firstWeekDays.map((d) => (
                        <td key={toLocalKey(d)} className="px-1 py-2 text-center font-semibold text-gray-100">
                            {formatHours((totalsPerDayMinutes[toLocalKey(d)] ?? 0) / 60)}
                        </td>
                    ))}
                    <td className={`px-2 py-2 text-center font-semibold text-gray-100 ${separatorClass}`}>
                        {formatHours(sumMinutesForDays(totalsPerDayMinutes, firstWeekDays) / 60)}
                    </td>
                    {secondWeekDays.map((d) => (
                        <td key={toLocalKey(d)} className="px-1 py-2 text-center font-semibold text-gray-100">
                            {formatHours((totalsPerDayMinutes[toLocalKey(d)] ?? 0) / 60)}
                        </td>
                    ))}
                    <td className="px-2 py-2 text-center font-semibold text-gray-100">
                        {formatHours(sumMinutesForDays(totalsPerDayMinutes, secondWeekDays) / 60)}
                    </td>
                </tr>
                </tbody>
            </table>
        </div>
    );
}
