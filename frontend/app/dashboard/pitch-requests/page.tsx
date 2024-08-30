import { PitchReqCardSection } from "@/components/pitchRequests/prCardSection";
export default async function Page() {
  return (
    <div className="flex flex-col text-sm md:text-base p-2 h-full w-full xl:w-11/12 relative">
      <div className="p-3">Pitch Requests</div>
      <div className="px-3 py-1 text-green-300">
        Oh great, no pending requests
      </div>
      <div className="flex flex-col gap-3  p-3 w-full grow h-auto">
        <PitchReqCardSection />
      </div>
    </div>
  );
}
