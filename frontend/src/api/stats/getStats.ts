import { API_ROUTER_URL } from "../../constants";

async function getStats() {
  const url = `${API_ROUTER_URL}/statistics`;

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

export default getStats;
