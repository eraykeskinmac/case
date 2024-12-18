import { createContext } from "react";
import { InvoiceContextType } from "../types";

export const InvoiceContext = createContext<InvoiceContextType | undefined>(
  undefined
);
