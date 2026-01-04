import { apiGet, apiPost } from "./client";

const LOGIN_URL = "/api/auth/login";
const ME_URL = "/api/me";

export const login = async ({ email, password }) =>
    apiPost(LOGIN_URL, {
        body: JSON.stringify({ email, password }),
        headers: { "Content-Type": "application/json" },
    });

export const getMe = async () => apiGet(ME_URL);
