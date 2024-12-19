import React from 'react';
import { Modal, Descriptions, Tag } from 'antd';
import { Invoice } from '../../types';

interface InvoiceModalProps {
  visible: boolean;
  onClose: () => void;
  invoice: Invoice | undefined;
}

export const InvoiceModal: React.FC<InvoiceModalProps> = ({ invoice, visible, onClose }) => {
  if (!invoice) return null;

  return (
    <Modal
      title="Fatura Detayı"
      open={visible}
      onCancel={onClose}
      footer={null}
      width={600}
    >
      <Descriptions bordered column={1}>
        <Descriptions.Item label="Servis Adı">{invoice.service_name}</Descriptions.Item>
        <Descriptions.Item label="Fatura Numarası">{invoice.invoice_number}</Descriptions.Item>
        <Descriptions.Item label="Tarih">
          {new Date(invoice.date).toLocaleDateString('tr-TR')}
        </Descriptions.Item>
        <Descriptions.Item label="Tutar">
          {new Intl.NumberFormat('tr-TR', {
            style: 'currency',
            currency: 'USD'
          }).format(invoice.amount)}
        </Descriptions.Item>
        <Descriptions.Item label="Durum">
          <Tag color={
            invoice.status === 'Paid' ? 'success' :
            invoice.status === 'Pending' ? 'warning' : 'error'
          }>
            {invoice.status === 'Paid' ? 'Ödendi' :
             invoice.status === 'Pending' ? 'Bekliyor' : 'Ödenmedi'}
          </Tag>
        </Descriptions.Item>
      </Descriptions>
    </Modal>
  );
};