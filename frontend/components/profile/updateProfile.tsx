"use client";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";

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
  userId: z.number(),
  username: z.string(),
  fullname: z.string(),
  email: z.string(),
});

type UpdateProfileProps = {
  userId: number;
  fullname: string;
  username: string;
  email: string;
};
export function UpdateProfileForm(pr: UpdateProfileProps) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      userId: pr.userId,
      username: pr.username,
      fullname: pr.fullname,
      email: pr.email,
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
        {/* userId */}
        <FormField
          control={form.control}
          name="userId"
          render={({ field }) => (
            <FormItem className="hidden">
              <FormLabel>{"user ID"}</FormLabel>
              <FormControl>
                <Input type="number" placeholder="user id" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* full name*/}
        <FormField
          control={form.control}
          name="fullname"
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

        {/* username*/}
        <FormField
          control={form.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{"Username"}</FormLabel>
              <FormControl>
                <Input placeholder="input new username" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* email */}
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{"Email"}</FormLabel>
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
     
        <Button type="submit">Update</Button>
      </form>
    </Form>
  );
}
