import { cn } from "@/lib/utils";
import Link from "next/link";

interface PrProps {
  props: PitchReq
}
export function PitchReqCard({props}: PrProps) {
  const date = new Date(props.request_deadline);

  // Extract month, day, and last two digits of the year
  const month = ("0" + (date.getMonth() + 1)).slice(-2);
  const day = ("0" + date.getDate()).slice(-2);
  const year = date.getFullYear().toString().slice(-2);
  const formattedDate = `${month}/${day}/${year}`; // "12/24/24"
  const duration =  Math.floor((Date.now() - Date.parse(props.created_at))/(60 * 60 * 1000 * 24))

  return (
    <div className={cn(props.status == "ongoing" ? "border-green-200": "border-red-200", "bg-gray-100 w-full min-h-20 h-auto rounded-lg border p-2  gap-1 text-sm xl:text-base flex")}>
      <div className="flex flex-grow flex-col sm:flex-row gap-1">
        <div className="rounded border flex justify-center items-center border-slate-300 p-2 w-full h-full">
          {props.customer_name}
        </div>
        <div className="rounded border border-slate-300 p-1 w-full h-full flex justify-center items-center">
          {props.customer_request.join(", ")}
        </div>
      </div>
      <div className="flex flex-grow flex-col sm:flex-row gap-1 ">
        <div className={cn(props.admin_viewed == true ? "bg-green-300": "bg-red-100","rounded border border-slate-300 p-1 flex justify-center items-center w-full h-full") }>
          Admin viewed
        </div>
        <p className="rounded border border-slate-300 p-1 flex justify-center items-center w-full h-full">
          {duration > 0 ? `days:${duration}` : "new"}
        </p>
      </div>
      <div className="flex flex-grow flex-col sm:flex-row gap-1 ">
        <p className="rounded border flex justify-center items-center border-slate-300 p-1 w-full h-full">
          {formattedDate}
        </p>
        <Link
          href={`/dashboard/pitch-requests/update-pr?${props.pitch_id}`}
          className="bg-slate-300 w-full h-full rounded flex flex-col items-center justify-center border border-slate-300"
        >
          <p className="text-xs p-1">Update</p>
        </Link>

        <button
          className="bg-sky-100 w-full h-full rounded flex flex-col items-center justify-center border border-slate-300 p-1"
        >
          <p className="text-xs">create</p>
          <p className="text-xs">deal</p>
          </button>
      </div>
    </div>
  );
}
