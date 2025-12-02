import axios from "axios";

const API_URL = "/api/entries";
const SUMMARY_URL = "/api/tickets/summary";

export async function getEntries({ start, end } = {}) {
    const res = await axios.get(API_URL, {
        params: start && end ? { start, end } : {},
    });
    return res.data;
}

export async function createEntry(entry) {
    const dateStr = entry.date
        ? entry.date
        : new Date().toISOString().split("T")[0]; // YYYY-MM-DD in UTC
    const payload = { ...entry, date: dateStr };
    const res = await axios.post(API_URL, payload);
    return res.data;
}

export async function updateEntry(id, entry) {
    const payload = { ...entry };
    if (entry.date) {
        payload.date = entry.date;
    }
    const res = await axios.put(`${API_URL}/${id}`, payload);
    return res.data;
}

export async function deleteEntry(id) {
    await axios.delete(`${API_URL}/${id}`);
}

export async function getTicketSummary({ start, end }) {
    const res = await axios.get(SUMMARY_URL, { params: { start, end } });
    return res.data;
}
