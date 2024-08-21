import { CreateUserForm } from "@/components/users/createUserForm";
export default function Home() {
  return (
    <div className="flex  bg-slate-50 text-sm md:text-base font-normal flex-col p-2 h-auto w-full sm:w-5/6 md:w-10/12 lg:w-8/12 xl:w-6/12 relative">
      <div className="p-3">Create User</div>
      <div className="p-3 mt-1 flex flex-col gap-3 w-full rounded-md bg-gray-100">
        <CreateUserForm />
      </div>
    </div>
  );
}
