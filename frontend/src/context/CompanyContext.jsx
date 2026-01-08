import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";
import { apiGet } from "../api/client";
import { useAuth } from "../auth/AuthContext";

const CompanyContext = createContext(null);
const SELECTED_COMPANY_KEY = "company.selected_id";
const COMPANIES_URL = "/api/companies";

const getStoredSelectedCompanyId = () => {
    const storedValue = localStorage.getItem(SELECTED_COMPANY_KEY);
    if (!storedValue) {
        return null;
    }
    const parsed = Number(storedValue);
    return Number.isNaN(parsed) ? null : parsed;
};

export function CompanyProvider({ children }) {
    const { isAuthenticated, isInitializing } = useAuth();
    const [companies, setCompanies] = useState([]);
    const [selectedCompanyId, setSelectedCompanyId] = useState(() => getStoredSelectedCompanyId());

    const handleSetSelectedCompanyId = useCallback((id) => {
        if (id === null) {
            setSelectedCompanyId(null);
            return;
        }
        const nextId = Number(id);
        setSelectedCompanyId(Number.isNaN(nextId) ? null : nextId);
    }, []);

    const resetCompany = useCallback(() => {
        localStorage.removeItem(SELECTED_COMPANY_KEY);
        setSelectedCompanyId(null);
    }, []);

    useEffect(() => {
        if (isInitializing) {
            return;
        }
        if (!isAuthenticated) {
            setCompanies([]);
            return;
        }
        apiGet(COMPANIES_URL)
            .then((data) => {
                const nextCompanies = Array.isArray(data) ? data : data?.companies || [];
                setCompanies(nextCompanies);
                const storedId = getStoredSelectedCompanyId();
                const hasStoredCompany =
                    storedId !== null &&
                    nextCompanies.some((company) => Number(company.id) === storedId);
                setSelectedCompanyId(hasStoredCompany ? storedId : null);
            })
            .catch(() => {
                setCompanies([]);
            });
    }, [isAuthenticated, isInitializing]);

    useEffect(() => {
        if (selectedCompanyId) {
            localStorage.setItem(SELECTED_COMPANY_KEY, String(selectedCompanyId));
        } else {
            localStorage.removeItem(SELECTED_COMPANY_KEY);
        }
    }, [selectedCompanyId]);

    const value = useMemo(
        () => ({
            companies,
            selectedCompanyId,
            setSelectedCompanyId: handleSetSelectedCompanyId,
            resetCompany,
        }),
        [companies, handleSetSelectedCompanyId, resetCompany, selectedCompanyId]
    );

    return <CompanyContext.Provider value={value}>{children}</CompanyContext.Provider>;
}

export const useCompany = () => {
    const context = useContext(CompanyContext);
    if (!context) {
        throw new Error("useCompany must be used within a CompanyProvider");
    }
    return context;
};
