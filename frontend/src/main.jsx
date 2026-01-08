import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.jsx";
import "./index.css";
import { AuthProvider } from "./auth/AuthContext";
import { CompanyProvider } from "./context/CompanyContext";

ReactDOM.createRoot(document.getElementById("root")).render(
    <React.StrictMode>
        <AuthProvider>
            <CompanyProvider>
                <App />
            </CompanyProvider>
        </AuthProvider>
    </React.StrictMode>
);
