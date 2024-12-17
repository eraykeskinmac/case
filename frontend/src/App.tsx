import { Button } from 'antd';
import { useEffect, useState } from 'react';
import './App.css';

function App() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const testApi = async () => {
      setLoading(true);
      try {
        const response = await fetch('http://localhost:3000/api/health');
        const data = await response.json();
        console.log('API Health Check Response:', data);
        setError(null);
      } catch (err) {
        console.error('API Connection Error:', err);
        setError('Failed to connect to API');
      } finally {
        setLoading(false);
      }
    };

    testApi();
  }, []);

  return (
    <div style={{ padding: '20px' }}>
      <Button loading={loading}>
        {loading ? 'Checking API...' : error || 'API Test'}
      </Button>
    </div>
  );
}

export default App;