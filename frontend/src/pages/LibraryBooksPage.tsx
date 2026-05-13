import { ChangeEvent, useMemo, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { Button, Card, Input, Table, Typography, Tag } from "antd";
import { Book, Library, User } from "../types";

type Props = {
  libraries: Library[];
  books: Book[];
  currentUser: User | null;
  onReserve: (bookId: string) => void;
};

const { Title, Text } = Typography;

function LibraryBooksPage({ libraries, books, currentUser, onReserve }: Props) {
  const { libraryUid } = useParams<{ libraryUid: string }>();
  const [search, setSearch] = useState("");

  const library = useMemo(
    () => libraries.find((item) => item.libraryUid === libraryUid),
    [libraryUid, libraries],
  );
  const libraryBooks = useMemo(
    () => books.filter((book) => book.libraryUid === libraryUid),
    [books, libraryUid],
  );

  const filteredBooks = useMemo(
    () =>
      libraryBooks.filter((book) =>
        book.name.toLowerCase().includes(search.toLowerCase()),
      ),
    [libraryBooks, search],
  );

  if (!library) {
    return (
      <div className="page-card">
        <Title level={2}>Библиотека не найдена</Title>
        <Link to="/libraries">
          <Button>Вернуться к списку библиотек</Button>
        </Link>
      </div>
    );
  }

  const columns = [
    { title: "Название", dataIndex: "name", key: "name" },
    { title: "Автор", dataIndex: "author", key: "author" },
    { title: "Жанр", dataIndex: "genre", key: "genre" },
    {
      title: "Количество",
      dataIndex: "availableCount",
      key: "availableCount",
    },
    {
      title: "Состояние",
      dataIndex: "condition",
      key: "condition",
      render: (status: string) => {
        const color =
          status === "EXCELLENT"
            ? "green"
            : status === "GOOD"
              ? "blue"
              : "orange";
        const label =
          status === "EXCELLENT"
            ? "Отлично"
            : status === "GOOD"
              ? "Хорошо"
              : "Плохо";
        return <Tag color={color}>{label}</Tag>;
      },
    },
    {
      title: "Статус",
      dataIndex: "status",
      key: "status",
      render: (status: string) => {
        const color =
          status === "available"
            ? "green"
            : status === "reserved"
              ? "orange"
              : "blue";
        const label =
          status === "available"
            ? "Доступна"
            : status === "reserved"
              ? "Забронирована"
              : "Выдана";
        return <Tag color={color}>{label}</Tag>;
      },
    },
    {
      title: "Действие",
      key: "action",
      render: (_: unknown, record: Book) => {
        if (
          (currentUser?.role === "user" || currentUser?.role === "admin") &&
          record.status === "available"
        ) {
          return (
            <Button type="primary" onClick={() => onReserve(record.bookUid)}>
              Забронировать
            </Button>
          );
        }
        if (record.status === "reserved" || record.status === "issued") {
          return <Text type="secondary">{record.reservedBy}</Text>;
        }
        return <Text type="secondary">Требуется авторизация</Text>;
      },
    },
  ];

  return (
    <div className="page-card">
      <Title level={2}>{library.name}</Title>
      <Text type="secondary">{library.address}</Text>
      <Input.Search
        placeholder="Поиск по книге"
        value={search}
        onChange={(event: ChangeEvent<HTMLInputElement>) =>
          setSearch(event.target.value)
        }
        allowClear
        enterButton="Поиск"
        style={{ margin: "24px 0" }}
      />
      <Table
        columns={columns}
        dataSource={filteredBooks}
        rowKey="id"
        pagination={false}
      />
      <div style={{ marginTop: 24 }}>
        <Link to="/libraries">
          <Button>Вернуться к библиотекам</Button>
        </Link>
      </div>
    </div>
  );
}

export default LibraryBooksPage;
