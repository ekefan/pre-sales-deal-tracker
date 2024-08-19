"use client";

import Link from "next/link";
import clsx from "clsx";
import { usePathname } from "next/navigation";
import {
  House,
  LucideIcon,
  UsersRound,
  Handshake,
  GitPullRequestArrow,
  User,
  LogOut,
} from "lucide-react";

const [AdminRole, SalesRole, MangerRole] = ["admin", "sales", "manager"];
type Links = {
  name: string;
  href: string;
  icon: LucideIcon;
};

const adminLinks: Links[] = [
  { name: "Home", href: "/dashboard", icon: House },
  { name: "Deals", href: "/dashboard/deals", icon: Handshake },
  { name: "Users", href: "/dashboard/users", icon: UsersRound },
  { name: "PRs", href: "/dashboard/pitch-requests", icon: GitPullRequestArrow },
];
const salesLinks: Links[] = [
  { name: "Home", href: "/dashboard", icon: House },
  { name: "Deals", href: "/dashboard/deals", icon: Handshake },
  { name: "PRs", href: "/dashboard/pitch-requests", icon: GitPullRequestArrow },
];
const managerLinks: Links[] = [
  { name: "Home", href: "/dashboard", icon: House },
  { name: "Deals", href: "/dashboard/deals", icon: Handshake },
];

type Prop = {
  user: string;
};
export default function NavLinks({ user }: Prop) {
  let links: Links[];
  const pathname = usePathname();
  switch (user) {
    case MangerRole:
      links = managerLinks;
      break;
    case SalesRole:
      links = salesLinks;
      break;
    default:
      links = adminLinks;
  }

  return (
    <>
      {links.map((link) => {
        const LinkIcon = link.icon;
        return (
          <Link
            key={link.name}
            href={link.href}
            className={clsx(
              " w-full bg-gray-100 flex h-[32px] md:h-[48px] items-center justify-center md:justify-normal gap-2 rounded-md p-3 text-sm font-medium hover:bg-sky-100 hover:text-blue-500",
              {
                "bg-sky-100 text-blue-600": pathname === link.href,
              }
            )}
          >
            <LinkIcon className="size-4 md:size-6" />
            <p className="hidden md:block">{link.name}</p>
          </Link>
        );
      })}
      <div className="hidden h-auto w-full grow rounded-md bg-gray-100 md:block"></div>
      <Link
        key="profile"
        href="/dashboard/profile"
        className={clsx(
          " w-full bg-gray-100 flex h-[32px] md:h-[48px] items-center justify-center md:justify-normal gap-2 rounded-md p-3 text-sm font-medium hover:bg-sky-100 hover:text-blue-500",
          {
            "bg-sky-100 text-blue-600": pathname === "/dashboard/profile",
          }
        )}
      >
        <User className="size-4 md:size-6" />
        <p className="hidden md:block">Profile</p>
      </Link>
      <form className="w-full bg-gray-100 flex h-[32px] md:h-[48px] items-center justify-center md:justify-normal gap-2 rounded-md p-3 text-sm font-medium hover:bg-green-100 hover:text-green-400">
        <LogOut className="size-4 md:size-6" />
        <button className="hidden md:flex  items-start grow pl-1">
          Log out
        </button>
      </form>
    </>
  );
}
