"use client";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { MultiSelect } from "@/components/ui/mutiSelect";
import { services } from "./filterForm";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import axios from "axios";
import { BASE_URL } from "@/lib/utils";
import { useUser } from "@/context/userContext";
const formSchema = z.object({
  serviceToRender: z.array(z.string()),
  status: z.string(),
  statusTag: z.string(),
  currentPitchRequest: z.string(),
  netTotalCost: z.number(),
  profit: z.number(),
  awarded: z.boolean(),
});
type CreateDealProps = {
  deal_id: number;
  status: string;
  department: string;
  servicesToRender: string[];
  netTotalCost: number;
  profit: number;
  CurrentPitchRequest: string;
  awarded: boolean;
};
export function UpdateDealForm(props: CreateDealProps) {
  const {usr} = useUser()
  //1. Define your form
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      serviceToRender: props.servicesToRender,
      status: props.status,
      statusTag: props.department,
      currentPitchRequest: props.CurrentPitchRequest,
      netTotalCost: props.netTotalCost,
      profit: props.profit,
      awarded: props.awarded,
    },
  });

  //2. Define a submit handler
  async function onSubmit(values: z.infer<typeof formSchema>) {

    try {
      const resp  = await axios({
        method:"put",
        baseURL: BASE_URL,
        url: "/admin/deals/update",
        data: {
          id: props.deal_id,
          services_to_render: values.serviceToRender,
          status: values.status,
          department: values.statusTag,
          current_pitch_request: values.currentPitchRequest,
          net_total_cost: values.netTotalCost,
          profit: values.profit,
          awarded: values.awarded
        },
        headers: {
          Authorization: `Bearer ${usr?.access_token}`
        }
      })
      const updatedDeal = await resp.data
      console.log(updatedDeal, "the updated deal")
    } catch (error){
      console.log(error)
    }
  }
  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        {/* servicesToRender */}
        <FormField
          control={form.control}
          name="serviceToRender"
          render={({ field }) => (
            <FormItem>
              <FormLabel>
                {"Services to render (no edit if not new service is requested)"}
              </FormLabel>
              <FormControl>
                <MultiSelect
                  options={[
                    ...services,
                    ...field.value.map((service) => {
                      return { label: service, value: service };
                    }),
                  ]}
                  value={field.value || []}
                  onChange={(newValue) => field.onChange(newValue)}
                  placeholder="Select services"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        {/* status */}
        <FormField
          control={form.control}
          name="status"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Current deal status</FormLabel>
              <FormControl>
                <Input placeholder="ongoing or pending or closed" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        {/* statusTag */}
        <FormField
          control={form.control}
          name="statusTag"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Current Department</FormLabel>
              <FormControl>
                <Input placeholder="pre-sales/sales/manager" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <div className="flex justify-around gap-3 w-full flex-grow">
          {/* netTotalCost */}
          <FormField
            control={form.control}
            name="netTotalCost"
            render={({ field }) => (
              <FormItem className="flex flex-col grow">
                <FormLabel>Net total cost</FormLabel>
                <FormControl>
                  <Input placeholder="expenditure" {...field}/>
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          {/* profit */}
          <FormField
            control={form.control}
            name="profit"
            render={({ field }) => (
              <FormItem className="flex flex-col grow">
                <FormLabel>Profit</FormLabel>
                <FormControl>
                  <Input placeholder="profit" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>
        {/* currentPitchRequest */}
        <FormField
          control={form.control}
          name="currentPitchRequest"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{"Current pitch request"}</FormLabel>
              <FormControl>
                <Input
                  placeholder="what's the sales-rep asking for"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="awarded"
          render={({ field }) => (
            <FormItem className="flex gap-4 items-end">
              <FormLabel>Contract awarded</FormLabel>
              <FormControl>
                <Checkbox
                  checked={field.value}
                  onCheckedChange={field.onChange}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Update</Button>
      </form>
    </Form>
  );
}

///Creating the Create Deals Form
