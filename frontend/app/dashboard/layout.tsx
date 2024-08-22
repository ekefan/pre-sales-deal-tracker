'use client';
import SideNav from "@/components/main/sidenav";
import { UserProvider } from "@/context/userContext";

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <UserProvider>
      <div className="bg-slate-200 flex w-full h-screen md:w-11/12 md:rounded-r-2xl items-center p-2">
        <div className="bg-slate-100 w-full h-full rounded md:absolute md:right-0 md:h-5/6 md:rounded-l-2xl md:w-11/12 flex gap-2 p-3 flex-col md:flex-row">
          <div className="md:rounded-2xl md:w-1/6 flex flex-col gap-2">
            <SideNav />
          </div>
          <main className="bg-slate-50 rounded flex-grow px-2 pb-2 overflow-y-auto">
            {children}
          </main>
        </div>
      </div>
    </UserProvider>
  );
}
