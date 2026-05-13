import { API_ROUTER_URL, IDP_URL } from "../../constants";

type GetBooksByLibraryType = {
  libraryUid: string;
  page?: number;
  size?: number;
};

const DEFAULT_PAGE = 1;
const DEFAULT_PAGE_SIZE = 10;

async function getBooksByLibrary({
  libraryUid,
  page,
  size,
}: GetBooksByLibraryType) {
  const url = `${API_ROUTER_URL}/libraries/${libraryUid}/books/?page=${page || DEFAULT_PAGE}&size=${size || DEFAULT_PAGE_SIZE}&showAll=true`;

  try {
    const response: Response = await fetch(url, {
      method: "GET",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    });

    return { ...(await response.json()), libraryUid };
  } catch (e) {
    console.error(e);
  }
}

export default getBooksByLibrary;
