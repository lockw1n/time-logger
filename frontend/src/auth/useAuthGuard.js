import { useEffect } from "react";
import { useAuth } from "./AuthContext";
import { redirectToLogin } from "./navigation";

export const useAuthGuard = (path) => {
    const { isAuthenticated, isInitializing } = useAuth();

    useEffect(() => {
        if (!isInitializing && !isAuthenticated && !path.startsWith("/login")) {
            redirectToLogin(path);
        }
    }, [isAuthenticated, isInitializing, path]);

    if (isInitializing) return false;
    return isAuthenticated;
};
