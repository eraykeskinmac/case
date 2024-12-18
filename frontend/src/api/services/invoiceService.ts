import { api } from "../config";
import { Invoice, InvoiceFilters, ApiResponse } from "../../types";

export const invoiceService = {
  getInvoices: async (
    filters: InvoiceFilters
  ): Promise<ApiResponse<Invoice[]>> => {
    try {
      const queryParams = new URLSearchParams();

      if (filters.page) queryParams.append("page", String(filters.page));
      if (filters.limit) queryParams.append("limit", String(filters.limit));
      if (filters.search) {
        queryParams.append("search", filters.search);
      }
      if (filters.sort_by) queryParams.append("sort_by", filters.sort_by);
      if (filters.sort_dir) queryParams.append("sort_dir", filters.sort_dir);

      const response = await api.get(`/invoices?${queryParams}`);

      return {
        data: response.data.data.data || [],
        meta: response.data.data.meta || {
          total: 0,
          page: 1,
          limit: filters.limit,
          total_pages: 1,
        },
      };
    } catch (error) {
      console.error("Error fetching invoices:", error);
      throw error;
    }
  },

  getInvoiceById: async (id: number): Promise<ApiResponse<Invoice>> => {
    return api.get(`/invoices/${id}`);
  },

  createInvoice: async (
    data: Omit<Invoice, "id" | "created_at" | "updated_at">
  ): Promise<ApiResponse<Invoice>> => {
    return api.post("/invoices", data);
  },

  updateInvoice: async (
    id: number,
    data: Partial<Invoice>
  ): Promise<ApiResponse<Invoice>> => {
    return api.put(`/invoices/${id}`, data);
  },

  deleteInvoice: async (id: number): Promise<ApiResponse<void>> => {
    return api.delete(`/invoices/${id}`);
  },
};
