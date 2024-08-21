import { UpdatePitchRequestForm } from "@/components/pitchRequests/updatePitchRequestForm";
export default function Home() {
  //get user
  const pitchReqDetails = {
    ID: 1,
    status: "ongoing",
    pitchRequest: "costing",
    customerRequests: ["service 1"],
    adminViewed: false,
    requestDeadline: "14/11/2024",
  };
  return (
    <div className="flex  bg-slate-50 text-sm md:text-base font-normal flex-col p-2 h-auto w-full sm:w-5/6 md:w-10/12 lg:w-8/12 xl:w-6/12 relative">
      <div className="p-3">update pitch-request</div>
      <div className="p-3">
        <div className="p-2 flex flex-col gap-3 w-full rounded-md bg-gray-100">
          <UpdatePitchRequestForm
            pitchId={pitchReqDetails.ID}
            status={pitchReqDetails.status}
            pitchRequest={pitchReqDetails.pitchRequest}
            customerRequests={pitchReqDetails.customerRequests}
            adminViewed={pitchReqDetails.adminViewed}
            requestDeadline={pitchReqDetails.requestDeadline}
          />
        </div>
      </div>
    </div>
  );
}
    //   pitchId: pr.pitchId,
    //   status: pr.status,
    //   pitchRequest: pr.pitchRequest,
    //   customerRequests: pr.customerRequests,
    //   adminViewed: pr.adminViewed,
    //   requestDeadline: pr.requestDeadline,