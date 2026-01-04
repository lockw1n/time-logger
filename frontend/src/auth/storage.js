const TOKEN_KEY = "auth.token";
const EXPIRES_AT_KEY = "auth.expires_at";

export const getStoredToken = () => localStorage.getItem(TOKEN_KEY);

export const getStoredExpiresAt = () => localStorage.getItem(EXPIRES_AT_KEY);

export const setStoredAuth = ({ token, expiresAt }) => {
    localStorage.setItem(TOKEN_KEY, token);
    localStorage.setItem(EXPIRES_AT_KEY, expiresAt);
};

export const clearStoredAuth = () => {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(EXPIRES_AT_KEY);
};

export const isTokenExpired = (expiresAt) => {
    if (!expiresAt) return true;
    const expiryTime = new Date(expiresAt).getTime();
    if (Number.isNaN(expiryTime)) return true;
    return Date.now() >= expiryTime;
};
