import React, { useState } from "react";
import InvoiceGeneratorModal from "./InvoiceGeneratorModal";

export default function InvoiceGenerator() {
    const [open, setOpen] = useState(false);

    return (
        <div className="w-full max-w-6xl mt-8 flex justify-end">
            <button
                className="px-4 py-2 rounded bg-blue-600 hover:bg-blue-700 text-white font-medium shadow-sm transition"
                onClick={() => setOpen(true)}
            >
                Generate invoice
            </button>
            <InvoiceGeneratorModal open={open} onCancel={() => setOpen(false)} />
        </div>
    );
}
