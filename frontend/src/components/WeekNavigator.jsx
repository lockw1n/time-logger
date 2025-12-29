import React from "react";

export default function WeekNavigator({ label, onPrev, onNext }) {
    return (
        <div className="flex items-center gap-3">
            <button
                className="px-3 py-1 rounded bg-gray-800 border border-gray-700 hover:bg-gray-700 transition"
                onClick={onPrev}
            >
                ← Previous week
            </button>
            <div className="px-4 py-1 rounded bg-gray-800 border border-gray-700 text-sm min-w-[220px] text-center">
                {label}
            </div>
            <button
                className="px-3 py-1 rounded bg-gray-800 border border-gray-700 hover:bg-gray-700 transition"
                onClick={onNext}
            >
                Next week →
            </button>
        </div>
    );
}
