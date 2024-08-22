"use client";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { MultiSelect } from "@/components/ui/mutiSelect";

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
  username: z.string(),
  role: z.string(),
  fullname: z.string(),
  email: z.string(),
});

type UpdateUserProps = {
    username: string,
    fullname: string,
    email: string,
}

export function UpdateUserForm(user: UpdateUserProps) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: user.username,
      fullname:  user.fullname,
      email: user.email,
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
        {/* username */}
        <FormField
          control={form.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{"username"}</FormLabel>
              <FormControl>
                <Input placeholder="enter username" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        {/* fullname*/}
        <FormField
          control={form.control}
          name="fullname"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{"fullname"}</FormLabel>
              <FormControl>
                <Input placeholder="enter user full-name" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        {/* email*/}
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{"user email"}</FormLabel>
              <FormControl>
                <Input placeholder="example@vastech.com" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        {/* role
        <FormField
          control={form.control}
          name="role"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{"user role"}</FormLabel>
              <FormControl>
                <Input
                  placeholder="admin, sales or manager"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        /> */}
        <Button type="submit">Update</Button>
      </form>
    </Form>
  );
}
