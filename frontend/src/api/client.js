import { clearStoredAuth, getStoredExpiresAt, getStoredToken, isTokenExpired } from "../auth/storage";
import { redirectToLogin } from "../auth/navigation";

const buildQuery = (params = {}) => {
    const entries = Object.entries(params).filter(([, value]) => value !== undefined);
    if (!entries.length) return "";
    return `?${new URLSearchParams(entries).toString()}`;
};

const getAuthHeaders = () => {
    const token = getStoredToken();
    const expiresAt = getStoredExpiresAt();
    if (!token) return {};
    if (isTokenExpired(expiresAt)) {
        clearStoredAuth();
        redirectToLogin();
        return {};
    }
    return { Authorization: `Bearer ${token}` };
};

const handleUnauthorized = () => {
    clearStoredAuth();
    redirectToLogin();
};

const request = async (path, { method = "GET", params, body, headers, responseType = "json" } = {}) => {
    const url = `${path}${buildQuery(params)}`;
    const response = await fetch(url, {
        method,
        headers: {
            ...getAuthHeaders(),
            ...headers,
        },
        body,
    });

    if (response.status === 401) {
        handleUnauthorized();
    }

    if (!response.ok) {
        let errorPayload = null;
        try {
            errorPayload = await response.json();
        } catch {
            errorPayload = null;
        }
        const error = new Error(errorPayload?.error || "request failed");
        error.status = response.status;
        error.payload = errorPayload;
        throw error;
    }

    if (responseType === "blob") {
        return {
            data: await response.blob(),
            headers: response.headers,
        };
    }

    if (response.status === 204) return null;
    return response.json();
};

export const apiGet = (path, options) => request(path, { ...options, method: "GET" });
export const apiPost = (path, options) => request(path, { ...options, method: "POST" });
export const apiPut = (path, options) => request(path, { ...options, method: "PUT" });
export const apiDelete = (path, options) => request(path, { ...options, method: "DELETE" });
