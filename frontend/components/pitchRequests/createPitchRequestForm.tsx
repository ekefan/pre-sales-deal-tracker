"use client";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { MultiSelect } from "@/components/ui/mutiSelect";
import { services } from "@/components/deals/filterForm";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";

const formSchema = z.object({
  salesRepId: z.number(),
  salesRepName: z.string(),
  status: z.string(),
  customerName: z.string(),
  pitchRequest: z.string(), //pitchTag
  customerRequests: z.array(z.string()),
  newRequest: z.string(),
  requestDeadline: z.string(), //convert to date please
});

type PitchRequestProps = {
  salesRepId: number;
  salesRepName: string;
  status: string;
};
export function CreatePitchRequestForm(pr: PitchRequestProps) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      salesRepId: pr.salesRepId,
      salesRepName: pr.salesRepName,
      status: pr.status,
      newRequest: "",
    },
  });

  function onSubmit(values: z.infer<typeof formSchema>) {
    // Do something with the form values.
    // ✅ This will be type-safe and validated.
    console.log(values);
  }
  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        {/* salesRepId */}
        <FormField
          control={form.control}
          name="salesRepId"
          render={({ field }) => (
            <FormItem className="hidden">
              <FormLabel>{"username"}</FormLabel>
              <FormControl>
                <Input type="number" placeholder="sales rep id" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* sales rep name*/}
        <FormField
          control={form.control}
          name="salesRepName"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{"Full Name"}</FormLabel>
              <FormControl>
                <Input placeholder="Enter your full name" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* status*/}
        <FormField
          control={form.control}
          name="status"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{"Pitch request Status"}</FormLabel>
              <FormControl>
                <Input placeholder="ongoing by default" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/*customer name*/}
        <FormField
          control={form.control}
          name="customerName"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{"Customer name"}</FormLabel>
              <FormControl>
                <Input
                  placeholder="Enter fullname or General Abbreviation"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        {/*pitch request*/}
        <FormField
          control={form.control}
          name="pitchRequest"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{"pitch request"}</FormLabel>
              <FormControl>
                <Input
                  placeholder="costing/proposal/presentation..."
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        {/*ccustomerRequests*/}
        <FormField
          control={form.control}
          name="customerRequests"
          render={({ field }) => (
            <FormItem className="w-full">
              <FormLabel>{"Customer Request"}</FormLabel>
              <MultiSelect
                options={services}
                value={field.value || []}
                onChange={(newValue) => field.onChange(newValue)}
                placeholder="Select services we offer"
              />
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="newRequest"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{"Customer request not on services list?"}</FormLabel>
              <FormControl>
                <Input placeholder="input custom request" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        {/*Request Deadline*/}
        <FormField
          control={form.control}
          name="requestDeadline"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{"When is the request due"}</FormLabel>
              <FormControl>
                <Input placeholder="day/month/year" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type="submit">create p-r</Button>
      </form>
    </Form>
  );
}
