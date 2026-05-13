import { API_ROUTER_URL } from "../../constants";

type ReturnReservationType = {
  reservationUid: string;
  reservationBody: {
    condition: string;
    date: string;
  };
};

async function returnReservation({
  reservationUid,
  reservationBody,
}: ReturnReservationType) {
  const url = `${API_ROUTER_URL}/reservations/${reservationUid}/return`;
  const body = JSON.stringify(reservationBody);

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

    return {};
  } catch (e) {
    console.error(e);
    return undefined;
  }
}

export default returnReservation;
