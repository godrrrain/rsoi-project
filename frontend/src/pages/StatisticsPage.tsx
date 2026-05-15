import { useEffect, useState } from "react";
import { Card, List, Typography } from "antd";
import getStats from "../api/stats/getStats";

const { Title, Text } = Typography;

type Stats = Record<string, number>;

function StatisticsPage() {
  const [stats, setStats] = useState<Stats | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    getStats()
      .then((response) => {
        if (!response || typeof response !== "object") {
          throw new Error("Не удалось загрузить статистику");
        }
        setStats(response as Stats);
      })
      .catch((e) => {
        setError((e as Error).message || "Ошибка загрузки статистики");
      });
  }, []);

  return (
    <div className="page-card">
      <Title level={2}>Статистика</Title>
      <Card>
        {error ? (
          <Text type="danger">{error}</Text>
        ) : !stats || Object.keys(stats).length === 0 ? (
          <Text>Статистика пока отсутствует.</Text>
        ) : (
          <List
            dataSource={Object.entries(stats)}
            renderItem={([metric, value]) => (
              <List.Item
                style={{
                  display: "flex",
                  justifyContent: "flex-start",
                  gap: "16px",
                }}
              >
                <Text strong>{metric}</Text>
                <Text>{value}</Text>
              </List.Item>
            )}
          />
        )}
      </Card>
    </div>
  );
}

export default StatisticsPage;
