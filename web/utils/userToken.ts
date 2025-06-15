import { jwtDecode } from "jwt-decode";

interface CustomJwtPayload {
	role?: string;
	name?: string;
	email?: string;
	username?: string;
	id?: number;
  }

  
export const decodeJwt = (): { role?: string; name?: string; email?: string; username?: string; id?: number } | null => {
  if (typeof document !== "undefined") {
    const jwt = document.cookie
      .split("; ")
      .find((row) => row.startsWith("jwt="))
      ?.split("=")[1];

    if (jwt) {
      try {
        const decodedToken = jwtDecode<CustomJwtPayload>(jwt);
        return {
          role: decodedToken.role,
          name: decodedToken.name,
          email: decodedToken.email,
		  username: decodedToken.username,
		  id: decodedToken.id,
        };
      } catch (error) {
        console.error("Failed to decode JWT:", error);
        return null;
      }
    }
  }
  return null;
};