import React, {
  useState,
  useCallback,
  ReactNode,
  useEffect,
} from "react";
import { Invoice, InvoiceFilters } from "../types";
import { invoiceService } from "../api/services/invoiceService";
import { InvoiceContext } from "./InvoiceContext";


export const InvoiceProvider: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  const [invoices, setInvoices] = useState<Invoice[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [total, setTotal] = useState(0);
  const [filters, setFilters] = useState<InvoiceFilters>({
    page: 1,
    limit: 10,
  });

  const fetchInvoices = useCallback(async () => {
    setLoading(true);
    try {
      const response = await invoiceService.getInvoices(filters);

      if (Array.isArray(response.data)) {
        setInvoices(response.data);
        setTotal(response.meta?.total || response.data.length);
      } else {
        console.error('Invalid data format received:', response);
        setInvoices([]);
        setTotal(0);
      }
      setError(null);
    } catch (err) {
      console.error('Fetch error:', err);
      setError(err instanceof Error ? err.message : 'Failed to fetch invoices');
      setInvoices([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    fetchInvoices();
  }, [fetchInvoices]);

  const updateFilters = useCallback((newFilters: Partial<InvoiceFilters>) => {
    setFilters(prev => {
      const updatedFilters = {
        ...prev,
        ...newFilters,
        page: newFilters.page || 1,
      };

      return updatedFilters;
    });
  }, []);

  return (
    <InvoiceContext.Provider
      value={{
        invoices,
        loading,
        error,
        filters,
        total,
        setFilters: updateFilters,
        fetchInvoices,
      }}
    >
      {children}
    </InvoiceContext.Provider>
  );
};
