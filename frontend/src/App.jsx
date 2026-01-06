import React, { useEffect, useState } from "react";
import Login from "./pages/Login";
import { useAuth } from "./auth/AuthContext";
import ProtectedApp from "./components/ProtectedApp";
import TimesheetPage from "./components/TimesheetPage";
import MePage from "./pages/Me";

export default function App() {
    const { isAuthenticated, isInitializing } = useAuth();
    const [path, setPath] = useState(
        `${window.location.pathname}${window.location.search}`
    );

    useEffect(() => {
        const onPopState = () => {
            setPath(`${window.location.pathname}${window.location.search}`);
        };
        window.addEventListener("popstate", onPopState);
        return () => window.removeEventListener("popstate", onPopState);
    }, []);

    if (path.startsWith("/login")) {
        if (isAuthenticated) {
            window.history.replaceState(null, "", "/");
            window.dispatchEvent(new PopStateEvent("popstate"));
            return null;
        }
        return <Login />;
    }

    if (isInitializing) {
        return (
            <div className="min-h-screen flex items-center justify-center text-gray-400">
                Loading…
            </div>
        );
    }

    const isProfileRoute = path.startsWith("/me");
    return (
        <ProtectedApp path={path}>
            {isProfileRoute ? <MePage /> : <TimesheetPage />}
        </ProtectedApp>
    );
}
