import axios from "axios";

const INVOICE_URL = "/api/invoices/monthly";

const extractFileMeta = (headers = {}) => {
    const disposition = String(headers["content-disposition"] || "");
    const isInline = disposition.toLowerCase().includes("inline");

    let filename = "invoice.pdf";
    const filenameRegex = /filename\*?=(?:UTF-8''|")?([^";]+)/i;
    const match = disposition.match(filenameRegex);
    if (match && match[1]) {
        try {
            filename = decodeURIComponent(match[1].trim());
        } catch {
            filename = match[1].trim();
        }
    }

    return { filename, isInline };
};

export async function generateMonthlyInvoicePdf({ month, consultantId = 1, companyId = 1 } = {}) {
    const res = await axios.post(
        `${INVOICE_URL}?format=pdf`,
        {
            month,
            consultant_id: consultantId,
            company_id: companyId,
        },
        { responseType: "blob" }
    );

    const { filename, isInline } = extractFileMeta(res.headers);

    return {
        blob: res.data,
        filename,
        isInline,
    };
}

export async function generateMonthlyInvoiceExcel({ month, consultantId = 1, companyId = 1,} = {}) {
    const res = await axios.post(
        `${INVOICE_URL}?format=excel`,
        {
            month,
            consultant_id: consultantId,
            company_id: companyId,
        },
        { responseType: "blob" }
    );

    const disposition = res.headers["content-disposition"] || "";
    let filename = "invoice.xlsx";

    const match = disposition.match(/filename\*?=(?:UTF-8''|")?([^";]+)/i);
    if (match && match[1]) {
        try {
            filename = decodeURIComponent(match[1].trim());
        } catch {
            filename = match[1].trim();
        }
    }

    return {
        blob: res.data,
        filename,
    };
}
