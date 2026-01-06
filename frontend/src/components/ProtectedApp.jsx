import React from "react";
import { useAuthGuard } from "../auth/useAuthGuard";
import { useAuth } from "../auth/AuthContext";

export default function ProtectedApp({ path, children }) {
    const canRenderProtected = useAuthGuard(path);
    if (!canRenderProtected) return null;
    const { isAuthenticated, logout, user } = useAuth();
    const firstName = (user?.first_name || user?.firstName || "").trim();
    const lastName = (user?.last_name || user?.lastName || "").trim();
    const fullName = [firstName, lastName].filter(Boolean).join(" ");
    return (
        <div className="min-h-screen">
            {isAuthenticated ? (
                <div className="flex items-center justify-between px-6 pt-6">
                    <button
                        className="text-gray-300 hover:text-white transition"
                        onClick={() => {
                            window.history.pushState(null, "", "/");
                            window.dispatchEvent(new PopStateEvent("popstate"));
                        }}
                    >
                        Timesheet
                    </button>
                    <div className="flex items-center gap-4">
                        {fullName ? (
                            <button
                                className="text-gray-300 hover:text-white transition"
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
                            onClick={logout}
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
