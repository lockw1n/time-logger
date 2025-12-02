import { useCallback, useEffect, useMemo, useState } from "react";
import { getEntries, getTicketSummary } from "../api/entries";
import { getStartOfWeek, getWeekDays, toYMD } from "../utils/date";

const DAYS_WINDOW = 14;

const buildRangeParams = (anchor) => {
    const start = new Date(anchor);
    start.setHours(0, 0, 0, 0);
    const end = new Date(anchor);
    end.setDate(end.getDate() + DAYS_WINDOW - 1);
    end.setHours(0, 0, 0, 0);

    return {
        startStr: toYMD(start),
        endStr: toYMD(end),
    };
};

const groupEntries = (entries) => {
    const grouped = {};
    entries.forEach((e) => {
        const key = toYMD(e.date || e.created_at);
        if (!grouped[e.ticket]) {
            grouped[e.ticket] = { ticket: e.ticket, label: e.label || "feature", cells: {} };
        }
        grouped[e.ticket].cells[key] = e;
    });
    return Object.values(grouped);
};

export function useTimesheet() {
    const initialAnchor = (() => {
        const currentWeek = getStartOfWeek(new Date());
        currentWeek.setDate(currentWeek.getDate() - 7); // show previous week on the left, current on the right
        return currentWeek;
    })();
    const [anchorDate, setAnchorDate] = useState(initialAnchor);
    const [days, setDays] = useState(() => getWeekDays(initialAnchor, DAYS_WINDOW));
    const [rows, setRows] = useState([]);
    const [totals, setTotals] = useState({});

    useEffect(() => {
        setDays(getWeekDays(anchorDate, DAYS_WINDOW));
    }, [anchorDate]);

    const refresh = useCallback(
        async (anchor = anchorDate) => {
            const { startStr, endStr } = buildRangeParams(anchor);
            const [data, summary] = await Promise.all([
                getEntries({ start: startStr, end: endStr }),
                getTicketSummary({ start: startStr, end: endStr }),
            ]);
            setRows(groupEntries(data));
            const totalsMap = {};
            summary.forEach((s) => {
                totalsMap[s.ticket] = s.total_hours;
            });
            setTotals(totalsMap);
        },
        [anchorDate]
    );

    useEffect(() => {
        refresh(anchorDate);
    }, [anchorDate, refresh]);

    const goToPreviousWeek = useCallback(() => {
        setAnchorDate((prev) => {
            const d = new Date(prev);
            d.setDate(d.getDate() - 7);
            return getStartOfWeek(d);
        });
    }, []);

    const goToNextWeek = useCallback(() => {
        setAnchorDate((prev) => {
            const d = new Date(prev);
            d.setDate(d.getDate() + 7);
            return getStartOfWeek(d);
        });
    }, []);

    const rangeLabel = useMemo(() => {
        if (!days.length) return "";
        const start = days[0];
        const end = days[days.length - 1];
        const format = (date) => {
            const day = String(date.getDate()).padStart(2, "0");
            const month = String(date.getMonth() + 1).padStart(2, "0");
            const year = date.getFullYear();
            return `${day}.${month}.${year}`;
        };
        return `${format(start)} – ${format(end)}`;
    }, [days]);

    return {
        anchorDate,
        days,
        rows,
        totals,
        rangeLabel,
        refresh,
        setAnchorDate,
        goToNextWeek,
        goToPreviousWeek,
    };
}
