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
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
const formSchema = z.object({
  serviceToRender: z.array(z.string()),
  status: z.string(),
  statusTag: z.string(),
  currentPitchRequest: z.string(),
  netTotalCost: z.number(),
  profit: z.number(),
});
type CreateDealProps = {
  status: string;
  statusTag: string;
  servicesToRender: string[];
  netTotalCost: number;
  profit: number;
  CurrentPitchRequest: string; //eg. costing
};
export function UpdateDealForm(props: CreateDealProps) {
  //1. Define your form
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      serviceToRender: props.servicesToRender,
      status: props.status,
      statusTag: props.statusTag,
      currentPitchRequest: props.CurrentPitchRequest,
      netTotalCost: props.netTotalCost,
      profit: props.profit,
    },
  });

  //2. Define a submit handler
  function onSubmit(values: z.infer<typeof formSchema>) {
    // Do something with the form values.
    // ✅ This will be type-safe and validated.
    console.log(values);
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

        {/* netTotalCost */}
        <FormField
          control={form.control}
          name="netTotalCost"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Net total cost</FormLabel>
              <FormControl>
                <Input placeholder="expenditure" {...field} />
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
            <FormItem>
              <FormLabel>Profit</FormLabel>
              <FormControl>
                <Input placeholder="profit" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
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
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  );
}

///Creating the Create Deals Form
