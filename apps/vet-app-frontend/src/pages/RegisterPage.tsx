// src/pages/RegisterPage.jsx
import { useEffect, useState } from 'react';
import { Tabs, Form, Input, Button, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { authApi } from '../api/auth';

const RegisterPage = () => {
  const [activeTab, setActiveTab] = useState('user');
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const [messageApi, contextHolder] = message.useMessage();

  const onFinish = async (values) => {
    setLoading(true);

    try {
      const registerFunction =
        activeTab === 'user' ? authApi.registerUser : authApi.registerClinic;

      const data = await registerFunction(values);

      messageApi.success({
        content: data.message || 'Регистрация прошла успешно!',
        duration: 3,
      });

      // navigate('/auth/login');
    } catch (error) {
      messageApi.error({
        content: error.message || 'Ошибка регистрации',
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

      <Form form={form} layout="vertical" onFinish={onFinish}>
        {activeTab === 'user' ? (
          <>
            <Form.Item
              name="firstName"
              label="Имя"
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="lastName"
              label="Фамилия"
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="telegram"
              label="Телега"
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
          </>
        ) : (
          <>
            <Form.Item
              name="name"
              label="Название клиники"
              rules={[{ required: true }]}
            >
              <Input />
            </Form.Item>
          </>
        )}
        <Form.Item name="phone" label="Телефон" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
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
          rules={[
            { required: true, message: 'Введите пароль' },
            { min: 6, message: 'Минимум 6 символов' },
          ]}
        >
          <Input.Password />
        </Form.Item>

        <Button type="primary" htmlType="submit" block size="large">
          Зарегистрироваться
        </Button>
      </Form>
    </div>
  );
};

export default RegisterPage;
