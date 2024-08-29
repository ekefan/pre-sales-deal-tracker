import { NextResponse } from 'next/server';
import axios from 'axios';
import { BASE_URL } from '@/lib/utils';

export async function POST(request: Request) {
  try {
    const { username, password } = await request.json();
        const response = await axios({
        method: "post",
        baseURL: BASE_URL,
        url:"/users/login",
        data: {
        username: username,
        password: password,
        }
    })
    return NextResponse.json(response.data);
  } catch (error: any) {
    if (error.response) {
        return NextResponse.json(error.response.data, {status: error.response.status });
    }
    else {
        console.log("unexpected Error", error);
        return NextResponse.json({ error: "unexpected error occurred"}, {status:500})
    }
  }
}