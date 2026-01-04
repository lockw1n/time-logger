import React from "react";
import { useAuthGuard } from "../auth/useAuthGuard";
import { useAuth } from "../auth/AuthContext";
import TimesheetPage from "./TimesheetPage";

export default function ProtectedApp({ path }) {
    const canRenderProtected = useAuthGuard(path);
    if (!canRenderProtected) return null;
    const { isAuthenticated, logout } = useAuth();
    return (
        <div className="min-h-screen">
            {isAuthenticated ? (
                <div className="flex items-center justify-end px-6 pt-6">
                    <button
                        className="px-4 py-2 rounded bg-gray-700 hover:bg-gray-600 text-white font-medium shadow-sm transition"
                        onClick={logout}
                    >
                        Logout
                    </button>
                </div>
            ) : null}
            <TimesheetPage />
        </div>
    );
}
