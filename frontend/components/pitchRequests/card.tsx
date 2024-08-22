import Link from "next/link";

export function PitchReqCard() {
  return (
    <div className="bg-gray-100 w-full min-h-20 h-auto rounded-lg border p-2  gap-1 text-sm xl:text-base flex">
      <div className="flex flex-grow flex-col sm:flex-row gap-1">
        <div className="rounded border flex justify-center items-center border-slate-300 p-2 w-full h-full">
          {"<customer name>"}
        </div>
        <div className="rounded border border-slate-300 p-1 w-full h-full flex justify-center items-center">
          {"service1, service2, service3"}
        </div>
      </div>
      <div className="flex flex-grow flex-col sm:flex-row gap-1 ">
        <div className="rounded border border-slate-300 p-1 flex justify-center items-center w-full h-full bg-red-100">
          request received
        </div>
        <p className="rounded border border-slate-300 p-1 flex justify-center items-center w-full h-full">
          {"5 days"}
        </p>
      </div>
      <div className="flex flex-grow flex-col sm:flex-row gap-1 ">
        <p className="rounded border flex justify-center items-center border-slate-300 p-1 w-full h-full">
          {"12/14/24"}
        </p>
        <Link
          href="/dashboard/deals/update"
          className="bg-slate-300 w-full h-full rounded flex flex-col items-center justify-center border border-slate-300"
        >
          <p className="text-xs p-1">Update</p>
        </Link>

        <Link
          href="/dashboard/deals/update"
          className="bg-sky-100 w-full h-full rounded flex flex-col items-center justify-center border border-slate-300 p-1"
        >
          <p className="text-xs">create</p>
          <p className="text-xs">deal</p>
        </Link>
      </div>
    </div>
  );
}
