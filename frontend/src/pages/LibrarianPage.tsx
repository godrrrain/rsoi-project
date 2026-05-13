import { ChangeEvent, useState } from "react";
import {
  Button,
  Form,
  Input,
  Modal,
  Space,
  Table,
  Typography,
  Tag,
} from "antd";
import { Reservation } from "../types";

type Props = {
  reservations: Reservation[];
  onReturn: (reservationUid: string, condition: string) => void;
};

const { Title } = Typography;

function LibrarianPage({ reservations, onReturn }: Props) {
  const [search, setSearch] = useState("");
  const [returnModalOpen, setReturnModalOpen] = useState(false);
  const [returnReservationUid, setReturnReservationUid] = useState<string>("");

  const columns = [
    { title: "Читатель", dataIndex: ["username", "name"], key: "username" },
    { title: "Название", dataIndex: ["book", "name"], key: "name" },
    { title: "Автор", dataIndex: ["book", "author"], key: "author" },
    { title: "Библиотека", dataIndex: ["library", "name"], key: "library" },
    {
      title: "Статус",
      dataIndex: "status",
      key: "status",
      render: (status: string) => {
        const color =
          status === "RENTED"
            ? "orange"
            : status === "RETURNED"
              ? "green"
              : "blue";
        const label =
          status === "RENTED"
            ? "Забронирована"
            : status === "RETURNED"
              ? "Возвращена"
              : "Выдана";
        return <Tag color={color}>{label}</Tag>;
      },
    },
    {
      title: "Действие",
      key: "action",
      render: (_: unknown, record: Reservation) =>
        record.status !== "RETURNED" ? (
          <Space>
            <Button
              onClick={() => {
                setReturnReservationUid(record.reservationUid);
                setReturnModalOpen(true);
              }}
            >
              Вернуть
            </Button>
          </Space>
        ) : null,
    },
  ];

  return (
    <div className="page-card">
      <Title level={2}>Страница библиотекаря</Title>
      <Input.Search
        placeholder="Поиск читателя по имени"
        value={search}
        onChange={(event: ChangeEvent<HTMLInputElement>) =>
          setSearch(event.target.value)
        }
        enterButton="Найти"
        allowClear
        style={{ marginBottom: 24 }}
      />
      <div style={{ marginTop: 24 }}>
        <Title level={4}>Забронированные книги</Title>
        <Table
          columns={columns}
          dataSource={reservations}
          rowKey="reservationUid"
          pagination={false}
          locale={{ emptyText: "Нет забронированных книг" }}
        />
      </div>
      <Modal
        title="Вернуть книгу"
        open={returnModalOpen}
        onCancel={() => setReturnModalOpen(false)}
        footer={null}
      >
        <Form
          layout="vertical"
          onFinish={(values: { condition: string; date: string }) => {
            onReturn(returnReservationUid, values.condition);
            setReturnModalOpen(false);
          }}
        >
          <Form.Item
            label="Состояние книги"
            name="condition"
            rules={[{ required: true, message: "Укажите состояние" }]}
          >
            <Input placeholder="EXCELLENT, GOOD, BAD" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit">
              Подтвердить возврат
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}

export default LibrarianPage;
