declare interface LoginResp {
    user_id: string;
    username: string;
    role: string;
    fullname: string;
    email: string;
    updatedAt: number;
    createdAt: number;
  }

  declare interface User {
    access_token: string;
    user: LoginResp
}


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