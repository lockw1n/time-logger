import React, { useEffect } from "react";

const LABEL_COLORS = {
    feature: "bg-blue-500",
    bug: "bg-red-500",
    meeting: "bg-green-500",
    research: "bg-yellow-400",
};

export default function TimeLogModal({
    open,
    ticket,
    label,
    hours,
    date,
    errors = {},
    locked = false,
    onChangeTicket,
    onChangeLabel,
    onChangeHours,
    onChangeDate,
    onCancel,
    onSave,
    onDelete,
    canDelete = false,
}) {
    useEffect(() => {
        if (!open) return;
        const onKeyDown = (e) => {
            if (e.key === "Escape") {
                onCancel();
            }
        };
        window.addEventListener("keydown", onKeyDown);
        return () => {
            window.removeEventListener("keydown", onKeyDown);
        };
    }, [open, onCancel]);

    if (!open) return null;

    return (
        <div
            className="fixed inset-0 bg-black/60 flex items-center justify-center z-50 px-4"
            onClick={onCancel}
        >
            <div
                className="bg-gray-800 p-6 rounded-lg shadow-2xl border border-gray-700 w-full max-w-lg"
                onClick={(e) => e.stopPropagation()}
            >
                <h2 className="text-xl font-semibold text-gray-100 mb-4">Log time</h2>
                <div className="flex flex-col gap-3">
                    <input
                        type="text"
                        placeholder="Ticket key (e.g. ABC-123)"
                        value={ticket}
                        onChange={(e) => onChangeTicket(e.target.value)}
                        readOnly={locked}
                        className={`px-3 py-2 rounded border text-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-400 w-full ${
                            errors.ticket ? "border-red-500 ring-red-400/60" : "border-gray-600"
                        } ${locked ? "bg-gray-800 cursor-not-allowed" : "bg-gray-900"}`}
                    />
                    <div className="relative w-full">
                        <select
                            value={label}
                            onChange={(e) => onChangeLabel(e.target.value)}
                            disabled={locked}
                            className={`label-select px-3 py-2 pr-14 rounded border text-gray-100 focus:outline-none w-full ${
                                errors.label ? "border-red-500 ring-1 ring-red-400/60" : "border-gray-600"
                            } ${locked ? "bg-gray-800 cursor-not-allowed" : "bg-gray-900"}`}
                        >
                            <option value="feature">Feature</option>
                            <option value="bug">Bug</option>
                            <option value="meeting">Meeting</option>
                            <option value="research">Research</option>
                        </select>
                        <span
                            className={`absolute right-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 rounded-sm border border-gray-600 ${
                                LABEL_COLORS[label] || "bg-gray-500"
                            }`}
                            aria-hidden="true"
                        />
                    </div>
                    <div className="grid grid-cols-2 gap-3">
                        <input
                            type="text"
                            inputMode="decimal"
                            placeholder="Hours (e.g. 2.75 or 2h 15m)"
                            value={hours}
                            onChange={(e) => onChangeHours(e.target.value)}
                            className={`px-3 py-2 rounded bg-gray-900 border text-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-400 ${
                                errors.hours ? "border-red-500 ring-red-400/60" : "border-gray-600"
                            }`}
                        />
                        <input
                            type="date"
                            value={date}
                            onChange={(e) => onChangeDate(e.target.value)}
                            className={`px-3 py-2 rounded bg-gray-900 border text-gray-100 focus:outline-none focus:ring-1 focus:ring-blue-400 date-input ${
                                errors.date ? "border-red-500 ring-red-400/60" : "border-gray-600"
                            }`}
                        />
                    </div>
                </div>
                <div className="flex items-center justify-between gap-3 mt-5">
                    {canDelete ? (
                        <button
                            className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded font-medium"
                            onClick={onDelete}
                        >
                            Delete
                        </button>
                    ) : (
                        <span />
                    )}
                    <div className="flex gap-3">
                        <button className="text-gray-300 hover:text-white" onClick={onCancel}>
                            Cancel
                        </button>
                        <button
                            className="bg-blue-600 px-4 py-2 rounded text-white font-medium hover:bg-blue-700"
                            onClick={onSave}
                        >
                            Save
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
