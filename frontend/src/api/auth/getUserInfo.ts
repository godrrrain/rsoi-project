import { IDP_URL } from "../../constants";

type GetUserInfoType = {
  token: string;
};

async function getUserInfo({ token }: GetUserInfoType) {
  const url = `${IDP_URL}/oauth2/userinfo`;

  try {
    const response = await fetch(url, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      throw new Error("Ошибка входа");
    }

    return await response.json();
  } catch (e) {
    console.error(e);

    localStorage.removeItem("access_token");
    localStorage.removeItem("id_token");
    localStorage.removeItem("refresh_token");

    return undefined;
  }
}

export default getUserInfo;
