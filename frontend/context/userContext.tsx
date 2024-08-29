"use client";

import React, { createContext, useContext, useState } from "react";

type UserContextType = {
  usr: UserResp | null;
  setUser: (usr: UserResp) => void;
};

const UserContext = createContext<UserContextType | undefined>(undefined);

export const UserProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [usr, setUser] = useState<UserContextType["usr"]>(() => {
    //check if user exists in localstorage
    if (typeof window !== "undefined") {
      const storedUser = localStorage.getItem("user");
      return storedUser ? JSON.parse(storedUser) : null;
    }
    return null;
  });

  const updateUser = (user: UserResp) => {
    setUser(user);
    if (user) {
      localStorage.setItem("user", JSON.stringify(user));
    } else {
      localStorage.removeItem("user");
    }
  };
  return (
    <UserContext.Provider value={{ usr, setUser: updateUser }}>
      {children}
    </UserContext.Provider>
  );
};

export const useUser = () => {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error("useUser must be used within a UserProvider");
  }
  return context;
};
