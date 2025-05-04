import { useEffect, useState } from 'react';
import {
  Button,
  Form,
  Input,
  Modal,
  Select,
  Table,
  message,
  Tabs,
  InputNumber,
  DatePicker,
} from 'antd';
import {
  SearchOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ScheduleOutlined,
  MedicineBoxOutlined,
} from '@ant-design/icons';
import type { TableProps, TabsProps } from 'antd';
import { apiInstance } from '../api/auth';
import dayjs from 'dayjs';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import {
  createSchedule,
  deleteSchedule,
  fetchSchedules,
  updateSchedule,
} from '../store/slices/shedulesSlice';
import {
  autoOrder,
  createInventoryItem,
  deleteInventoryItem,
  fetchInventory,
  updateInventoryItem,
  updateQuantity,
} from '../store/slices/inventorySlice';

interface Schedule {
  id: number;
  doctor_id: number;
  doctor_name: string;
  start_time: string;
  end_time: string;
}

interface Inventory {
  id: number;
  medicine_name: string;
  quantity: number;
  threshold: number;
}

interface Doctor {
  id: number;
  firstName: string;
  lastName: string;
}

const ManagementPage = () => {
  const [scheduleForm] = Form.useForm();
  const [inventoryForm] = Form.useForm();
  const [editScheduleForm] = Form.useForm();
  const [editInventoryForm] = Form.useForm();

  const [isScheduleModalOpen, setIsScheduleModalOpen] = useState(false);
  const [isInventoryModalOpen, setIsInventoryModalOpen] = useState(false);
  const [isEditScheduleModalOpen, setIsEditScheduleModalOpen] = useState(false);
  const [isEditInventoryModalOpen, setIsEditInventoryModalOpen] =
    useState(false);

  const dispatch = useAppDispatch();
  const { schedules, loading: schedulesLoading } = useAppSelector(
    (state) => state.schedules
  );
  const { items: inventory, loading: inventoryLoading } = useAppSelector(
    (state) => state.inventory
  );
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState('schedules');
  const [editingSchedule, setEditingSchedule] = useState<Schedule | null>(null);
  const [editingInventory, setEditingInventory] = useState<Inventory | null>(
    null
  );

  useEffect(() => {
    dispatch(fetchSchedules());
    dispatch(fetchInventory());
  }, []);

  //   const fetchDoctors = async () => {
  //     try {
  //       const response = await apiInstance.get('auth/doctors');
  //       setDoctors(response.data);
  //     } catch (error) {
  //       message.error('Ошибка при загрузке врачей');
  //     }
  //   };

  const handleCreateSchedule = async (values: any) => {
    dispatch(
      createSchedule({
        doctor_id: values.doctor_id,
        start_time: values.start_time.toISOString(),
        end_time: values.end_time.toISOString(),
      })
    );
  };

  const handleUpdateSchedule = async (values: any) => {
    if (!editingSchedule) return;
    dispatch(
      updateSchedule({
        id: editingSchedule.id,
        scheduleData: {
          doctor_id: values.doctor_id,
          start_time: values.start_time.toISOString(),
          end_time: values.end_time.toISOString(),
        },
      })
    );
  };

  const handleDeleteSchedule = async (id: number) => {
    dispatch(deleteSchedule(id));
  };

  const handleCreateInventory = async (values: any) => {
    dispatch(createInventoryItem(values));
  };

  const handleUpdateInventory = async (values: any) => {
    if (!editingInventory) return;
    dispatch(
      updateInventoryItem({
        id: editingInventory.id,
        itemData: values,
      })
    );
  };

  const handleUpdateQuantity = async (id: number, quantity: number) => {
    dispatch(
      updateQuantity({
        id: id,
        quantity: quantity,
      })
    );
  };

  const handleDeleteInventory = async (id: number) => {
    dispatch(deleteInventoryItem(id));
  };

  const handleAutoOrder = async () => {
    dispatch(autoOrder());
  };

  const scheduleColumns: TableProps<Schedule>['columns'] = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: 'Врач',
      dataIndex: 'doctor_name',
      key: 'doctor_name',
      render: (_, record) => `${record.doctor_name}`,
    },
    {
      title: 'Начало приема',
      dataIndex: 'start_time',
      key: 'start_time',
      render: (text) => dayjs(text).format('DD.MM.YYYY HH:mm'),
    },
    {
      title: 'Окончание приема',
      dataIndex: 'end_time',
      key: 'end_time',
      render: (text) => dayjs(text).format('DD.MM.YYYY HH:mm'),
    },
    {
      title: 'Действия',
      key: 'actions',
      render: (_, record) => (
        <div>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={(e) => {
              e.stopPropagation();
              setEditingSchedule(record);
              editScheduleForm.setFieldsValue({
                doctor_id: record.doctor_id,
                start_time: dayjs(record.start_time),
                end_time: dayjs(record.end_time),
              });
              setIsEditScheduleModalOpen(true);
            }}
          />
          <Button
            type="link"
            danger
            icon={<DeleteOutlined />}
            onClick={(e) => {
              e.stopPropagation();
              handleDeleteSchedule(record.id);
            }}
          />
        </div>
      ),
    },
  ];

  const inventoryColumns: TableProps<Inventory>['columns'] = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: 'Название',
      dataIndex: 'medicine_name',
      key: 'medicine_name',
    },
    {
      title: 'Количество',
      dataIndex: 'quantity',
      key: 'quantity',
      render: (text, record) => (
        <InputNumber
          min={0}
          value={text}
          onChange={(value) => handleUpdateQuantity(record.id, value as number)}
        />
      ),
    },
    {
      title: 'Порог',
      dataIndex: 'threshold',
      key: 'threshold',
    },
    {
      title: 'Действия',
      key: 'actions',
      render: (_, record) => (
        <div>
          {/* <Button
            type="link"
            icon={<EditOutlined />}
            onClick={(e) => {
              e.stopPropagation();
              setEditingInventory(record);
              editInventoryForm.setFieldsValue(record);
              setIsEditInventoryModalOpen(true);
            }}
          /> */}
          <Button
            type="link"
            danger
            icon={<DeleteOutlined />}
            onClick={(e) => {
              e.stopPropagation();
              handleDeleteInventory(record.id);
            }}
          />
        </div>
      ),
    },
  ];

  const items: TabsProps['items'] = [
    {
      key: 'schedules',
      label: (
        <span>
          <ScheduleOutlined />
          Расписания
        </span>
      ),
      children: (
        <div>
          <div style={{ marginBottom: 16, textAlign: 'right' }}>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => setIsScheduleModalOpen(true)}
            >
              Добавить расписание
            </Button>
          </div>
          <Table
            columns={scheduleColumns}
            dataSource={schedules}
            rowKey="id"
            loading={loading}
          />
        </div>
      ),
    },
    {
      key: 'inventory',
      label: (
        <span>
          <MedicineBoxOutlined />
          Инвентарь
        </span>
      ),
      children: (
        <div>
          <div
            style={{
              marginBottom: 16,
              display: 'flex',
              justifyContent: 'space-between',
            }}
          >
            <Button type="primary" onClick={handleAutoOrder} loading={loading}>
              Выполнить автозаказ
            </Button>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => setIsInventoryModalOpen(true)}
            >
              Добавить лекарство
            </Button>
          </div>
          <Table
            columns={inventoryColumns}
            dataSource={inventory}
            rowKey="id"
            loading={loading}
          />
        </div>
      ),
    },
  ];

  return (
    <div style={{ padding: 24 }}>
      <Tabs
        activeKey={activeTab}
        items={items}
        onChange={(key) => setActiveTab(key)}
      />

      {/* Модальное окно добавления расписания */}
      <Modal
        title="Добавить расписание"
        open={isScheduleModalOpen}
        onCancel={() => setIsScheduleModalOpen(false)}
        footer={null}
      >
        <Form
          form={scheduleForm}
          layout="vertical"
          onFinish={handleCreateSchedule}
        >
          <Form.Item
            name="doctor_id"
            label="Врач"
            rules={[{ required: true, message: 'Пожалуйста, выберите врача' }]}
          >
            <Select placeholder="Выберите врача">
              {[].map((doctor) => (
                <Select.Option key={doctor.id} value={doctor.id}>
                  {doctor.firstName} {doctor.lastName}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="start_time"
            label="Начало приема"
            rules={[
              { required: true, message: 'Пожалуйста, укажите время начала' },
            ]}
          >
            <DatePicker
              showTime
              format="DD.MM.YYYY HH:mm"
              style={{ width: '100%' }}
            />
          </Form.Item>
          <Form.Item
            name="end_time"
            label="Окончание приема"
            rules={[
              {
                required: true,
                message: 'Пожалуйста, укажите время окончания',
              },
            ]}
          >
            <DatePicker
              showTime
              format="DD.MM.YYYY HH:mm"
              style={{ width: '100%' }}
            />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              Сохранить
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      {/* Модальное окно редактирования расписания */}
      <Modal
        title="Редактировать расписание"
        open={isEditScheduleModalOpen}
        onCancel={() => setIsEditScheduleModalOpen(false)}
        footer={null}
      >
        <Form
          form={editScheduleForm}
          layout="vertical"
          onFinish={handleUpdateSchedule}
        >
          <Form.Item
            name="doctor_id"
            label="Врач"
            rules={[{ required: true, message: 'Пожалуйста, выберите врача' }]}
          >
            <Select placeholder="Выберите врача">
              {[].map((doctor) => (
                <Select.Option key={doctor.id} value={doctor.id}>
                  {doctor.firstName} {doctor.lastName}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="start_time"
            label="Начало приема"
            rules={[
              { required: true, message: 'Пожалуйста, укажите время начала' },
            ]}
          >
            <DatePicker
              showTime
              format="DD.MM.YYYY HH:mm"
              style={{ width: '100%' }}
            />
          </Form.Item>
          <Form.Item
            name="end_time"
            label="Окончание приема"
            rules={[
              {
                required: true,
                message: 'Пожалуйста, укажите время окончания',
              },
            ]}
          >
            <DatePicker
              showTime
              format="DD.MM.YYYY HH:mm"
              style={{ width: '100%' }}
            />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              Обновить
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      {/* Модальное окно добавления лекарства */}
      <Modal
        title="Добавить лекарство"
        open={isInventoryModalOpen}
        onCancel={() => setIsInventoryModalOpen(false)}
        footer={null}
      >
        <Form
          form={inventoryForm}
          layout="vertical"
          onFinish={handleCreateInventory}
        >
          <Form.Item
            name="medicine_name"
            label="Название"
            rules={[
              { required: true, message: 'Пожалуйста, введите название' },
            ]}
          >
            <Input placeholder="Введите название лекарства" />
          </Form.Item>
          <Form.Item
            name="quantity"
            label="Количество"
            rules={[
              { required: true, message: 'Пожалуйста, укажите количество' },
            ]}
          >
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item
            name="threshold"
            label="Порог"
            rules={[{ required: true, message: 'Пожалуйста, укажите порог' }]}
          >
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              Сохранить
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      {/* Модальное окно редактирования лекарства */}
      <Modal
        title="Редактировать лекарство"
        open={isEditInventoryModalOpen}
        onCancel={() => setIsEditInventoryModalOpen(false)}
        footer={null}
      >
        <Form
          form={editInventoryForm}
          layout="vertical"
          onFinish={handleUpdateInventory}
        >
          <Form.Item
            name="medicine_name"
            label="Название"
            rules={[
              { required: true, message: 'Пожалуйста, введите название' },
            ]}
          >
            <Input placeholder="Введите название лекарства" />
          </Form.Item>
          <Form.Item
            name="quantity"
            label="Количество"
            rules={[
              { required: true, message: 'Пожалуйста, укажите количество' },
            ]}
          >
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item
            name="threshold"
            label="Порог"
            rules={[{ required: true, message: 'Пожалуйста, укажите порог' }]}
          >
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              Обновить
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default ManagementPage;
