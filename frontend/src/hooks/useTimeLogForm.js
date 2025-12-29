import { useCallback, useEffect, useMemo, useState } from "react";
import { createEntry, updateEntry, deleteEntry } from "../api/entries";
import { toYMD } from "../utils/date";
import { parseHoursInput } from "../utils/time";

const getFallbackDate = (defaultDate) => {
    const base = defaultDate ?? new Date();
    return toYMD(base);
};

export function useTimeLogForm({ activities, onSaved, defaultDate }) {
    const [isOpen, setIsOpen] = useState(false);
    const [activeEntryId, setActiveEntryId] = useState(null);
    const [ticket, setTicket] = useState("");
    const [activityId, setActivityId] = useState(null);
    const [hours, setHours] = useState("");
    const [date, setDate] = useState(getFallbackDate(defaultDate));
    const [errors, setErrors] = useState({});
    const [locked, setLocked] = useState(false);

    const activityOptions = activities || [];
    const defaultActivityId = useMemo(
        () =>
            activityOptions.find((option) => option.id !== null && option.id !== undefined)?.id ??
            null,
        [activityOptions]
    );

    const resetFormState = useCallback(() => {
        setActiveEntryId(null);
        setTicket("");
        setActivityId(defaultActivityId);
        setHours("");
        setDate(getFallbackDate(defaultDate));
        setErrors({});
        setLocked(false);
    }, [defaultActivityId, defaultDate]);

    const openNew = useCallback(() => {
        resetFormState();
        setIsOpen(true);
    }, [resetFormState]);

    const openFromCell = useCallback(
        ({ ticket: selectedTicket, activity, date: selectedDate, entry }) => {
            const resolvedActivityId = activity?.id ?? defaultActivityId;
            const resolvedTicketCode = selectedTicket?.code || "";
            setActiveEntryId(entry?.id ?? null);
            setTicket(resolvedTicketCode);
            setActivityId(resolvedActivityId);
            setHours(entry ? String(entry.hours ?? "") : "");
            setDate(toYMD(selectedDate || new Date()));
            setErrors({});
            setLocked(true);
            setIsOpen(true);
        },
        [defaultActivityId]
    );

    const close = useCallback(() => {
        setIsOpen(false);
        resetFormState();
    }, [resetFormState]);

    useEffect(() => {
        if (!isOpen) return;
        if (activityId != null) return;
        if (defaultActivityId == null) return;
        setActivityId(defaultActivityId);
    }, [activityId, defaultActivityId, isOpen]);

    const handleSave = useCallback(async () => {
        const nextErrors = {};
        const trimmedTicket = (ticket || "").trim();
        const selectedActivity = activityOptions.find(
            (option) => String(option.id) === String(activityId)
        );
        const hoursNum = parseHoursInput(hours);
        const hasDate = Boolean(date);
        const hasActivity = Boolean(selectedActivity?.id);

        if (!trimmedTicket) nextErrors.ticket = true;
        if (!hasActivity) nextErrors.activity = true;
        if (Number.isNaN(hoursNum) || hoursNum < 0 || hoursNum > 24) nextErrors.hours = true;
        if (!hasDate) nextErrors.date = true;

        if (Object.keys(nextErrors).length) {
            setErrors(nextErrors);
            return;
        }

        setErrors({});
        const dateISO = (date && date.trim()) || toYMD(new Date());
        const durationMinutes = Math.round(hoursNum * 60);
        if (durationMinutes % 15 !== 0) {
            setErrors({ hours: true });
            return;
        }

        if (activeEntryId) {
            await updateEntry(activeEntryId, { duration_minutes: durationMinutes });
        } else {
            await createEntry({
                ticket_code: trimmedTicket,
                activity_id: selectedActivity?.id,
                date: dateISO,
                duration_minutes: durationMinutes,
            });
        }

        setIsOpen(false);
        resetFormState();
        if (onSaved) {
            onSaved();
        }
    }, [
        activityId,
        activityOptions,
        activeEntryId,
        date,
        hours,
        onSaved,
        resetFormState,
        ticket,
    ]);

    const handleDelete = useCallback(async () => {
        if (!activeEntryId) return;
        await deleteEntry(activeEntryId);
        setIsOpen(false);
        resetFormState();
        if (onSaved) {
            onSaved();
        }
    }, [activeEntryId, onSaved, resetFormState]);

    return {
        openNew,
        openFromCell,
        close,
        isOpen,
        isEdit: activeEntryId !== null,
        activeEntryId,
        modalProps: {
            open: isOpen,
            ticket,
            activityId,
            activityOptions,
            hours,
            date,
            errors,
            locked,
            canDelete: Boolean(activeEntryId),
            onChangeTicket: (val) => {
                setTicket(val);
                setErrors((prev) => ({ ...prev, ticket: false }));
            },
            onChangeActivity: (val) => {
                setActivityId(val);
                setErrors((prev) => ({ ...prev, activity: false }));
            },
            onChangeHours: (val) => {
                setHours(val);
                setErrors((prev) => ({ ...prev, hours: false }));
            },
            onChangeDate: (val) => {
                setDate(val);
                setErrors((prev) => ({ ...prev, date: false }));
            },
            onCancel: close,
            onSave: handleSave,
            onDelete: handleDelete,
        },
    };
}
