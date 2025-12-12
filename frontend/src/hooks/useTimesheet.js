import { useCallback, useEffect, useMemo, useState } from "react";
import { getTimesheet } from "../api/timesheet";
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

const buildRowsFromTimesheet = (rows) => {
    const tableRows = [];
    const totals = {};
    const labelMap = new Map();
    const ticketMap = new Map();

    rows.forEach((row) => {
        const ticketId = row.ticket?.id ?? null;
        const ticketCode = row.ticket?.code || "—";
        const labelId = row.label?.id ?? null;
        const labelName = row.label?.name || "";
        const labelColor = (row.label?.color || "").toLowerCase();
        const cells = {};

        if (ticketId !== null) {
            const key = String(ticketId);
            if (!ticketMap.has(key)) {
                ticketMap.set(key, { id: ticketId, code: ticketCode });
            }
        }

        if (labelId !== null) {
            const key = String(labelId);
            if (!labelMap.has(key)) {
                labelMap.set(key, { id: labelId, name: labelName, color: labelColor });
            }
        }

        (row.entries || []).forEach((entry) => {
            const dateKey = toYMD(entry.date);
            cells[dateKey] = {
                id: entry.id,
                date: entry.date,
                hours: (entry.duration_minutes || 0) / 60,
                comment: entry.comment,
            };
        });

        tableRows.push({
            ticket: { id: ticketId, code: ticketCode },
            label: { id: labelId, name: labelName, color: labelColor },
            color: labelColor,
            cells,
        });

        totals[ticketCode] = (row.total || 0) / 60;
    });

    const labelOptions = Array.from(labelMap.values());
    if (!labelOptions.length) {
        labelOptions.push({ id: null, name: "No label", color: "gray" });
    }

    const ticketOptions = Array.from(ticketMap.values());

    return { tableRows, totals, labelOptions, ticketOptions };
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
    const [labelOptions, setLabelOptions] = useState([]);
    const [ticketOptions, setTicketOptions] = useState([]);

    useEffect(() => {
        setDays(getWeekDays(anchorDate, DAYS_WINDOW));
    }, [anchorDate]);

    const refresh = useCallback(
        async (anchor = anchorDate) => {
            const { startStr, endStr } = buildRangeParams(anchor);
            const data = await getTimesheet({ start: startStr, end: endStr });
            const { tableRows, totals, labelOptions, ticketOptions } = buildRowsFromTimesheet(data?.rows || []);
            setRows(tableRows);
            setTotals(totals);
            setLabelOptions(labelOptions);
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
        labelOptions,
        ticketOptions,
        setAnchorDate,
        goToNextWeek,
        goToPreviousWeek,
    };
}
