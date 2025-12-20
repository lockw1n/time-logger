import express from "express";
import {chromium} from "playwright";

const app = express();
app.use(express.json({limit: "10mb"}));

const PORT = Number(process.env.PORT || 3001);
const INTERNAL_TOKEN = process.env.INTERNAL_TOKEN || ""; // optional

// Launch one shared browser, reuse it across requests.
const browserPromise = chromium.launch({
    headless: true,
    args: [
        "--no-sandbox",
        "--disable-setuid-sandbox",
        "--disable-dev-shm-usage"
    ]
});

app.get("/health", async (_req, res) => {
    res.json({ok: true});
});

app.post("/render", async (req, res) => {
    const started = Date.now();

    // Optional internal auth
    if (INTERNAL_TOKEN) {
        const token = req.header("X-Internal-Token");
        if (token !== INTERNAL_TOKEN) {
            return res.status(401).json({error: "unauthorized"});
        }
    }

    const {html, footerHtml, baseUrl, pdf, timeoutMs} = req.body ?? {};
    if (typeof html !== "string" || html.trim().length === 0) {
        return res.status(400).json({error: "Field 'html' is required (non-empty string)."});
    }

    const effectiveTimeout = Number.isFinite(timeoutMs) ? timeoutMs : 30000;

    const pdfOptions = {
        format: "A4",
        printBackground: true,
        preferCSSPageSize: true,
        margin: {
            top: "10mm",
            right: "10mm",
            bottom: footerHtml ? "20mm" : "10mm",
            left: "10mm"
        },
        ...((pdf && typeof pdf === "object") ? pdf : {})
    };

    if (typeof footerHtml === "string" && footerHtml.trim() !== "") {
        pdfOptions.displayHeaderFooter = true;
        pdfOptions.headerTemplate = "<div></div>";
        pdfOptions.footerTemplate = footerHtml;
    }

    const browser = await browserPromise;
    const context = await browser.newContext();

    try {
        const page = await context.newPage();

        // If baseUrl is provided, allow only same-origin requests (plus data/blob).
        if (typeof baseUrl === "string" && baseUrl.trim() !== "") {
            let allowedOrigin = "";
            try {
                allowedOrigin = new URL(baseUrl).origin;
            } catch {
                return res.status(400).json({error: "Invalid 'baseUrl' (must be a valid URL)."});
            }

            await page.route("**/*", (route) => {
                const url = route.request().url();

                // Always allow inline resources
                if (url.startsWith("data:") || url.startsWith("blob:")) {
                    return route.continue();
                }

                let origin;
                try {
                    origin = new URL(url).origin;
                } catch {
                    return route.abort();
                }

                if (origin === allowedOrigin) {
                    return route.continue();
                }

                return route.abort();
            });

        }

        await page.setContent(html, {waitUntil: "networkidle", timeout: effectiveTimeout});

        // Best-effort: wait for fonts
        await page.evaluate(async () => {
            try {
                await document.fonts.ready;
            } catch {
            }
        });

        const buffer = await page.pdf({...pdfOptions});

        res.setHeader("Content-Type", "application/pdf");
        res.setHeader("Content-Disposition", 'inline; filename="document.pdf"');
        res.setHeader("X-Render-Time-Ms", String(Date.now() - started));
        return res.status(200).send(buffer);
    } catch (e) {
        return res.status(500).json({
            error: "render_failed",
            message: e?.message ?? String(e)
        });
    } finally {
        await context.close().catch(() => {});
    }
});

app.listen(PORT, () => {
    console.log(`pdf-renderer listening on :${PORT}`);
});

// Graceful shutdown
process.on("SIGTERM", async () => {
    try {
        const browser = await browserPromise;
        await browser.close();
    } finally {
        process.exit(0);
    }
});
