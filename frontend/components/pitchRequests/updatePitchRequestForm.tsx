"use client";
import { boolean, z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { MultiSelect } from "@/components/ui/mutiSelect";
import { services } from "@/components/deals/filterForm";
import { Checkbox } from "@/components/ui/checkbox";
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
  pitchId: z.number(),
  status: z.string(),
  pitchRequest: z.string(), //pitchTag
  customerRequests: z.array(z.string()),
  newRequest: z.string(),
  adminViewed: z.boolean(),
  requestDeadline: z.string(), //convert to date please
});

type PitchRequestProps = {
  pitchId: number;
  status: string;
  pitchRequest: string;
  customerRequests: string[];
  adminViewed: boolean;
  requestDeadline: string;
};
export function UpdatePitchRequestForm(pr: PitchRequestProps) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      pitchId: pr.pitchId,
      status: pr.status,
      pitchRequest: pr.pitchRequest,
      customerRequests: pr.customerRequests,
      adminViewed: pr.adminViewed,
      requestDeadline: pr.requestDeadline,
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
        {/* pitchID*/}
        <FormField
          control={form.control}
          name="pitchId"
          render={({ field }) => (
            <FormItem className="hidden">
              <FormLabel>{"pitch id"}</FormLabel>
              <FormControl>
                <Input type="number" placeholder="sales rep id" {...field} />
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
              <FormLabel>{"Pitch Status"}</FormLabel>
              <FormControl>
                <Input placeholder="ongoing on closed" {...field} />
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
                options={[...services, ...field.value.map((service) => {
                    return { label: service, value: service };
                  }),]}
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
                <Input placeholder="input custom request (leave empty if non)" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/*admin viewed*/}
        <FormField
          control={form.control}
          name="adminViewed"
          render={({ field }) => (
            <FormItem className="">
              <div className="flex gap-4">
              <FormLabel>{"Admin Viewed"}</FormLabel>
              <FormControl>
              <Checkbox
                    checked={field.value || false}
                    onCheckedChange={(checked) => field.onChange(checked)}
                  />
              </FormControl>
              </div>
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
