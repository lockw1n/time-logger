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
    hourly_rate: 100,
    currency: "€",
    invoice_number: "20251201",
    order_number: "0000000000",
    consultant_name: "Oleksii Kotsiuba",
    consultant_address: "Ukraine, full address here…",
    consultant_tax_number: "0000000000",
    company_name: "Company Name",
    company_uid: "AT U000000000",
    company_street: "Industriestr. 1",
    company_city: "0000 city",
    company_country: "country",
    bank_name: "UKRSIBBANK",
    bank_address: "ANDRIIVSKA STREET 21/2 KYIV, UKRAINE",
    iban: "UA00000000000000000000000000",
    bic: "KHABUA2K",
    bank_country: "UKRAINE",
    payment_condition: "Net 14 days",
};

export default function InvoiceGenerator() {
    const [form, setForm] = useState(defaultForm);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    const payload = useMemo(() => {
        const hourly = Number(form.hourly_rate) || 0;
        return { ...form, hourly_rate: hourly };
    }, [form]);

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
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Hourly rate</label>
                    <input
                        type="number"
                        min="0"
                        step="0.01"
                        value={form.hourly_rate}
                        onChange={handleChange("hourly_rate")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Invoice number</label>
                    <input
                        type="text"
                        value={form.invoice_number}
                        onChange={handleChange("invoice_number")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Order number</label>
                    <input
                        type="text"
                        value={form.order_number}
                        onChange={handleChange("order_number")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Consultant name</label>
                    <input
                        type="text"
                        value={form.consultant_name}
                        onChange={handleChange("consultant_name")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Consultant address</label>
                    <input
                        type="text"
                        value={form.consultant_address}
                        onChange={handleChange("consultant_address")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Company name</label>
                    <input
                        type="text"
                        value={form.company_name}
                        onChange={handleChange("company_name")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Company UID</label>
                    <input
                        type="text"
                        value={form.company_uid}
                        onChange={handleChange("company_uid")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Company street</label>
                    <input
                        type="text"
                        value={form.company_street}
                        onChange={handleChange("company_street")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Company city</label>
                    <input
                        type="text"
                        value={form.company_city}
                        onChange={handleChange("company_city")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Company country</label>
                    <input
                        type="text"
                        value={form.company_country}
                        onChange={handleChange("company_country")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Payment condition</label>
                    <input
                        type="text"
                        value={form.payment_condition}
                        onChange={handleChange("payment_condition")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Bank name</label>
                    <input
                        type="text"
                        value={form.bank_name}
                        onChange={handleChange("bank_name")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Bank address</label>
                    <input
                        type="text"
                        value={form.bank_address}
                        onChange={handleChange("bank_address")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">IBAN</label>
                    <input
                        type="text"
                        value={form.iban}
                        onChange={handleChange("iban")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">BIC</label>
                    <input
                        type="text"
                        value={form.bic}
                        onChange={handleChange("bic")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
                <div className="flex flex-col gap-2">
                    <label className="text-gray-300">Bank country</label>
                    <input
                        type="text"
                        value={form.bank_country}
                        onChange={handleChange("bank_country")}
                        className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100"
                    />
                </div>
            </div>

            {error && <p className="text-red-400 text-sm mt-3">{error}</p>}
            {success && <p className="text-green-400 text-sm mt-3">{success}</p>}
        </div>
    );
}
