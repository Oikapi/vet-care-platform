import { Tabs } from 'antd';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';

const AuthLayout = () => {
  const navigate = useNavigate();
  const location = useLocation();

  const activeTab = location.pathname.includes('login') ? 'login' : 'register';

  const handleTabChange = (key) => {
    navigate(`/auth/${key}`);
  };

  return (
    <div style={{ maxWidth: 500, margin: '50px auto', padding: 20 }}>
      <Tabs
        activeKey={activeTab}
        onChange={handleTabChange}
        items={[
          { key: 'register', label: 'Регистрация' },
          { key: 'login', label: 'Вход' },
        ]}
        centered
      />
      <Outlet />
    </div>
  );
};

export default AuthLayout;
