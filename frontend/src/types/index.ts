export interface Invoice {
  id: number;
  service_name: string;
  invoice_number: number;
  date: string;
  amount: number;
  status: 'Paid' | 'Pending' | 'Unpaid';
  created_at: string;
  updated_at: string;
}

export type CreateInvoiceData = Omit<Invoice, 'id' | 'created_at' | 'updated_at'>;

export interface InvoiceFilters {
  page: number;
  limit: number;
  search?: string;
  sort_by?: string;
  sort_dir?: 'asc' | 'desc';
}

export interface ApiResponse<T> {
  data: T;
  message?: string;
  meta?: {
    total: number;
    page: number;
    limit: number;
    total_pages: number;
  };
}

export interface InvoiceContextType {
  invoices: Invoice[];
  loading: boolean;
  error: string | null;
  filters: InvoiceFilters;
  total: number;
  setFilters: (filters: Partial<InvoiceFilters>) => void;
  fetchInvoices: () => Promise<void>;
}