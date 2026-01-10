import React from "react";
import { useCompany } from "../context/CompanyContext";

export default function CompanySelectionGate() {
    const { companies, selectedCompanyId, switchCompany } = useCompany();

    return (
        <div className="w-full max-w-lg rounded-xl border border-gray-700 bg-gray-800 p-6 text-gray-100 shadow-lg">
            <h2 className="text-3xl font-bold mb-2 text-center">Select a company</h2>
            <p className="text-sm text-gray-400 mb-4 text-center">
                Choose a company to load your timesheet and start logging time.
            </p>
            <select
                id="company-selection"
                name="company-selection"
                className="w-full bg-gray-900 border border-gray-700 rounded px-3 py-2 pr-10 text-gray-100"
                value={selectedCompanyId ?? ""}
                onChange={(event) =>
                    switchCompany(event.target.value ? Number(event.target.value) : null)
                }
                disabled={companies.length === 0}
            >
                {companies.length === 0 ? (
                    <option value="" disabled>
                        No companies available
                    </option>
                ) : (
                    <option value="">Select company</option>
                )}
                {companies.map((company) => (
                    <option key={company.id} value={company.id}>
                        {company.name || company.name_short || company.nameShort || "Company"}
                    </option>
                ))}
            </select>
        </div>
    );
}
