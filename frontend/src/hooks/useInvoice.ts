import { useContext } from 'react';
import { InvoiceContextType } from '../types';
import { InvoiceContext } from '../context/InvoiceContext';


export const useInvoices = (): InvoiceContextType => {
  const context = useContext(InvoiceContext);
  if (!context) {
    throw new Error("useInvoices must be used within an InvoiceProvider");
  }
  return context;
};