import { Layout as AntLayout } from 'antd';
import { Sidebar } from './Sidebar';

const { Sider, Content: AntContent } = AntLayout;

export function Layout({ children }: { children: React.ReactNode }) {
    return (
        <AntLayout style={{ minHeight: '100vh' }}>
            <Sider
                theme="light"
                width={250}
                style={{
                  padding: '24px 0',
                  borderRight: '1px solid #f0f0f0',
                  background: '#fafafa',
              }}
            >
                <Sidebar />
            </Sider>
            <AntContent style={{ background: '#fff', padding: '24px' }}>
                {children}
            </AntContent>
        </AntLayout>
    );
}