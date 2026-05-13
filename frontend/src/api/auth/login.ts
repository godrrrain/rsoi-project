import { CLIENT_ID, IDP_URL, REDIRECT_URI } from "../../constants";

type LoginType = {
  code: string;
};

async function login({ code }: LoginType) {
  try {
    const params = new URLSearchParams({
      grant_type: "authorization_code",
      code: code,
      redirect_uri: REDIRECT_URI,
      client_id: CLIENT_ID,
    });
    console.log("[login] Sending to token endpoint:", {
      code: code,
      redirect_uri: REDIRECT_URI,
      client_id: CLIENT_ID,
      grant_type: "authorization_code",
    });

    const response = await fetch(`${IDP_URL}/oauth2/token`, {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: params.toString(),
    });

    console.log("[login] Token response status:", response.status);
    const responseText = await response.text();
    console.log("[login] Token response body:", responseText);

    if (!response.ok) {
      throw new Error("Failed to get token: " + responseText);
    }

    const data = JSON.parse(responseText);
    localStorage.setItem("access_token", data.access_token);
    localStorage.setItem("id_token", data.id_token);
    if (data.refresh_token) {
      localStorage.setItem("refresh_token", data.refresh_token);
    }

    // Clear URL and reload
    // window.history.replaceState({}, document.title, "/");

    return data;
  } catch (e) {
    console.error(e);
    return undefined;
  }
}

export default login;
