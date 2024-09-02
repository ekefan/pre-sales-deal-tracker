import { ProfileSection } from "@/components/profile/profileSection";
export default async function Page() {
  return (
    <div className="flex text-sm md:text-base flex-col p-2 h-full w-full sm:w-5/6 md:w-10/12 lg:w-8/12 xl:w-6/12 relative">
      <div className="p-3">Profile</div>
      <ProfileSection/>
    </div>
  );
}
