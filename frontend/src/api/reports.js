import axios from "axios";

export async function downloadInvoicePdf(payload) {
    const res = await axios.post("/api/reports/invoice/pdf", payload, { responseType: "blob" });
    return res.data;
}
