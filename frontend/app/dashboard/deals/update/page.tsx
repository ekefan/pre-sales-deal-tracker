import { UpdateDealForm } from "@/components/deals/updateDealForm";
import { UpdatePitchRequestSection } from "@/components/pitchRequests/updatePrSection";

export default function Home() {
  return (
    <div className="flex  bg-slate-50 text-sm md:text-base font-normal flex-col p-2 h-auto w-full sm:w-5/6 md:w-10/12 lg:w-8/12 xl:w-6/12 relative">
      <div className="p-3">Update deal</div>
      <UpdatePitchRequestSection />
    </div>
  );
}
