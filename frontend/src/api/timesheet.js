import { apiGet } from "./client";

const TIMESHEET_URL = "/api/timesheet";

export async function getTimesheet({ start, end, companyId } = {}) {
    return apiGet(TIMESHEET_URL, {
        params: {
            company_id: companyId,
            start,
            end,
        },
    });
}
