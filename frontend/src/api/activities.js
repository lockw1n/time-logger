import { apiGet } from "./client";

const ACTIVITIES_URL = "/api/activities";

export async function listActivitiesForCompany(companyId) {
    return apiGet(ACTIVITIES_URL, {
        params: { company_id: companyId },
    });
}
