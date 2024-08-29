"use client";
import { DealCard } from "./card";
import { useEffect, useState } from "react";
import { useUser } from "@/context/userContext";
import axios from "axios";
import { BASE_URL, DEFAULT_PAGE_SIZE } from "@/lib/utils";
import { useDeal } from "@/context/dealContext";

type Props = {
  page: string;
};

// CardSection: card component for deals and home page of the dashboard
export function CardSection({ page }: Props) {
  const { usr } = useUser();
  const [deals, setDeals] = useState<Deal[]>([]);
  const { dealParam } = useDeal();

  /* to fetch deals from client side, using use effect */
  useEffect(() => {
    async function getDeals(
      token: string | undefined,
      url: string,
      param: Status | DealFilter
    ) {
      /* defined an async function */
      try {
        if (!token) {
          throw new Error("user not signed in, no access token");
        }
        let resp = await axios({
          method: "get",
          baseURL: BASE_URL,
          url: url,
          params: param,
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        let deals: Deal[] = await resp.data;
        setDeals(deals);
      } catch (error: any) {
        console.log(error);
      }
    }
    if (page === "home") {
      getDeals(usr?.access_token, "a/deals/vas", { status: "ongoing" });
    } else if (page === "deals") {
      const param =
        dealParam != null
          ? dealParam
          : {
              customer_name: null,
              service_to_render: null,
              status: null,
              max_profit: null,
              min_profit: null,
              awarded: null,
              sales_rep_name: null,
              page_size: DEFAULT_PAGE_SIZE,
              page_id: 1,
            };
      getDeals(usr?.access_token, "a/deals/filtered", param);
    }
  }, [usr, page, dealParam]);
  return (
    <div className="flex flex-col gap-3  p-3 w-full grow h-auto">
      <section className="flex flex-col w-full h-full gap-3 ">
        {deals.length < 1 ? (
          <>no deals</>
        ) : (
          deals.map((deal) => {
            return <DealCard key={deal.deal_id} deal={deal} />;
          })
        )}
      </section>
    </div>
  );
}
