import { useCallback, useEffect, useState } from "react";
import {
  Link,
  Navigate,
  Route,
  Routes,
  useLocation,
  useNavigate,
} from "react-router-dom";
import { Layout, Menu, Button, Space, Typography, message } from "antd";
import {
  HomeOutlined,
  LoginOutlined,
  LogoutOutlined,
  TeamOutlined,
  UserOutlined,
  BarChartOutlined,
} from "@ant-design/icons";
import { Book, Library, Reservation, User } from "./types";
import LibrariesPage from "./pages/LibrariesPage";
import LibraryBooksPage from "./pages/LibraryBooksPage";
import ReaderPage from "./pages/ReaderPage";
import LibrarianPage from "./pages/LibrarianPage";
import StatisticsPage from "./pages/StatisticsPage";
import { CLIENT_ID, DEFAULT_CITY, IDP_URL, REDIRECT_URI } from "./constants";
import getLibraries from "./api/libraries/getLibraries";
import getBooksByLibrary from "./api/libraries/getBooksByLibrary";
import getUserReservations from "./api/reservations/getUserReservations";
import returnReservation from "./api/reservations/returnReservation";
import createReservations from "./api/reservations/createReservations";
import { generateRandomString } from "./helpers/generateRandomString";
import getUserInfo from "./api/auth/getUserInfo";
import login from "./api/auth/login";
import { parseJWT } from "./helpers/parseJWT";
import getUserReservationsAll from "./api/reservations/getUserReservationsAll";
import getRating from "./api/auth/getRating";

const { Header, Content } = Layout;
const { Text } = Typography;

