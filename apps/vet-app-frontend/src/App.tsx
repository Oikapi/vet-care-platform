import { Outlet } from 'react-router-dom';
import { ConfigProvider } from 'antd';

function App() {
  return (
    <ConfigProvider theme={{ token: { colorPrimary: '#00b96b' } }}>
      <Outlet />
    </ConfigProvider>
  );
}

export default App;
