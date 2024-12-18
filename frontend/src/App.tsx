import { InvoiceProvider } from './context/InvoiceProvider';
import { Layout } from './components/layout/Layout';
import { InvoicePage } from './pages/InvoicePage';

function App() {
    return (
        <InvoiceProvider>
            <Layout>
                <InvoicePage />
            </Layout>
        </InvoiceProvider>
    );
}

export default App;