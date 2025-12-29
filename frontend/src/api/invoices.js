import axios from "axios";

const INVOICE_URL = "/api/invoices/generate";

const buildMonthRange = (month) => {
    if (!month) return null;
    const [yearStr, monthStr] = String(month).split("-");
    const year = Number(yearStr);
    const monthIndex = Number(monthStr) - 1;
    if (!Number.isInteger(year) || !Number.isInteger(monthIndex) || monthIndex < 0 || monthIndex > 11) {
        return null;
    }
    const start = new Date(year, monthIndex, 1);
    const end = new Date(year, monthIndex + 1, 0);
    const pad = (n) => String(n).padStart(2, "0");
    const toYMD = (date) => `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;
    return { start: toYMD(start), end: toYMD(end) };
};

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

const invalidMonthError = new Error("invalid month");
invalidMonthError.code = "invalid_month";

export async function generateMonthlyInvoicePdf({ month, consultantId = 1, companyId = 1 } = {}) {
    const range = buildMonthRange(month);
    if (!range) {
        throw invalidMonthError;
    }
    const res = await axios.post(
        INVOICE_URL,
        null,
        {
            params: {
                consultant_id: consultantId,
                company_id: companyId,
                start: range.start,
                end: range.end,
                format: "pdf",
            },
            responseType: "blob",
        }
    );

    const { filename, isInline } = extractFileMeta(res.headers);

    return {
        blob: res.data,
        filename,
        isInline,
    };
}

export async function generateMonthlyInvoiceExcel({ month, consultantId = 1, companyId = 1 } = {}) {
    const range = buildMonthRange(month);
    if (!range) {
        throw invalidMonthError;
    }
    const res = await axios.post(
        INVOICE_URL,
        null,
        {
            params: {
                consultant_id: consultantId,
                company_id: companyId,
                start: range.start,
                end: range.end,
                format: "excel",
            },
            responseType: "blob",
        }
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