function App() {
  const [libraries, setLibraries] = useState<Library[]>([]);
  const [books, setBooks] = useState<Book[]>([]);
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [selectedCity, setSelectedCity] = useState<string | undefined>(
    DEFAULT_CITY,
  );

  const location = useLocation();
  const navigate = useNavigate();

  const urlParams = new URLSearchParams(window.location.search);
  const code = urlParams.get("code");
  const state = urlParams.get("state");

  const checkToken = useCallback(async (token: string) => {
    getUserInfo({ token }).then((userInfoResponse) => {
      if (userInfoResponse) {
        getRating().then((ratingResponse) => {
          setCurrentUser({
            id: userInfoResponse.sub,
            name: userInfoResponse.name,
            email: userInfoResponse.email,
            password: "",
            role: parseJWT(token)?.role || "user",
            rating: ratingResponse?.stars || 0,
          });
          message.success(`Вход выполнен как ${userInfoResponse.name}`);

          navigate(
            userInfoResponse?.role === "admin" ? "/librarian" : "/reader",
          );
        });
      } else {
        message.error("Ошибка входа");
      }
    });
  }, []);

  // Load user from localStorage
  useEffect(() => {
    if (code && state) {
      login({ code }).then((response) => {
        if (response) {
          checkToken(response.access_token);
        } else {
          message.error("Ошибка входа");
        }
      });
    } else {
      const token = localStorage.getItem("access_token");
      if (token) {
        checkToken(token);
      }
    }
  }, []);

  // Load libraries on mount or city change
  useEffect(() => {
    getLibraries({ city: selectedCity }).then((response) => {
      if (response) {
        setLibraries(response.items);
      }
    });
  }, [selectedCity]);

  // Load reservations for current user
  useEffect(() => {
    if (currentUser) {
      if (currentUser.role === "admin") {
        getUserReservationsAll().then((response) => {
          if (response) {
            setReservations(response);
          }
        });
      } else {
        getUserReservations().then((response) => {
          if (response) {
            setReservations(response);
          }
        });
      }
    } else {
      setReservations([]);
    }
  }, [currentUser]);

  const selectedKey = location.pathname.startsWith("/libraries")
    ? "/libraries"
    : location.pathname;

  const menuItems = [
    {
      label: <Link to="/libraries">Библиотеки</Link>,
      key: "/libraries",
      icon: <HomeOutlined />,
    },
    currentUser?.role === "user" && {
      label: <Link to="/reader">Читатель</Link>,
      key: "/reader",
      icon: <UserOutlined />,
    },
    currentUser?.role === "admin" && {
      label: <Link to="/librarian">Библиотекарь</Link>,
      key: "/librarian",
      icon: <TeamOutlined />,
    },
    currentUser?.role === "admin" && {
      label: <Link to="/statistics">Статистика</Link>,
      key: "/statistics",
      icon: <BarChartOutlined />,
    },
  ].filter(Boolean) as any[];

  const handleLogin = () => {
    const state = generateRandomString();
    sessionStorage.setItem("oauth_state", state);

    const params = new URLSearchParams({
      client_id: CLIENT_ID,
      redirect_uri: REDIRECT_URI,
      response_type: "code",
      scope: "openid profile email",
      state: state,
    });

    window.location.href = `${IDP_URL}/oauth2/authorize?${params.toString()}`;
  };

  const handleLogout = () => {
    localStorage.removeItem("access_token");
    localStorage.removeItem("id_token");
    localStorage.removeItem("refresh_token");

    setCurrentUser(null);
    message.info("Вы вышли из системы");
    navigate("/libraries");
  };

  const handleReserve = async (bookUid: string) => {
    if (
      !currentUser ||
      (currentUser.role !== "user" && currentUser.role !== "admin")
    ) {
      message.warning("Только авторизованный читатель может бронировать");
      return;
    }

    try {
      const libraryUid =
        books.find((b) => b.bookUid === bookUid)?.libraryUid || "";
      const tillDate = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000)
        .toISOString()
        .split("T")[0];

      const reservation = await createReservations({
        bookUid,
        libraryUid,
        tillDate,
      });

      if (reservation) {
        setReservations((prev) => [
          ...prev,
          { ...reservation, username: { name: currentUser.name } },
        ]);

        // Update books state
        setBooks((prev) =>
          prev.map((book) =>
            book.bookUid === bookUid
              ? {
                  ...book,
                  status: "reserved",
                  reservedBy: currentUser.name,
                  reservedById: currentUser.id,
                  reservationUid: reservation.reservationUid,
                }
              : book,
          ),
        );
        message.success("Книга забронирована");
      } else {
        throw new Error("Ошибка бронирования");
      }
    } catch (error: any) {
      message.error(error.message || "Ошибка бронирования");
    }
  };

  const handleReturn = async (reservationUid: string, condition: string) => {
    if (!currentUser) return;

    try {
      await returnReservation({
        reservationUid,
        reservationBody: {
          condition,
          date: new Date(Date.now()).toISOString().split("T")[0],
        },
      });

      setReservations((prev) =>
        prev.map((r) =>
          r.reservationUid === reservationUid
            ? { ...r, status: "RETURNED" }
            : r,
        ),
      );

      // Update books state
      const bookId = reservationUid.replace("res-", "");
      setBooks((prev) =>
        prev.map((book) =>
          book.bookUid === bookId
            ? {
                ...book,
                status: "available",
                reservedBy: undefined,
                reservedById: undefined,
                reservationUid: undefined,
              }
            : book,
        ),
      );
      message.success("Книга возвращена");
    } catch (error: any) {
      message.error(error.message || "Ошибка возврата");
    }
  };

  // Load books for libraries
  useEffect(() => {
    if (libraries.length > 0) {
      Promise.all(
        libraries.map((lib) =>
          getBooksByLibrary({ libraryUid: lib.libraryUid }),
        ),
      ).then((responses) => {
        const allBooks = responses.flatMap((response) => {
          if (!response) return [];

          const { items, libraryUid } = response;
          if (!items) return [];
          return items.map((item: Book) => ({
            bookUid: item.bookUid,
            libraryUid,
            name: item.name,
            author: item.author,
            genre: item.genre,
            condition: item.condition,
            availableCount: item.availableCount,
            status:
              item.availableCount > 0
                ? "available"
                : ("reserved" as Book["status"]),
          }));
        });

        setBooks(allBooks);
      });
    }
  }, [libraries, reservations]);

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Header className="app-header">
        <div className="app-title">Library Booking</div>
        <Menu
          theme="dark"
          mode="horizontal"
          selectedKeys={[selectedKey]}
          items={menuItems}
          style={{ flex: 1 }}
        />
        <Space>
          {currentUser ? (
            <>
              <Text style={{ color: "#fff" }}>{currentUser.name}</Text>
              <Button
                icon={<LogoutOutlined />}
                type="primary"
                onClick={handleLogout}
              >
                Выйти
              </Button>
            </>
          ) : (
            <>
              <Button
                icon={<LoginOutlined />}
                type="primary"
                onClick={handleLogin}
              >
                Войти через Identity Provider
              </Button>
            </>
          )}
        </Space>
      </Header>

      <Content
        style={{
          padding: "24px",
          maxWidth: 1200,
          margin: "0 auto",
          width: "100%",
        }}
      >
        <Routes>
          <Route path="/" element={<Navigate to="/libraries" replace />} />
          <Route
            path="/libraries"
            element={
              <LibrariesPage
                libraries={libraries}
                selectedCity={selectedCity}
                onCityChange={setSelectedCity}
              />
            }
          />
          <Route
            path="/libraries/:libraryUid/books"
            element={
              <LibraryBooksPage
                libraries={libraries}
                books={books}
                currentUser={currentUser}
                onReserve={handleReserve}
              />
            }
          />
          <Route
            path="/reader"
            element={
              currentUser?.role === "user" ? (
                <ReaderPage
                  currentUser={currentUser}
                  reservations={reservations}
                />
              ) : (
                <Navigate to="/libraries" replace />
              )
            }
          />
          <Route
            path="/librarian"
            element={
              currentUser?.role === "admin" ? (
                <LibrarianPage
                  reservations={reservations}
                  onReturn={handleReturn}
                />
              ) : (
                <Navigate to="/libraries" replace />
              )
            }
          />
          <Route
            path="/statistics"
            element={
              currentUser?.role === "admin" ? (
                <StatisticsPage />
              ) : (
                <Navigate to="/libraries" replace />
              )
            }
          />
          <Route path="*" element={<Navigate to="/libraries" replace />} />
        </Routes>
      </Content>
    </Layout>
  );
}

export default App;
