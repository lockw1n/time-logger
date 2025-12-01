import axios from "axios";

const API_URL = "/api/entries";

export async function getEntries({ start, end } = {}) {
    const res = await axios.get(API_URL, {
        params: start && end ? { start, end } : {},
    });
    return res.data;
}

export async function createEntry(entry) {
    const payload = {
        ...entry,
        date: entry.date ? new Date(entry.date).toISOString() : new Date().toISOString(),
    };
    const res = await axios.post(API_URL, payload);
    return res.data;
}

export async function updateEntry(id, entry) {
    const payload = {
        ...entry,
        date: entry.date ? new Date(entry.date).toISOString() : new Date().toISOString(),
    };
    const res = await axios.put(`${API_URL}/${id}`, payload);
    return res.data;
}

export async function deleteEntry(id) {
    await axios.delete(`${API_URL}/${id}`);
}
