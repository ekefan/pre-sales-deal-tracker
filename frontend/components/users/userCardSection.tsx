"use client";
import { useEffect, useState } from "react";
import { UserCard } from "./card";
import { useUser } from "@/context/userContext";
import axios from "axios";
import { BASE_URL, DEFAULT_PAGE_SIZE } from "@/lib/utils";

export function UserCardSection() {
  const { usr } = useUser();
  const [usrs, setUsers] = useState<User[]>([]);

  useEffect(() => {
    async function getUsers(
      token: string | undefined,
      param: UserParam,
      url: string
    ) {
      try {
        if (!token) {
          throw new Error("user not signed in, no access token");
        }
        const resp = await axios({
          method: "get",
          baseURL: BASE_URL,
          url: url,
          params: param,
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        const users: User[] = await resp.data;
        setUsers(users);
      } catch (error) {
        console.log(error);
      }
    }

    getUsers(
      usr?.access_token,
      { page_id: 1, page_size: DEFAULT_PAGE_SIZE },
      "/a/users"
    );
  }, [usr]);
  return (
    <section className="flex flex-col w-full h-full gap-3 ">
      {usrs.length === 0 ? (
        <>Create your first User</>
      ) : (
        usrs.map((user) => {
          return (
            <UserCard
              key={user.username}
              userId={user.user_id}
              username={user.username}
              fullname={user.fullname}
              email={user.email}
            />
          );
        })
      )}
    </section>
  );
}
