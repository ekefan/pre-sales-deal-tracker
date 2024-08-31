// LoginResp -- props for user in login Respo
declare interface LoginResp {
    user_id: number;
    username: string;
    role: string;
    fullname: string;
    email: string;
    updatedAt: number;
    createdAt: number;
  }

  // UserResp --- props or Login Response
  declare interface UserResp {
    access_token: string;
    user: LoginResp
}

//Deal --- props for Deal model
declare interface Deal {
  deal_id: number;
  pitch_id: {num: number; valid: boolean};
  sales_rep_name: string;
  customer_name: string;
  services_to_render: string[];
  status: string;
  department: string;
  current_pitch_request: string;
  net_total_cost: string;
  profit: string;
  created_at: string;
  updated_at: string;
  closed_at: stirng;
  awarded: boolean;
}

// DealFilter --- params needed to get filtered deals at /a/deals/filtered
declare interface DealFilter {
  customer_name: string | null;
  service_to_render: string[] | null;
  status: string | null;
  max_profit: string | null;
  min_profit: string | null;
  awarded: bool | null
  sales_rep_name: string | null;
  page_size: number;
  page_id: number;
}

declare interface Status {
  status: string;
}

//User --- props for user model
declare interface User {
  fullname: string;
  username: string;
  email: string;
  user_id: number;
}


// UserParams --- params needed to get list of users at /a/users
declare interface UserParam {
  page_id: number;
  page_size: number;
}

//PitchReq --- props for pitch request model
declare interface PitchReq {
  pitch_id: number;
  sales_rep_id: number;
  sales_rep_name: string;
  status: string;
  customer_name: string;
  department: string;
  customer_request: string[];
  request_deadline: string;
  admin_viewed: boolean;
  created_at: string;
  updated_at: string;
}
declare interface SalesPitchReqParams {
  sales_rep_id: number;
  page_id: number;
  page_size: number;
}

declare interface AdminPitchReqParams {
  admin_viewed: boolean;
}