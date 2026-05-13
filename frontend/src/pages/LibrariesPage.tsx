import { ChangeEvent, useMemo, useState } from "react";
import { Button, Card, Col, Input, Row, Select, Typography } from "antd";
import { Library } from "../types";
import { Link } from "react-router-dom";
import { CITIES } from "../constants";

type Props = {
  libraries: Library[];
  selectedCity: string | undefined;
  onCityChange: (city: string | undefined) => void;
};

const { Title } = Typography;
const { Option } = Select;

function LibrariesPage({ libraries, selectedCity, onCityChange }: Props) {
  const [search, setSearch] = useState("");

  const visibleLibraries = useMemo(
    () =>
      libraries.filter(
        ({ name, city }) =>
          name.toLowerCase().includes(search.toLowerCase()) &&
          (!selectedCity || city === selectedCity),
      ),
    [libraries],
  );

  return (
    <div className="page-card">
      <Title level={2}>Все библиотеки</Title>
      <div style={{ marginBottom: 24, display: "flex", gap: 16 }}>
        <Select
          value={selectedCity}
          onChange={(value) => onCityChange(value || undefined)}
          style={{ width: 200 }}
          placeholder="Выберите город"
          allowClear
        >
          {CITIES.map((city) => (
            <Option key={city} value={city}>
              {city}
            </Option>
          ))}
        </Select>
        <Input.Search
          placeholder="Введите библиотеку"
          value={search}
          onChange={(event: ChangeEvent<HTMLInputElement>) =>
            setSearch(event.target.value)
          }
          allowClear
          enterButton="Поиск"
          style={{ flex: 1 }}
        />
      </div>
      <Row gutter={[16, 16]}>
        {visibleLibraries.map(({ libraryUid, name, address }) => (
          <Col xs={24} sm={12} md={8} key={libraryUid}>
            <Card
              title={name}
              actions={[
                <Link key="select" to={`/libraries/${libraryUid}/books`}>
                  <Button type="link">Выбрать</Button>
                </Link>,
              ]}
            >
              <Card.Meta description={address} />
            </Card>
          </Col>
        ))}
      </Row>
      {visibleLibraries.length === 0 && (
        <Typography.Text>Ничего не найдено.</Typography.Text>
      )}
    </div>
  );
}

export default LibrariesPage;
