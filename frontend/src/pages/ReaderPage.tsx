import { Button, Card, Rate, Table, Typography, Tag } from "antd";
import { Reservation, User } from "../types";

type Props = {
  currentUser: User;
  reservations: Reservation[];
};

const { Title, Text } = Typography;

function ReaderPage({ currentUser, reservations }: Props) {
  const columns = [
    { title: "Название", dataIndex: ["book", "name"], key: "title" },
    { title: "Автор", dataIndex: ["book", "author"], key: "author" },
    { title: "Библиотека", dataIndex: ["library", "name"], key: "library" },
    {
      title: "Статус",
      dataIndex: "status",
      key: "status",
      render: (status: string) => {
        const color = status === "RENTED" ? "orange" : "blue";
        const label = status === "RENTED" ? "Забронирована" : "Выдана";
        return <Tag color={color}>{label}</Tag>;
      },
    },
  ];

  return (
    <div className="page-card">
      <Title level={2}>Страница читателя</Title>
      <Card style={{ marginBottom: 24 }}>
        <Title level={4}>{currentUser.name}</Title>
        <Text>Рейтинг читателя:</Text>
        <div style={{ marginTop: 8 }}>
          <Rate disabled value={currentUser.rating} />
          <Text style={{ marginLeft: 12 }}>
            {currentUser.rating.toFixed(1)}
          </Text>
        </div>
      </Card>
      <Title level={4}>Все взятые книги</Title>
      <Table
        columns={columns}
        dataSource={reservations}
        rowKey="reservationUid"
        pagination={false}
        locale={{ emptyText: "У вас нет взятых книг" }}
      />
    </div>
  );
}

export default ReaderPage;
