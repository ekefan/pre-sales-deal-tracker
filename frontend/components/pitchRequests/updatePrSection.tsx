"use client";
import { UpdateDealForm } from "../deals/updateDealForm";
import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { useUser } from "@/context/userContext";
import { BASE_URL } from "@/lib/utils";
import axios from "axios";
export function UpdatePitchRequestSection() {
  const { usr } = useUser();
  const searchParams = useSearchParams();
  const id = searchParams.get("deal_id");
  const [deal, setDeal] = useState<Deal | null>(null);

  useEffect(() => {
    const getDeal = async (
      token: string | undefined,
      urlPath: string,
      param: { deal_id: number }
    ) => {
      try {
        if (!token) {
          throw new Error("user not signed in, no access token");
        }
        const response = await axios({
          method: "get",
          baseURL: BASE_URL,
          url: urlPath,
          params: param,
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        const dealFromServer = response.data;
        setDeal(dealFromServer);
      } catch (error) {
        console.log(error);
      }
    };
    getDeal(usr?.access_token, "a/admin/getdeal", {
      deal_id: id != null ? Number(id) : 0,
    });
  }, [id, usr]);
  if (deal != null) {
    return (
      <div className="p-3 mt-1 flex flex-col gap-3 w-full border rounded-md bg-gray-100">
        <UpdateDealForm
          deal_id={deal.deal_id}
          servicesToRender={
            deal.services_to_render != undefined ? deal.services_to_render : []
          }
          CurrentPitchRequest={
            deal.current_pitch_request != undefined
              ? deal.current_pitch_request
              : ""
          }
          status={deal.deal_status != undefined ? deal.deal_status : ""}
          department={deal.department != undefined ? deal.department : ""}
          profit={deal.profit != undefined ? Number(deal.profit) : 0}
          netTotalCost={
            deal.net_total_cost != undefined ? Number(deal.net_total_cost) : 0
          }
          awarded={deal.awarded != undefined ? deal.awarded : false}
        />
      </div>
    );
  } else {
    return <div>{`can't update deal`}</div>;
  }
}
