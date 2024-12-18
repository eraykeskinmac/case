import { Menu } from 'antd';
import { FileTextOutlined, CreditCardOutlined, SettingOutlined, CustomerServiceOutlined } from '@ant-design/icons';

export function Sidebar() {
    return (
        <Menu
            mode="inline"
            defaultSelectedKeys={['invoices']}
            style={{ height: '100%', borderRight: 0 }}
            items={[
                {
                    key: 'invoices',
                    icon: <FileTextOutlined />,
                    label: 'Faturalar',
                },
                {
                    key: 'payment-methods',
                    icon: <CreditCardOutlined />,
                    label: 'Ödeme Yöntemleri',
                },
                {
                    key: 'services',
                    icon: <CustomerServiceOutlined />,
                    label: 'Hizmetler',
                },
                {
                    key: 'settings',
                    icon: <SettingOutlined />,
                    label: 'Ayarlar',
                },
            ]}
        />
    );
}