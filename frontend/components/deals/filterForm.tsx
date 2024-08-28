"use client";
import { Button } from "@/components/ui/button";
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
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { MultiSelect } from "@/components/ui/mutiSelect";
import { useDeal } from "@/context/dealContext";
import { DEFAULT_PAGE_SIZE } from "@/lib/utils";

export const services = [
  { label: "Service 1", value: "service1" },
  { label: "Service 2", value: "service2" },
  { label: "Service 3", value: "service3" },
  { label: "Service 4", value: "service4" },
  { label: "Service 5", value: "service5" },
  { label: "goat meat", value: "goat" },
  { label: "beans", value: "beans" },
  { label: "rice", value: "rice" },

  // Add more services as needed
];

const FormSchema = z.object({
  customerName: z.string().nullable(),
  serviceToRender: z.array(z.string()).optional().default([]),
  status: z.enum(["ongoing", "closed"]).nullable(),
  maxProfit: z
    .number()
    .gt(-1, { message: "please enter a profit greater than 0" })
    .nullable(),
  minProfit: z
    .number()
    .gt(-1, { message: "please enter a profit greater than 0" })
    .nullable(),
  awarded: z.boolean().nullable(),
  salesRepName: z.string().nullable(),
});

/**
 * Processes form data used to filter deals requested by user
 *
 * @param {object} values - The form data object... based on formSchema
 * @param {string} values.customerName - The name of the customer from a list of customers
 * @param {string} values.servicesToRender - The service rendered for the deal
 * @param {"ongoing" | "closed"} value.status - The current status of the deal
 * @param {number} values.maxProfit - Deals where profit greater than or equal to this number
 * @param {number} values.minProfit - Deals where profit less than or equal to this number
 * @param {boolean} values.awarded - Deals where the contract was awarded
 * @param {salesRepName} values.salesRepName - Name of the sales rep who brought in the deal
 */
export function FilterDealsForm() {
  const { setDealParam } = useDeal();
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      customerName: null,
      serviceToRender: [], // Default to an empty array
      status: null,
      maxProfit: null,
      minProfit: null,
      awarded: null,
      salesRepName: null,
    },
  });

  function onSubmit(values: z.infer<typeof FormSchema>) {
    console.log(values);
    const dealfilterParams: DealFilter = {
      customer_name: values?.customerName,
      service_to_render: values?.serviceToRender,
      status: values?.status,
      max_profit: values?.maxProfit != null ? String(values?.maxProfit) : null,
      min_profit: values?.minProfit != null ? String(values?.minProfit) : null,
      awarded: values?.awarded,
      sales_rep_name: values?.salesRepName,
      page_size: DEFAULT_PAGE_SIZE,
      page_id: 1,
    };

    /// set deal params here//////////////////////////////
    setDealParam(dealfilterParams);
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="border p-2 rounded-md flex flex-col gap-2"
      >
        <FormField
          control={form.control}
          name="customerName"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <Input
                  placeholder="Customer Name"
                  aria-describedby="customer name error"
                  {...field}
                  value={field.value || ""} // Default to empty string
                  onChange={(e) => field.onChange(e.target.value || null)} // Convert empty string to null
                />
              </FormControl>
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="serviceToRender"
          render={({ field }) => (
            <FormItem className="w-full">
              <MultiSelect
                options={services}
                value={field.value || []}
                onChange={(newValue) => field.onChange(newValue)}
                placeholder="Select services"
              />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="maxProfit"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <Input
                  type="number"
                  placeholder="Max Profit"
                  {...field}
                  aria-describedby="profit error"
                  value={field.value !== null ? String(field.value) : ""} // Handle number values
                  onChange={(e) =>
                    field.onChange(
                      e.target.value ? Number(e.target.value) : null
                    )
                  } // Convert empty string to null
                />
              </FormControl>
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="minProfit"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <Input
                  type="number"
                  placeholder="Min Profit"
                  {...field}
                  value={field.value !== null ? String(field.value) : ""} // Handle number values
                  onChange={(e) =>
                    field.onChange(
                      e.target.value ? Number(e.target.value) : null
                    )
                  } // Convert empty string to null
                />
              </FormControl>
            </FormItem>
          )}
        />

        <div className="flex gap-2">
          <FormField
            control={form.control}
            name="awarded"
            render={({ field }) => (
              <FormItem className="flex flex-row items-center space-x-3 text-slate-600 space-y-0 rounded-md border p-2 w-2/6">
                <div className="space-y-1 leading-none">
                  <FormLabel>Awarded</FormLabel>
                  <FormDescription>
                    (tick if contract was awarded)
                  </FormDescription>
                </div>
                <FormControl>
                  <Checkbox
                    checked={field.value || false}
                    onCheckedChange={(checked) => field.onChange(checked)}
                  />
                </FormControl>
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="status"
            render={({ field }) => (
              <FormItem className="w-4/6">
                <Select
                  onValueChange={(value) => field.onChange(value || null)}
                  value={field.value || ""}
                >
                  <FormControl>
                    <SelectTrigger className="w-full h-full text-sm sm:text-base text-slate-600">
                      <SelectValue placeholder="select deal status" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value="ongoing">ongoing</SelectItem>
                    <SelectItem value="closed">closed</SelectItem>
                  </SelectContent>
                </Select>
              </FormItem>
            )}
          />
        </div>

        <FormField
          control={form.control}
          name="salesRepName"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <Input
                  placeholder="Sales Rep Name"
                  {...field}
                  value={field.value || ""} // Default to empty string
                  onChange={(e) => field.onChange(e.target.value || null)} // Convert empty string to null
                />
              </FormControl>
            </FormItem>
          )}
        />
        <Button
          className="mt-4 w-1/6 bg-slate-800 hover:bg-slate-900"
          type="submit"
        >
          Filter
        </Button>
      </form>
    </Form>
  );
}
