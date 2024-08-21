import { UpdateDealForm } from "@/components/deals/updateDealForm";

export default function Home() {
  return (
    <div className="flex  bg-slate-50 text-sm md:text-base font-normal flex-col p-2 h-auto w-full sm:w-5/6 md:w-10/12 lg:w-8/12 xl:w-6/12 relative">
      <div className="p-3">Update deal</div>
      <div className="p-3 mt-1 flex flex-col gap-3 w-full border rounded-md bg-gray-100">
        <UpdateDealForm
          servicesToRender={["Service 1", "Service 2"]}
          CurrentPitchRequest="costing"
          status="ongoing"
          statusTag="pre-sales"
          profit={0}
          netTotalCost={0}
        />
      </div>
    </div>
  );
}
