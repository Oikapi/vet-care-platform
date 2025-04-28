import { useState } from 'react';
import { Tabs, Form, Input, Button, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { authApi } from '../api/auth';

const LoginPage = () => {
  const [activeTab, setActiveTab] = useState('user');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const [messageApi, contextHolder] = message.useMessage();

  const onFinish = async (values) => {
    setLoading(true);
    try {
      const loginFunction =
        activeTab === 'user' ? authApi.loginUser : authApi.loginClinic;

      const data = await loginFunction(values);

      localStorage.setItem('token', data.access_token);
      messageApi.success({
        content: data.message || 'Вход выполнен успешно!',
        duration: 3,
      });
    } catch (error) {
      messageApi.error({
        content: error.message || 'Ошибка входа',
        duration: 3,
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: '20px 0' }}>
      {contextHolder}
      <Tabs
        activeKey={activeTab}
        onChange={setActiveTab}
        items={[
          { key: 'user', label: 'Я пользователь' },
          { key: 'clinic', label: 'Я клиника' },
        ]}
        centered
      />

      <Form layout="vertical" onFinish={onFinish}>
        <Form.Item
          name="email"
          label="Email"
          rules={[
            { required: true, message: 'Введите email' },
            { type: 'email', message: 'Некорректный email' },
          ]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          name="password"
          label="Пароль"
          rules={[{ required: true, message: 'Введите пароль' }]}
        >
          <Input.Password />
        </Form.Item>

        <Button
          type="primary"
          htmlType="submit"
          block
          size="large"
          loading={loading}
        >
          Войти
        </Button>
      </Form>
    </div>
  );
};

export default LoginPage;
