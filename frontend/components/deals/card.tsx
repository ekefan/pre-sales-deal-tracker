import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/components/ui/hover-card";
import { cn } from "@/lib/utils";
import Link from "next/link";
import { ADMIN_ROLE } from "@/lib/utils";
import { useRouter } from "next/navigation";
import { useUser } from "@/context/userContext";
import { revalidatePath } from "next/cache";
interface DealProp {
  deal: Deal;
}
export function DealCard({ deal }: DealProp) {
  const router = useRouter();
  const {usr} = useUser();
  const margin =
    Number(deal.profit) != 0
      ? String(
          (parseFloat(deal.profit) * 100) / parseFloat(deal.net_total_cost)
        )
      : "0.00";
  const duration =  Math.floor((Date.now() - Date.parse(deal.created_at))/(60 * 60 * 1000 * 24))
  

  const onSubmit = (userRole: string | undefined, deal_id: number) => {
    if (userRole !== ADMIN_ROLE) {
      revalidatePath("/dashboard/deals/")
    }

    const url = `/dashboard/deals/update?id=${deal_id}`
    router.push(url)
  }
  return (
    <div
      className={cn(
        `bg-gray-100 w-full min-h-20 h-auto rounded-lg border p-2  gap-1 text-sm xl:text-base flex`,
        deal.status === "ongoing"
          ? "border-red-200"
          : deal.awarded === true
          ? "border-yellow-200"
          : "border-red-200"
      )}
    >
      <div className="flex flex-grow flex-col sm:flex-row gap-1">
        <div className="rounded border flex justify-center items-center border-slate-300 p-2 w-full h-full">
          {deal.customer_name}
        </div>
        <div className="rounded border border-slate-300 p-1 w-full h-full flex justify-center items-center">
          <HoverCard>
            <HoverCardTrigger className="hover:cursor-pointer">
              deal details
            </HoverCardTrigger>

            <HoverCardContent>
              <div className="flex flex-col gap-2">
                <p className="rounded border border-slate-300 p-1 w-full h-full">
                  brought by:{deal.sales_rep_name}
                </p>
                <p className="rounded border border-slate-300 p-1 w-full h-full">
                  Services:&nbsp;
                  {deal.services_to_render.join(", ")}
                </p>
              </div>
            </HoverCardContent>
          </HoverCard>
        </div>
      </div>
      <div className="flex flex-grow flex-col sm:flex-row gap-1 ">
        <div className="rounded border border-slate-300 p-1 flex justify-center items-center w-full h-full">
          <HoverCard>
            <HoverCardTrigger className="hover:cursor-pointer">
              deal finance
            </HoverCardTrigger>
            <HoverCardContent>
              <div className="flex flex-col gap-2">
                <p className="rounded p-1 border border-slate-300 block">
                  Margin: {margin}
                </p>
                <p className="rounded p-1 border border-slate-300 block">
                  Net cost: {deal.net_total_cost}
                </p>
                <p className="rounded p-1 border border-slate-300 block">
                  Profit: {deal.profit}
                </p>
              </div>
            </HoverCardContent>
          </HoverCard>
        </div>
        <p className="rounded border border-slate-300 p-1 flex justify-center items-center w-full h-full">
          {`days: ${duration}`}
        </p>
      </div>
      <div className="flex flex-grow flex-col sm:flex-row gap-1 ">
        <p className="rounded border flex justify-center items-center border-slate-300 p-1 w-full h-full">
          {deal.department}
        </p>

        <button onClick={() => onSubmit(usr?.user?.role, deal.deal_id)} className="bg-green-100 w-full h-full rounded flex flex-col items-center justify-center border border-slate-300">
          {/* <Link
            href="/dashboard/deals/update"
            className=
          > */}
            <p>Update</p>
            <p className="text-xs text-gray-500">admin only</p>
          {/* </Link> */}
        </button>
      </div>
    </div>
  );
}
