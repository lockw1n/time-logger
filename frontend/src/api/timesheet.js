import axios from "axios";

const TIMESHEET_URL = "/api/timesheet";

export async function getTimesheet({ start, end } = {}) {
    // Temporarily hardcode consultant/company while backend is single-tenant.
    const res = await axios.get(TIMESHEET_URL, {
        params: {
            consultant_id: 1,
            company_id: 1,
            start,
            end,
        },
    });
    return res.data;
}
