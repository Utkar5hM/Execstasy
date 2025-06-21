"use client";
import React, { createContext, useContext, useState, useEffect } from "react";

import { jwtDecode } from "jwt-decode";
import CryptoJS from "crypto-js";
import Cookies from 'js-cookie';

type AuthJwtPayload = {
  name?: string;
  username?: string;
  email?: string;
};

type User = {
  name: string;
  username: string;
  avatar: string;
};

const AuthContext = createContext<{
  user: User;
  setUser: React.Dispatch<React.SetStateAction<User>>;
  refreshUser: () => void;
}>({
  user: {
    name: "",
    username: "",
    avatar: "",
  },
  setUser: () => {},
  refreshUser: () => {},
});

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<User>({
    name: "",
    username: "",
    avatar: "",
  });

  const refreshUser = () => {
    const jwt = Cookies.get('jwt');
    if (jwt) {
      try {
        const decoded = jwtDecode<AuthJwtPayload>(jwt);
        const avatar = decoded.email
          ? `https://www.gravatar.com/avatar/${CryptoJS.MD5(
              decoded.email.trim().toLowerCase()
            ).toString()}`
          : "";
        setUser({
          name: decoded.name || "",
          username: decoded.username || "",
          avatar,
        });
      } catch {
        setUser({
          name: "",
          username: "",
          avatar: "",
        });
      }
    } else {
      setUser({
        name: "",
        username: "",
        avatar: "",
      });
    }
  };

  useEffect(() => {
    refreshUser();
  }, []);

  return (
    <AuthContext.Provider value={{ user, setUser, refreshUser }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);