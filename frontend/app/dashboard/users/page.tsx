import { UserCard } from "@/components/users/card";
import { Plus, User } from "lucide-react";
import Link from "next/link";
type usr = {
  fullname: string;
  username: string;
  email: string;
  userId: number;
};
export default async function Page() {
  const usrs: usr[] = [
    {
      fullname: "Joshua Vastech",
      username: "josh2funny",
      email: "eebenezer949@gmail.com",
      userId:5, 
    },
    {
      fullname: "Josh Vastech",
      username: "joshds2funny",
      email: "eebeneze49@gmail.com",
      userId: 4,
    },
    {
      fullname: "JoVastech",
      username: "jo2funny",
      email: "eebr949@gmail.com",
      userId: 3,
    },
    {
      fullname: "Joshua Vatech",
      username: "josnny",
      email: "eebenezer@gmail.com",
      userId: 1,
    },
  ];
 
  return (
    <div className="flex text-sm md:text-base flex-col p-2 h-full w-full xl:w-11/12 relative">
      <div className="p-3">Users</div>
      <div className="flex px-3 py-1 justify-end">
        <Link
          href="/dashboard/users/create"
          className="bg-slate-800 rounded-lg p-2 text-sm md:text-base text-white flex items-center gap-2 hover:bg-sky-200 hover:text-slate-900"
        >
          <p>create user</p>
          <Plus size={20} />
        </Link>
      </div>
      <div className="flex flex-col gap-3  p-3 w-full grow h-auto">
        <section className="flex flex-col w-full h-full gap-3 ">
          {usrs.map((usr) => {
            return (
              <UserCard
                key={usr.username}
                userId={usr.userId}
                username={usr.username}
                fullname={usr.fullname}
                email={usr.email}
              />
            );
          })}
        </section>
      </div>
    </div>
  );
}
