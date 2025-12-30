import React, { useEffect, useState } from "react";
import { generateMonthlyInvoiceExcel, generateMonthlyInvoicePdf } from "../api/invoices";

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

export default function InvoiceGeneratorModal({ open, onCancel }) {
    const [form, setForm] = useState(defaultForm);
    const [loadingPdf, setLoadingPdf] = useState(false);
    const [loadingExcel, setLoadingExcel] = useState(false);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    const isBusy = loadingPdf || loadingExcel;
    const isMonthValid = Boolean(form.month);

    const mapInvoiceError = (err) => {
        if (err?.code === "invalid_month") {
            return "Select a valid month to generate the invoice.";
        }
        const status = err?.response?.status;
        if (status === 400) return "Invoice request is invalid. Check the selected month.";
        if (status === 409) return "Invoice could not be generated for this period.";
        return "Failed to generate invoice. Check required fields and try again.";
    };

    const createBlobUrl = (blob, filename, type) => {
        const file = new File([blob], filename, { type });
        return window.URL.createObjectURL(file);
    };

    const revokeLater = (url, delayMs) => {
        window.setTimeout(() => window.URL.revokeObjectURL(url), delayMs);
    };

    const handleChange = (key) => (e) => {
        const value = e.target.value;
        setForm((prev) => ({ ...prev, [key]: value }));
        setError("");
        setSuccess("");
    };

    const handleDownloadPdf = async () => {
        setLoadingPdf(true);
        setError("");
        setSuccess("");
        try {
            if (!isMonthValid) {
                throw { code: "invalid_month" };
            }
            const { blob, filename, isInline } = await generateMonthlyInvoicePdf({ month: form.month });
            const url = createBlobUrl(blob, filename, "application/pdf");

            if (isInline) {
                const newTab = window.open(url, "_blank");
                if (!newTab) {
                    triggerDownload(url, filename);
                }
                setSuccess("Invoice opened in a new tab.");
                revokeLater(url, 60_000);
            } else {
                triggerDownload(url, filename);
                revokeLater(url, 5_000);
                setSuccess("Invoice download started.");
            }
        } catch (err) {
            console.error("invoice pdf generation failed", err);
            setError(mapInvoiceError(err));
        } finally {
            setLoadingPdf(false);
        }
    };

    const handleDownloadExcel = async () => {
        setLoadingExcel(true);
        setError("");
        setSuccess("");
        try {
            if (!isMonthValid) {
                throw { code: "invalid_month" };
            }
            const { blob, filename } = await generateMonthlyInvoiceExcel({ month: form.month });
            const url = createBlobUrl(
                blob,
                filename,
                "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
            );
            triggerDownload(url, filename);
            revokeLater(url, 5_000);
            setSuccess("Invoice Excel download started.");
        } catch (err) {
            console.error("invoice excel generation failed", err);
            setError(mapInvoiceError(err));
        } finally {
            setLoadingExcel(false);
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
                <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 mb-4">
                    <div>
                        <h2 className="text-xl font-semibold text-gray-100">Invoice</h2>
                        <p className="text-sm text-gray-400">Generate an invoice for a selected month.</p>
                    </div>
                    <div className="flex flex-wrap gap-3 justify-end md:ml-auto mt-5 md:mt-0">
                        <button
                            onClick={handleDownloadPdf}
                            disabled={isBusy || !isMonthValid}
                            className="px-4 py-2 rounded bg-blue-600 hover:bg-blue-700 disabled:opacity-60 text-white font-medium shadow-sm transition w-[160px]"
                        >
                            {loadingPdf ? "Generating…" : "Generate PDF"}
                        </button>
                        <button
                            onClick={handleDownloadExcel}
                            disabled={isBusy || !isMonthValid}
                            className="px-4 py-2 rounded bg-emerald-600 hover:bg-emerald-700 disabled:opacity-60 text-white font-medium shadow-sm transition w-[160px]"
                        >
                            {loadingExcel ? "Generating…" : "Generate Excel"}
                        </button>
                    </div>
                </div>

                <div className="grid grid-cols-1 gap-4 text-sm max-w-xs">
                    <div className="flex flex-col gap-2">
                        <label className="text-gray-300">Month</label>
                        <input
                            type="month"
                            value={form.month}
                            onChange={handleChange("month")}
                            disabled={isBusy}
                            className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100 date-input month-input w-full"
                        />
                    </div>
                </div>

                {error && <p className="text-red-400 text-sm mt-3">{error}</p>}
                {success && <p className="text-green-400 text-sm mt-3">{success}</p>}

                <div className="flex justify-end mt-5">
                    <button className="text-gray-300 hover:text-white" onClick={onCancel}>
                        Cancel
                    </button>
                </div>
            </div>
        </div>
    );
}
