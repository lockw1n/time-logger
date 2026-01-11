import React from "react";

export default function InvoiceGeneratorSkeleton() {
    return (
        <div className="animate-pulse">
            <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 mb-4">
                <div>
                    <div className="h-5 w-24 bg-gray-700 rounded" />
                    <div className="h-4 w-56 bg-gray-700 rounded mt-2" />
                </div>
                <div className="flex flex-wrap gap-3 justify-end md:ml-auto mt-5 md:mt-0">
                    <div className="h-10 w-[160px] bg-gray-700 rounded" />
                    <div className="h-10 w-[160px] bg-gray-700 rounded" />
                </div>
            </div>
            <div className="grid grid-cols-1 gap-4 text-sm max-w-xs">
                <div className="flex flex-col gap-2">
                    <div className="h-4 w-16 bg-gray-700 rounded" />
                    <div className="h-10 w-full bg-gray-700 rounded" />
                </div>
            </div>
            <div className="mt-5 h-9 w-20 bg-gray-700 rounded ml-auto" />
        </div>
    );
}
