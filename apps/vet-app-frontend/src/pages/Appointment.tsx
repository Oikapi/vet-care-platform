import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  createAppointment,
  fetchAvailableSlots,
  fetchAppointment,
  createSlot,
  clearAvailableSlots,
  resetAppointmentsError,
} from '../store/slices/appointmentSlice';
import { RootState, AppDispatch } from '../store/store';
import {
  Button,
  Table,
  Modal,
  Form,
  Input,
  DatePicker,
  TimePicker,
  Select,
  message,
  Card,
  Tag,
  Space,
  Tabs,
} from 'antd';
import { PlusOutlined, SearchOutlined, CloseOutlined } from '@ant-design/icons';
import dayjs, { Dayjs } from 'dayjs';
import type { ColumnsType } from 'antd/es/table';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { Appointment, Slot } from '../types';

const { RangePicker } = DatePicker;

const AppointmentsPage = () => {
  const dispatch = useAppDispatch();
  const { appointments, availableSlots, loading, error } = useAppSelector(
    (state) => state.appointments
  );

  const [selectedClinicId, setSelectedClinicId] = useState<number>(1); // Можно сделать выбор клиники
  const [isSlotModalVisible, setIsSlotModalVisible] = useState(false);
  const [isAppointmentModalVisible, setIsAppointmentModalVisible] =
    useState(false);
  const [selectedSlot, setSelectedSlot] = useState<Slot | null>(null);
  const [form] = Form.useForm();
  const [slotForm] = Form.useForm();

  // Загрузка доступных слотов при изменении выбранной клиники
  useEffect(() => {
    dispatch(fetchAvailableSlots(selectedClinicId));
    return () => {
      dispatch(clearAvailableSlots());
    };
  }, [dispatch, selectedClinicId]);

  // Обработка ошибок
  useEffect(() => {
    if (error) {
      message.error(error);
      dispatch(resetAppointmentsError());
    }
  }, [error, dispatch]);

  // Создание нового слота
  const handleCreateSlot = async (values: {
    doctor_id: number;
    timeRange: [Dayjs, Dayjs];
  }) => {
    try {
      await dispatch(
        createSlot({
          doctor_id: values.doctor_id,
          slot_time: values.timeRange[0],
        })
      ).unwrap();

      message.success('Слот успешно создан');
      setIsSlotModalVisible(false);
      slotForm.resetFields();
    } catch (err) {
      message.error('Ошибка при создании слота');
    }
  };

  // Создание записи на прием
  const handleCreateAppointment = async (values: { patient_id: number }) => {
    if (!selectedSlot) return;

    try {
      await dispatch(
        createAppointment({
          slot_id: selectedSlot.id,
          patient_id: values.patient_id,
        })
      ).unwrap();

      message.success('Запись успешно создана');
      setIsAppointmentModalVisible(false);
      setSelectedSlot(null);
    } catch (err) {
      message.error('Ошибка при создании записи');
    }
  };

  // Колонки для таблицы слотов
  const slotColumns: ColumnsType<Slot> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: 'Доктор',
      dataIndex: ['doctor', 'first_name'],
      key: 'doctor',
      render: (_, record) => `${record.doctor_id} ${record.doctor_id}`,
    },
    {
      title: 'Время',
      key: 'time',
      render: (_, record) => (
        <>
          {dayjs(record.start_time).format('DD.MM.YYYY HH:mm')} -
          {dayjs(record.end_time).format('HH:mm')}
        </>
      ),
    },
    {
      title: 'Статус',
      key: 'status',
      render: (_, record) => (
        <Tag color={record.IsBooked ? 'green' : 'red'}>
          {record.is_available ? 'Доступен' : 'Занят'}
        </Tag>
      ),
    },
    {
      title: 'Действия',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            disabled={!record.is_available}
            onClick={() => {
              setSelectedSlot(record);
              setIsAppointmentModalVisible(true);
            }}
          >
            Записать
          </Button>
        </Space>
      ),
    },
  ];

  // Колонки для таблицы записей
  const appointmentColumns: ColumnsType<Appointment> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: 'Доктор',
      key: 'doctor',
      render: (_, record) => `${record.slot.doctor_id}`,
    },
    {
      title: 'Время приема',
      key: 'time',
      render: (_, record) => (
        <>
          {dayjs(record.slot.start_time).format('DD.MM.YYYY HH:mm')} -
          {dayjs(record.slot.end_time).format('HH:mm')}
        </>
      ),
    },
    {
      title: 'Статус',
      key: 'status',
      render: (_, record) => {
        let color = '';
        switch (record.status) {
          case 'completed':
            color = 'green';
            break;
          case 'cancelled':
            color = 'red';
            break;
          default:
            color = 'blue';
        }
        return <Tag color={color}>{record.status}</Tag>;
      },
    },
  ];

  // Моковые данные для демонстрации (замените на реальные)
  const mockDoctors = [
    { id: 1, first_name: 'Иван', last_name: 'Петров' },
    { id: 2, first_name: 'Мария', last_name: 'Иванова' },
  ];

  const mockPatients = [
    { id: 1, first_name: 'Алексей', last_name: 'Сидоров' },
    { id: 2, first_name: 'Елена', last_name: 'Кузнецова' },
  ];

  return (
    <div className="appointments-page">
      <Card title="Управление записями на прием" bordered={false}>
        <div style={{ marginBottom: 16 }}>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => setIsSlotModalVisible(true)}
            style={{ marginRight: 8 }}
          >
            Создать слот
          </Button>
          <Select
            value={selectedClinicId}
            onChange={(value) => setSelectedClinicId(value)}
            style={{ width: 200 }}
          >
            <Select.Option value={1}>Клиника №1</Select.Option>
            <Select.Option value={2}>Клиника №2</Select.Option>
          </Select>
        </div>

        <Tabs defaultActiveKey="slots">
          <Tabs.TabPane tab="Доступные слоты" key="slots">
            <Table
              columns={slotColumns}
              dataSource={availableSlots}
              rowKey="id"
              loading={loading}
              pagination={{ pageSize: 10 }}
            />
          </Tabs.TabPane>
          <Tabs.TabPane tab="Записи на прием" key="appointments">
            <Table
              columns={appointmentColumns}
              dataSource={appointments}
              rowKey="id"
              loading={loading}
              pagination={{ pageSize: 10 }}
            />
          </Tabs.TabPane>
        </Tabs>
      </Card>

      {/* Модальное окно создания слота */}
      <Modal
        title="Создать новый слот"
        visible={isSlotModalVisible}
        onCancel={() => setIsSlotModalVisible(false)}
        footer={null}
      >
        <Form form={slotForm} layout="vertical" onFinish={handleCreateSlot}>
          <Form.Item
            name="doctor_id"
            label="Доктор"
            rules={[{ required: true, message: 'Выберите доктора' }]}
          >
            <Select placeholder="Выберите доктора">
              {mockDoctors.map((doctor) => (
                <Select.Option key={doctor.id} value={doctor.id}>
                  {doctor.first_name} {doctor.last_name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="timeRange"
            label="Временной интервал"
            rules={[{ required: true, message: 'Укажите время приема' }]}
          >
            <RangePicker showTime format="DD.MM.YYYY HH:mm" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              Создать
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      {/* Модальное окно создания записи */}
      <Modal
        title={`Запись на прием - ${selectedSlot?.doctor_id}`}
        visible={isAppointmentModalVisible}
        onCancel={() => {
          setIsAppointmentModalVisible(false);
          setSelectedSlot(null);
        }}
        footer={null}
      >
        {selectedSlot && (
          <div style={{ marginBottom: 16 }}>
            <p>
              <strong>Время: </strong>
              {dayjs(selectedSlot.start_time).format('DD.MM.YYYY HH:mm')} -
              {dayjs(selectedSlot.end_time).format('HH:mm')}
            </p>
          </div>
        )}
        <Form form={form} layout="vertical" onFinish={handleCreateAppointment}>
          <Form.Item
            name="patient_id"
            label="Пациент"
            rules={[{ required: true, message: 'Выберите пациента' }]}
          >
            <Select placeholder="Выберите пациента">
              {mockPatients.map((patient) => (
                <Select.Option key={patient.id} value={patient.id}>
                  {patient.first_name} {patient.last_name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              Создать запись
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default AppointmentsPage;
