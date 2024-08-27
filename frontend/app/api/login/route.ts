// app/api/login/route.ts
import { NextResponse } from 'next/server';
import axios from 'axios';

export async function POST(request: Request) {
  try {
    const { username, password } = await request.json();
    console.log("program got here too", username, password)
    const response = await axios({
        method: "post",
        baseURL: "http://localhost:8080",
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