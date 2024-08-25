'use client';

import React, { createContext, useContext, useState } from 'react';

type UserContextType = {
  user: { userId: number; username: string; fullname: string; email: string } | null;
  setUser: (user: { userId: number; username: string; fullname: string; email: string }) => void;
};

const UserContext = createContext<UserContextType | undefined>(undefined);

export const UserProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<UserContextType['user']>(null);

  return (
    <UserContext.Provider value={{ user, setUser }}>
      {children}
    </UserContext.Provider>
  );
};

export const useUser = () => {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error('useUser must be used within a UserProvider');
  }
  return context;
};