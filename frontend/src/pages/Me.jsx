import React from "react";
import { useAuth } from "../auth/AuthContext";

const renderField = (label, value) => (
    <div className="flex flex-col gap-1">
        <span className="text-xs uppercase tracking-wide text-gray-400">{label}</span>
        <span className="text-sm text-gray-100">{value || "—"}</span>
    </div>
);

export default function MePage() {
    const { user } = useAuth();
    const regionValue = user?.region || user?.region_name || user?.regionName;

    return (
        <div className="min-h-screen p-6 flex flex-col items-center relative">
            <h1 className="text-3xl font-bold mb-6 text-gray-100 flex items-center gap-2">
                ⏱️ Time Logger <span className="text-gray-400">– Profile</span>
            </h1>

            <div className="w-full max-w-6xl">
                <div className="bg-gray-800 border border-gray-700 rounded-xl shadow-lg p-6 max-w-lg">
                    <div className="flex flex-col gap-4">
                        {renderField("First name", user?.first_name || user?.firstName)}
                        {renderField("Last name", user?.last_name || user?.lastName)}
                        {renderField("Email", user?.email)}
                        {renderField("Country", user?.country)}
                        {renderField("City", user?.city)}
                        {regionValue ? renderField("Region", regionValue) : null}
                        {renderField("ID", user?.id != null ? String(user.id) : "")}
                    </div>
                </div>
            </div>
        </div>
    );
}
