import { useCallback, useEffect, useMemo, useState } from "react";
import { getTimesheet } from "../api/timesheet";
import { getStartOfWeek, getWeekDays, parseDMY, toYMD } from "../utils/date";

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

const buildRowsFromTimesheet = (rows) => {
    const tableRows = [];
    const totals = {};
    const ticketMap = new Map();

    rows.forEach((row) => {
        const ticketId = row.ticket?.id ?? null;
        const ticketCode = row.ticket?.code || "—";
        const cells = {};

        if (ticketId !== null) {
            const key = String(ticketId);
            if (!ticketMap.has(key)) {
                ticketMap.set(key, { id: ticketId, code: ticketCode });
            }
        }

        (row.entries || []).forEach((entry) => {
            const parsedDate = parseDMY(entry.date);
            if (!parsedDate) return;
            const dateKey = toYMD(parsedDate);
            cells[dateKey] = {
                id: entry.id,
                date: parsedDate,
                hours: (entry.duration_minutes || 0) / 60,
                comment: entry.comment,
            };
        });

        tableRows.push({
            ticket: { id: ticketId, code: ticketCode },
            activity: {
                id: row.activity?.id ?? null,
                name: row.activity?.name || "",
                color: (row.activity?.color || "").toLowerCase(),
            },
            color: (row.activity?.color || "").toLowerCase(),
            cells,
        });

        totals[ticketCode] = (totals[ticketCode] || 0) + (row.total || 0) / 60;
    });

    const ticketOptions = Array.from(ticketMap.values());

    return { tableRows, totals, ticketOptions };
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
    const [ticketOptions, setTicketOptions] = useState([]);

    useEffect(() => {
        setDays(getWeekDays(anchorDate, DAYS_WINDOW));
    }, [anchorDate]);

    const refresh = useCallback(
        async (anchor = anchorDate) => {
            const { startStr, endStr } = buildRangeParams(anchor);
            const data = await getTimesheet({ start: startStr, end: endStr });
            const { tableRows, totals, ticketOptions } = buildRowsFromTimesheet(data?.rows || []);
            setRows(tableRows);
            setTotals(totals);
            setTicketOptions(ticketOptions);
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
        ticketOptions,
        setAnchorDate,
        goToNextWeek,
        goToPreviousWeek,
    };
}
