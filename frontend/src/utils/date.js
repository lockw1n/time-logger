export const formatDate = (date) =>
    date.toISOString().slice(0, 10); // YYYY-MM-DD

export const toYMD = (date) => {
    const d = new Date(date);
    const pad = (n) => String(n).padStart(2, "0");
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
};

export const getStartOfWeek = (date) => {
    const d = new Date(date);
    const day = d.getDay(); // 0=Sun
    const diff = (day + 6) % 7; // Monday = 0
    d.setDate(d.getDate() - diff);
    d.setHours(0, 0, 0, 0);
    return d;
};

export const getWeekDays = (startDate, totalDays = 7) => {
    const days = [];
    for (let i = 0; i < totalDays; i++) {
        const d = new Date(startDate);
        d.setDate(d.getDate() + i);
        days.push(d);
    }
    return days;
};
