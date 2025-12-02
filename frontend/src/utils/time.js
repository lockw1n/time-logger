export const isQuarterHour = (hours) => {
    if (!Number.isFinite(hours)) return false;
    const quarters = Math.round(hours * 4);
    return Math.abs(hours * 4 - quarters) < 1e-8;
};

export const parseHoursInput = (input) => {
    if (!input) return NaN;
    const val = String(input).trim().toLowerCase();
    if (!val) return NaN;

    // Match patterns like "2h", "15m", "2h 30m"
    const timeRegex = /(\d+(?:\.\d+)?)\s*(h|m)/g;
    let match;
    let totalHours = 0;
    let found = false;
    while ((match = timeRegex.exec(val)) !== null) {
        found = true;
        const num = parseFloat(match[1]);
        if (Number.isNaN(num)) continue;
        if (match[2] === "h") totalHours += num;
        if (match[2] === "m") totalHours += num / 60;
    }
    if (found) return isQuarterHour(totalHours) ? totalHours : NaN;

    // Fallback: plain numeric hours (supports dot)
    if (/^-?\d+(\.\d+)?$/.test(val) || /^-?\d+(,\d+)?$/.test(val)) {
        const numeric = parseFloat(val.replace(",", "."));
        if (Number.isNaN(numeric)) return NaN;
        return isQuarterHour(numeric) ? numeric : NaN;
    }

    return NaN;
};
