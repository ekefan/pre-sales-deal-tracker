import { FilterDealsForm } from "@/components/deals/filterForm";
import { SlidersHorizontal } from "lucide-react";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

export default function Page() {
  return (
    <>
      <div className="flex flex-col text-sm md:text-base p-2 h-full w-full xl:w-11/12">
        <div className="p-3">
          Explore Deals
        </div>
        <div className="flex justify-end px-3 py-1">
          <Popover>
            <PopoverTrigger className="w-20 sm:w-24 md:w-28 lg:w-32">
              <div className="flex border p-2 justify-between gap-2 rounded-lg items-center bg-slate-800 hover:bg-sky-200 hover:text-slate-900 text-white">
                <div className=" ">
                  Filter
                </div>
                <SlidersHorizontal size={18}/>
              </div>
            </PopoverTrigger>
            <PopoverContent className="w-full">
              <FilterDealsForm />
            </PopoverContent>
          </Popover>
        </div>
        <div className="flex flex-col gap-3  p-3 w-full grow h-auto">
        <section className="flex flex-col w-full h-full gap-3 ">
          <div className="bg-green-200 w-full h-32 rounded-lg border"></div>
          <div className="bg-yellow-200 w-full h-32 rounded-lg border"></div>
          <div className="bg-indigo-300 w-full h-32 rounded-lg border"></div>
          <div className="bg-pink-300 w-full h-32 rounded-lg border"></div>
        </section>
      </div>
      </div>
    </>
  );
}
