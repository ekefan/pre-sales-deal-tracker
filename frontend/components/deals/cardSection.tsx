// import { useEffect } from "react";
"use client";
import { useUser } from "@/context/userContext";
import { DealCard } from "./card";
import axios from "axios";
import { useEffect, useState} from "react";
import { useRouter } from "next/navigation";


type Props = {
  page: string;
};

// async function getdeals fetches deals from the database if the user is signed
async function getDeals(token: string | undefined, url: string) {
    try {
      if (!token) {
        throw new Error("user not signed in, no access token");
      }
      const response = await axios({
        method: "get",
        baseURL: "http://localhost:8080",
        url: url,
        params: {
          status: "closed",
        },
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      return response
    } catch (error: any) {
      console.log(error);
    }
  }

// CardSection: card component for deals and home page of the dashboard
export function CardSection({ page }: Props) {
  const { usr } = useUser();
  const router = useRouter();
//   const {deal, setDeal} = useState([])

    if (page === "home") {
        const deals = getDeals(usr?.access_token, "a/deals/vas");
        console.log(deals)
      } else if (page === "deals") {
        console.log("fetch all deals");
      }
  return (
    <div className="flex flex-col gap-3  p-3 w-full grow h-auto">
      <section className="flex flex-col w-full h-full gap-3 ">
        {}
        <DealCard />
        <DealCard />
        <DealCard />
        <DealCard />
      </section>
    </div>
  );
}

