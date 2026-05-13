import { API_ROUTER_URL } from "../../constants";

type CreateReservationsType = {
  bookUid: string;
  libraryUid: string;
  tillDate: string;
};

async function createReservations(reservation: CreateReservationsType) {
  const url = `${API_ROUTER_URL}/reservations`;
  const body = JSON.stringify(reservation);

  const token = localStorage.getItem("access_token");

  try {
    const response: Response = await fetch(url, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: body,
    });

    if (response.ok) {
      return await response.json();
    }

    return undefined;
  } catch (e) {
    console.error(e);
    return undefined;
  }
}

export default createReservations;
