import axios from "axios";
import { toYMD } from "../utils/date";

const API_URL = "/api/entries";

export async function createEntry(entry) {
    const dateStr = entry.date ? entry.date : toYMD(new Date());
    const payload = {
        consultant_id: 1,
        company_id: 1,
        ...entry,
        date: dateStr,
    };
    const res = await axios.post(API_URL, payload);
    return res.data;
}

export async function updateEntry(id, entry) {
    const payload = {};
    if (entry.duration_minutes !== undefined) payload.duration_minutes = entry.duration_minutes;
    if (entry.comment !== undefined) payload.comment = entry.comment;
    const res = await axios.put(`${API_URL}/${id}`, payload);
    return res.data;
}

export async function deleteEntry(id) {
    await axios.delete(`${API_URL}/${id}`);
}
