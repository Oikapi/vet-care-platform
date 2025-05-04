import React, { useEffect, useState } from 'react';
import { Card, Col, Row, Modal, Button, Form, Input, Select } from 'antd';
import axios from 'axios';
import { apiInstance } from '../api/auth';

const { Option } = Select;

const DashboardPage: React.FC = () => {
  const [clinics, setClinics] = useState<any[]>([]);
  const [doctors, setDoctors] = useState<any[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form] = Form.useForm();

  useEffect(() => {
    apiInstance
      .get('/auth/clinics')
      .then((res) => setClinics(res.data))
      .catch((err) => console.error(err));

    apiInstance
      .get('/auth/doctors')
      .then((res) => setDoctors(res.data))
      .catch((err) => console.error(err));
  }, []);

  const handleCreateAppointment = (values: any) => {
    apiInstance
      .post('/appointment/appointments', values)
      .then(() => {
        form.resetFields();
        setIsModalOpen(false);
      })
      .catch((err) => console.error(err));
  };

  return (
    <div style={{ padding: 24 }}>
      <Row gutter={24}>
        <Col span={12}>
          <Card title="Список клиник" bordered={false}>
            {clinics.map((clinic) => (
              <p key={clinic.id}>
                {clinic.name} — {clinic.email}
              </p>
            ))}
          </Card>
        </Col>
        <Col span={12}>
          <Card title="Список Докторов" bordered={false}>
            {doctors.map((doctor) => (
              <Card>
                <p key={doctor.id}>
                  {doctor.firstName} — {doctor.lastName}
                </p>
                <p>{doctor.specialization}</p>
              </Card>
            ))}
          </Card>
        </Col>
        <Col span={12}>
          <Card title="Список Докторов" bordered={false}>
            {doctors.map((doctor) => (
              <Card>
                <p key={doctor.id}>
                  {doctor.firstName} — {doctor.lastName}
                </p>
                <p>{doctor.specialization}</p>
              </Card>
            ))}
          </Card>
        </Col>
      </Row>

      <Button
        type="primary"
        style={{ marginTop: 24 }}
        onClick={() => setIsModalOpen(true)}
      >
        Создать запись
      </Button>

      <Modal
        title="Создание записи"
        open={isModalOpen}
        onCancel={() => setIsModalOpen(false)}
        onOk={() => form.submit()}
      >
        <Form form={form} layout="vertical" onFinish={handleCreateAppointment}>
          <Form.Item
            name="doctor_id"
            label="Доктор"
            rules={[{ required: true }]}
          >
            <Select placeholder="Выберите доктора">
              {doctors.map((user) => (
                <Option key={user.id} value={user.id}>
                  {user.firstName} {user.lastName}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="clinic_id"
            label="Клиника"
            rules={[{ required: true }]}
          >
            <Select placeholder="Выберите клинику">
              {clinics.map((clinic) => (
                <Option key={clinic.id} value={clinic.id}>
                  {clinic.name}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="doctor_id"
            label="ID доктора"
            rules={[{ required: true }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            name="slot_id"
            label="ID слота"
            rules={[{ required: true }]}
          >
            <Input />
          </Form.Item>

          <Form.Item name="telegram_id" label="Telegram ID">
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default DashboardPage;
