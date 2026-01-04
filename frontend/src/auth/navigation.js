export const buildLoginRedirectUrl = (fromPath) => {
    const redirectParam = fromPath ? `?redirect=${encodeURIComponent(fromPath)}` : "";
    return `/login${redirectParam}`;
};

export const redirectToLogin = (fromPath) => {
    const next = fromPath || `${window.location.pathname}${window.location.search}`;
    const loginUrl = buildLoginRedirectUrl(next);
    if (window.location.pathname === "/login") return;
    window.history.replaceState(null, "", loginUrl);
    window.dispatchEvent(new PopStateEvent("popstate"));
};
