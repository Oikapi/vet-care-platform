import { useEffect, useState } from 'react';
import { Button, Form, Input, Modal, Select, Table, message } from 'antd';
import { SearchOutlined, PlusOutlined, EditOutlined } from '@ant-design/icons';
import type { TableProps } from 'antd';
import { apiInstance } from '../api/auth';

interface Pet {
  id: number;
  species: string;
  name: string;
  breed: string;
  gender: string;
  age: number;
}

const PetsPage = () => {
  const [form] = Form.useForm();
  const [editForm] = Form.useForm();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [searchId, setSearchId] = useState('');
  const [pet, setPet] = useState<Pet | null>(null);
  const [pets, setPets] = useState<Pet[]>([]);
  const [loading, setLoading] = useState(false);
  const [users, setUsers] = useState([]);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [editingPet, setEditingPet] = useState<Pet | null>(null);

  useEffect(() => {
    const getUsers = async () => {
      try {
        const response = await apiInstance.get<any[]>('auth/users');
        const data = response.data;
        setUsers(data);
      } catch (e) {
        console.log(e);
      }
    };

    getUsers();
  }, []);

  useEffect(() => {
    getAllPets();
  }, []);

  const getAllPets = async () => {
    try {
      const response = await apiInstance.get('pets/pets');

      const data = response.data;
      setPets(data);
    } catch (e) {
      console.log(e);
    }
  };

  const showModal = () => {
    setIsModalOpen(true);
  };

  const handleCancel = () => {
    setIsModalOpen(false);
    form.resetFields();
  };

  const handleEditCancel = () => {
    setIsEditModalOpen(false);
    editForm.resetFields();
    setEditingPet(null);
  };

  const showEditModal = (pet: Pet) => {
    setEditingPet(pet);
    editForm.setFieldsValue({
      ...pet,
      age: pet.age.toString(),
    });
    setIsEditModalOpen(true);
  };

  const handleCreatePet = async (values: Omit<Pet, 'id'>) => {
    try {
      setLoading(true);
      const response = await apiInstance.post<Pet>('pets/register', {
        ...values,
      });
      const data = response.data;

      // Имитация ответа от сервера
      //   const newPet = {
      //     id: Math.floor(Math.random() * 1000),
      //     ...values,
      //     age: Number(values.age),
      //   };

      setPets([...pets, data]);
      message.success('Питомец успешно создан!');
      setIsModalOpen(false);
      form.resetFields();
    } catch (error) {
      message.error('Ошибка при создании питомца');
    } finally {
      setLoading(false);
    }
  };

  const handleUpdatePet = async (values: any) => {
    if (!editingPet) return;

    try {
      setLoading(true);
      const response = await apiInstance.put<Pet>(
        `pets/update/${editingPet.id}`,
        {
          ...values,
          age: Number(values.age),
        }
      );
      const updatedPet = response.data;

      setPets(pets.map((p) => (p.id === updatedPet.id ? updatedPet : p)));
      message.success('Питомец успешно обновлен!');
      setIsEditModalOpen(false);
      editForm.resetFields();
      setEditingPet(null);
    } catch (error) {
      message.error('Ошибка при обновлении питомца');
    } finally {
      setLoading(false);
    }
  };

  const handleFindPet = async () => {
    if (!searchId) {
      message.warning('Введите ID питомца');
      return;
    }

    try {
      setLoading(true);
      const response = await apiInstance.get(`/find_pet/${searchId}`);
      const foundPet = response.data;

      if (foundPet) {
        setPet(foundPet);
        message.success('Питомец найден!');
      } else {
        setPet(null);
        message.warning('Питомец не найден');
      }
    } catch (error) {
      message.error('Ошибка при поиске питомца');
    } finally {
      setLoading(false);
    }
  };

  const columns: TableProps<Pet>['columns'] = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: 'Вид',
      dataIndex: 'species',
      key: 'species',
    },
    {
      title: 'Кличка',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Порода',
      dataIndex: 'breed',
      key: 'breed',
    },
    {
      title: 'Пол',
      dataIndex: 'gender',
      key: 'gender',
    },
    {
      title: 'Возраст',
      dataIndex: 'age',
      key: 'age',
    },
    {
      title: 'Действия',
      key: 'actions',
      render: (_, record) => (
        <Button
          type="link"
          icon={<EditOutlined />}
          onClick={() => showEditModal(record)}
        />
      ),
    },
  ];

  return (
    <div style={{ padding: 24 }}>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          marginBottom: 16,
        }}
      >
        <div style={{ display: 'flex', gap: 8 }}>
          <Input
            placeholder="Введите ID питомца"
            value={searchId}
            onChange={(e) => setSearchId(e.target.value)}
            style={{ width: 200 }}
          />
          <Button
            type="primary"
            icon={<SearchOutlined />}
            onClick={handleFindPet}
            loading={loading}
          >
            Найти
          </Button>
        </div>

        <Button type="primary" icon={<PlusOutlined />} onClick={showModal}>
          Добавить питомца
        </Button>
      </div>

      {pet && (
        <div style={{ marginBottom: 24 }}>
          <h3>Найденный питомец:</h3>
          <p>ID: {pet.id}</p>
          <p>Вид: {pet.species}</p>
          <p>Кличка: {pet.name}</p>
          <p>Порода: {pet.breed}</p>
          <p>Пол: {pet.gender}</p>
          <p>Возраст: {pet.age}</p>
        </div>
      )}

      <Table
        columns={columns}
        dataSource={pets}
        rowKey="id"
        loading={loading}
      />

      <Modal
        title="Добавить нового питомца"
        open={isModalOpen}
        onCancel={handleCancel}
        footer={null}
      >
        <Form form={form} layout="vertical" onFinish={handleCreatePet}>
          <Form.Item
            name="species"
            label="Вид животного"
            rules={[{ required: true, message: 'Пожалуйста, укажите вид' }]}
          >
            <Input placeholder="Собака, кошка и т.д." />
          </Form.Item>

          <Form.Item
            name="name"
            label="Кличка"
            rules={[{ required: true, message: 'Пожалуйста, укажите кличку' }]}
          >
            <Input placeholder="Введите кличку" />
          </Form.Item>

          <Form.Item
            name="breed"
            label="Порода"
            rules={[{ required: true, message: 'Пожалуйста, укажите породу' }]}
          >
            <Input placeholder="Введите породу" />
          </Form.Item>

          <Form.Item
            name="gender"
            label="Пол"
            rules={[{ required: true, message: 'Пожалуйста, укажите пол' }]}
          >
            <Select placeholder="Выберите пол">
              <Select.Option value="male">Самец</Select.Option>
              <Select.Option value="female">Самка</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="user_id"
            label="Пользователь"
            rules={[
              { required: true, message: 'Пожалуйста, укажите пользователя' },
            ]}
          >
            <Select placeholder="Выберите пользователя">
              {users.map((user) => (
                <Select.Option value={user.id}>{user.firstName}</Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="age"
            label="Возраст (лет)"
            rules={[{ required: true, message: 'Пожалуйста, укажите возраст' }]}
          >
            <Input type="number" min={0} placeholder="Введите возраст" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              Создать
            </Button>
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        title={`Редактировать питомца: ${editingPet?.name || ''}`}
        open={isEditModalOpen}
        onCancel={handleEditCancel}
        footer={null}
      >
        <Form form={editForm} layout="vertical" onFinish={handleUpdatePet}>
          <Form.Item
            name="species"
            label="Вид животного"
            rules={[{ required: true, message: 'Пожалуйста, укажите вид' }]}
          >
            <Input placeholder="Собака, кошка и т.д." />
          </Form.Item>

          <Form.Item
            name="name"
            label="Кличка"
            rules={[{ required: true, message: 'Пожалуйста, укажите кличку' }]}
          >
            <Input placeholder="Введите кличку" />
          </Form.Item>

          <Form.Item
            name="breed"
            label="Порода"
            rules={[{ required: true, message: 'Пожалуйста, укажите породу' }]}
          >
            <Input placeholder="Введите породу" />
          </Form.Item>

          <Form.Item
            name="gender"
            label="Пол"
            rules={[{ required: true, message: 'Пожалуйста, укажите пол' }]}
          >
            <Select placeholder="Выберите пол">
              <Select.Option value="male">Самец</Select.Option>
              <Select.Option value="female">Самка</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="user_id"
            label="Пользователь"
            rules={[
              { required: true, message: 'Пожалуйста, укажите пользователя' },
            ]}
          >
            <Select placeholder="Выберите пользователя">
              {users.map((user) => (
                <Select.Option key={user.id} value={user.id}>
                  {user.firstName} {user.lastName}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="age"
            label="Возраст (лет)"
            rules={[{ required: true, message: 'Пожалуйста, укажите возраст' }]}
          >
            <Input type="number" min={0} placeholder="Введите возраст" />
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

export default PetsPage;
