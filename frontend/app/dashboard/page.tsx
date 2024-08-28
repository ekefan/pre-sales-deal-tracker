import { CardSection } from "@/components/deals/cardSection";
import Welcome from "@/components/main/welcome";
import { DealProvider } from "@/context/dealContext";
import { revalidatePath } from "next/cache";

export default async function Page() {
  revalidatePath("/dashboard");
  return (
    <div className="flex flex-col p-2 h-full w-full xl:w-11/12 relative text-sm">
      <div className="flex flex-col p-3 bg-slate-50 gap-4 sm:text-base md:text-lg text-slate-700">
        <div>
          <Welcome />
        </div>
        <p className="text-slate-600">Ongoing Deals</p>
      </div>
      <DealProvider>
        <CardSection page="home" />
      </DealProvider>
    </div>
  );
}
