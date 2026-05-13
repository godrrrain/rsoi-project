import { API_ROUTER_URL } from "../../constants";

const MOCK_DATA = [
  {
    reservationUid: "f464ca3a-fcf7-4e3f-86f0-76c7bba96f72",
    status: "RENTED",
    startDate: "2021-10-09",
    tillDate: "2021-10-11",
    book: {
      bookUid: "f7cdc58f-2caf-4b15-9727-f89dcc629b27",
      name: "Краткий курс C++ в 7 томах",
      author: "Бьерн Страуструп",
      genre: "Научная фантастика",
    },
    library: {
      libraryUid: "83575e12-7ce0-48ee-9931-51919ff3c9ee",
      name: "Библиотека имени 7 Непьющих",
      address: "2-я Бауманская ул., д.5, стр.1",
      city: "Москва",
    },
    username: {
      name: "Алексей",
      rating: 75,
    },
  },
  {
    reservationUid: "f464ca3a-fcf7-4e3f-86f0-76c7bba96f73",
    status: "RENTED",
    startDate: "2021-10-09",
    tillDate: "2021-10-11",
    book: {
      bookUid: "f7cdc58f-2caf-4b15-9727-f89dcc629b27",
      name: "Краткий курс C++ в 7 томах",
      author: "Бьерн Страуструп",
      genre: "Научная фантастика",
    },
    library: {
      libraryUid: "83575e12-7ce0-48ee-9931-51919ff3c9ee",
      name: "Библиотека имени 7 Непьющих",
      address: "2-я Бауманская ул., д.5, стр.1",
      city: "Москва",
    },
    username: {
      name: "Мария",
      rating: 75,
    },
  },
];

async function getUserReservations() {
  const url = `${API_ROUTER_URL}/reservations`;

  const token = localStorage.getItem("access_token");

  try {
    const response: Response = await fetch(url, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${token}`,
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    });

    return await response.json();
  } catch (e) {
    console.error(e);
  }
}

export default getUserReservations;
