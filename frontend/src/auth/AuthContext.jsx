import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";
import { getMe, login as loginRequest } from "../api/auth";
import { clearStoredAuth, getStoredExpiresAt, getStoredToken, isTokenExpired, setStoredAuth } from "./storage";
import { redirectToLogin } from "./navigation";

const AuthContext = createContext(null);

const buildAuthState = (token, expiresAt, user = null) => ({
    token,
    expiresAt,
    user,
});

export function AuthProvider({ children }) {
    const [auth, setAuth] = useState(() => {
        const token = getStoredToken();
        const expiresAt = getStoredExpiresAt();
        if (token && !isTokenExpired(expiresAt)) {
            return buildAuthState(token, expiresAt);
        }
        clearStoredAuth();
        return buildAuthState(null, null);
    });
    const [isInitializing, setIsInitializing] = useState(true);

    const logout = useCallback(() => {
        clearStoredAuth();
        setAuth(buildAuthState(null, null));
        redirectToLogin();
    }, []);

    const login = useCallback(async ({ email, password }) => {
        const result = await loginRequest({ email, password });
        const token = result?.access_token;
        const expiresAt = result?.expires_at;
        if (!token || !expiresAt) {
            throw new Error("invalid login response");
        }
        setStoredAuth({ token, expiresAt });
        const user = await getMe();
        setAuth(buildAuthState(token, expiresAt, user));
        return user;
    }, []);

    useEffect(() => {
        const token = getStoredToken();
        const expiresAt = getStoredExpiresAt();
        if (!token) {
            setIsInitializing(false);
            return;
        }
        if (isTokenExpired(expiresAt)) {
            logout();
            setIsInitializing(false);
            return;
        }
        getMe()
            .then((user) => setAuth(buildAuthState(token, expiresAt, user)))
            .catch((error) => {
                if (error?.status === 401) {
                    logout();
                }
            })
            .finally(() => setIsInitializing(false));
    }, [logout]);

    useEffect(() => {
        if (!auth.token) return;
        if (isTokenExpired(auth.expiresAt)) {
            logout();
            return;
        }
        const expiryTime = new Date(auth.expiresAt).getTime();
        if (Number.isNaN(expiryTime)) {
            logout();
            return;
        }
        const timeoutMs = Math.max(0, expiryTime - Date.now());
        const timer = setTimeout(() => logout(), timeoutMs);
        return () => clearTimeout(timer);
    }, [auth.expiresAt, auth.token, logout]);

    const value = useMemo(
        () => ({
            login,
            logout,
            isAuthenticated: Boolean(auth.token && !isTokenExpired(auth.expiresAt)),
            user: auth.user,
            isInitializing,
        }),
        [auth.expiresAt, auth.token, auth.user, isInitializing, login, logout]
    );

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
};
