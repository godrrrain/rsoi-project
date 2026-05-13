import { DEFAULT_CITY, API_ROUTER_URL } from "../../constants";

type GetLibrariesType = {
  page?: number;
  size?: number;
  city?: string;
};

const DEFAULT_PAGE = 1;
const DEFAULT_PAGE_SIZE = 10;

async function getLibraries({ page, size, city }: GetLibrariesType) {
  const url = `${API_ROUTER_URL}/libraries?page=${page || DEFAULT_PAGE}&size=${size || DEFAULT_PAGE_SIZE}&city=${city || DEFAULT_CITY}`;

  try {
    const response: Response = await fetch(url, {
      method: "GET",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    });

    return await response.json();
  } catch (e) {
    console.error(e);
  }
}

export default getLibraries;
