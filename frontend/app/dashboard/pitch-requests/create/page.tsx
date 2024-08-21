import { CreatePitchRequestForm } from "@/components/pitchRequests/createPitchRequestForm";
export default function Home() {
    //get user
    const salesRepReq = {id: 10, name: "Lily Joshua", status: "ongoing"}
  return (
    <div className="flex  bg-slate-50 text-sm md:text-base font-normal flex-col p-2 h-auto w-full sm:w-5/6 md:w-10/12 lg:w-8/12 xl:w-6/12 relative">
      <div className="p-3">Create pitch-request</div>
      <div className="p-3">
        <div className="p-2 flex flex-col gap-3 w-full rounded-md bg-gray-100">
          <CreatePitchRequestForm  salesRepId={salesRepReq.id} salesRepName={salesRepReq.name} status={salesRepReq.status}/>
        </div>
      </div>
    </div>
  );
}