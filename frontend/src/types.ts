export type Role = "guest" | "user" | "admin";

export type User = {
  id: string;
  name: string;
  email: string;
  password: string;
  role: Role;
  rating: number;
};

export type BookStatus = "available" | "reserved" | "issued";

export type Book = {
  bookUid: string;
  libraryUid: string;
  name: string;
  author: string;
  genre: string;
  condition: string;
  availableCount: number;
  status: BookStatus;
  reservedBy?: string;
  reservedById?: string;
  reservationUid?: string;
};

export type Library = {
  libraryUid: string;
  name: string;
  address: string;
  city: string;
};

export type Reservation = {
  reservationUid: string;
  status: string;
  startDate: string;
  tillDate: string;
  book: { bookUid: string; name: string; author: string; genre: string };
  library: Library;
  username: {
    name: string;
  };
};
