import React from "react";

const COLORS = {
    bug: "bg-red-500/60 border-red-400/40",
    feature: "bg-blue-500/60 border-blue-400/40",
    meeting: "bg-green-500/60 border-green-400/40",
    research: "bg-yellow-500/60 border-yellow-400/40",
};

export default function TimesheetCell({ entry, label, onOpen, weekend = false, extraClass = "" }) {
    const isFilled = entry && entry.hours > 0;
    const emptyBg = weekend ? "bg-slate-600/60 border-slate-500" : "bg-gray-800/40 border-gray-700";
    const filledBg = COLORS[label] || "bg-gray-700/60 border-gray-600";
    const bgColor = isFilled ? filledBg : emptyBg;

    const weekendTint = weekend ? "ring-1 ring-slate-400/50 shadow-inner shadow-slate-300/10" : "";

    return (
        <td
            className={`${bgColor} ${weekendTint} text-center px-2 py-2 border cursor-pointer transition hover:brightness-110 w-20 ${extraClass}`}
            onClick={onOpen}
        >
            <span className="font-medium">
                {entry && entry.hours > 0 ? entry.hours : ""}
            </span>
        </td>
    );
}
