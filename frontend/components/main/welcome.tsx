"use client";
import { useUser } from "@/context/userContext";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation"

export default function Welcome() {
  const newDate = new Date(Date.now());
  const dateString = newDate.toDateString();
  const router = useRouter()

  const { usr } = useUser();
  const [firstName, setFirstName] = useState("");

  useEffect(() => {
    if (!usr) {
        router.push("/")
    } else if (usr != null && usr.user != null && usr.user.fullname != null) {
      const firstname = usr.user.fullname.split(" ")[0];
      setFirstName(firstname);
    }
  }, [usr, router]);
  return (
    <div>
      <p className="font-medium">Welcome {firstName}</p>
      <p className="text-sm md:text-base">{dateString}</p>
    </div>
  );
}
