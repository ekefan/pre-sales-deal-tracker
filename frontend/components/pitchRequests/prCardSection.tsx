"use client";

import { useEffect, useState } from "react";
import { useUser } from "@/context/userContext";
import axios from "axios";
import { BASE_URL, DEFAULT_PAGE_SIZE } from "@/lib/utils";
import { PitchReqCard } from "@/components/pitchRequests/card";

import Link from "next/link";
export function PitchReqCardSection() {
  const { usr } = useUser();
  const [pitchRequests, setPitchRequests] = useState<PitchReq[]>([]);

  useEffect(() => {
    async function getPitchRequests(
      token: string | undefined,
      url: string,
      param: SalesPitchReqParams
    ) {
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
        const pitchReqs: PitchReq[] = await resp.data;
        console.log(pitchReqs);
        setPitchRequests(pitchReqs);
      } catch (error) {
        console.log(error);
      }
    }
    if (usr?.user.role == "admin") {
      console.log("waiting to write api for that");
    } else if (usr?.user.role == "sales") {
      getPitchRequests(usr?.access_token, "a/sales/pitchrequest/", {
        sales_rep_id: usr?.user.user_id,
        page_id: 1,
        page_size: DEFAULT_PAGE_SIZE,
      });
    }
  }, [usr]);
  return (
    <section className="flex flex-col w-full h-full gap-3 ">
      {pitchRequests.map((pr) => {
        return <PitchReqCard key={pr.pitch_id} props={pr}/>
      })}
    </section>
  );
}

/// get pitchRequests for admin
/// get pitchRequests for sales
//// if usr.role not equal to admin or sales push to login
