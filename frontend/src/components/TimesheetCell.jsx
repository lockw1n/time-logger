import React, { useState } from "react";

const COLORS = {
    bug: "bg-red-500/60 border-red-400/40",
    feature: "bg-blue-500/60 border-blue-400/40",
    meeting: "bg-green-500/60 border-green-400/40",
    research: "bg-yellow-500/60 border-yellow-400/40",
};

export default function TimesheetCell({ entry, label, onChange }) {
    const [editing, setEditing] = useState(false);
    const [value, setValue] = useState(entry ? entry.hours : "");

    const isFilled = entry && entry.hours > 0;
    const bgColor = isFilled ? COLORS[label] || "bg-gray-700/60 border-gray-600"
        : "bg-gray-800/40 border-gray-700";

    const handleBlur = () => {
        setEditing(false);
        const num = parseFloat(value);
        if (!isNaN(num) && num >= 0) onChange(num);
    };

    return (
        <td
            className={`${bgColor} text-center px-2 py-2 border cursor-pointer transition hover:brightness-110`}
            onClick={() => setEditing(true)}
        >
            {editing ? (
                <input
                    type="number"
                    step="0.25"
                    value={value}
                    onChange={(e) => setValue(e.target.value)}
                    onBlur={handleBlur}
                    onKeyDown={(e) => e.key === "Enter" && e.target.blur()}
                    className="w-14 rounded bg-gray-900 text-gray-100 text-center border border-gray-600 focus:outline-none focus:ring-1 focus:ring-blue-400"
                    autoFocus
                />
            ) : (
                <span className="font-medium">
                    {entry && entry.hours > 0 ? entry.hours : ""}
                </span>
            )}
        </td>
    );
}
