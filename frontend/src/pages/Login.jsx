import React, { useMemo, useState } from "react";
import { useAuth } from "../auth/AuthContext";
import ledvixWordmark from "../assets/brand/ledvix-wordmark.svg";

const getRedirectTarget = (search) => {
    const params = new URLSearchParams(search);
    return params.get("redirect") || "/";
};

export default function Login() {
    const { login } = useAuth();
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    const redirectTo = useMemo(
        () => getRedirectTarget(window.location.search),
        []
    );

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError("");
        setLoading(true);
        try {
            await login({ email, password });
            window.history.replaceState(null, "", redirectTo);
            window.dispatchEvent(new PopStateEvent("popstate"));
        } catch (err) {
            setError(err?.payload?.error || err?.message || "Login failed");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen p-6 flex flex-col items-center">
            <img
                src={ledvixWordmark}
                alt="Ledvix"
                className="h-10 w-auto mb-8"
            />

            <div className="w-full max-w-6xl flex-1 flex items-center justify-center">
                <div className="bg-gray-800 border border-gray-700 rounded-xl shadow-lg p-6 max-w-lg w-full">
                    <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
                        <div className="flex flex-col gap-2">
                            <label className="text-gray-300" htmlFor="login-email">
                                Email
                            </label>
                            <input
                                id="login-email"
                                name="email"
                                autoComplete="email"
                                type="email"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                disabled={loading}
                                className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100 w-full"
                                required
                            />
                        </div>
                        <div className="flex flex-col gap-2">
                            <label className="text-gray-300" htmlFor="login-password">
                                Password
                            </label>
                            <input
                                id="login-password"
                                name="password"
                                autoComplete="current-password"
                                type="password"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                disabled={loading}
                                className="bg-gray-900 border border-gray-700 rounded px-3 py-2 text-gray-100 w-full"
                                required
                            />
                        </div>
                        {error && <p className="text-red-400 text-sm">{error}</p>}
                        <button
                            type="submit"
                            disabled={loading}
                            className="px-4 py-2 rounded bg-blue-600 hover:bg-blue-700 disabled:opacity-60 text-white font-medium shadow-sm transition"
                        >
                            {loading ? "Signing in…" : "Sign in"}
                        </button>
                    </form>
                </div>
            </div>
        </div>
    );
}
