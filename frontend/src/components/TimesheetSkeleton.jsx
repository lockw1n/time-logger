import React from "react";

export default function TimesheetSkeleton() {
    return (
        <div className="shadow-lg rounded-xl bg-gray-800 border border-gray-700 min-h-[260px] animate-pulse">
            <div className="px-4 py-3 border-b border-gray-700">
                <div className="flex items-center gap-4">
                    <div className="h-4 w-24 bg-gray-700 rounded" />
                    <div className="h-4 w-8 bg-gray-700 rounded" />
                    <div className="flex-1 grid grid-cols-7 gap-2">
                        {Array.from({ length: 7 }).map((_, idx) => (
                            <div key={`th-${idx}`} className="h-4 bg-gray-700 rounded" />
                        ))}
                    </div>
                </div>
            </div>
            <div className="px-4 py-3 space-y-3">
                {Array.from({ length: 6 }).map((_, rowIdx) => (
                    <div key={`row-${rowIdx}`} className="flex items-center gap-4">
                        <div className="h-4 w-24 bg-gray-700 rounded" />
                        <div className="h-4 w-6 bg-gray-700 rounded" />
                        <div className="flex-1 grid grid-cols-7 gap-2">
                            {Array.from({ length: 7 }).map((_, cellIdx) => (
                                <div
                                    key={`cell-${rowIdx}-${cellIdx}`}
                                    className="h-4 bg-gray-700 rounded"
                                />
                            ))}
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}
