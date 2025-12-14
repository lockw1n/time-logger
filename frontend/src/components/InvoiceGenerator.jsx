import React, { useState } from "react";
import { generateMonthlyInvoicePdf } from "../api/invoices";

const toMonthInput = (date) => {
    const d = new Date(date);
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, "0");
    return `${year}-${month}`;
};

const defaultPreviousMonth = () => {
    const d = new Date();
    d.setDate(1);
    d.setMonth(d.getMonth() - 1);
    return toMonthInput(d);
};

const defaultForm = {
    month: defaultPreviousMonth(),
};

export default function InvoiceGenerator() {
    const [form, setForm] = useState(defaultForm);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    const handleChange = (key) => (e) => {
        const value = e.target.value;
        setForm((prev) => ({ ...prev, [key]: value }));
        setError("");
        setSuccess("");
    };

    const handleDownload = async () => {
        setLoading(true);
        setError("");
        setSuccess("");
        try {
            const { blob, filename, isInline } = await generateMonthlyInvoicePdf({ month: form.month });
            const file = new File([blob], filename, { type: "application/pdf" });
            const url = window.URL.createObjectURL(file);

            const cleanup = () => window.URL.revokeObjectURL(url);

            if (isInline) {
                const newTab = window.open(url, "_blank");
                if (!newTab) {
                    triggerDownload(url, filename);
                }
                setSuccess("Invoice opened in a new tab.");
                setTimeout(cleanup, 60_000);
            } else {
                triggerDownload(url, filename);
                setTimeout(cleanup, 5_000);
                setSuccess("Invoice download started.");
            }
        } catch (err) {
            console.error("invoice pdf generation failed", err);
            setError("Failed to generate invoice. Check required fields and try again.");
        } finally {
            setLoading(false);
        }
    };

    const triggerDownload = (href, filename) => {
        const a = document.createElement("a");
        a.href = href;
        a.download = filename;
        a.target = "_blank";
        document.body.appendChild(a);
        a.click();
        a.remove();
    };

    return (
        <div className="w-full max-w-6xl mt-8 flex justify-end">
            <div className="w-full sm:w-[47%] max-w-3xl bg-gray-800 border border-gray-700 rounded-xl p-4 shadow-lg">
            <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 mb-4">
                <div>
                    <h2 className="text-xl font-semibold text-gray-100">Invoice</h2>
                    <p className="text-sm text-gray-400">Generate a PDF for a selected month.</p>
                </div>
                <button
                    onClick={handleDownload}
                    disabled={loading}
                    className="px-4 py-2 rounded bg-blue-600 hover:bg-blue-700 disabled:opacity-60 text-white font-medium shadow-sm transition"
                >
                    {loading ? "Generating…" : "Generate PDF"}
                </button>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-[1fr_auto] gap-4 text-sm max-w-md">
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Month</label>
                    <input
                        type="month"
                        value={form.month}
                        onChange={handleChange("month")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100 date-input month-input"
                    />
                </div>
            </div>

            {error && <p className="text-red-400 text-sm mt-3">{error}</p>}
            {success && <p className="text-green-400 text-sm mt-3">{success}</p>}
            </div>
        </div>
    );
}
