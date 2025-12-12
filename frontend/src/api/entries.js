import axios from "axios";

const API_URL = "/api/entries";

export async function createEntry(entry) {
    const dateStr = entry.date ? entry.date : new Date().toISOString().split("T")[0]; // YYYY-MM-DD in UTC
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
