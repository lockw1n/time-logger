import React from "react";
import { useAuthGuard } from "../auth/useAuthGuard";
import { useAuth } from "../auth/AuthContext";
import { useCompany } from "../context/CompanyContext";
import ledvixWordmark from "../assets/brand/ledvix-wordmark.svg";

export default function ProtectedApp({ path, children }) {
    const canRenderProtected = useAuthGuard(path);
    if (!canRenderProtected) return null;
    const { isAuthenticated, logout, user } = useAuth();
    const { companies, selectedCompanyId, isSwitchingCompany, switchCompany, resetCompany } = useCompany();
    const firstName = (user?.first_name || user?.firstName || "").trim();
    const lastName = (user?.last_name || user?.lastName || "").trim();
    const fullName = [firstName, lastName].filter(Boolean).join(" ");
    return (
        <div className="min-h-screen">
            {isAuthenticated ? (
                <div className="flex items-center px-6 py-4">
                    <div className="flex items-center gap-3">
                        <button
                            className="flex items-center cursor-pointer transition hover:opacity-90"
                            onClick={() => {
                                window.history.pushState(null, "", "/");
                                window.dispatchEvent(new PopStateEvent("popstate"));
                            }}
                        >
                            <img
                                src={ledvixWordmark}
                                alt="Ledvix"
                                className="h-8 w-auto"
                            />
                        </button>
                        {selectedCompanyId !== null ? (
                            <select
                                id="header-company-selection"
                                name="header-company-selection"
                                className="bg-gray-900 border border-gray-700 rounded px-2 py-0.5 text-[11px] text-gray-400"
                                value={selectedCompanyId || ""}
                                disabled={isSwitchingCompany}
                                onChange={(event) =>
                                    switchCompany(event.target.value ? Number(event.target.value) : null)
                                }
                            >
                                <option value="">Select company</option>
                                {companies.map((company) => (
                                    <option key={company.id} value={company.id}>
                                        {company.name || company.name_short || company.nameShort || "Company"}
                                    </option>
                                ))}
                            </select>
                        ) : null}
                    </div>
                    <div className="flex-1" />
                    <div className="flex items-center gap-4">
                        {fullName ? (
                            <button
                                className="text-gray-400 hover:text-gray-200 transition"
                                onClick={() => {
                                    window.history.pushState(null, "", "/me");
                                    window.dispatchEvent(new PopStateEvent("popstate"));
                                }}
                            >
                                {fullName}
                            </button>
                        ) : null}
                        <button
                            className="px-4 py-2 rounded bg-gray-700 hover:bg-gray-600 text-white font-medium shadow-sm transition"
                            onClick={() => {
                                resetCompany();
                                logout();
                            }}
                        >
                            Logout
                        </button>
                    </div>
                </div>
            ) : null}
            {children}
        </div>
    );
}
