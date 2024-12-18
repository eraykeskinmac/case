import React from "react";
import { Table, Tag, Button, Space, Popconfirm } from "antd";
import { EditOutlined, DeleteOutlined, EyeOutlined } from "@ant-design/icons";
import type {
  ColumnsType,
  TablePaginationConfig,
  SorterResult,
} from "antd/es/table/interface";
import { Invoice } from "../../types";
import { useInvoices } from "../../hooks/useInvoice";

interface InvoiceTableProps {
  onShowInvoice: (invoice: Invoice) => void;
  onEditInvoice: (invoice: Invoice) => void;
  onDeleteInvoice: (invoice: Invoice) => void;
}

export const InvoiceTable: React.FC<InvoiceTableProps> = ({
  onShowInvoice,
  onEditInvoice,
  onDeleteInvoice,
}) => {
  const { invoices, loading, filters, total, setFilters } = useInvoices();

  const formatDate = (dateString: string) => {
    return new Intl.DateTimeFormat("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
    }).format(new Date(dateString));
  };

  const columns: ColumnsType<Invoice> = [
    {
      title: "Servis Adı",
      dataIndex: "service_name",
      key: "service_name",
      sorter: true,
    },
    {
      title: "Fatura Numarası",
      dataIndex: "invoice_number",
      key: "invoice_number",
    },
    {
      title: "Tarih",
      dataIndex: "date",
      key: "date",
      sorter: true,
      render: (date: string) => formatDate(date),
    },
    {
      title: "Tutar",
      dataIndex: "amount",
      key: "amount",
      sorter: true,
      align: "right",
      render: (amount: number) =>
        new Intl.NumberFormat("tr-TR", {
          style: "currency",
          currency: "USD",
        }).format(amount),
    },
    {
      title: "Durum",
      dataIndex: "status",
      key: "status",
      sorter: true,
      align: "center",
      render: (status: string) => {
        let color = "";
        let text = "";

        switch (status) {
          case "Paid":
            color = "success";
            text = "Ödendi";
            break;
          case "Pending":
            color = "warning";
            text = "Bekliyor";
            break;
          case "Unpaid":
            color = "error";
            text = "Ödenmedi";
            break;
          default:
            color = "default";
            text = status;
        }

        return (
          <Tag
            color={color}
            style={{
              padding: "0 12px",
              minWidth: "80px",
              textAlign: "center",
            }}
          >
            {text}
          </Tag>
        );
      },
    },
    {
      title: "İşlemler",
      key: "actions",
      align: "center",
      render: (_, record) => (
        <Space size="middle">
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => onShowInvoice(record)}
          />
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => onEditInvoice(record)}
          />
          <Popconfirm
            title="Faturayı silmek istediğinize emin misiniz?"
            okText="Evet"
            cancelText="Hayır"
            onConfirm={() => onDeleteInvoice(record)}
          >
            <Button type="link" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Table
      columns={columns}
      dataSource={invoices}
      loading={loading}
      onChange={(
        pagination: TablePaginationConfig,
        _,
        sorter: SorterResult<Invoice> | SorterResult<Invoice>[]
      ) => {
        const sortInfo = Array.isArray(sorter) ? sorter[0] : sorter;

        setFilters({
          ...filters,
          page: pagination.current || 1,
          limit: pagination.pageSize || 10,
          sort_by: sortInfo.field?.toString(),
          sort_dir: sortInfo.order === "descend" ? "desc" : "asc",
        });
      }}
      pagination={{
        total,
        current: filters.page,
        pageSize: filters.limit,
        showSizeChanger: true,
        showTotal: (total, range) =>
          `${range[0]}-${range[1]} / ${total} fatura`,
      }}
      rowKey="id"
    />
  );
};
