"use client";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { useUser } from "@/context/userContext";

import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useSearchParams } from "next/navigation";
import axios from "axios";
import { BASE_URL } from "@/lib/utils";

const formSchema = z.object({
  username: z.string().min(1),
  fullname: z.string().min(1),
  email: z.string().email(),
});

export function UpdateUserForm() {
  const { usr } = useUser();
  const searchParams = useSearchParams()
  const user_id = searchParams.get("user_id")
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),

    defaultValues: {
      username: "",
      fullname: "",
      email: "",
    },
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    try {
      if (!usr?.access_token) {
        return
      }

      const resp = await axios({
        method: "put",
        baseURL: BASE_URL,
        url: "/a/users/update",
        headers: {
          Authorization: `Bearer ${usr.access_token}`
        },
        data: {
          user_id: user_id != null ? Number(user_id) : 0,
          fullname: values.fullname,
          username: values.username,
          email: values.email
        }
      })

      const updatedUsr = await resp.data
      console.log(updatedUsr)
    } catch(error) {
      console.log(error)
    }
    console.log(values, user_id);
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
                <Input placeholder="enter username" {...field} defaultValue={usr?.user.username} />
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
                <Input placeholder="enter user full-name" {...field}  defaultValue={usr?.user.fullname}/>
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
                <Input placeholder="example@vastech.com" {...field} defaultValue={usr?.user.email}/>
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
