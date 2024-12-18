import React from 'react';
import { Modal, Form, Input, DatePicker, Select, InputNumber } from 'antd';
import { CreateInvoiceData, Invoice } from '../../types';
import moment from 'moment';

interface InvoiceFormModalProps {
  visible: boolean;
  onClose: () => void;
  onSubmit: (values: CreateInvoiceData) => Promise<void>;
  initialValues?: Invoice;
  title: string;
}

export const InvoiceFormModal: React.FC<InvoiceFormModalProps> = ({
  visible,
  onClose,
  onSubmit,
  initialValues,
  title,
}) => {
  const [form] = Form.useForm();

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      await onSubmit({
        ...values,
        date: values.date.toISOString(),
      });
      form.resetFields();
      onClose();
    } catch (error) {
      console.error('Form validation failed:', error);
    }
  };

  React.useEffect(() => {
    if (visible && initialValues) {
      form.setFieldsValue({
        ...initialValues,
        date: moment(initialValues.date),
      });
    }
  }, [visible, initialValues, form]);

  return (
    <Modal
      title={title}
      open={visible}
      onOk={handleSubmit}
      onCancel={onClose}
      okText={initialValues ? 'Güncelle' : 'Oluştur'}
      cancelText="İptal"
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={initialValues ? {
          ...initialValues,
          date: moment(initialValues.date),
        } : undefined}
      >
        <Form.Item
          name="service_name"
          label="Servis Adı"
          rules={[{ required: true, message: 'Servis adı zorunludur' }]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="invoice_number"
          label="Fatura Numarası"
          rules={[{ required: true, message: 'Fatura numarası zorunludur' }]}
        >
          <InputNumber style={{ width: '100%' }} />
        </Form.Item>

        <Form.Item
          name="date"
          label="Tarih"
          rules={[{ required: true, message: 'Tarih zorunludur' }]}
        >
          <DatePicker style={{ width: '100%' }} />
        </Form.Item>

        <Form.Item
          name="amount"
          label="Tutar"
          rules={[{ required: true, message: 'Tutar zorunludur' }]}
        >
          <InputNumber
            style={{ width: '100%' }}
            precision={2}
            formatter={value => `$ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
            parser={value => value!.replace(/\$\s?|(,*)/g, '')}
          />
        </Form.Item>

        <Form.Item
          name="status"
          label="Durum"
          rules={[{ required: true, message: 'Durum zorunludur' }]}
        >
          <Select>
            <Select.Option value="Paid">Ödendi</Select.Option>
            <Select.Option value="Pending">Bekliyor</Select.Option>
            <Select.Option value="Unpaid">Ödenmedi</Select.Option>
          </Select>
        </Form.Item>
      </Form>
    </Modal>
  );
};