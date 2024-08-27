import { DealCard } from "@/components/deals/card";
import Welcome from "@/components/main/welcome";
import { revalidatePath } from "next/cache";

import { Suspense } from "react";

export default function Page() {
  revalidatePath("/dashboard")
  return (
    <div className="flex flex-col p-2 h-full w-full xl:w-11/12 relative text-sm">
      <div className="flex flex-col p-3 bg-slate-50 gap-4 sm:text-base md:text-lg text-slate-700">
        <div>
          <Suspense>
            <Welcome />
          </Suspense>
        </div>
        <p className="text-slate-600">Ongoing Deals</p>
      </div>
      <div className="flex flex-col gap-3  p-3 w-full grow h-auto">
        <section className="flex flex-col w-full h-full gap-3 ">
          <DealCard />
          <DealCard />
          <DealCard />
          <DealCard />
        </section>
      </div>
    </div>
  );
}
