"use client";
import { DEFAULT_PAGE_SIZE } from "@/lib/utils";
import React, { createContext, useContext, useState } from "react";

// Define the type you want the context to hold...
type DealContextType = {
  dealParam: DealFilter | null;
  setDealParam: (dealParam: DealFilter) => void;
};
const DealContext = createContext<DealContextType | undefined>(undefined);

export const DealProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
    const customParam: DealFilter = {
        customer_name: null,
        service_to_render: null,
        status: null,
        max_profit: null,
        min_profit: null,
        awarded: null,
        sales_rep_name: null,
        page_size: DEFAULT_PAGE_SIZE,
        page_id: 1,

      }
    const [dealParam, setDealParam] = useState<DealFilter>(customParam);

  return <DealContext.Provider value={{dealParam, setDealParam}}>{children}</DealContext.Provider>;
};


export const useDeal = () => {
    const context = useContext(DealContext)
    if (!context) {
        throw new Error("useDeal must be used within a DealProvider");
    }
    return context;
}