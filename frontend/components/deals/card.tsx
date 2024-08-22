import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/components/ui/hover-card";
import Link from "next/link";

export function DealCard() {
  return (
    <div className="bg-gray-100 w-full min-h-20 h-auto rounded-lg border p-2  gap-1 text-sm xl:text-base flex">
      <div className="flex flex-grow flex-col sm:flex-row gap-1">
        <div className="rounded border flex justify-center items-center border-slate-300 p-2 w-full h-full">
          {"<customer name>"}
        </div>
        <div className="rounded border border-slate-300 p-1 w-full h-full flex justify-center items-center">
          <HoverCard>
            <HoverCardTrigger className="hover:cursor-pointer">Customer details</HoverCardTrigger>

            <HoverCardContent>
              <div className="flex flex-col gap-2">
                <p className="rounded border border-slate-300 p-1 w-full h-full">
                  brought by:{" <salesrep name>"}
                </p>
                <p className="rounded border border-slate-300 p-1 w-full h-full">
                  Services:
                  {" service1, service2"}
                </p>
              </div>
            </HoverCardContent>
          </HoverCard>
        </div>
      </div>
      <div className="flex flex-grow flex-col sm:flex-row gap-1 ">
        <div className="rounded border border-slate-300 p-1 flex justify-center items-center w-full h-full">
          <HoverCard>
            <HoverCardTrigger className="hover:cursor-pointer">deal finance</HoverCardTrigger>
            <HoverCardContent>
              <div className="flex flex-col gap-2">
                <p className="rounded p-1 border border-slate-300 block">
                  Margin: {"<margin>"}
                </p>
                <p className="rounded p-1 border border-slate-300 block">
                  Net cost: {"<net cost>"}
                </p>
                <p className="rounded p-1 border border-slate-300 block">
                  Profit: {"<profit>"}
                </p>
              </div>
            </HoverCardContent>
          </HoverCard>
        </div>
        <p className="rounded border border-slate-300 p-1 flex justify-center items-center w-full h-full">
          {"5 days"}
        </p>
      </div>
      <div className="flex flex-grow flex-col sm:flex-row gap-1 ">
        <p className="rounded border flex justify-center items-center border-slate-300 p-1 w-full h-full">
          {"<pre-sales>"}
        </p>

        <Link
          href="/dashboard/deals/update"
          className="bg-green-100 w-full h-full rounded flex flex-col items-center justify-center border border-slate-300"
        >
          <p>Update</p>
          <p className="text-xs text-gray-500">admin only</p>
        </Link>
      </div>
    </div>
  );
}

