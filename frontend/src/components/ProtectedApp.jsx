import React from "react";
import { useAuthGuard } from "../auth/useAuthGuard";
import TimesheetPage from "./TimesheetPage";

export default function ProtectedApp({ path }) {
    const canRenderProtected = useAuthGuard(path);
    if (!canRenderProtected) return null;
    return <TimesheetPage />;
}
