# Invoice Management - Frontend

A modern **React** and **TypeScript** application for managing invoices with features like search, filter, sort, and pagination.

## 🚀 Tech Stack

- **React** 18
- **TypeScript**
- **Vite** (Development Tool)
- **Ant Design** v4 (UI Components)
- **Axios** (HTTP Client)
- **Context API** (State Management)

---

## 🧩 Features

- **Invoice List**: View a table of invoices with details like Service Name, Invoice Number, Date, Amount, and Status.
- **Search & Filter**: Search invoices by Service Name or Invoice Number.
- **Sorting**: Sort invoices by Date, Amount, or Status.
- **Pagination**: Navigate between pages of invoices.
- **CRUD Operations**:
   - **Create**: Add a new invoice.
   - **Read**: View detailed invoice information.
   - **Update**: Edit existing invoices.
   - **Delete**: Remove invoices.
- **Status Badges**: Color-coded status indicators for "Paid," "Pending," and "Unpaid."
- **Reusable Components**: Modularized and reusable components for scalability.

---

## 🔧 Project Structure

```plaintext
frontend/
├── src/
│   ├── api/
│   │   ├── services/
│   │   │   └── invoiceService.ts      # API methods for invoices
│   │   └── config.ts                  # Axios configuration
│   ├── components/
│   │   ├── invoices/
│   │   │   ├── InvoiceFormModal.tsx   # Modal for creating/editing invoices
│   │   │   ├── InvoiceModal.tsx       # Modal for viewing invoice details
│   │   │   └── InvoiceTable.tsx       # Table component to display invoices
│   │   ├── layout/
│   │   │   ├── Layout.tsx             # Page layout
│   │   │   └── Sidebar.tsx            # Sidebar navigation
│   ├── context/
│   │   ├── InvoiceContext.tsx         # Invoice context
│   │   └── InvoiceProvider.tsx        # Invoice provider for state management
│   ├── hooks/
│   │   └── useInvoice.ts              # Custom hook for invoices
│   ├── pages/
│   │   └── InvoicePage.tsx            # Main invoice management page
│   ├── types/                         # Type definitions
│   ├── App.tsx                        # Main app entry point
│   ├── main.tsx                       # React entry file
│   └── vite-env.d.ts                  # Vite TypeScript config
│
├── public/
│   └── index.html                     # Root HTML file
│
├── .env                               # Environment variables
├── Dockerfile                         # Frontend Dockerfile
├── tsconfig.json                      # TypeScript configuration
└── vite.config.ts                     # Vite configuration
```

---

## 🌐 API Integration

The frontend uses **Axios** to communicate with the API. API calls are organized in the `invoiceService.ts` file.

---

## 🎨 UI Components

The application uses **Ant Design** for a modern and responsive UI. Components include:
- **InvoiceTable**: Displays the list of invoices.
- **InvoiceFormModal**: Modal for creating and editing invoices.
- **InvoiceModal**: Modal to view invoice details.
- **Status Badges**: Color-coded badges for invoice statuses.

---

## 📋 Notes

- The frontend is optimized with **Vite** for fast development and hot module replacement.
- State management is handled using the **Context API**.
- Ensure the backend API is running and accessible for proper integration.

