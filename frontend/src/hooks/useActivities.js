import { useCallback, useEffect, useState } from "react";
import { listActivitiesForCompany } from "../api/activities";

export function useActivities(companyId) {
    const [activities, setActivities] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    const refresh = useCallback(async () => {
        if (!companyId) return;
        setLoading(true);
        setError("");
        try {
            const data = await listActivitiesForCompany(companyId);
            const normalized = (data || []).map((activity) => ({
                id: activity.id,
                name: activity.name,
                color: activity.color || "",
                priority: activity.priority ?? 0,
            }));
            normalized.sort((a, b) => a.priority - b.priority || a.name.localeCompare(b.name));
            setActivities(normalized);
        } catch (err) {
            setActivities([]);
            setError("Failed to load activities.");
            // eslint-disable-next-line no-console
            console.error("activities fetch failed", err);
        } finally {
            setLoading(false);
        }
    }, [companyId]);

    useEffect(() => {
        refresh();
    }, [refresh]);

    return {
        activities,
        loading,
        error,
        refresh,
    };
}
