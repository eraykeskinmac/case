import { useState } from "react";
import { Input, Button, Row, Col, message } from "antd";
import { SearchOutlined, PlusOutlined } from "@ant-design/icons";
import { InvoiceTable } from "../components/invoices/InvoiceTable";
import { InvoiceFormModal } from "../components/invoices/InvoiceFormModal";
import { InvoiceModal } from "../components/invoices/InvoiceModal";
import { CreateInvoiceData, Invoice } from "../types";
import { invoiceService } from "../api/services/invoiceService";
import { useInvoices } from "../hooks/useInvoice";

export function InvoicePage() {
  const { setFilters, fetchInvoices } = useInvoices();
  const [searchValue, setSearchValue] = useState("");
  const [selectedInvoice, setSelectedInvoice] = useState<Invoice | undefined>();
  const [showViewModal, setShowViewModal] = useState(false);
  const [showFormModal, setShowFormModal] = useState(false);
  const [editingInvoice, setEditingInvoice] = useState<Invoice | undefined>();

  const handleSearch = (value: string) => {
    setSearchValue(value);
    const trimmedValue = value.trim();

    setFilters({
      page: 1,
      search: trimmedValue || undefined,
    });
  };

  const handleCreateInvoice = async (values: CreateInvoiceData) => {
    try {
      await invoiceService.createInvoice(values);
      message.success("Fatura başarıyla oluşturuldu");
      fetchInvoices();
      setShowFormModal(false);
    } catch (err) {
      const errorMessage =
        err instanceof Error
          ? err.message
          : "Fatura oluşturulurken bir hata oluştu";
      message.error(errorMessage);
    }
  };

  const handleUpdateInvoice = async (values: CreateInvoiceData) => {
    if (!editingInvoice?.id) return;

    try {
      await invoiceService.updateInvoice(editingInvoice.id, values);
      message.success("Fatura başarıyla güncellendi");
      fetchInvoices();
      setEditingInvoice(undefined);
      setShowFormModal(false);
    } catch (err) {
      const errorMessage =
        err instanceof Error
          ? err.message
          : "Fatura güncellenirken bir hata oluştu";
      message.error(errorMessage);
    }
  };

  const handleDeleteInvoice = async (invoice: Invoice) => {
    try {
      await invoiceService.deleteInvoice(invoice.id);
      message.success("Fatura başarıyla silindi");
      fetchInvoices();
    } catch (err) {
      const errorMessage =
        err instanceof Error
          ? err.message
          : "Fatura silinirken bir hata oluştu";
      message.error(errorMessage);
    }
  };

  const handleEdit = (invoice: Invoice) => {
    setEditingInvoice(invoice);
    setShowFormModal(true);
  };

  const handleView = (invoice: Invoice) => {
    setSelectedInvoice(invoice);
    setShowViewModal(true);
  };

  return (
    <div style={{ padding: "24px" }}>
      <Row
        gutter={[16, 16]}
        align="middle"
        justify="space-between"
        style={{ marginBottom: 16 }}
      >
        <Col>
          <Input
            placeholder="Servis adı veya fatura numarası ile ara..."
            prefix={<SearchOutlined />}
            value={searchValue}
            onChange={(e) => handleSearch(e.target.value)}
            style={{ width: 300 }}
            allowClear
          />
        </Col>
        <Col>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => setShowFormModal(true)}
          >
            Yeni Fatura
          </Button>
        </Col>
      </Row>

      <InvoiceTable
        onShowInvoice={handleView}
        onEditInvoice={handleEdit}
        onDeleteInvoice={handleDeleteInvoice}
      />

      <InvoiceModal
        invoice={selectedInvoice}
        visible={showViewModal}
        onClose={() => {
          setShowViewModal(false);
          setSelectedInvoice(undefined);
        }}
      />
      <InvoiceFormModal
        visible={showFormModal}
        onClose={() => {
          setShowFormModal(false);
          setEditingInvoice(undefined);
        }}
        onSubmit={editingInvoice ? handleUpdateInvoice : handleCreateInvoice}
        initialValues={editingInvoice}
        title={editingInvoice ? "Fatura Düzenle" : "Yeni Fatura"}
      />
    </div>
  );
}
