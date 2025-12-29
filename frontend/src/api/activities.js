import axios from "axios";

const ACTIVITIES_URL = "/api/activities";

export async function listActivitiesForCompany(companyId) {
    const res = await axios.get(ACTIVITIES_URL, {
        params: { company_id: companyId },
    });
    return res.data;
}
