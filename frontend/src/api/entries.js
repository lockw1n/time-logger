import { apiDelete, apiPost, apiPut } from "./client";
import { toYMD } from "../utils/date";

const API_URL = "/api/entries";

export async function createEntry(entry) {
    const dateStr = entry.date ? entry.date : toYMD(new Date());
    const payload = {
        company_id: 1,
        ...entry,
        date: dateStr,
    };
    return apiPost(API_URL, {
        body: JSON.stringify(payload),
        headers: { "Content-Type": "application/json" },
    });
}

export async function updateEntry(id, entry) {
    const payload = {};
    if (entry.duration_minutes !== undefined) payload.duration_minutes = entry.duration_minutes;
    if (entry.comment !== undefined) payload.comment = entry.comment;
    return apiPut(`${API_URL}/${id}`, {
        body: JSON.stringify(payload),
        headers: { "Content-Type": "application/json" },
    });
}

export async function deleteEntry(id) {
    await apiDelete(`${API_URL}/${id}`);
}
