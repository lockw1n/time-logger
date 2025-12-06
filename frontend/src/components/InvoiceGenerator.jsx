import React, { useMemo, useState } from "react";
import { downloadInvoicePdf } from "../api/reports";

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
    currency: "€",
};

export default function InvoiceGenerator() {
    const [form, setForm] = useState(defaultForm);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    const payload = useMemo(() => ({ ...form }), [form]);

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
            const blobData = await downloadInvoicePdf(payload);
            const url = window.URL.createObjectURL(new Blob([blobData], { type: "application/pdf" }));
            const a = document.createElement("a");
            a.href = url;
            a.download = `invoice-${payload.month}.pdf`;
            document.body.appendChild(a);
            a.click();
            a.remove();
            window.URL.revokeObjectURL(url);
            setSuccess("Invoice PDF downloaded.");
        } catch (err) {
            console.error("invoice pdf generation failed", err);
            setError("Failed to generate invoice. Check required fields and try again.");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="w-full max-w-6xl mt-8 bg-gray-800 border border-gray-700 rounded-xl p-4 shadow-lg">
            <div className="flex items-center justify-between mb-4">
                <div>
                    <h2 className="text-xl font-semibold text-gray-100">Invoice</h2>
                    <p className="text-sm text-gray-400">Generate a PDF for a selected month.</p>
                </div>
                <button
                    onClick={handleDownload}
                    disabled={loading}
                    className="px-4 py-2 rounded bg-blue-600 hover:bg-blue-700 disabled:opacity-60 text-white font-medium shadow-sm transition"
                >
                    {loading ? "Generating…" : "Download PDF"}
                </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Month</label>
                    <input
                        type="month"
                        value={form.month}
                        onChange={handleChange("month")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100 date-input month-input"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Currency</label>
                    <input
                        type="text"
                        value={form.currency}
                        onChange={handleChange("currency")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
            </div>

            {error && <p className="text-red-400 text-sm mt-3">{error}</p>}
            {success && <p className="text-green-400 text-sm mt-3">{success}</p>}
        </div>
    );
}
