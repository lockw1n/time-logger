import { apiGet } from "./client";

const TIMESHEET_URL = "/api/timesheet";

export async function getTimesheet({ start, end } = {}) {
    return apiGet(TIMESHEET_URL, {
        params: {
            company_id: 1,
            start,
            end,
        },
    });
}
